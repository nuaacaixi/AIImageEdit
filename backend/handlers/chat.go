package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/vibeCoding/AIImageEdit/backend/ai"
	"github.com/vibeCoding/AIImageEdit/backend/models"
	"github.com/vibeCoding/AIImageEdit/backend/tools"
)

type ChatHandler struct {
	store       *models.Store
	storagePath string
	llmGateway  *ai.LLMGateway
	toolReg     *tools.Registry
	systemPrompt string
}

type ChatRequest struct {
	SessionID    string `json:"sessionId"`
	Message      string `json:"message"`
	BaseImageID  string `json:"baseImageId,omitempty"`
}

type ChatResponse struct {
	TurnID        string                `json:"turnId"`
	SessionID     string                `json:"sessionId"`
	Reasoning     string                `json:"reasoning"`
	Intent        string                `json:"intent"`
	ToolName      string                `json:"toolName"`
	ResultURL     string                `json:"resultUrl"`
	Status        string                `json:"status"`
	Message       string                `json:"message"`
}

func NewChatHandler(store *models.Store, storagePath string, llmGateway *ai.LLMGateway, toolReg *tools.Registry) *ChatHandler {
	sp := buildSystemPrompt(toolReg)
	return &ChatHandler{
		store:        store,
		storagePath:  storagePath,
		llmGateway:   llmGateway,
		toolReg:      toolReg,
		systemPrompt: sp,
	}
}

func buildSystemPrompt(reg *tools.Registry) string {
	return `你是一个专业的 AI 修图助手，名字叫"小修"。你可以使用以下工具帮助用户处理图片：

` + reg.SystemPromptSection() + `
当用户向你发送修图请求时，你需要：
1. 理解用户的中文（或英文）意图
2. 选择最合适的工具
3. 将用户指令改写为适合该工具的英文提示词（详细、具体、包含视觉描述）
4. 用温暖、简洁的中文告知用户正在执行什么操作（一句话即可，像"好的，正在为您调整天空的色调..."这样）

请严格按以下 JSON 格式回复（不要包含其他文字，不要用 markdown 包装）：
{"intent":"意图类型","toolName":"工具名称","rewrittenPrompt":"改写后的英文提示词","reasoning":"给用户的简短中文反馈"}`
}

func (h *ChatHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.SessionID == "" || req.Message == "" {
		http.Error(w, "sessionId and message are required", http.StatusBadRequest)
		return
	}

	// Get session
	sess, err := h.store.GetSession(req.SessionID)
	if err != nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	// Determine which image to use
	baseImageURL := sess.CurrentImageURL
	if baseImageURL == "" {
		baseImageURL = sess.OriginalImageURL
	}
	baseImageID := sess.CurrentImageID
	if baseImageID == "" {
		baseImageID = sess.OriginalImageID
	}

	// Add user message to context
	userMsgID := generateUUID()
	h.store.AddContextMessage(&models.ContextMessage{
		ID:        userMsgID,
		SessionID: req.SessionID,
		Role:      "user",
		Content:   req.Message,
		ImageRef:  "",
	})

	// Build LLM context
	llmMessages, err := models.BuildLLMContext(h.store, req.SessionID, h.systemPrompt, req.Message, baseImageURL)
	if err != nil {
		http.Error(w, "failed to build context: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Call LLM for intent parsing
	var llmResp *ai.LLMResponsePayload
	if h.llmGateway != nil {
		llmResp, err = h.llmGateway.ParseIntent(llmMessages)
		if err != nil {
			log.Printf("LLM intent parse failed: %v, falling back to default edit", err)
			llmResp = &ai.LLMResponsePayload{
				Intent:          "edit",
				ToolName:        "edit_image",
				RewrittenPrompt: req.Message,
				Reasoning:       "正在为您处理图片...",
			}
		}
	} else {
		// No LLM configured, use default
		llmResp = &ai.LLMResponsePayload{
			Intent:          "edit",
			ToolName:        "edit_image",
			RewrittenPrompt: req.Message,
			Reasoning:       "正在为您处理图片...",
		}
	}

	// Create turn
	turnID := generateUUID()
	turn := &models.Turn{
		ID:        turnID,
		SessionID: req.SessionID,
		UserInput: req.Message,
		LLMResponse: models.LLMResponse{
			Intent:          llmResp.Intent,
			ToolName:        llmResp.ToolName,
			RewrittenPrompt: llmResp.RewrittenPrompt,
			Reasoning:       llmResp.Reasoning,
		},
		ToolName:     llmResp.ToolName,
		InputImageID: baseImageID,
		Status:       string(models.TurnStatusProcessing),
	}
	turn.SetToolParams(map[string]any{"prompt": llmResp.RewrittenPrompt})

	if err := h.store.CreateTurn(turn); err != nil {
		http.Error(w, "failed to create turn: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute tool
	tool, err := h.toolReg.Get(llmResp.ToolName)
	if err != nil {
		h.store.UpdateTurnResult(turnID, string(models.TurnStatusFailed), "", "")
		http.Error(w, "unknown tool: "+llmResp.ToolName, http.StatusBadRequest)
		return
	}

	inputPath := filepath.Join(h.storagePath, baseImageID)
	result, err := tool.Execute(tools.Params{
		ImagePath: inputPath,
		Prompt:    llmResp.RewrittenPrompt,
	})
	if err != nil {
		log.Printf("tool %s failed: %v", llmResp.ToolName, err)
		h.store.UpdateTurnResult(turnID, string(models.TurnStatusFailed), "", "")
		http.Error(w, "failed to process image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save result image
	resultID := models.GenerateResultID(turnID)
	resultFilename := fmt.Sprintf("%s.png", resultID)
	resultPath := filepath.Join(h.storagePath, resultFilename)
	if err := models.SaveBytesToFile(result.ImageBytes, resultPath); err != nil {
		h.store.UpdateTurnResult(turnID, string(models.TurnStatusFailed), "", "")
		http.Error(w, "failed to save result: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resultURL := "/api/images/" + resultFilename

	// Update turn with result
	h.store.UpdateTurnResult(turnID, string(models.TurnStatusDone), resultFilename, resultURL)

	// Update session current image
	h.store.UpdateSessionCurrent(req.SessionID, resultFilename, resultURL)

	// Add assistant context message
	assistantMsgID := generateUUID()
	h.store.AddContextMessage(&models.ContextMessage{
		ID:        assistantMsgID,
		SessionID: req.SessionID,
		Role:      "assistant",
		Content:   llmResp.Reasoning,
		ImageRef:  resultURL,
		TurnID:    turnID,
	})

	resp := ChatResponse{
		TurnID:    turnID,
		SessionID: req.SessionID,
		Reasoning: llmResp.Reasoning,
		Intent:    llmResp.Intent,
		ToolName:  llmResp.ToolName,
		ResultURL: resultURL,
		Status:    "success",
		Message:   llmResp.Reasoning,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
