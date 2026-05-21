# AI 修图助手 — 产品架构设计文档

## Context

当前项目是一个简单的"上传图片 → 输入指令 → 调用 OpenAI Images Edits API → 返回结果"的单轮修图工具。用户希望将其升级为一个**有温度的 AI 修图助手**，核心变化：

1. 从"表单式操作"变为"对话式交互"——用户与 AI 对话完成修图
2. 从"单一 API 调用"变为"多工具路由"——LLM 理解意图后选择不同的修图工具
3. 从"无状态"变为"会话持久化"——图文上下文跨轮次保留
4. 从"粗糙展示"变为"沉浸式体验"——大图预览、丝滑动画、流式反馈

---

## 一、整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (Vue 3)                          │
│  ┌──────────┐  ┌──────────┐  ┌───────────┐  ┌──────────────┐  │
│  │ Immersive │  │  AI Chat  │  │  Result   │  │  History     │  │
│  │  Viewer   │  │  Panel    │  │  Stream   │  │  Panel       │  │
│  └──────────┘  └──────────┘  └───────────┘  └──────────────┘  │
│                       │ Pinia Store │                             │
│                       └─────────────┘                             │
└───────────────────────────┬─────────────────────────────────────┘
                            │ HTTP / SSE
┌───────────────────────────┴─────────────────────────────────────┐
│                       Backend (Go)                                │
│  ┌──────────┐  ┌──────────────┐  ┌───────────┐  ┌───────────┐  │
│  │ Router   │  │  LLM Gateway │  │  Tool     │  │  Session  │  │
│  │ (upload, │  │  (intent     │  │  Registry │  │  Store    │  │
│  │  chat,   │  │   parse +    │  │  (edit,   │  │  (DB)     │  │
│  │  images) │  │   prompt     │  │  generate,│  │           │  │
│  │          │  │   rewrite)   │  │  upscale) │  │           │  │
│  └──────────┘  └──────────────┘  └───────────┘  └───────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### 核心原则
- **图文混排上下文**：LLM 看到的上下文包含图片引用 + 文字对话
- **LLM 作为调度中心**：所有用户输入先经过 LLM 解析意图，再路由到具体工具
- **会话与轮次分离**：Session → Turns → Results，每轮独立存储

---

## 二、数据模型设计

### 2.1 核心实体

```
Session (会话)
├── sessionId: string (UUID)
├── originalImageId: string      // 用户最初上传的原始图
├── originalImageUrl: string
├── turns: Turn[]                // 按时间顺序的轮次列表
├── contextWindow: ContextMessage[] // 给 LLM 的上下文窗口
├── createdAt: timestamp
└── updatedAt: timestamp

Turn (轮次 - 一次 AI 对话 + 修图操作)
├── turnId: string (UUID)
├── sessionId: string
├── userInput: string            // 用户原始输入
├── llmResponse: {
│     intent: string             // 意图类型 (edit / generate / upscale / remove_bg / ...)
│     rewrittenPrompt: string    // LLM 改写后的英文提示词
│     reasoning: string          // LLM 简短的中文解释，如 "正在把背景色调调整为暖色..."
│     toolName: string           // 被路由到的工具名称
│   }
├── toolName: string             // 实际调用的工具
├── toolParams: object           // 传给工具的完整参数
├── resultImageId: string        // 生成的结果图 ID
├── resultImageUrl: string       // 结果图访问 URL
├── status: enum (pending | processing | done | failed)
├── inputImageId: string         // 本轮的输入图（可能是原始图或上一轮结果图）
└── createdAt: timestamp

ContextMessage (上下文消息 - 给 LLM 的)
├── role: string                 // "user" | "assistant" | "system"
├── content: string              // 文字内容
├── imageRef: string | null      // 引用的图片 URL
└── turnId: string | null        // 关联的轮次
```

### 2.2 数据库选型
- 开发阶段使用 **SQLite**（零配置，单文件，足够轻量）
- 生产可迁移至 PostgreSQL
- 表结构：sessions, turns, context_messages

---

## 三、后端 API 设计

### 3.1 新增/修改接口

| Method | Path | 说明 |
|--------|------|------|
| POST | `/api/upload` | **(保留)** 上传图片，创建 Session |
| POST | `/api/chat` | **(新增)** 对话式修图入口 |
| GET | `/api/session/:id` | **(新增)** 获取会话完整状态 |
| GET | `/api/session/:id/turns` | **(新增)** 获取所有轮次 |
| GET | `/api/session/:id/context` | **(新增)** 获取上下文窗口 |
| POST | `/api/session/:id/reset` | **(新增)** 重置到某轮 / 回到原图 |
| GET | `/api/images/:filename` | **(保留)** 静态图片服务 |

### 3.2 核心接口：`POST /api/chat`

**请求体：**
```json
{
  "sessionId": "uuid",
  "message": "能不能让天空看起来像日落的时候",
  "baseImageId": "optional, 默认为上一轮结果或原图"
}
```

**处理流程：**
```
1. 接收用户消息
2. 构建 LLM 上下文（包含历史图文对话 + System Prompt）
3. 调用 LLM，传入：
   - System Prompt（定义工具列表、输出格式）
   - 上下文消息（图文混排）
   - 用户当前输入
4. LLM 返回结构化 JSON：
   {
     "intent": "edit",
     "toolName": "openai_image_edit",
     "rewrittenPrompt": "change the sky to a warm sunset with orange and pink tones",
     "reasoning": "正在为您把天空调整为温暖的日落色调..."
   }
5. 根据 toolName 路由到对应的工具处理器
6. 工具处理器调用第三方 API
7. 保存结果，更新上下文窗口
8. 返回结果给前端
```

**响应体（SSE 流式）：**
```
event: reasoning
data: {"text": "正在为您把天空调整为温暖的日落色调..."}

event: processing
data: {"status": "calling_openai", "progress": 50}

event: result
data: {"turnId": "uuid", "resultUrl": "/api/images/xxx.png", "status": "done"}

event: done
data: {"sessionId": "uuid", "turnsCount": 3}
```

> 前期阶段可用普通 JSON 响应替代 SSE，先落地核心逻辑再优化流式体验。

---

## 四、LLM 编排引擎

### 4.1 架构

```
用户输入 → LLM Gateway → 意图解析 + Prompt 改写 → 工具路由 → 第三方 API
              ↑                                        │
              │ 上下文窗口（图文混排）                    ├─ OpenAI Images Edits
              │                                        ├─ DALL-E Image Gen
              │                                        ├─ Remove.bg API
              │                                        ├─ Real-ESRGAN 超分
              │                                        └─ ...扩展工具
              │
         ┌────┴────┐
         │ Session  │
         │ Context  │
         └─────────┘
```

### 4.2 LLM System Prompt 设计要点

```
你是一个专业的AI修图助手。你可以使用以下工具：

1. edit_image - 编辑现有图片（改背景、调色、加滤镜、去物体等）
2. generate_image - 从文字生成新图片
3. remove_background - 去除图片背景
4. upscale - 提升图片分辨率/清晰度
5. style_transfer - 风格迁移（水彩、油画、动漫等）

你需要：
1. 理解用户的中文输入意图
2. 选择合适的工具
3. 将中文指令改写为适合对应工具的英文提示词
4. 用温暖、简洁的中文告知用户正在执行什么操作

返回 JSON 格式：
{
  "intent": "工具对应的意图类型",
  "toolName": "具体工具名",
  "rewrittenPrompt": "改写后的英文提示词",
  "reasoning": "给用户的简短中文反馈"
}
```

### 4.3 上下文窗口构建

```
System Prompt (固定)
  ↓
[user]: "我想修这张图"  (imageRef: original_xxx.png)
[assistant]: "好的，我看到了您的图片，有什么想调整的吗？"
[user]: "让天空变蓝一点"
[assistant]: "正在为您调整天空的颜色..." (imageRef: result_turn1.png)
[user]: "再亮一点"
  ↓
当前用户输入
```

上下文窗口限制：最近 N 条消息（如最近 10 条），超出部分做摘要或丢弃。

### 4.4 LLM 选型

- **主模型**：OpenAI GPT-4o / GPT-4.1（视觉+文本理解能力强）
- **备选**：Claude Sonnet（成本更低，且同样支持图片理解）
- **API 配置**：通过环境变量切换 `LLM_PROVIDER` / `LLM_API_KEY` / `LLM_MODEL`

---

## 五、工具注册表 (Tool Registry)

### 5.1 工具接口定义

```go
type Tool interface {
    Name() string
    Description() string
    Execute(params ToolParams) (*ToolResult, error)
}

type ToolParams struct {
    ImagePath string
    Prompt    string
    Options   map[string]any
}

type ToolResult struct {
    ImageBytes []byte
    ImageURL   string
    Metadata   map[string]string
}
```

### 5.2 首批工具实现

| 工具 | 底层 API | 参数 |
|------|----------|------|
| `openai_image_edit` | OpenAI `/images/edits` + gpt-image-1 | image + prompt |
| `openai_image_generate` | OpenAI `/images/generations` + dall-e-3 | prompt, size, quality |
| `remove_background` | remove.bg API | image |
| `upscale` | Real-ESRGAN (自部署或第三方) | image, scale |

### 5.3 工具注册

启动时自动注册所有可用工具（根据环境变量中是否有对应 API Key 决定是否启用），LLM 的 System Prompt 中只列出已启用的工具。

---

## 六、前端架构设计

### 6.1 页面布局

```
┌──────────────────────────────────────────┐
│              沉浸式图片预览区              │
│         (原图 / 当前结果 大图展示)         │
│                                          │
│       [← 上一版]  [对比]  [下一版→]       │
│                                          │
├──────────────────────────────────────────┤
│  功能按钮栏                                │
│  [🎨 AI编辑] [✨ 生成] [🖼 去背景] [🔍 超分]│
├──────────────────────────────────────────┤
│  AI 对话面板 (可折叠)                       │
│  ┌────────────────────────────────────┐  │
│  │ AI: 好的，为您调整天空色调...        │  │
│  │                           (已处理)  │  │
│  │ 用户: 让天空变暖一点                 │  │
│  │ AI: 正在优化画面，请稍候...          │  │
│  └────────────────────────────────────┘  │
│  ┌──────────────────────┐ [发送]         │
│  │ 输入修图指令...       │               │
│  └──────────────────────┘               │
├──────────────────────────────────────────┤
│  历史版本缩略图条 (横向滚动)               │
│  [原图] [v1] [v2] [v3(当前)]             │
└──────────────────────────────────────────┘
```

### 6.2 核心交互流程

```
1. 用户进入 → 空白画布 + 上传按钮
2. 上传图片 → 大图展示，底部出现功能栏和对话按钮
3. 点击功能按钮或 AI 对话图标 → 打开对话面板
4. 输入指令 → LLM 处理 → 显示 reasoning（"正在为您..."）
5. 生成完成 → 大图区域丝滑过渡到结果图
6. 可以继续对话修图，也可以切换历史版本
```

### 6.3 组件树

```
App.vue
├── ImmersiveViewer.vue        // 沉浸式大图预览
│   ├── ImageCanvas.vue        // 图片画布（支持缩放、拖拽）
│   └── VersionControl.vue     // 版本切换控制
├── ToolBar.vue                // 功能按钮栏
├── ChatPanel.vue              // AI 对话面板（底部弹出）
│   ├── ChatMessage.vue        // 单条消息（文字+图片）
│   ├── ReasoningBubble.vue    // AI 思考中气泡
│   └── ChatInput.vue          // 输入框
├── HistoryStrip.vue           // 历史版本缩略图条
└── UploadOverlay.vue          // 上传引导页（首次进入）
```

### 6.4 Pinia Store 重构

```typescript
// stores/session.ts
interface SessionState {
  sessionId: string | null
  originalImage: ImageInfo | null
  currentImage: ImageInfo | null        // 当前展示的图片
  turns: Turn[]                         // 所有轮次
  messages: ChatMessage[]               // 对话消息
  isStreaming: boolean                  // 是否正在流式响应
  activeTool: string | null             // 当前选中的工具
}

// stores/viewer.ts
interface ViewerState {
  zoom: number
  panX: number
  panY: number
  compareMode: boolean                  // 是否对比模式
  compareImage: ImageInfo | null        // 对比的图片
}
```

### 6.5 动画设计

- **上传后过渡**：上传区缩放消失 → 大图从缩略图位置展开
- **版本切换**：图片交叉淡入淡出 (crossfade)，300ms ease-out
- **结果生成**：新图片从模糊到清晰 (blur-in)，配合微妙的缩放脉冲
- **对话面板**：从底部滑入 (slide-up)，带弹性缓动
- **AI 思考中**：工具按钮发出呼吸光效，提示 AI 正在处理

---

## 七、开发阶段规划

### Phase 1：核心架构搭建
- 后端：SQLite 数据库 + Session/Turn 模型 + 基础 CRUD
- 后端：LLM Gateway 模块（调用 OpenAI Chat Completions，返回结构化 JSON）
- 前端：ImmersiveViewer + ChatPanel 基础布局
- **产出**：上传图片 → 输入文字 → LLM 返回 reasoning → 调用现有 edit API → 展示结果

### Phase 2：多工具支持
- 后端：Tool Registry 接口 + openai_image_edit 实现
- 后端：新增 generate / remove_bg / upscale 工具
- 前端：ToolBar 功能按钮栏切换 activeTool
- **产出**：点击不同工具按钮 → LLM 改写 prompt → 路由到对应 API

### Phase 3：体验优化
- SSE 流式响应（reasoning 实时展示 → processing 进度 → result）
- 图片动画过渡效果
- 历史版本管理（缩略图条、对比模式）
- 会话持久化（刷新页面后恢复）

### Phase 4：打磨上线
- 错误处理与重试
- 移动端适配
- 性能优化（图片懒加载、缓存）
- 生产环境部署配置

---

## 八、文件结构变更

```
github/AIImageEdit/
├── backend/
│   ├── main.go                     # 入口，路由注册
│   ├── go.mod
│   ├── ai/
│   │   ├── llm_gateway.go          # LLM 调用封装（新增）
│   │   └── thirdparty.go           # 现有 OpenAI Images Edits（保留）
│   ├── tools/                      # 工具注册表（新增目录）
│   │   ├── registry.go             # 工具注册中心
│   │   ├── edit.go                 # edit_image 工具
│   │   ├── generate.go             # generate_image 工具
│   │   ├── remove_bg.go            # remove_background 工具
│   │   └── upscale.go              # upscale 工具
│   ├── handlers/
│   │   ├── upload.go               # 保留，改为创建 Session
│   │   ├── edit.go                 # 保留（内部调用）
│   │   ├── chat.go                 # 新增：对话式修图入口
│   │   └── session.go              # 新增：会话查询接口
│   ├── models/
│   │   ├── session.go              # 扩充：Session + Turn + ContextMessage
│   │   ├── db.go                   # 新增：SQLite 初始化与迁移
│   │   └── context.go              # 新增：上下文窗口构建
│   └── storage/
│       └── images/                 # 图片存储目录
├── frontend/
│   ├── src/
│   │   ├── App.vue                 # 新布局
│   │   ├── main.ts
│   │   ├── style.css               # 更新全局样式
│   │   ├── components/
│   │   │   ├── ImmersiveViewer.vue  # 新增
│   │   │   ├── ImageCanvas.vue     # 新增
│   │   │   ├── VersionControl.vue   # 新增
│   │   │   ├── ToolBar.vue         # 新增
│   │   │   ├── ChatPanel.vue       # 新增（重组原 ChatEditor）
│   │   │   ├── ChatMessage.vue     # 新增
│   │   │   ├── ChatInput.vue       # 新增
│   │   │   ├── HistoryStrip.vue    # 新增
│   │   │   ├── UploadOverlay.vue   # 新增
│   │   │   └── ResultCard.vue      # 保留（用于历史列表）
│   │   ├── pages/
│   │   │   └── EditorPage.vue      # 新增：整合所有组件的页面
│   │   ├── stores/
│   │   │   ├── session.ts          # 重构
│   │   │   └── viewer.ts           # 新增
│   │   └── api/
│   │       └── index.ts            # 新增：Axios 封装 + SSE
│   └── ...
└── ARCHITECTURE.md                 # 架构文档（本文件）
```

---

## 九、验证方案

1. **Phase 1 验证**：
   - 上传图片 → 创建 Session → 发送 chat 消息 → LLM 返回 reasoning → 调用 edit → 前端展示结果
   - 检查 SQLite 中 sessions / turns / context_messages 表数据正确
2. **Phase 2 验证**：
   - 切换不同工具按钮 → LLM 改写 prompt 并路由到不同 API → 结果正确返回
   - 工具注册表根据 API Key 可用性正确启用/禁用
3. **Phase 3 验证**：
   - SSE 事件流正确推送 reasoning → processing → result → done
   - 图片切换动画流畅，版本切换正常
   - 刷新页面后会话恢复
4. **集成测试**：端到端上传 → 多轮对话修图 → 历史回退 → 下载结果
