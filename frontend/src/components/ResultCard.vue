<template>
  <div class="result-card">
    <div class="result-body">
      <div class="result-image">
        <img :src="resultUrl" alt="修图结果" loading="lazy" />
        <div class="image-actions">
          <a :href="resultUrl" target="_blank" class="img-btn" title="查看大图">
            &#128269;
          </a>
          <a :href="resultUrl" download class="img-btn" title="下载">
            &#128229;
          </a>
        </div>
      </div>
      <div class="result-info">
        <div class="prompt-text">
          <span class="prompt-label">指令</span>
          <p>"{{ prompt }}"</p>
        </div>
        <div class="result-actions">
          <button class="action-btn reuse-btn" @click="$emit('reuse', prompt)">
            复用指令
          </button>
          <a :href="resultUrl" download class="action-btn download-btn">
            下载结果
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ prompt: string; resultUrl: string }>()
defineEmits<{ reuse: [prompt: string] }>()
</script>

<style scoped>
.result-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  transition: all var(--transition);
}

.result-card:hover {
  border-color: var(--border-focus);
  box-shadow: var(--shadow);
}

.result-body {
  display: grid;
  grid-template-columns: 200px 1fr;
  gap: 0;
}

.result-image {
  position: relative;
  aspect-ratio: 1;
  background: var(--bg-input);
  overflow: hidden;
}

.result-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-actions {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  display: flex;
  gap: 0.3rem;
  opacity: 0;
  transition: opacity var(--transition);
}

.result-image:hover .image-actions {
  opacity: 1;
}

.img-btn {
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  border-radius: var(--radius-sm);
  color: #fff;
  text-decoration: none;
  font-size: 0.8rem;
  transition: background var(--transition);
}

.img-btn:hover {
  background: rgba(0, 0, 0, 0.85);
}

.result-info {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.prompt-label {
  font-size: 0.65rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--accent);
  margin-bottom: 0.3rem;
  display: block;
}

.prompt-text p {
  font-size: 0.85rem;
  color: var(--text-primary);
  line-height: 1.5;
  font-style: italic;
}

.result-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.8rem;
}

.action-btn {
  padding: 0.4rem 0.9rem;
  font-size: 0.75rem;
  font-weight: 600;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition);
  text-decoration: none;
  display: inline-flex;
  align-items: center;
}

.reuse-btn {
  background: transparent;
  color: var(--accent-hover);
  border: 1px solid var(--border);
}

.reuse-btn:hover {
  border-color: var(--accent);
  background: rgba(99, 102, 241, 0.08);
}

.download-btn {
  background: var(--accent-gradient);
  color: #fff;
  border: none;
}

.download-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 12px var(--accent-glow);
}

@media (max-width: 600px) {
  .result-body {
    grid-template-columns: 1fr;
  }
  .result-image {
    aspect-ratio: 16 / 9;
  }
}
</style>
