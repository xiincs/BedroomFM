<template>
  <div class="profile-page">
    <!-- Topbar -->
    <div class="profile-nav">
      <router-link to="/" class="nav-back">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="15,18 9,12 15,6"/>
        </svg>
        返回首页
      </router-link>
      <div class="nav-title">个人中心</div>
      <button class="nav-logout" @click="handleLogout">退出登录</button>
    </div>

    <!-- Not logged in -->
    <div v-if="!auth.isLoggedIn" class="not-logged-in">
      <div class="nli-icon">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="var(--text3)" stroke-width="1">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/>
        </svg>
      </div>
      <div class="nli-title">请先登录</div>
      <div class="nli-sub">登录账号以查看个人中心和 VIP 成长体系</div>
      <router-link to="/" class="btn btn-primary" style="margin-top:8px">去登录</router-link>
    </div>

    <template v-else>
      <!-- Hero Card -->
      <div class="hero-section">
        <div class="hero-glow" :style="heroGlowStyle"></div>
        <div class="hero-card">
          <!-- Avatar -->
          <div class="hero-avatar-wrap">
            <div class="avatar-ring" :style="avatarRingStyle">
              <img class="hero-avatar" :src="auth.user.avatar" :alt="auth.user.nickname" />
            </div>
            <div class="avatar-level" :style="levelBadgeStyle">Lv.{{ auth.user.level }}</div>
          </div>

          <!-- Info -->
          <div class="hero-info">
            <div class="hero-nickname">{{ auth.user.nickname }}</div>
            <div class="hero-username">@{{ auth.user.username }}</div>
            <div class="hero-tier-row">
              <span class="hero-tier-badge" :style="tierBadgeStyle">
                {{ auth.tier?.name || '听众' }}
              </span>
              <span class="hero-joined">加入于 {{ joinedDate }}</span>
            </div>
          </div>

          <!-- XP Bar -->
          <div class="xp-section">
            <div class="xp-labels">
              <span class="xp-current">{{ auth.xpProgress }} / {{ auth.xpToNext }} XP</span>
              <span class="xp-total">累计 {{ auth.user.xp }} XP</span>
            </div>
            <div class="xp-track">
              <div class="xp-fill" :style="xpFillStyle">
                <div class="xp-shimmer"></div>
              </div>
            </div>
            <div class="xp-next-label" v-if="auth.user.level < 100">
              距 Lv.{{ auth.user.level + 1 }} 还需 {{ auth.xpToNext - auth.xpProgress }} XP
            </div>
            <div class="xp-next-label" v-else style="color:var(--yellow)">✦ 已达最高等级</div>
          </div>
        </div>
      </div>

      <!-- Stats Row -->
      <div class="stats-row">
        <div class="stat-card">
          <div class="stat-icon" style="color:#7C5CFA">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/>
            </svg>
          </div>
          <div class="stat-num">{{ auth.user.totalSongs }}</div>
          <div class="stat-label">点歌数</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon" style="color:#FFD166">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="currentColor" stroke="none">
              <polygon points="12,2 15.09,8.26 22,9.27 17,14.14 18.18,21.02 12,17.77 5.82,21.02 7,14.14 2,9.27 8.91,8.26"/>
            </svg>
          </div>
          <div class="stat-num">{{ auth.user.totalVotes }}</div>
          <div class="stat-label">顶歌次数</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon" style="color:#22D97A">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/>
            </svg>
          </div>
          <div class="stat-num">{{ auth.user.totalRooms }}</div>
          <div class="stat-label">加入房间</div>
        </div>
      </div>

      <!-- Current Tier Showcase -->
      <div class="tier-showcase" v-if="auth.tier">
        <div class="section-label">当前段位</div>
        <div class="tier-card" :style="tierCardStyle">
          <div class="tier-card-glow" :style="tierCardGlowStyle"></div>
          <div class="tier-header">
            <div class="tier-icon" :style="tierIconStyle">
              <svg v-if="auth.tier.id === 0" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><polygon points="10,8 16,12 10,16"/></svg>
              <svg v-else-if="auth.tier.id === 1" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><polygon points="12,2 15.09,8.26 22,9.27 17,14.14 18.18,21.02 12,17.77 5.82,21.02 7,14.14 2,9.27 8.91,8.26"/></svg>
              <svg v-else-if="auth.tier.id === 2" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/></svg>
              <svg v-else-if="auth.tier.id === 3" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2" stroke="white" stroke-width="2" stroke-linecap="round" fill="none"/></svg>
              <svg v-else-if="auth.tier.id === 4" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><polygon points="12,2 22,8.5 22,15.5 12,22 2,15.5 2,8.5"/></svg>
              <svg v-else-if="auth.tier.id === 5" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><polygon points="12,2 22,8.5 22,15.5 12,22 2,15.5 2,8.5"/><polygon points="12,6 18,9.5 18,16.5 12,20 6,16.5 6,9.5"/></svg>
              <svg v-else-if="auth.tier.id === 6" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2l2.4 7.4H22l-6.2 4.5 2.4 7.4L12 17l-6.2 4.3 2.4-7.4L2 9.4h7.6z"/></svg>
              <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><path d="M12 8v4l3 3" stroke-linecap="round"/><path d="M12 2v2M12 20v2M2 12h2M20 12h2" stroke-linecap="round"/></svg>
            </div>
            <div>
              <div class="tier-name" :style="tierNameStyle">{{ auth.tier.name }}</div>
              <div class="tier-level-range">Lv.{{ auth.tier.minLevel }} – {{ auth.tier.maxLevel === 100 ? '100 (满级)' : auth.tier.maxLevel }}</div>
            </div>
          </div>
          <div class="tier-perks">
            <div v-for="perk in auth.tier.perks" :key="perk" class="tier-perk">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <polyline points="20,6 9,17 4,12"/>
              </svg>
              {{ perk }}
            </div>
          </div>
        </div>
      </div>

      <!-- VIP Tier Roadmap -->
      <div class="roadmap-section">
        <div class="section-label">VIP 成长之路</div>
        <div class="roadmap-scroll">
          <div class="roadmap-track">
            <div
              v-for="t in (auth.allTiers.length ? auth.allTiers : FALLBACK_TIERS)"
              :key="t.id"
              class="roadmap-node"
              :class="{
                'node-current': auth.tier && t.id === auth.tier.id,
                'node-unlocked': auth.user.level >= t.minLevel && (!auth.tier || t.id !== auth.tier.id),
                'node-locked': auth.user.level < t.minLevel,
              }"
            >
              <div class="node-connector" v-if="t.id > 0"></div>
              <div class="node-dot" :style="nodeDotStyle(t)">
                <svg v-if="auth.user.level >= t.minLevel" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                  <polyline points="20,6 9,17 4,12"/>
                </svg>
                <svg v-else width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
              </div>
              <div class="node-info">
                <div class="node-name" :style="nodeNameStyle(t)">{{ t.name }}</div>
                <div class="node-range">Lv.{{ t.minLevel }}{{ t.maxLevel !== t.minLevel ? '–' + t.maxLevel : '' }}</div>
              </div>
              <div v-if="auth.tier && t.id === auth.tier.id" class="node-current-arrow">▲</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Achievements -->
      <div class="achievements-section">
        <div class="section-label">成就徽章</div>
        <div class="achievements-grid">
          <div v-for="a in achievements" :key="a.id" class="achievement-item" :class="{ locked: !a.unlocked }">
            <div class="ach-icon" :style="a.unlocked ? { color: a.color } : {}">{{ a.icon }}</div>
            <div class="ach-name">{{ a.name }}</div>
            <div class="ach-desc">{{ a.desc }}</div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()

onMounted(() => {
  auth.fetchMe()
})

const FALLBACK_TIERS = [
  { id:0, name:'听众',      minLevel:1,  maxLevel:4,   color:'#6b7280', gradient:['#6b7280','#9ca3af'] },
  { id:1, name:'小歌迷',    minLevel:5,  maxLevel:9,   color:'#3b82f6', gradient:['#3b82f6','#60a5fa'] },
  { id:2, name:'节奏达人',  minLevel:10, maxLevel:19,  color:'#06b6d4', gradient:['#06b6d4','#10b981'] },
  { id:3, name:'黄金DJ',    minLevel:20, maxLevel:34,  color:'#f59e0b', gradient:['#f59e0b','#ef4444'] },
  { id:4, name:'白金制作人',minLevel:35, maxLevel:49,  color:'#8b5cf6', gradient:['#8b5cf6','#ec4899'] },
  { id:5, name:'钻石传说',  minLevel:50, maxLevel:74,  color:'#06b6d4', gradient:['#06b6d4','#8b5cf6','#ec4899'] },
  { id:6, name:'传奇之声',  minLevel:75, maxLevel:99,  color:'#a855f7', gradient:['#f59e0b','#8b5cf6','#06b6d4'] },
  { id:7, name:'音乐之神',  minLevel:100,maxLevel:100, color:'#fbbf24', gradient:['#fbbf24','#ef4444','#8b5cf6'] },
]

function tierGrad(t) {
  const g = t?.gradient ?? ['#6b7280']
  return g.length > 1 ? `linear-gradient(135deg, ${g.join(', ')})` : g[0]
}

const heroGlowStyle = computed(() => ({
  background: `radial-gradient(ellipse 70% 50% at 50% 0%, ${(auth.tier?.color ?? '#7C5CFA')}22 0%, transparent 70%)`
}))

const avatarRingStyle = computed(() => ({
  background: `linear-gradient(135deg, ${(auth.tier?.gradient ?? ['#7C5CFA','#9B7DFF']).join(', ')})`,
  boxShadow: `0 0 24px ${(auth.tier?.color ?? '#7C5CFA')}44`
}))

const levelBadgeStyle = computed(() => ({
  background: tierGrad(auth.tier),
  boxShadow: `0 0 12px ${(auth.tier?.color ?? '#7C5CFA')}66`,
}))

const tierBadgeStyle = computed(() => ({
  background: tierGrad(auth.tier),
  boxShadow: `0 0 10px ${(auth.tier?.color ?? '#7C5CFA')}33`,
}))

const xpPercent = computed(() => {
  if (!auth.xpToNext) return 100
  return Math.min(100, Math.round((auth.xpProgress / auth.xpToNext) * 100))
})

const xpFillStyle = computed(() => ({
  width: xpPercent.value + '%',
  background: tierGrad(auth.tier),
  boxShadow: `0 0 12px ${(auth.tier?.color ?? '#7C5CFA')}66`,
}))

const tierCardStyle = computed(() => ({
  borderColor: `${(auth.tier?.color ?? '#7C5CFA')}40`,
}))

const tierCardGlowStyle = computed(() => ({
  background: `radial-gradient(ellipse 80% 80% at 50% 0%, ${(auth.tier?.color ?? '#7C5CFA')}14 0%, transparent 70%)`
}))

const tierIconStyle = computed(() => ({
  color: auth.tier?.color ?? '#7C5CFA',
  background: `${(auth.tier?.color ?? '#7C5CFA')}15`,
  borderColor: `${(auth.tier?.color ?? '#7C5CFA')}30`,
}))

const tierNameStyle = computed(() => ({
  backgroundImage: tierGrad(auth.tier),
  WebkitBackgroundClip: 'text',
  WebkitTextFillColor: 'transparent',
  backgroundClip: 'text',
}))

const joinedDate = computed(() => {
  if (!auth.user?.createdAt) return '—'
  return new Date(auth.user.createdAt * 1000).toLocaleDateString('zh-CN', { year:'numeric', month:'long', day:'numeric' })
})

function nodeDotStyle(t) {
  const unlocked = auth.user?.level >= t.minLevel
  if (!unlocked) return { background: 'var(--bg4)', borderColor: 'var(--border)' }
  const isCurrent = auth.tier && t.id === auth.tier.id
  return {
    background: tierGrad(t),
    borderColor: 'transparent',
    boxShadow: isCurrent ? `0 0 14px ${t.color}66` : 'none',
    transform: isCurrent ? 'scale(1.3)' : 'none',
  }
}

function nodeNameStyle(t) {
  const unlocked = auth.user?.level >= t.minLevel
  if (!unlocked) return { color: 'var(--text3)' }
  return {
    backgroundImage: tierGrad(t),
    WebkitBackgroundClip: 'text',
    WebkitTextFillColor: 'transparent',
    backgroundClip: 'text',
  }
}

const achievements = computed(() => {
  const u = auth.user
  if (!u) return []
  return [
    { id:'first_song', icon:'🎵', name:'第一首歌', desc:'加入第一首歌到队列', color:'#7C5CFA', unlocked: u.totalSongs >= 1 },
    { id:'song_10', icon:'🎶', name:'点歌达人', desc:'累计点歌 10 首', color:'#3b82f6', unlocked: u.totalSongs >= 10 },
    { id:'song_50', icon:'🎸', name:'点歌王', desc:'累计点歌 50 首', color:'#06b6d4', unlocked: u.totalSongs >= 50 },
    { id:'vote_first', icon:'⭐', name:'顶歌初体验', desc:'第一次为歌曲投票', color:'#FFD166', unlocked: u.totalVotes >= 1 },
    { id:'vote_20', icon:'🌟', name:'品味独到', desc:'累计顶歌 20 次', color:'#f59e0b', unlocked: u.totalVotes >= 20 },
    { id:'room_3', icon:'🏠', name:'房间常客', desc:'加入 3 个房间', color:'#22D97A', unlocked: u.totalRooms >= 3 },
    { id:'lv5', icon:'🔵', name:'初出茅庐', desc:'达到 Lv.5', color:'#3b82f6', unlocked: u.level >= 5 },
    { id:'lv20', icon:'🟡', name:'黄金时代', desc:'达到 Lv.20', color:'#f59e0b', unlocked: u.level >= 20 },
    { id:'lv50', icon:'💎', name:'钻石音乐人', desc:'达到 Lv.50', color:'#06b6d4', unlocked: u.level >= 50 },
  ]
})

function handleLogout() {
  auth.logout()
  router.push('/')
}
</script>

<style scoped>
.profile-page {
  height: 100vh;
  overflow-y: auto;
  background: var(--bg0);
  display: flex;
  flex-direction: column;
  gap: 0;
}

/* Nav */
.profile-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 24px;
  background: var(--bg1);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 10;
  flex-shrink: 0;
}
.nav-back {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text2);
  text-decoration: none;
  transition: color 0.15s;
}
.nav-back:hover { color: var(--text0); }
.nav-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text0);
}
.nav-logout {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 12px;
  color: var(--text3);
  font-family: inherit;
  transition: color 0.15s;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
}
.nav-logout:hover { color: var(--red); background: rgba(255,75,110,0.08); }

/* Not logged in */
.not-logged-in {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px;
}
.nli-icon { opacity: 0.4; }
.nli-title { font-size: 18px; font-weight: 600; color: var(--text1); }
.nli-sub { font-size: 13px; color: var(--text3); text-align: center; }

/* Hero */
.hero-section {
  position: relative;
  padding: 32px 24px 0;
}
.hero-glow {
  position: absolute;
  inset: 0;
  pointer-events: none;
}
.hero-card {
  position: relative;
  background: var(--bg2);
  border: 1px solid var(--border-active);
  border-radius: var(--radius-xl);
  padding: 28px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  overflow: hidden;
}

/* Avatar */
.hero-avatar-wrap {
  display: flex;
  align-items: flex-end;
  gap: 16px;
}
.avatar-ring {
  width: 88px;
  height: 88px;
  border-radius: 50%;
  padding: 3px;
  flex-shrink: 0;
  transition: box-shadow 0.3s;
}
.hero-avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 3px solid var(--bg2);
  display: block;
  background: var(--bg3);
}
.avatar-level {
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 700;
  color: #fff;
  text-shadow: 0 1px 3px rgba(0,0,0,0.4);
  margin-bottom: 4px;
}

/* Info */
.hero-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.hero-nickname {
  font-size: 24px;
  font-weight: 800;
  color: var(--text0);
  letter-spacing: -0.5px;
}
.hero-username {
  font-size: 13px;
  color: var(--text3);
}
.hero-tier-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 4px;
}
.hero-tier-badge {
  padding: 3px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 700;
  color: #fff;
  text-shadow: 0 1px 2px rgba(0,0,0,0.4);
}
.hero-joined {
  font-size: 11px;
  color: var(--text3);
}

/* XP */
.xp-section { display: flex; flex-direction: column; gap: 6px; }
.xp-labels {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
}
.xp-current { color: var(--text1); font-weight: 500; }
.xp-total { color: var(--text3); }
.xp-track {
  width: 100%;
  height: 8px;
  background: var(--bg4);
  border-radius: 8px;
  overflow: hidden;
}
.xp-fill {
  height: 100%;
  border-radius: 8px;
  position: relative;
  transition: width 1.2s cubic-bezier(0.4,0,0.2,1);
  min-width: 4px;
  overflow: hidden;
}
.xp-shimmer {
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, transparent 0%, rgba(255,255,255,0.3) 50%, transparent 100%);
  animation: shimmer 2s infinite;
}
@keyframes shimmer {
  from { transform: translateX(-100%); }
  to { transform: translateX(100%); }
}
.xp-next-label { font-size: 11px; color: var(--text3); }

/* Stats */
.stats-row {
  display: flex;
  gap: 12px;
  padding: 20px 24px 0;
}
.stat-card {
  flex: 1;
  background: var(--bg2);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 18px 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  text-align: center;
  transition: border-color 0.2s;
}
.stat-card:hover { border-color: var(--border-active); }
.stat-icon { opacity: 0.8; }
.stat-num {
  font-size: 26px;
  font-weight: 800;
  color: var(--text0);
  line-height: 1;
}
.stat-label { font-size: 11px; color: var(--text3); }

/* Tier showcase */
.tier-showcase { padding: 20px 24px 0; }
.section-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--text3);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  margin-bottom: 12px;
}
.tier-card {
  position: relative;
  background: var(--bg2);
  border: 1px solid;
  border-radius: var(--radius-xl);
  padding: 24px;
  overflow: hidden;
}
.tier-card-glow {
  position: absolute;
  inset: 0;
  pointer-events: none;
}
.tier-header {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 18px;
}
.tier-icon {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  border: 1px solid;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.tier-name {
  font-size: 22px;
  font-weight: 800;
  letter-spacing: -0.3px;
}
.tier-level-range { font-size: 12px; color: var(--text3); margin-top: 2px; }
.tier-perks { display: flex; flex-direction: column; gap: 10px; }
.tier-perk {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: var(--text1);
}
.tier-perk svg { color: var(--green); flex-shrink: 0; }

/* Roadmap */
.roadmap-section { padding: 20px 24px 0; }
.roadmap-scroll { overflow-x: auto; padding-bottom: 8px; }
.roadmap-track {
  display: flex;
  align-items: flex-start;
  gap: 0;
  min-width: max-content;
  padding: 8px 0 4px;
}
.roadmap-node {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  width: 96px;
  position: relative;
  flex-shrink: 0;
}
.node-connector {
  position: absolute;
  top: 18px;
  right: 50%;
  width: 100%;
  height: 2px;
  background: var(--border);
  z-index: 0;
}
.node-unlocked .node-connector,
.node-current .node-connector {
  background: var(--bg4);
}
.node-dot {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 2px solid;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 1;
  color: #fff;
  transition: transform 0.3s, box-shadow 0.3s;
}
.node-info { text-align: center; }
.node-name {
  font-size: 11px;
  font-weight: 600;
  line-height: 1.3;
}
.node-range { font-size: 10px; color: var(--text3); margin-top: 2px; }
.node-current-arrow {
  font-size: 10px;
  color: var(--accent-light);
  margin-top: -4px;
}

/* Achievements */
.achievements-section { padding: 20px 24px 32px; }
.achievements-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}
.achievement-item {
  background: var(--bg2);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 16px 12px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  transition: border-color 0.2s, transform 0.2s;
}
.achievement-item:not(.locked):hover {
  border-color: var(--border-active);
  transform: translateY(-2px);
}
.achievement-item.locked { opacity: 0.35; filter: grayscale(1); }
.ach-icon { font-size: 26px; line-height: 1; }
.ach-name { font-size: 11px; font-weight: 600; color: var(--text1); }
.ach-desc { font-size: 10px; color: var(--text3); line-height: 1.4; }

/* Responsive */
@media (max-width: 640px) {
  .hero-section, .stats-row, .tier-showcase, .roadmap-section, .achievements-section {
    padding-left: 16px;
    padding-right: 16px;
  }
  .hero-card { padding: 20px; }
  .hero-nickname { font-size: 20px; }
  .avatar-ring { width: 72px; height: 72px; }
  .achievements-grid { grid-template-columns: repeat(3, 1fr); gap: 8px; }
  .achievement-item { padding: 12px 8px; }
}
</style>
