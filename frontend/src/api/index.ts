import axios from 'axios'

export interface UploadResponse {
  imageId: string
  originalUrl: string
  sessionId: string
}

export interface ChatRequest {
  sessionId: string
  message: string
  baseImageId?: string
}

export interface Turn {
  turnId: string
  sessionId: string
  userInput: string
  llmResponse: {
    intent: string
    toolName: string
    rewrittenPrompt: string
    reasoning: string
  }
  toolName: string
  resultImageId: string
  resultImageUrl: string
  inputImageId: string
  status: 'pending' | 'processing' | 'done' | 'failed'
  createdAt: string
}

export interface ChatResponse {
  turnId: string
  sessionId: string
  reasoning: string
  intent: string
  toolName: string
  resultUrl: string
  status: string
  message: string
}

export interface SessionInfo {
  sessionId: string
  originalImageId: string
  originalImageUrl: string
  currentImageId: string
  currentImageUrl: string
  createdAt: string
  updatedAt: string
}

export interface ContextMessage {
  id: string
  sessionId: string
  role: 'user' | 'assistant' | 'system'
  content: string
  imageRef: string
  turnId: string
  createdAt: string
}

const api = axios.create({
  baseURL: '/api',
})

export async function uploadImage(file: File): Promise<UploadResponse> {
  const formData = new FormData()
  formData.append('image', file)
  const res = await api.post<UploadResponse>('/upload', formData)
  return res.data
}

export async function sendChatMessage(req: ChatRequest): Promise<ChatResponse> {
  const res = await api.post<ChatResponse>('/chat', req)
  return res.data
}

export async function getSession(id: string): Promise<SessionInfo> {
  const res = await api.get<SessionInfo>(`/session/${id}`)
  return res.data
}

export async function getTurns(sessionId: string): Promise<Turn[]> {
  const res = await api.get<Turn[]>(`/session/${sessionId}/turns`)
  return res.data
}

export async function getContextMessages(sessionId: string): Promise<ContextMessage[]> {
  const res = await api.get<ContextMessage[]>(`/session/${sessionId}/context`)
  return res.data
}
