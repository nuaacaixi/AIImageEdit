## Plan: AI Image Editing Website Architecture

TL;DR - 构建一个 Go 后端 + Vue 前端的对话式 AI 修图网站，后端负责图像上传、任务管理、第三方图片编辑 API 调用和静态资源服务，前端负责上传、对话式修图界面、结果预览和下载。

**Steps**
1. 设计整体系统架构和模块边界。
   - 后端：Go HTTP 服务，图像文件存储，第三方图像编辑 API 代理层，任务状态管理，静态文件服务。
   - 前端：Vue 3 SPA，上传组件、聊天修图组件、结果预览组件。
   - AI 层：集成第三方图像编辑 API（如 OpenAI Image Edit、Azure OpenAI、或其他可用图像编辑服务）。
2. 定义核心后端接口。
   - POST /api/upload：上传原图并返回 imageId 和原图 URL。
   - POST /api/edit：提交修图指令（文本 + imageId），后端调用第三方图像编辑 API，返回任务状态、结果图片 URL。
   - GET /api/images/:id：获取存储的原图或修图结果图。
   - GET /api/session/:id：获取当前会话状态和历史修图记录（可选，用于多轮交互）。
3. 定义前端页面和组件结构。
   - UploadPage：图像选择、上传按钮、上传状态反馈。
   - ChatEditor：对话式指令输入、提交按钮、AI 修图响应列表、当前结果预览。
   - ResultCard：展示原图/修图结果、下载按钮、重新编辑按钮。
   - HistoryList：展示当前会话的指令历史和修图版本（可选）。
4. 规划数据流与状态管理。
   - 使用 Pinia 管理当前会话状态、已上传图片、编辑历史、AI 修图结果。
   - 上传后保存 imageId，前端发送修图指令到 /api/edit。
   - 后端收到指令后调用第三方 API：传入原图、提示词/指令，保存返回的修图结果图并返回给前端。
5. 规划文件存储与部署方案。
   - 开发阶段使用本地 filesystem 存储 images/ 目录，如 `backend/storage/images/`。
   - 生产可扩展为 S3/MinIO 或第三方云存储。
   - 前端构建产物可由 Go 静态文件服务，或单独部署在静态站点平台。
6. 验证设计。
   - 确保接口契约：上传返回 `imageId` 和 `originalUrl`；修图返回 `status`、`resultUrl`、`message`。
   - 通过第三方 API 调用完整流程：上传图 -> 提交修图指令 -> 展示修图结果。
   - 验证前端 UI 是否支持多次“继续修图”操作。

**Relevant files**
- `backend/main.go` — Go 服务入口，路由与静态资源配置。
- `backend/handlers/upload.go` — 处理图像上传与存储。
- `backend/handlers/edit.go` — 处理修图指令、调用第三方图像编辑 API。
- `backend/ai/thirdparty.go` — 封装第三方图像编辑 API 调用逻辑，包括请求构建、结果解析和错误处理。
- `backend/models/session.go` — 定义会话状态、修图记录数据结构。
- `frontend/src/App.vue` — 应用入口。
- `frontend/src/pages/UploadPage.vue` — 处理图像上传。
- `frontend/src/pages/ChatEditor.vue` — 对话式修图界面。
- `frontend/src/components/ResultCard.vue` — 结果展示组件。
- `frontend/src/stores/session.ts` — Pinia 状态管理。

**Verification**
1. 手动验证：上传样图后，提交修图指令能通过第三方 API 返回修图结果，并在前端显示。
2. 检查接口契约：`/api/upload` 返回 `imageId`、`originalUrl`；`/api/edit` 返回 `status`、`resultUrl`、`message`。
3. 代码结构验证：后端清晰分离上传、编辑接口和第三方 API 调用逻辑，前端具备上传-提交-预览流程。

**Decisions**
- 使用第三方图像编辑 API 实现核心修图功能，后端作为代理层负责上传、调用和结果存储。
- 首先实现匿名上传与对话式指令交互，不引入用户登录和权限管理。
- 对话式修图采用“多轮文本指令 + 返回图像”模式，暂不做实时流式更新。

**Further Considerations**
1. 第三方 API 的选择：OpenAI Image Edit / Azure OpenAI / 其他图像编辑服务。若后续需要更高吞吐可改为可插拔的 API 适配层。
2. 如果希望支持多人会话或用户隔离，可加入简单的 sessionId 或用户 token。
3. 若需要“撤销/比较”，可在后端增加 editVersion 和 resultHistory 字段。
