<template>
  <div class="history-strip" v-if="items.length > 0">
    <div class="strip-label">版本历史</div>
    <div class="strip-list">
      <button
        v-for="(item, index) in items"
        :key="item.id"
        class="strip-item"
        :class="{ current: index === currentIndex }"
        @click="$emit('select', index)"
      >
        <img :src="item.url" :alt="`版本 ${index + 1}`" />
        <span class="strip-version">v{{ index + 1 }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useSessionStore } from '@/stores/session'

const store = useSessionStore()

const emit = defineEmits<{
  select: [index: number]
}>()

const items = computed(() => {
  const list: { id: string; url: string }[] = []
  if (store.originalImage) {
    list.push(store.originalImage)
  }
  for (const turn of store.doneTurns) {
    if (turn.resultImageUrl) {
      list.push({ id: turn.turnId, url: turn.resultImageUrl })
    }
  }
  return list
})

const currentIndex = computed(() => {
  if (!store.currentImage) return 0
  const idx = items.value.findIndex(i => i.url === store.currentImage?.url)
  return idx >= 0 ? idx : items.value.length - 1
})
</script>

<style scoped>
.history-strip {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.6rem 1rem;
  background: var(--bg-card);
  border-radius: var(--radius);
  border: 1px solid var(--border);
}

.strip-label {
  font-size: 0.72rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
  flex-shrink: 0;
}

.strip-list {
  display: flex;
  gap: 0.5rem;
  overflow-x: auto;
  padding-bottom: 2px;
}

.strip-item {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-sm);
  overflow: hidden;
  border: 2px solid var(--border);
  cursor: pointer;
  transition: all var(--transition);
  position: relative;
  flex-shrink: 0;
  padding: 0;
  background: var(--bg-input);
}

.strip-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.strip-item:hover {
  border-color: var(--accent);
  transform: translateY(-2px);
}

.strip-item.current {
  border-color: var(--accent);
  box-shadow: 0 0 12px var(--accent-glow);
}

.strip-version {
  position: absolute;
  bottom: 2px;
  right: 2px;
  font-size: 0.55rem;
  font-weight: 700;
  background: rgba(0, 0, 0, 0.7);
  color: #fff;
  padding: 1px 4px;
  border-radius: 3px;
}
</style>
