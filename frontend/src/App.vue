<template>
  <div class="app-shell">
    <header class="app-header">
      <div class="header-left">
        <div class="logo">
          <span class="logo-icon">&#10038;</span>
          <span class="logo-text">AI Image Editor</span>
        </div>
      </div>
      <div class="header-right">
        <span class="badge">Beta</span>
      </div>
    </header>

    <main class="app-main">
      <aside class="sidebar">
        <UploadPage />
      </aside>
      <section class="editor-area">
        <ChatEditor />
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import UploadPage from './pages/UploadPage.vue'
import ChatEditor from './pages/ChatEditor.vue'
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1.8rem;
  height: 3.6rem;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
  backdrop-filter: blur(12px);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.logo-icon {
  font-size: 1.5rem;
  background: var(--accent-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.logo-text {
  font-size: 1.1rem;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--text-primary);
}

.badge {
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  padding: 0.25rem 0.6rem;
  border-radius: 100px;
  background: var(--accent-gradient);
  color: #fff;
}

.app-main {
  flex: 1;
  display: grid;
  grid-template-columns: 380px 1fr;
  gap: 1.5rem;
  padding: 1.5rem;
  max-width: 1440px;
  width: 100%;
  margin: 0 auto;
  align-items: start;
}

.sidebar {
  position: sticky;
  top: calc(3.6rem + 1.5rem);
}

.editor-area {
  min-height: calc(100vh - 3.6rem - 3rem - 3rem);
}

@media (max-width: 900px) {
  .app-main {
    grid-template-columns: 1fr;
  }
  .sidebar {
    position: static;
  }
}
</style>
