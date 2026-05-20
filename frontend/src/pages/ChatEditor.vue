<template>
  <div class="editor-panel">
    <div class="panel-header">
      <h2 class="panel-title">
        <span class="icon">&#10033;</span>
        AI 修图指令
      </h2>
      <p class="panel-desc">用自然语言描述你想要的效果</p>
    </div>

    <div v-if="!session.imageId" class="empty-state">
      <div class="empty-icon">&#128444;</div>
      <p class="empty-title">等待上传图片</p>
      <p class="empty-desc">请先在左侧上传一张图片，然后在这里输入修图指令</p>
    </div>

    <div v-else class="editor-content">
      <div class="image-comparison">
        <div class="compare-card">
          <div class="compare-label">原图</div>
          <div class="compare-image">
            <img :src="session.originalUrl" alt="原始图片" />
          </div>
        </div>

        <div class="compare-arrow">
          <span>&#10132;</span>
        </div>

        <div class="compare-card" :class="{ 'compare-empty': !latestResult }">
          <div class="compare-label">结果</div>
          <div v-if="latestResult" class="compare-image">
            <img :src="latestResult" alt="修图结果" />
          </div>
          <div v-else class="compare-placeholder">
            <span>&#10033;</span>
            <p>等待修图</p>
          </div>
        </div>
      </div>

      <div class="prompt-area">
        <div class="prompt-input-wrap">
          <textarea
            v-model="prompt"
            placeholder="描述你想要的修图效果，例如：把背景换成傍晚的天空，增强照片对比度和色彩饱和度..."
            rows="3"
            @keydown.enter.ctrl="submitPrompt"
          ></textarea>
          <button
            class="send-btn"
            :disabled="loading || !prompt.trim()"
            @click="submitPrompt"
          >
            <span v-if="loading" class="spinner"></span>
            <span v-else>&#8594;</span>
          </button>
        </div>
        <p class="prompt-hint">Ctrl + Enter 快速提交</p>
      </div>

      <p v-if="session.statusMessage" class="status" :class="{ 'status-error': isError }">
        <span v-if="loading" class="spinner-sm"></span>
        {{ session.statusMessage }}
      </p>

      <div class="examples" v-if="!session.results.length">
        <p class="examples-title">试试这些指令</p>
        <div class="example-chips">
          <button
            v-for="ex in examplePrompts"
            :key="ex"
            class="example-chip"
            @click="prompt = ex"
          >{{ ex }}</button>
        </div>
      </div>

      <div class="results" v-if="session.results.length">
        <div class="results-header">
          <h3 class="results-title">修图历史</h3>
          <span class="results-count">{{ session.results.length }} 条</span>
        </div>
        <ResultCard
          v-for="result in session.results"
          :key="result.resultUrl"
          :prompt="result.prompt"
          :result-url="result.resultUrl"
          @reuse="prompt = $event"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import axios from 'axios'
import { useSessionStore } from '@/stores/session'
import ResultCard from '@/components/ResultCard.vue'

const session = useSessionStore()
const prompt = ref('')
const loading = ref(false)
const isError = ref(false)

const examplePrompts = [
  '把背景变成模糊的森林',
  '增强色彩饱和度，让画面更鲜艳',
  '把照片变成黑白风格',
  '去掉背景中多余的人物',
  '添加柔和的光晕效果',
]

const latestResult = computed(() => {
  if (!session.results.length) return null
  return session.results[0].resultUrl
})

async function submitPrompt() {
  if (!session.imageId || !prompt.value.trim()) return

  loading.value = true
  isError.value = false
  session.setStatus('AI 正在处理你的图片...')

  try {
    const response = await axios.post('/api/edit', {
      imageId: session.imageId,
      prompt: prompt.value.trim(),
    })
    const data = response.data as { status: string; resultUrl: string; message: string }
    session.addResult(prompt.value.trim(), data.resultUrl)
    session.setStatus(data.message)
    prompt.value = ''
    isError.value = false
  } catch {
    isError.value = true
    session.setStatus('修图失败，请稍后再试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.editor-panel {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
}

.panel-header {
  margin-bottom: 0.2rem;
}

.panel-title {
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.panel-title .icon {
  font-size: 1.1rem;
  color: var(--accent);
}

.panel-desc {
  font-size: 0.8rem;
  color: var(--text-secondary);
  margin-top: 0.2rem;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
  background: var(--bg-card);
  border-radius: var(--radius-xl);
  border: 2px dashed var(--border);
}

.empty-icon {
  font-size: 3.5rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text-secondary);
  margin-bottom: 0.4rem;
}

.empty-desc {
  font-size: 0.85rem;
  color: var(--text-muted);
  max-width: 320px;
}

.editor-content {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
}

.image-comparison {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 1rem;
  align-items: center;
}

.compare-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border);
  overflow: hidden;
}

.compare-label {
  font-size: 0.72rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
  padding: 0.6rem 1rem;
  border-bottom: 1px solid var(--border);
}

.compare-image {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-input);
}

.compare-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.compare-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.8rem;
  color: var(--accent);
  opacity: 0.7;
}

.compare-placeholder {
  aspect-ratio: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  gap: 0.4rem;
}

.compare-placeholder span {
  font-size: 2rem;
  opacity: 0.4;
}

.compare-placeholder p {
  font-size: 0.8rem;
  opacity: 0.5;
}

.compare-empty {
  border-style: dashed;
}

.prompt-area {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 1rem;
  border: 1px solid var(--border);
}

.prompt-input-wrap {
  display: flex;
  gap: 0.6rem;
  align-items: flex-end;
}

.prompt-input-wrap textarea {
  flex: 1;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  color: var(--text-primary);
  padding: 0.8rem 1rem;
  resize: vertical;
  min-height: 80px;
  font-size: 0.9rem;
  transition: border-color var(--transition);
}

.prompt-input-wrap textarea:focus {
  outline: none;
  border-color: var(--border-focus);
  box-shadow: 0 0 0 3px var(--accent-glow);
}

.prompt-input-wrap textarea::placeholder {
  color: var(--text-muted);
}

.send-btn {
  width: 46px;
  height: 46px;
  border-radius: var(--radius);
  background: var(--accent-gradient);
  color: #fff;
  font-size: 1.2rem;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.send-btn:hover:not(:disabled) {
  transform: scale(1.05);
  box-shadow: 0 4px 16px var(--accent-glow);
}

.prompt-hint {
  font-size: 0.7rem;
  color: var(--text-muted);
  margin-top: 0.5rem;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

.spinner-sm {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  vertical-align: middle;
  margin-right: 0.4rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.status {
  font-size: 0.82rem;
  padding: 0.6rem 0.9rem;
  border-radius: var(--radius-sm);
  background: rgba(34, 197, 94, 0.1);
  color: var(--success);
  border: 1px solid rgba(34, 197, 94, 0.2);
  display: flex;
  align-items: center;
}

.status-error {
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger);
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.examples {
  padding: 1rem 0;
}

.examples-title {
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--text-muted);
  margin-bottom: 0.6rem;
}

.example-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.example-chip {
  padding: 0.4rem 0.9rem;
  font-size: 0.78rem;
  font-weight: 500;
  background: var(--bg-card);
  color: var(--text-secondary);
  border: 1px solid var(--border);
  border-radius: 100px;
  cursor: pointer;
  transition: all var(--transition);
}

.example-chip:hover {
  border-color: var(--accent);
  color: var(--accent-hover);
  background: var(--bg-card-hover);
}

.results {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.results-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.4rem;
}

.results-title {
  font-size: 0.95rem;
  font-weight: 700;
}

.results-count {
  font-size: 0.75rem;
  color: var(--text-muted);
  background: var(--bg-card);
  padding: 0.2rem 0.6rem;
  border-radius: 100px;
  border: 1px solid var(--border);
}

@media (max-width: 900px) {
  .image-comparison {
    grid-template-columns: 1fr;
  }
  .compare-arrow {
    transform: rotate(90deg);
  }
}
</style>
