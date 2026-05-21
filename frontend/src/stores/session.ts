import { defineStore } from 'pinia'
import type { Turn, ContextMessage } from '@/api'
import * as api from '@/api'

export interface ImageInfo {
  id: string
  url: string
}

interface SessionState {
  sessionId: string | null
  originalImage: ImageInfo | null
  currentImage: ImageInfo | null
  turns: Turn[]
  messages: ContextMessage[]
  isStreaming: boolean
  activeTool: string | null
  statusMessage: string
}

export const useSessionStore = defineStore('session', {
  state: (): SessionState => ({
    sessionId: null,
    originalImage: null,
    currentImage: null,
    turns: [],
    messages: [],
    isStreaming: false,
    activeTool: null,
    statusMessage: '',
  }),

  getters: {
    hasImage: (state): boolean => state.originalImage !== null,
    latestResultUrl: (state): string | null => {
      if (!state.turns.length) return null
      const doneTurns = state.turns.filter(t => t.status === 'done' && t.resultImageUrl)
      if (!doneTurns.length) return null
      return doneTurns[doneTurns.length - 1].resultImageUrl
    },
    turnsCount: (state): number => state.turns.length,
    doneTurns: (state): Turn[] => state.turns.filter(t => t.status === 'done'),
  },

  actions: {
    async uploadAndCreateSession(file: File) {
      this.statusMessage = '正在上传图片...'
      try {
        const data = await api.uploadImage(file)
        this.sessionId = data.sessionId
        this.originalImage = { id: data.imageId, url: data.originalUrl }
        this.currentImage = { id: data.imageId, url: data.originalUrl }
        this.turns = []
        this.messages = []
        this.statusMessage = '图片上传成功'
      } catch {
        this.statusMessage = '上传失败，请重试'
        throw new Error('upload failed')
      }
    },

    async sendMessage(message: string) {
      if (!this.sessionId) return
      this.isStreaming = true
      this.statusMessage = 'AI 正在思考...'
      this.messages.push({
        id: crypto.randomUUID(),
        sessionId: this.sessionId,
        role: 'user',
        content: message,
        imageRef: '',
        turnId: '',
        createdAt: new Date().toISOString(),
      })

      try {
        const resp = await api.sendChatMessage({
          sessionId: this.sessionId,
          message,
        })

        this.statusMessage = resp.reasoning

        // Add turn
        this.turns.push({
          turnId: resp.turnId,
          sessionId: resp.sessionId,
          userInput: message,
          llmResponse: {
            intent: resp.intent,
            toolName: resp.toolName,
            rewrittenPrompt: '',
            reasoning: resp.reasoning,
          },
          toolName: resp.toolName,
          resultImageId: '',
          resultImageUrl: resp.resultUrl,
          inputImageId: this.currentImage?.id ?? '',
          status: 'done',
          createdAt: new Date().toISOString(),
        })

        // Update current image
        this.currentImage = {
          id: resp.turnId,
          url: resp.resultUrl,
        }

        // Add AI message
        this.messages.push({
          id: crypto.randomUUID(),
          sessionId: this.sessionId!,
          role: 'assistant',
          content: resp.reasoning,
          imageRef: resp.resultUrl,
          turnId: resp.turnId,
          createdAt: new Date().toISOString(),
        })

        this.statusMessage = '处理完成'
      } catch {
        this.statusMessage = '处理失败，请稍后再试'
      } finally {
        this.isStreaming = false
      }
    },

    async loadSession(id: string) {
      try {
        const [session, turns, messages] = await Promise.all([
          api.getSession(id),
          api.getTurns(id),
          api.getContextMessages(id),
        ])
        this.sessionId = session.sessionId
        this.originalImage = {
          id: session.originalImageId,
          url: session.originalImageUrl,
        }
        this.currentImage = {
          id: session.currentImageId,
          url: session.currentImageUrl,
        }
        this.turns = turns
        this.messages = messages
      } catch {
        console.warn('Failed to load session:', id)
      }
    },

    setActiveTool(toolName: string | null) {
      this.activeTool = toolName
    },

    resetSession() {
      this.sessionId = null
      this.originalImage = null
      this.currentImage = null
      this.turns = []
      this.messages = []
      this.isStreaming = false
      this.activeTool = null
      this.statusMessage = ''
    },
  },
})
