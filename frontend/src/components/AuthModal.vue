<template>
  <Transition name="modal">
    <div v-if="modelValue" class="auth-overlay" @click.self="close">
      <div class="auth-box">
        <!-- Close -->
        <button class="auth-close" @click="close">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>

        <!-- Brand -->
        <div class="auth-brand">
          <div class="auth-logo">
            <svg width="28" height="28" viewBox="0 0 32 32" fill="none">
              <rect width="32" height="32" rx="10" fill="#7C5CFA"/>
              <path d="M8 22V12l12-4v10M8 22a3 3 0 1 0 6 0 3 3 0 0 0-6 0zm12-4a3 3 0 1 0 6 0 3 3 0 0 0-6 0z"
                stroke="white" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <div>
            <div class="auth-brand-name">Bedroom FM</div>
            <div class="auth-brand-sub">音乐账号 · VIP 成长体系</div>
          </div>
        </div>

        <!-- Tabs -->
        <div class="auth-tabs">
          <button class="auth-tab" :class="{ active: tab === 'login' }" @click="tab = 'login'; error = ''">登录</button>
          <button class="auth-tab" :class="{ active: tab === 'register' }" @click="tab = 'register'; error = ''">注册</button>
          <div class="auth-tab-bar" :style="{ left: tab === 'login' ? '0' : '50%' }"></div>
        </div>

        <!-- Success state -->
        <div v-if="success" class="auth-success">
          <div class="success-icon">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#22D97A" stroke-width="1.5">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/>
            </svg>
          </div>
          <div class="success-title">{{ tab === 'register' ? '注册成功！' : '欢迎回来！' }}</div>
          <div class="success-sub">{{ successMsg }}</div>
          <div class="success-xp" v-if="tab === 'register'">
            <span class="xp-tag">+50 XP</span> 新人注册奖励已到账
          </div>
        </div>

        <!-- Login Form -->
        <div v-else-if="tab === 'login'" class="auth-form">
          <div class="auth-field">
            <label>用户名</label>
            <input
              v-model="loginForm.username"
              class="input"
              placeholder="输入用户名"
              maxlength="20"
              autocomplete="username"
              @keyup.enter="doLogin"
            />
          </div>
          <div class="auth-field">
            <label>密码</label>
            <input
              v-model="loginForm.password"
              type="password"
              class="input"
              placeholder="输入密码"
              autocomplete="current-password"
              @keyup.enter="doLogin"
            />
          </div>
          <div v-if="error" class="auth-error">{{ error }}</div>
          <button class="btn btn-primary auth-submit" @click="doLogin" :disabled="loading">
            <span v-if="loading" class="spinner"></span>
            {{ loading ? '登录中...' : '登录' }}
          </button>
          <div class="auth-switch">
            还没有账号？<button @click="tab = 'register'; error = ''">立即注册</button>
          </div>
        </div>

        <!-- Register Form -->
        <div v-else class="auth-form">
          <div class="auth-field">
            <label>用户名 <span class="field-hint">3–20 个字符，用于登录</span></label>
            <input
              v-model="regForm.username"
              class="input"
              placeholder="设置登录用户名"
              maxlength="20"
              autocomplete="username"
              @keyup.enter="doRegister"
            />
          </div>
          <div class="auth-field">
            <label>昵称 <span class="field-hint">在房间中显示</span></label>
            <input
              v-model="regForm.nickname"
              class="input"
              placeholder="给自己起个好听的名字"
              maxlength="16"
              @keyup.enter="doRegister"
            />
          </div>
          <div class="auth-field">
            <label>密码 <span class="field-hint">至少 6 位</span></label>
            <input
              v-model="regForm.password"
              type="password"
              class="input"
              placeholder="设置登录密码"
              autocomplete="new-password"
              @keyup.enter="doRegister"
            />
          </div>
          <div v-if="error" class="auth-error">{{ error }}</div>
          <button class="btn btn-primary auth-submit" @click="doRegister" :disabled="loading">
            <span v-if="loading" class="spinner"></span>
            {{ loading ? '注册中...' : '创建账号' }}
          </button>
          <div class="auth-perks">
            <span class="perk-tag">+50 XP 新人礼包</span>
            <span class="perk-tag">VIP 成长体系</span>
            <span class="perk-tag">专属头衔</span>
          </div>
          <div class="auth-switch">
            已有账号？<button @click="tab = 'login'; error = ''">立即登录</button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

defineProps({ modelValue: Boolean })
const emit = defineEmits(['update:modelValue'])

const auth = useAuthStore()

const tab = ref('login')
const error = ref('')
const loading = ref(false)
const success = ref(false)
const successMsg = ref('')

const loginForm = ref({ username: '', password: '' })
const regForm = ref({ username: '', password: '', nickname: '' })

function close() {
  if (success.value) {
    success.value = false
  }
  emit('update:modelValue', false)
}

async function doLogin() {
  if (!loginForm.value.username || !loginForm.value.password) {
    error.value = '请填写用户名和密码'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const data = await auth.login(loginForm.value.username, loginForm.value.password)
    successMsg.value = `${data.user.nickname}，欢迎回到 Bedroom FM`
    success.value = true
    setTimeout(close, 1800)
  } catch (e) {
    error.value = e?.response?.data?.error || '登录失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

async function doRegister() {
  if (!regForm.value.username || !regForm.value.password || !regForm.value.nickname) {
    error.value = '请完整填写所有字段'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const data = await auth.register(regForm.value.username, regForm.value.password, regForm.value.nickname)
    successMsg.value = `${data.user.nickname}，账号已创建！起步 Lv.${data.user.level} · ${data.tier?.name}`
    success.value = true
    setTimeout(close, 2200)
  } catch (e) {
    error.value = e?.response?.data?.error || '注册失败，请稍后重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.7);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 16px;
}

.auth-box {
  position: relative;
  width: 100%;
  max-width: 400px;
  background: var(--bg2);
  border: 1px solid var(--border-active);
  border-radius: var(--radius-xl);
  padding: 28px 28px 24px;
  box-shadow: 0 24px 80px rgba(0,0,0,0.6), 0 0 0 1px rgba(124,92,250,0.1);
}

.auth-close {
  position: absolute;
  top: 16px;
  right: 16px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--text3);
  padding: 4px;
  border-radius: 6px;
  display: flex;
  transition: color 0.15s;
}
.auth-close:hover { color: var(--text1); }

.auth-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}
.auth-logo {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: rgba(124,92,250,0.1);
  border: 1px solid rgba(124,92,250,0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.auth-brand-name {
  font-size: 16px;
  font-weight: 700;
  color: var(--text0);
}
.auth-brand-sub {
  font-size: 11px;
  color: var(--text3);
  margin-top: 2px;
}

.auth-tabs {
  position: relative;
  display: flex;
  border-bottom: 1px solid var(--border);
  margin-bottom: 20px;
}
.auth-tab {
  flex: 1;
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 10px 0;
  font-size: 13px;
  font-weight: 500;
  color: var(--text3);
  font-family: inherit;
  transition: color 0.2s;
}
.auth-tab.active { color: var(--accent-light); }
.auth-tab-bar {
  position: absolute;
  bottom: -1px;
  width: 50%;
  height: 2px;
  background: var(--accent);
  border-radius: 2px;
  transition: left 0.25s cubic-bezier(0.4,0,0.2,1);
}

.auth-form { display: flex; flex-direction: column; gap: 14px; }

.auth-field { display: flex; flex-direction: column; gap: 6px; }
.auth-field label {
  font-size: 11px;
  font-weight: 500;
  color: var(--text2);
  display: flex;
  align-items: center;
  gap: 6px;
}
.field-hint {
  color: var(--text3);
  font-weight: 400;
}

.auth-error {
  padding: 8px 12px;
  background: rgba(255,75,110,0.08);
  border: 1px solid rgba(255,75,110,0.2);
  border-radius: var(--radius-sm);
  font-size: 12px;
  color: var(--red);
}

.auth-submit {
  width: 100%;
  height: 40px;
  font-size: 14px;
  font-weight: 600;
  margin-top: 2px;
}

.auth-switch {
  text-align: center;
  font-size: 12px;
  color: var(--text3);
}
.auth-switch button {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--accent-light);
  font-size: 12px;
  font-family: inherit;
  padding: 0;
}
.auth-switch button:hover { text-decoration: underline; }

.auth-perks {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: center;
}
.perk-tag {
  padding: 3px 9px;
  border-radius: 20px;
  font-size: 10px;
  font-weight: 500;
  background: rgba(124,92,250,0.12);
  color: var(--accent-light);
  border: 1px solid rgba(124,92,250,0.2);
}

/* Success state */
.auth-success {
  text-align: center;
  padding: 12px 0 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}
.success-icon { margin-bottom: 4px; }
.success-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text0);
}
.success-sub {
  font-size: 13px;
  color: var(--text2);
}
.success-xp {
  font-size: 12px;
  color: var(--text2);
  display: flex;
  align-items: center;
  gap: 8px;
}
.xp-tag {
  padding: 2px 8px;
  border-radius: 20px;
  background: rgba(34,217,122,0.12);
  color: var(--green);
  font-size: 11px;
  font-weight: 600;
}

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  flex-shrink: 0;
}

.modal-enter-active, .modal-leave-active { transition: opacity 0.2s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
.modal-enter-active .auth-box, .modal-leave-active .auth-box { transition: transform 0.2s ease; }
.modal-enter-from .auth-box { transform: scale(0.96) translateY(8px); }
.modal-leave-to .auth-box { transform: scale(0.96) translateY(8px); }
</style>
