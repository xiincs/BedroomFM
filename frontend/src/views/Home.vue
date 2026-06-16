<template>
  <div class="home">
    <div class="home-bg"></div>

    <!-- User chip (top-right, logged-in only) -->
    <div class="auth-bar" v-if="auth.isLoggedIn">
      <div class="auth-user-chip" @click="router.push('/profile')">
        <img class="chip-avatar" :src="auth.user.avatar" :alt="auth.user.nickname" />
        <span class="chip-name">{{ auth.user.nickname }}</span>
        <span class="chip-badge" :style="chipBadgeStyle">Lv.{{ auth.user.level }}</span>
      </div>
    </div>

    <div class="home-content">
      <!-- Brand -->
      <div class="brand">
        <div class="brand-icon">
          <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
            <rect width="32" height="32" rx="10" fill="#7C5CFA"/>
            <path d="M8 22V12l12-4v10M8 22a3 3 0 1 0 6 0 3 3 0 0 0-6 0zm12-4a3 3 0 1 0 6 0 3 3 0 0 0-6 0z" stroke="white" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <div>
          <div class="brand-name">Bedroom FM</div>
          <div class="brand-sub">宿舍共享音乐房间</div>
        </div>
      </div>

      <!-- ── NOT logged in: login gate ── -->
      <template v-if="!auth.isLoggedIn">
        <div class="gate-card">
          <div class="gate-icon">
            <svg width="34" height="34" viewBox="0 0 24 24" fill="none" stroke="var(--accent-light)" stroke-width="1.5">
              <path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/>
            </svg>
          </div>
          <div class="gate-title">开始音乐之旅</div>
          <div class="gate-sub">登录后畅享实时点歌、弹幕互动与 VIP 成长体系</div>
          <button class="btn btn-primary gate-btn" @click="showAuth = true">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4M10 17l5-5-5-5M15 12H3"/>
            </svg>
            登录 / 注册
          </button>
        </div>
      </template>

      <!-- ── Logged in ── -->
      <template v-else>
        <!-- Checking -->
        <div v-if="checking" class="rejoin-checking">
          <div class="spinner-sm"></div>
          <span>检查房间状态…</span>
        </div>

        <!-- Active room card -->
        <div v-else-if="activeRoom" class="rejoin-card">
          <div class="rejoin-header">
            <span class="dot-pulse-green"></span>
            <span class="rejoin-label">进行中的房间</span>
          </div>
          <div class="rejoin-room-name">{{ activeRoom.room.name }}</div>
          <div class="rejoin-meta">
            <span>{{ activeRoom.room.members?.length || 0 }} 人在线</span>
            <span class="rejoin-code">{{ activeRoom.room.code }}</span>
            <span v-if="isMyHost" class="badge badge-orange" style="font-size:10px;padding:1px 6px;">DJ</span>
          </div>
          <div class="rejoin-actions">
            <button class="btn btn-primary" @click="rejoin">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="currentColor"><polygon points="5,3 19,12 5,21"/></svg>
              回到房间
            </button>
            <button class="btn btn-ghost btn-sm" @click="leaveRoom">离开此房间</button>
          </div>
        </div>

        <!-- Divider when rejoin card is present -->
        <div class="section-sep" v-if="!checking && activeRoom">
          <span>或加入 / 创建其他房间</span>
        </div>

        <!-- Create / Join cards -->
        <div class="home-cards" :class="{ compact: !checking && activeRoom }">
          <!-- Create Room -->
          <div class="home-card">
            <div class="home-card-header">
              <span class="home-card-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 5v14M5 12h14"/>
                </svg>
              </span>
              <span>创建房间</span>
            </div>
            <div class="home-card-body">
              <div class="field">
                <label>房间名称</label>
                <input v-model="createForm.name" class="input" placeholder="505宿舍音乐房" maxlength="30" @keyup.enter="create"/>
              </div>
              <div class="field">
                <label>你的昵称</label>
                <input v-model="createForm.nickname" class="input" placeholder="DJ小明" maxlength="16" @keyup.enter="create"/>
              </div>
              <button class="btn btn-primary" style="width:100%" :disabled="creating" @click="create">
                <span v-if="creating" class="spinner"></span>
                <span>{{ creating ? '创建中...' : '创建房间' }}</span>
              </button>
              <div v-if="createError" class="field-error">{{ createError }}</div>
            </div>
          </div>

          <div class="home-divider"><span>或</span></div>

          <!-- Join Room -->
          <div class="home-card">
            <div class="home-card-header">
              <span class="home-card-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4M10 17l5-5-5-5M15 12H3"/>
                </svg>
              </span>
              <span>加入房间</span>
            </div>
            <div class="home-card-body">
              <div class="field">
                <label>邀请码</label>
                <input v-model="joinForm.code" class="input code-input" placeholder="A7KD9" maxlength="5" @keyup.enter="join"/>
              </div>
              <div class="field">
                <label>你的昵称</label>
                <input v-model="joinForm.nickname" class="input" placeholder="阿豪" maxlength="16" @keyup.enter="join"/>
              </div>
              <button class="btn btn-ghost" style="width:100%" :disabled="joining" @click="join">
                <span v-if="joining" class="spinner"></span>
                <span>{{ joining ? '加入中...' : '加入房间' }}</span>
              </button>
              <div v-if="joinError" class="field-error">{{ joinError }}</div>
            </div>
          </div>
        </div>
      </template>

      <!-- Features (shown when no active room or not logged in) -->
      <div class="home-features" v-if="!auth.isLoggedIn || (!checking && !activeRoom)">
        <div class="feat"><span class="feat-dot green"></span>实时同步播放</div>
        <div class="feat"><span class="feat-dot purple"></span>轮流点歌机制</div>
        <div class="feat"><span class="feat-dot orange"></span>顶歌 / 切歌投票</div>
        <div class="feat"><span class="feat-dot red"></span>表情轰炸互动</div>
      </div>
    </div>

    <AuthModal v-model="showAuth" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useRoomStore } from '../stores/room'
import { useAuthStore } from '../stores/auth'
import AuthModal from '../components/AuthModal.vue'

const router = useRouter()
const store = useRoomStore()
const auth = useAuthStore()

const showAuth = ref(false)
const activeRoom = ref(null)
const checking = ref(false)

const chipBadgeStyle = computed(() => {
  const g = auth.tier?.gradient
  if (!g?.length) return { background: '#7C5CFA' }
  return { background: g.length > 1 ? `linear-gradient(90deg, ${g.join(', ')})` : g[0] }
})

const isMyHost = computed(() => {
  if (!activeRoom.value) return false
  return activeRoom.value.room.hostId === activeRoom.value.memberId
})

const createForm = ref({ name: '', nickname: '' })
const joinForm = ref({ code: '', nickname: '' })
const creating = ref(false)
const joining = ref(false)
const createError = ref('')
const joinError = ref('')

async function checkRoom() {
  checking.value = true
  activeRoom.value = await store.checkActiveRoom()
  checking.value = false
}

function syncNickname() {
  const nick = auth.user?.nickname || ''
  createForm.value.nickname = nick
  joinForm.value.nickname = nick
}

onMounted(async () => {
  if (auth.isLoggedIn) {
    await auth.fetchMe()
    syncNickname()
    await checkRoom()
  }
})

watch(() => auth.isLoggedIn, async (loggedIn) => {
  if (loggedIn) {
    syncNickname()
    await checkRoom()
  }
})

function rejoin() {
  store.roomId = activeRoom.value.room.id
  store.memberId = activeRoom.value.memberId
  router.push(`/room/${activeRoom.value.room.id}`)
}

async function leaveRoom() {
  await store.leaveRoom()  // removes member server-side, clears localStorage
  activeRoom.value = null
}

async function create() {
  if (!createForm.value.name || !createForm.value.nickname) {
    createError.value = '请填写房间名称和昵称'
    return
  }
  creating.value = true
  createError.value = ''
  try {
    const data = await store.createRoom(createForm.value.name, createForm.value.nickname)
    auth.gainXP('join_room')
    activeRoom.value = null
    router.push(`/room/${data.roomId}`)
  } catch (e) {
    createError.value = '创建失败，请检查服务是否启动'
  } finally {
    creating.value = false
  }
}

async function join() {
  if (!joinForm.value.code || !joinForm.value.nickname) {
    joinError.value = '请填写邀请码和昵称'
    return
  }
  joining.value = true
  joinError.value = ''
  try {
    const data = await store.joinRoom(joinForm.value.code, joinForm.value.nickname)
    auth.gainXP('join_room')
    activeRoom.value = null
    router.push(`/room/${data.roomId}`)
  } catch (e) {
    joinError.value = e?.response?.data?.error || '加入失败，邀请码错误'
  } finally {
    joining.value = false
  }
}
</script>

<style scoped>
/* Auth bar */
.auth-bar {
  position: absolute;
  top: 16px;
  right: 20px;
  z-index: 10;
}

.auth-user-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 5px 10px 5px 5px;
  background: var(--bg2);
  border: 1px solid var(--border-active);
  border-radius: 24px;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s;
}
.auth-user-chip:hover { background: var(--bg3); border-color: var(--accent); }
.chip-avatar { width: 26px; height: 26px; border-radius: 50%; object-fit: cover; }
.chip-name {
  font-size: 12px;
  font-weight: 500;
  color: var(--text1);
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.chip-badge {
  padding: 2px 7px;
  border-radius: 20px;
  font-size: 10px;
  font-weight: 700;
  color: #fff;
  text-shadow: 0 1px 2px rgba(0,0,0,0.3);
}

/* Page layout */
.home {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.home-bg {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(ellipse 60% 50% at 30% 40%, rgba(124,92,250,0.15) 0%, transparent 70%),
    radial-gradient(ellipse 40% 40% at 70% 60%, rgba(34,217,122,0.06) 0%, transparent 70%);
}

.home-content {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 28px;
  width: 100%;
  max-width: 720px;
  padding: 24px;
  max-height: 100vh;
  overflow-y: auto;
}

/* Brand */
.brand {
  display: flex;
  align-items: center;
  gap: 14px;
}
.brand-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  border-radius: 14px;
  background: rgba(124,92,250,0.1);
  border: 1px solid rgba(124,92,250,0.2);
}
.brand-name {
  font-size: 24px;
  font-weight: 700;
  color: var(--text0);
  letter-spacing: -0.5px;
}
.brand-sub {
  font-size: 13px;
  color: var(--text2);
  margin-top: 2px;
}

/* Gate card (not logged in) */
.gate-card {
  background: var(--bg2);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 32px 28px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  width: 100%;
  max-width: 380px;
  text-align: center;
}
.gate-icon {
  width: 64px;
  height: 64px;
  border-radius: var(--radius-lg);
  background: rgba(124,92,250,0.08);
  border: 1px solid rgba(124,92,250,0.18);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 4px;
}
.gate-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text0);
}
.gate-sub {
  font-size: 13px;
  color: var(--text2);
  line-height: 1.65;
  max-width: 280px;
}
.gate-btn {
  width: 100%;
  margin-top: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

/* Checking state */
.rejoin-checking {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: var(--text3);
}
.spinner-sm {
  width: 16px;
  height: 16px;
  border: 2px solid var(--bg4);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

/* Rejoin card */
.rejoin-card {
  width: 100%;
  background: linear-gradient(135deg, rgba(124,92,250,0.1) 0%, rgba(34,217,122,0.05) 100%);
  border: 1px solid rgba(124,92,250,0.28);
  border-radius: var(--radius-lg);
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.rejoin-header {
  display: flex;
  align-items: center;
  gap: 7px;
  font-size: 11px;
  font-weight: 600;
  color: var(--green);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}
.dot-pulse-green {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--green);
  animation: pulse 2s ease-in-out infinite;
  flex-shrink: 0;
}
.rejoin-room-name {
  font-size: 22px;
  font-weight: 700;
  color: var(--text0);
  letter-spacing: -0.3px;
}
.rejoin-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text2);
}
.rejoin-code {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.12em;
  padding: 2px 8px;
  border-radius: 4px;
  background: rgba(124,92,250,0.15);
  color: var(--accent-light);
}
.rejoin-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 6px;
}

/* Section separator */
.section-sep {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  font-size: 12px;
  color: var(--text3);
}
.section-sep::before,
.section-sep::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border);
}

/* Create/Join cards */
.home-cards {
  display: flex;
  align-items: stretch;
  gap: 0;
  width: 100%;
}
.home-cards.compact .home-card-body {
  padding: 14px 20px;
  gap: 10px;
}
.home-cards.compact .field {
  gap: 4px;
}

.home-card {
  flex: 1;
  background: var(--bg2);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.home-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 20px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text0);
  border-bottom: 1px solid var(--border);
  background: var(--bg3);
}
.home-card-icon {
  display: flex;
  align-items: center;
  color: var(--accent-light);
}

.home-card-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.field label {
  font-size: 12px;
  font-weight: 500;
  color: var(--text2);
  letter-spacing: 0.02em;
}
.field-error {
  font-size: 12px;
  color: var(--red);
}

.code-input {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.3em;
  text-transform: uppercase;
  text-align: center;
}

.home-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  flex-shrink: 0;
  position: relative;
}
.home-divider span {
  font-size: 12px;
  color: var(--text3);
  background: var(--bg0);
  padding: 4px 0;
  z-index: 1;
  writing-mode: vertical-rl;
}
.home-divider::before {
  content: '';
  position: absolute;
  top: 0;
  bottom: 0;
  left: 50%;
  width: 1px;
  background: var(--border);
}

/* Features */
.home-features {
  display: flex;
  align-items: center;
  gap: 24px;
}
.feat {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text2);
}
.feat-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}
.feat-dot.green  { background: var(--green); }
.feat-dot.purple { background: var(--accent); }
.feat-dot.orange { background: var(--orange); }
.feat-dot.red    { background: var(--red); }

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

/* ====== Mobile ====== */
@media (max-width: 768px) {
  .home {
    height: 100vh;
    height: 100dvh;
    align-items: flex-start;
  }

  .home-content {
    padding: 20px 16px 32px;
    gap: 20px;
    max-width: 100%;
  }

  .brand {
    flex-direction: column;
    text-align: center;
    gap: 10px;
  }
  .brand-icon {
    width: 44px;
    height: 44px;
    border-radius: 12px;
  }
  .brand-icon svg { width: 26px; height: 26px; }
  .brand-name { font-size: 20px; }
  .brand-sub  { font-size: 12px; }

  .gate-card { max-width: 100%; padding: 24px 20px; }
  .gate-title { font-size: 17px; }
  .gate-sub   { font-size: 12px; }

  .rejoin-card { padding: 16px; }
  .rejoin-room-name { font-size: 18px; }

  .home-cards {
    flex-direction: column;
  }
  .home-card { width: 100%; }
  .home-card-body { padding: 14px; gap: 10px; }

  .home-divider {
    width: 100%;
    height: auto;
    padding: 6px 0;
    justify-content: center;
  }
  .home-divider span {
    writing-mode: horizontal-tb;
    padding: 2px 8px;
  }
  .home-divider::before {
    top: 50%;
    bottom: auto;
    left: 0;
    right: 0;
    width: 100%;
    height: 1px;
  }

  .home-features {
    flex-wrap: wrap;
    justify-content: center;
    gap: 12px;
  }
  .feat { font-size: 11px; }

  .input { height: 40px; font-size: 14px; }
  .code-input { font-size: 18px; }
  .btn { height: 40px; font-size: 14px; }

  .auth-bar { top: 12px; right: 12px; }
}
</style>
