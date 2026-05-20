import { defineStore } from 'pinia'

export interface EditResult {
  prompt: string
  resultUrl: string
}

export const useSessionStore = defineStore('session', {
  state: () => ({
    imageId: '' as string,
    originalUrl: '' as string,
    loading: false as boolean,
    statusMessage: '' as string,
    results: [] as EditResult[],
  }),

  actions: {
    setImage(imageId: string, originalUrl: string) {
      this.imageId = imageId
      this.originalUrl = originalUrl
      this.results = []
      this.statusMessage = ''
    },
    addResult(prompt: string, resultUrl: string) {
      this.results.unshift({ prompt, resultUrl })
    },
    setLoading(value: boolean) {
      this.loading = value
    },
    setStatus(message: string) {
      this.statusMessage = message
    },
  },
})
