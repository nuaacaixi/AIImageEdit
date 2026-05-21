package tools

import (
	"fmt"

	"github.com/vibeCoding/AIImageEdit/backend/ai"
)

type EditImageTool struct {
	client *ai.Client
}

func NewEditImageTool(client *ai.Client) *EditImageTool {
	return &EditImageTool{client: client}
}

func (t *EditImageTool) Name() string {
	return "edit_image"
}

func (t *EditImageTool) Description() string {
	return "编辑现有图片——改变背景、调整颜色、添加滤镜、移除物体、风格转换等。需要提供图片和描述修改内容的提示词。"
}

func (t *EditImageTool) Execute(params Params) (*Result, error) {
	if t.client == nil {
		return nil, fmt.Errorf("edit_image tool: AI client is not configured (missing OPENAI_API_KEY)")
	}

	bytes, err := t.client.EditImage(params.ImagePath, params.Prompt)
	if err != nil {
		return nil, fmt.Errorf("edit_image: %w", err)
	}

	return &Result{
		ImageBytes: bytes,
	}, nil
}
