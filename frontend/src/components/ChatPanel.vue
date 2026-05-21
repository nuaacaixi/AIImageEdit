<template>
  <div class="chat-panel" :class="{ open: isOpen }">
    <button class="chat-toggle" @click="isOpen = !isOpen">
      <span class="toggle-icon">{{ isOpen ? '✕' : '💬' }}</span>
      <span v-if="!isOpen" class="toggle-label">AI 对话</span>
    </button>

    <Transition name="slide-up">
      <div v-if="isOpen" class="chat-body">
        <div class="chat-header">
          <span class="chat-title">AI 修图助手</span>
          <span v-if="store.isStreaming" class="thinking-dot"></span>
        </div>

        <div class="chat-messages" ref="messagesEl">
          <div v-if="store.messages.length === 0" class="chat-empty">
            <p>告诉我你想怎么修改这张图片吧！</p>
            <div class="quick-prompts">
              <button
                v-for="qp in quickPrompts"
                :key="qp"
                class="quick-prompt"
                @click="$emit('send', qp)"
              >{{ qp }}</button>
            </div>
          </div>

          <template v-for="msg in store.messages" :key="msg.id">
            <div class="message" :class="msg.role">
              <div class="msg-avatar">{{ msg.role === 'user' ? '👤' : '✦' }}</div>
              <div class="msg-body">
                <p class="msg-text">{{ msg.content }}</p>
                <img
                  v-if="msg.imageRef"
                  :src="msg.imageRef"
                  class="msg-image"
                  @click="$emit('viewImage', msg.imageRef)"
                />
              </div>
            </div>
          </template>

          <div v-if="store.isStreaming" class="message assistant">
            <div class="msg-avatar">✦</div>
            <div class="msg-body">
              <div class="typing-indicator">
                <span></span><span></span><span></span>
              </div>
            </div>
          </div>
        </div>

        <div class="chat-input-area">
          <div class="input-row">
            <textarea
              v-model="inputText"
              placeholder="描述你想如何修改图片..."
              rows="2"
              @keydown.enter.prevent.ctrl="sendMessage"
              :disabled="store.isStreaming"
            ></textarea>
            <button
              class="send-btn"
              :disabled="!inputText.trim() || store.isStreaming"
              @click="sendMessage"
            >
              ↑
            </button>
          </div>
          <p class="input-hint">Ctrl + Enter 发送</p>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { useSessionStore } from '@/stores/session'

const store = useSessionStore()
const isOpen = ref(false)
const inputText = ref('')
const messagesEl = ref<HTMLElement | null>(null)

const emit = defineEmits<{
  send: [message: string]
  viewImage: [url: string]
}>()

const quickPrompts = [
  '让天空更蓝一些',
  '把照片变成黑白风格',
  '增强色彩对比度',
  '添加柔和的光线效果',
  '把背景虚化',
]

watch(() => store.messages.length, async () => {
  await nextTick()
  if (messagesEl.value) {
    messagesEl.value.scrollTop = messagesEl.value.scrollHeight
  }
})

function sendMessage() {
  const text = inputText.value.trim()
  if (!text || store.isStreaming) return
  inputText.value = ''
  store.sendMessage(text)
}
</script>

<style scoped>
.chat-panel {
  position: fixed;
  bottom: 1.5rem;
  right: 1.5rem;
  z-index: 200;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.chat-toggle {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: var(--accent-gradient);
  box-shadow: 0 4px 20px var(--accent-glow);
  color: #fff;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: none;
  cursor: pointer;
  transition: all var(--transition);
  flex-shrink: 0;
}

.chat-toggle:hover {
  transform: scale(1.08);
  box-shadow: 0 6px 28px var(--accent-glow);
}

.toggle-icon {
  font-size: 1.2rem;
}

.toggle-label {
  font-size: 0.55rem;
  font-weight: 600;
  margin-top: -2px;
}

.chat-body {
  position: fixed;
  bottom: 5.5rem;
  right: 1.5rem;
  width: 380px;
  max-height: 520px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: var(--shadow-lg);
}

.chat-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.9rem 1.1rem;
  border-bottom: 1px solid var(--border);
  background: var(--bg-secondary);
}

.chat-title {
  font-size: 0.85rem;
  font-weight: 700;
}

.thinking-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--accent);
  animation: pulse 1s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 0.8rem;
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
  max-height: 320px;
}

.chat-empty {
  text-align: center;
  padding: 1.5rem 1rem;
  color: var(--text-muted);
  font-size: 0.85rem;
}

.quick-prompts {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
  margin-top: 0.8rem;
  justify-content: center;
}

.quick-prompt {
  padding: 0.35rem 0.8rem;
  font-size: 0.72rem;
  font-weight: 500;
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 100px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all var(--transition);
}

.quick-prompt:hover {
  border-color: var(--accent);
  color: var(--accent-hover);
}

.message {
  display: flex;
  gap: 0.6rem;
}

.msg-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  flex-shrink: 0;
  background: var(--bg-secondary);
}

.message.assistant .msg-avatar {
  background: var(--accent-gradient);
}

.msg-body {
  flex: 1;
  min-width: 0;
}

.msg-text {
  font-size: 0.82rem;
  line-height: 1.5;
  color: var(--text-primary);
  padding: 0.5rem 0.7rem;
  background: var(--bg-input);
  border-radius: var(--radius);
  border: 1px solid var(--border);
}

.message.assistant .msg-text {
  background: rgba(99, 102, 241, 0.08);
  border-color: rgba(99, 102, 241, 0.2);
}

.msg-image {
  margin-top: 0.4rem;
  border-radius: var(--radius);
  max-width: 160px;
  max-height: 120px;
  object-fit: cover;
  cursor: pointer;
  border: 1px solid var(--border);
  transition: transform var(--transition);
}

.msg-image:hover {
  transform: scale(1.05);
}

.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 0.5rem 0.7rem;
}

.typing-indicator span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--text-muted);
  animation: bounce 1.4s ease-in-out infinite;
}

.typing-indicator span:nth-child(2) { animation-delay: 0.2s; }
.typing-indicator span:nth-child(3) { animation-delay: 0.4s; }

@keyframes bounce {
  0%, 60%, 100% { transform: translateY(0); }
  30% { transform: translateY(-6px); }
}

.chat-input-area {
  padding: 0.7rem;
  border-top: 1px solid var(--border);
  background: var(--bg-secondary);
}

.input-row {
  display: flex;
  gap: 0.5rem;
  align-items: flex-end;
}

.input-row textarea {
  flex: 1;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  color: var(--text-primary);
  padding: 0.6rem 0.8rem;
  resize: none;
  font-size: 0.82rem;
  font-family: inherit;
  transition: border-color var(--transition);
}

.input-row textarea:focus {
  outline: none;
  border-color: var(--border-focus);
  box-shadow: 0 0 0 2px var(--accent-glow);
}

.input-row textarea::placeholder {
  color: var(--text-muted);
}

.send-btn {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  background: var(--accent-gradient);
  color: #fff;
  font-size: 1.1rem;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all var(--transition);
}

.send-btn:hover:not(:disabled) {
  transform: scale(1.1);
}

.send-btn:disabled {
  opacity: 0.4;
}

.input-hint {
  font-size: 0.65rem;
  color: var(--text-muted);
  margin-top: 0.3rem;
  padding-left: 0.3rem;
}

.slide-up-enter-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.slide-up-leave-active {
  transition: all 0.2s ease-in;
}
.slide-up-enter-from {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(10px) scale(0.98);
}
</style>
