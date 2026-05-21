import { defineStore } from 'pinia'
import type { ImageInfo } from './session'

interface ViewerState {
  zoom: number
  panX: number
  panY: number
  compareMode: boolean
  compareImage: ImageInfo | null
}

export const useViewerStore = defineStore('viewer', {
  state: (): ViewerState => ({
    zoom: 1,
    panX: 0,
    panY: 0,
    compareMode: false,
    compareImage: null,
  }),

  actions: {
    resetView() {
      this.zoom = 1
      this.panX = 0
      this.panY = 0
    },
    toggleCompare() {
      this.compareMode = !this.compareMode
      if (!this.compareMode) {
        this.compareImage = null
      }
    },
    setCompareImage(img: ImageInfo | null) {
      this.compareImage = img
    },
  },
})
