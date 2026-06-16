<template>
  <div class="room-layout" :class="'tab-' + mobileTab" v-if="store.room">
    <!-- Danmaku layer -->
    <div class="danmaku-layer" aria-hidden="true">
      <div
        v-for="d in danmaku"
        :key="d.id"
        class="danmaku-item"
        :style="{ top: d.top + 'px', color: d.color, animationDuration: d.duration + 's' }"
      >{{ d.text }}</div>
    </div>

    <!-- Floating reactions -->
    <div class="reaction-layer" aria-hidden="true">
      <TransitionGroup name="float">
        <div
          v-for="r in store.reactions"
          :key="r.id"
          class="float-reaction"
          :style="{ left: r.x + '%' }"
        >{{ r.emoji }}</div>
      </TransitionGroup>
    </div>

    <!-- LEFT: Queue + Members -->
    <aside class="sidebar left-sidebar">
      <!-- Room Header -->
      <div class="room-header">
        <div class="room-title-row">
          <div class="room-name">{{ store.room.name }}</div>
          <div class="room-code-wrap">
            <span class="room-code badge badge-purple" :title="codeCopied ? '已复制' : '点击复制'" @click="copyCode" style="cursor:pointer">{{ store.room.code }}</span>
            <span v-if="codeCopied" class="copy-tip">已复制</span>
          </div>
        </div>
        <div class="room-meta">
          <span class="dot-green"></span>
          <span>{{ memberList.length }} 人在线</span>
          <span v-if="store.isHost" class="badge badge-orange ml-2">DJ</span>
          <button class="vip-btn" :class="{ active: hasCookie }" @click="showVipSettings = true" title="VIP设置">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
            {{ hasCookie ? 'VIP' : '设置VIP' }}
          </button>
        </div>
      </div>
      <div class="divider"></div>

      <!-- Members -->
      <div class="section-header">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/></svg>
        成员
      </div>
      <div class="member-list">
        <div v-for="m in memberList" :key="m.id" class="member-item">
          <div class="member-avatar" :style="{ background: memberColor(m.nickname) }">
            {{ m.nickname[0] }}
          </div>
          <div class="member-info">
            <div class="member-name">
              {{ m.nickname }}
              <span v-if="m.id === store.room.hostId" class="badge badge-orange" style="font-size:10px;padding:1px 5px;">DJ</span>
            </div>
            <div class="member-persona">{{ m.persona }}</div>
          </div>
          <div class="member-coins">
            <svg width="10" height="10" viewBox="0 0 24 24" fill="var(--yellow)" stroke="none"><circle cx="12" cy="12" r="10"/></svg>
            {{ m.coins }}
          </div>
        </div>
      </div>
      <div class="divider"></div>

      <!-- Queue -->
      <div class="section-header" style="justify-content: space-between">
        <div class="flex-row">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>
          队列 <span class="text-muted">({{ store.room.queue?.length || 0 }})</span>
        </div>
        <button class="btn btn-xs btn-primary" @click="showSearch = true">+ 点歌</button>
      </div>

      <div class="queue-list scroll-y">
        <div v-if="!store.room.queue?.length" class="empty-state">
          <div>队列为空</div>
          <div class="text-muted">点击「+ 点歌」搜索歌曲</div>
        </div>
        <div v-for="(item, idx) in store.room.queue" :key="item.qid" class="queue-item fade-slide-up">
          <div class="queue-rank">{{ idx + 1 }}</div>
          <div class="queue-cover">
            <img :src="item.song.cover" :alt="item.song.name" />
          </div>
          <div class="queue-info">
            <div class="queue-song truncate">{{ item.song.name }}</div>
            <div class="queue-meta">
              <span class="text-muted">{{ item.song.artist }}</span>
              <span class="queue-adder">{{ item.adder }}</span>
            </div>
          </div>
          <div class="queue-actions">
            <div class="vote-score" v-if="item.votes > 0">
              <svg width="9" height="9" viewBox="0 0 24 24" fill="var(--yellow)"><polygon points="12,2 15.09,8.26 22,9.27 17,14.14 18.18,21.02 12,17.77 5.82,21.02 7,14.14 2,9.27 8.91,8.26"/></svg>
              {{ item.votes }}
            </div>
            <div class="vote-btns">
              <button class="btn btn-xs btn-ghost" @click.stop="voteUp(item.qid, 1)">+1</button>
              <button class="btn btn-xs btn-ghost" @click.stop="voteUp(item.qid, 3)">+3</button>
              <button class="btn btn-xs btn-ghost" @click.stop="voteUp(item.qid, 5)">+5</button>
              <button
                v-if="item.addedBy === store.memberId || store.isHost"
                class="btn btn-xs btn-danger"
                @click.stop="store.sendQueueRemove(item.qid)"
              >
                <svg width="9" height="9" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </aside>

    <!-- CENTER: Now Playing -->
    <main class="center-panel">
      <div class="now-playing" v-if="store.currentSong">
        <!-- Cover -->
        <div class="cover-wrap">
          <div class="cover-glow" :style="coverGlowStyle"></div>
          <div class="cover-img-wrap" :class="{ spinning: store.room.playback?.isPlaying }">
            <img class="cover-img" :src="store.currentSong.cover" :alt="store.currentSong.name" />
          </div>
          <!-- Skip vote overlay -->
          <div v-if="skipVoteActive" class="skip-overlay">
            <div class="skip-count">{{ skipVoteCount }}/{{ memberList.length }}</div>
            <div class="skip-label">人投票切歌</div>
          </div>
        </div>

        <!-- Song Info -->
        <div class="song-info">
          <div class="song-name">{{ store.currentSong.name }}</div>
          <div class="song-artist">{{ store.currentSong.artist }}</div>
          <div class="song-album text-muted">{{ store.currentSong.album }}</div>
        </div>

        <!-- Progress Bar -->
        <div class="progress-section">
          <span class="time-text">{{ formatTime(currentPosition) }}</span>
          <div class="progress-bar" @click="seekTo">
            <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
            <div class="progress-thumb" :style="{ left: progressPercent + '%' }"></div>
          </div>
          <span class="time-text">{{ formatTime(store.currentSong.duration) }}</span>
        </div>

        <!-- Controls -->
        <div class="controls">
          <button class="btn btn-ghost btn-icon" @click="store.sendVoteSkip()" title="发起切歌投票">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="5,4 15,12 5,20"/><line x1="19" y1="5" x2="19" y2="19"/>
            </svg>
          </button>

          <button v-if="store.isHost" class="ctrl-btn" :class="{ playing: isPlaying }" @click="togglePlay">
            <svg v-if="!isPlaying" width="22" height="22" viewBox="0 0 24 24" fill="currentColor"><polygon points="5,3 19,12 5,21"/></svg>
            <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="currentColor"><rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/></svg>
          </button>
          <div v-else class="ctrl-btn ctrl-btn-indicator">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/>
            </svg>
          </div>

          <button v-if="store.isHost" class="btn btn-ghost btn-icon" @click="store.sendNextSong()" title="下一首">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="5,4 15,12 5,20"/><line x1="19" y1="5" x2="19" y2="19"/>
            </svg>
          </button>

          <!-- Local audio toggle: non-host only -->
          <button
            v-if="!store.isHost"
            class="btn btn-ghost btn-icon local-audio-btn"
            :class="{ active: localAudio }"
            @click="toggleLocalAudio"
            :title="localAudio ? '关闭本机声音' : '开启本机声音'"
          >
            <svg v-if="localAudio" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="11,5 6,9 2,9 2,15 6,15 11,19"/><path d="M15.54 8.46a5 5 0 0 1 0 7.07"/><path d="M19.07 4.93a10 10 0 0 1 0 14.14"/>
            </svg>
            <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="11,5 6,9 2,9 2,15 6,15 11,19"/><line x1="23" y1="9" x2="17" y2="15"/><line x1="17" y1="9" x2="23" y2="15"/>
            </svg>
          </button>
        </div>

        <!-- Lyrics -->
        <div class="lyrics-section scroll-y">
          <div
            v-for="(line, i) in lyrics"
            :key="i"
            class="lyric-line"
            :class="{ active: activeLyricIdx === i, past: i < activeLyricIdx }"
            :ref="el => { if (activeLyricIdx === i && el) el.scrollIntoView({ behavior: 'smooth', block: 'center' }) }"
          >{{ line.text }}</div>
          <div v-if="!lyrics.length" class="lyric-line" style="color:var(--text3)">暂无歌词</div>
        </div>

      </div>

      <!-- Empty state -->
      <div v-else class="empty-playing">
        <div class="empty-disc">
          <div class="empty-disc-ring outer"></div>
          <div class="empty-disc-ring inner"></div>
          <div class="empty-disc-center"></div>
        </div>
        <div class="empty-title">等待点歌</div>
        <div class="empty-sub text-muted">在左侧点击「+ 点歌」添加第一首歌</div>
        <button class="btn btn-primary" @click="showSearch = true">立刻点歌</button>
      </div>

      <!-- Reaction Bar -->
      <div class="reaction-bar">
        <button
          v-for="emoji in EMOJIS"
          :key="emoji"
          class="emoji-btn"
          @click="store.sendReaction(emoji)"
        >{{ emoji }}</button>
      </div>
    </main>

    <!-- RIGHT: Chat -->
    <aside class="sidebar right-sidebar">
      <div class="section-header" style="border-bottom: 1px solid var(--border); margin-bottom:0; padding-bottom:12px;">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
        聊天
      </div>
      <div class="chat-messages scroll-y" ref="chatEl">
        <div v-if="!store.room.messages?.length" class="empty-state" style="padding:20px 0">
          <div class="text-muted">开始聊天吧</div>
        </div>
        <div
          v-for="msg in store.room.messages"
          :key="msg.id"
          class="chat-msg"
          :class="{ mine: msg.memberId === store.memberId }"
        >
          <div v-if="msg.memberId !== store.memberId" class="chat-nick">{{ msg.nickname }}</div>
          <div class="chat-bubble">{{ msg.content }}</div>
        </div>
      </div>
      <div class="chat-input-row">
        <input
          v-model="chatInput"
          class="input"
          placeholder="说点什么..."
          maxlength="200"
          @keyup.enter="sendChat"
        />
        <button class="btn btn-primary btn-icon" @click="sendChat">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="22" y1="2" x2="11" y2="13"/><polygon points="22,2 15,22 11,13 2,9"/></svg>
        </button>
      </div>
    </aside>

    <!-- Mobile Bottom Nav -->
    <nav class="mobile-nav">
      <button class="mobile-nav-btn" :class="{ active: mobileTab === 'queue' }" @click="mobileTab = 'queue'">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>
        <span>队列</span>
      </button>
      <button class="mobile-nav-btn" :class="{ active: mobileTab === 'now' }" @click="mobileTab = 'now'">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polygon points="10,8 16,12 10,16"/></svg>
        <span>播放</span>
      </button>
      <button class="mobile-nav-btn" :class="{ active: mobileTab === 'chat' }" @click="mobileTab = 'chat'">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
        <span>聊天</span>
      </button>
    </nav>

    <!-- Hidden audio element -->
    <audio
      ref="audioEl"
      :src="audioSrc"
      @timeupdate="onTimeUpdate"
      @ended="onEnded"
      @canplay="onCanPlay"
      @error="onError"
      preload="auto"
    ></audio>

    <!-- VIP Settings Modal -->
    <Transition name="modal">
      <div v-if="showVipSettings" class="modal-overlay" @click.self="closeVipModal">
        <div class="modal-box" style="max-width:420px">
          <div class="modal-header">
            <div class="modal-title" style="display:flex;align-items:center;gap:8px">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="var(--yellow)"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
              VIP 账户
            </div>
            <button class="btn btn-ghost btn-icon btn-sm" @click="closeVipModal">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
            </button>
          </div>

          <!-- Tabs -->
          <div class="vip-tabs">
            <button class="vip-tab" :class="{ active: vipTab === 'qr' }" @click="vipTab = 'qr'; startQR()">扫码登录</button>
            <button class="vip-tab" :class="{ active: vipTab === 'cookie' }" @click="vipTab = 'cookie'; stopQRPoll()">手动Cookie</button>
          </div>

          <!-- QR Tab -->
          <div v-if="vipTab === 'qr'" class="vip-modal-body">
            <div v-if="qrStatus === 'success'" class="vip-qr-success">
              <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="var(--green)" stroke-width="1.5"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
              <div class="qr-success-title">登录成功</div>
              <div class="qr-success-sub">{{ qrNickname }} · VIP 歌曲已解锁</div>
              <button class="btn btn-primary btn-sm" @click="closeVipModal">完成</button>
            </div>
            <template v-else>
              <div class="vip-qr-wrap">
                <div v-if="qrImg" class="qr-img-wrap" :class="{ expired: qrStatus === 'expired' }">
                  <img :src="qrImg" width="160" height="160" />
                  <div v-if="qrStatus === 'expired'" class="qr-expired-overlay" @click="startQR">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>
                    <span>已过期，点击刷新</span>
                  </div>
                </div>
                <div v-else class="qr-loading">
                  <div class="spinner-lg" style="width:24px;height:24px;border-width:2px"></div>
                </div>
              </div>
              <div class="vip-qr-hint">
                <div class="qr-status-dot" :class="qrStatusClass"></div>
                {{ qrStatusText }}
              </div>
              <div class="vip-tip" style="text-align:center;font-size:12px">
                用<strong style="color:var(--text1)">网易云音乐 App</strong> 扫码登录 · 无需密码
              </div>
            </template>
          </div>

          <!-- Manual Cookie Tab -->
          <div v-if="vipTab === 'cookie'" class="vip-modal-body">
            <div class="cookie-steps">
              <div class="cookie-step">
                <span class="step-num">1</span>
                <span>在 <strong>这个浏览器</strong> 打开 <a href="https://music.163.com" target="_blank" class="link">music.163.com</a> 并登录 VIP 账号</span>
              </div>
              <div class="cookie-step">
                <span class="step-num">2</span>
                <span>按 <code>F12</code> → 切到 <code>Network</code>（网络）标签</span>
              </div>
              <div class="cookie-step">
                <span class="step-num">3</span>
                <span>刷新页面，点击任意一个 <code>music.163.com</code> 的请求</span>
              </div>
              <div class="cookie-step">
                <span class="step-num">4</span>
                <span>右侧 <code>Request Headers</code> → 找到 <code>Cookie</code> 一行 → 右键复制值</span>
              </div>
              <div class="cookie-step">
                <span class="step-num">5</span>
                <span>粘贴到下方（整串，包含 <code>MUSIC_U</code>、<code>__csrf</code> 等）</span>
              </div>
            </div>
            <div class="vip-input-row">
              <textarea
                v-model="cookieInput"
                class="input cookie-textarea"
                placeholder="MUSIC_U=xxx; __csrf=xxx; NMTID=xxx; ..."
                rows="3"
              />
            </div>
            <div class="vip-actions">
              <button class="btn btn-danger btn-sm" @click="clearCookie" v-if="hasCookie">清除</button>
              <div style="flex:1"></div>
              <button class="btn btn-ghost btn-sm" @click="closeVipModal">取消</button>
              <button class="btn btn-primary btn-sm" @click="saveCookieAndClose" :disabled="!cookieInput.trim()">保存</button>
            </div>
            <div v-if="hasCookie" class="vip-saved-tip">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="var(--green)"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
              Cookie 已保存，VIP 歌曲已解锁
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Search Modal -->
    <Transition name="modal">
      <div v-if="showSearch" class="modal-overlay" @click.self="showSearch = false">
        <div class="modal-box">
          <div class="modal-header">
            <div class="modal-title">点歌</div>
            <button class="btn btn-ghost btn-icon btn-sm" @click="showSearch = false">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
            </button>
          </div>
          <div class="modal-search">
            <input
              v-model="searchQuery"
              class="input"
              placeholder="搜索歌曲、歌手..."
              ref="searchInputEl"
              @keyup.enter="doSearch"
              autofocus
            />
            <button class="btn btn-primary" @click="doSearch" :disabled="searching">
              <span v-if="searching" class="spinner" style="width:12px;height:12px;border-width:2px"></span>
              <span v-else>搜索</span>
            </button>
          </div>
          <div class="search-results scroll-y">
            <div v-if="searchError" class="empty-state text-muted" style="padding:24px 0">{{ searchError }}</div>
            <div v-else-if="searchResults.length === 0 && !searching" class="empty-state text-muted" style="padding:24px 0">
              {{ searchQuery ? '未找到结果' : '输入歌名开始搜索' }}
            </div>
            <div
              v-for="song in searchResults"
              :key="song.id"
              class="search-item"
              @click="addToQueue(song)"
            >
              <img class="search-cover" :src="song.cover" :alt="song.name" />
              <div class="search-info">
                <div class="search-name truncate">{{ song.name }}</div>
                <div class="search-meta text-muted">{{ song.artist }} · {{ song.album }}</div>
              </div>
              <div class="search-duration text-muted">{{ formatTime(song.duration) }}</div>
              <button class="btn btn-primary btn-xs">+ 加入</button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </div>

  <!-- Loading / Not joined -->
  <div v-else class="loading-screen">
    <div class="spinner-lg"></div>
    <div style="color:var(--text2);margin-top:12px;">连接中...</div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useRoomStore } from '../stores/room'
import axios from 'axios'

const API = `${import.meta.env.VITE_API_BASE}/api`

const EMOJIS = ['❤️', '🔥', '😭', '666', '牛', '顶']

const route = useRoute()
const router = useRouter()
const store = useRoomStore()

const audioEl = ref(null)
const chatEl = ref(null)
const searchInputEl = ref(null)

const isPlaying = ref(false)
const currentPosition = ref(0)
const positionTimer = ref(null)

const mobileTab = ref('now') // 'now' | 'queue' | 'chat'

const showSearch = ref(false)
const searchQuery = ref('')
const searchResults = ref([])
const searching = ref(false)
const searchError = ref('')

const showVipSettings = ref(false)
const vipTab = ref('cookie')
const cookieInput = ref(store.getCookie())
const hasCookie = computed(() => !!store.getCookie())

// QR login state
const qrImg = ref('')
const qrKey = ref('')
const qrStatus = ref('idle') // idle | waiting | scanned | success | expired
const qrNickname = ref('')
let qrPollTimer = null

const qrStatusText = computed(() => ({
  idle: '生成二维码中...',
  waiting: '等待扫码',
  scanned: '扫码成功，等待确认...',
  success: '登录成功',
  expired: '二维码已过期',
}[qrStatus.value] || ''))

const qrStatusClass = computed(() => ({
  waiting: 'dot-pulse',
  scanned: 'dot-yellow',
  success: 'dot-green-solid',
  expired: 'dot-red',
}[qrStatus.value] || ''))

async function startQR() {
  stopQRPoll()
  qrImg.value = ''
  qrStatus.value = 'idle'
  try {
    const keyRes = await axios.get(`${API}/music/login/qr/key`)
    const key = keyRes.data?.data?.unikey
    if (!key) return
    qrKey.value = key

    const imgRes = await axios.get(`${API}/music/login/qr/create`, { params: { key } })
    qrImg.value = imgRes.data?.data?.qrimg || ''
    qrStatus.value = 'waiting'

    qrPollTimer = setInterval(pollQR, 2000)
  } catch (e) {
    console.error('QR start error', e)
  }
}

async function pollQR() {
  if (!qrKey.value) return
  try {
    const res = await axios.get(`${API}/music/login/qr/check`, { params: { key: qrKey.value } })
    const { code, cookie, profile } = res.data
    if (code === 800) {
      qrStatus.value = 'expired'
      stopQRPoll()
    } else if (code === 801) {
      qrStatus.value = 'waiting'
    } else if (code === 802) {
      qrStatus.value = 'scanned'
    } else if (code === 803) {
      qrStatus.value = 'success'
      stopQRPoll()
      // cookie is returned as a string; store it
      store.saveCookie(cookie || '')
      cookieInput.value = cookie || ''
      qrNickname.value = profile?.nickname || '用户'
    }
  } catch (e) {}
}

function stopQRPoll() {
  clearInterval(qrPollTimer)
  qrPollTimer = null
}

function closeVipModal() {
  stopQRPoll()
  showVipSettings.value = false
}

function saveCookieAndClose() {
  store.saveCookie(cookieInput.value.trim())
  closeVipModal()
}
function clearCookie() {
  cookieInput.value = ''
  store.saveCookie('')
  closeVipModal()
}

watch(showVipSettings, (v) => {
  if (v && vipTab.value === 'qr') startQR()
  if (!v) stopQRPoll()
})
watch(vipTab, (t) => {
  if (t === 'qr') startQR()
  else stopQRPoll()
})

const chatInput = ref('')
const audioSrc = ref('')

// Danmaku
const danmaku = ref([])
const DANMAKU_LANES = 6
const danmakuLaneCooldown = Array.from({ length: DANMAKU_LANES }, () => 0)

function addDanmaku(msg) {
  const now = Date.now()
  let lane = 0
  for (let i = 1; i < DANMAKU_LANES; i++) {
    if (danmakuLaneCooldown[i] < danmakuLaneCooldown[lane]) lane = i
  }
  const duration = 9 + Math.random() * 3
  const id = now + Math.random()
  const top = 60 + lane * 42
  const color = memberColor(msg.nickname)
  danmakuLaneCooldown[lane] = now + duration * 500
  danmaku.value.push({ id, text: msg.nickname + '  ' + msg.content, top, color, duration })
  setTimeout(() => { danmaku.value = danmaku.value.filter(d => d.id !== id) }, (duration + 1) * 1000)
}
const lyrics = ref([])
const codeCopied = ref(false)
// localAudio: true = play audio on this device, false = silent follower (uses startedAt for position)
// Host defaults to true (they control the room speaker); guests default to false
const localAudio = ref(false)

const shouldPlayAudio = computed(() => store.isHost || localAudio.value)

function copyCode() {
  navigator.clipboard?.writeText(store.room?.code || '').then(() => {
    codeCopied.value = true
    setTimeout(() => { codeCopied.value = false }, 1500)
  })
}

function toggleLocalAudio() {
  if (store.isHost) return  // host always plays
  localAudio.value = !localAudio.value
  const pb = store.room?.playback
  if (!pb?.song) return
  if (localAudio.value && pb.isPlaying && audioEl.value) {
    // Just turned on: seek to current position and play
    const elapsed = (Date.now() - pb.startedAt) / 1000
    const seekPos = (pb.position || 0) + elapsed
    audioEl.value.currentTime = Math.min(seekPos, audioEl.value.duration || seekPos)
    audioEl.value.play().catch(() => {})
  } else if (!localAudio.value && audioEl.value) {
    // Just turned off: pause local audio, keep position timer running
    audioEl.value.pause()
  }
}

const memberList = computed(() => store.room?.members || [])

const skipVoteActive = computed(() => {
  const v = store.room?._skipVoteStatus
  return v && v.votes > 0
})
const skipVoteCount = computed(() => store.room?._skipVoteStatus?.votes || 0)

const progressPercent = computed(() => {
  if (!store.currentSong?.duration) return 0
  return Math.min(100, (currentPosition.value / store.currentSong.duration) * 100)
})

const coverGlowStyle = computed(() => {
  return { backgroundImage: `url(${store.currentSong?.cover})` }
})

const activeLyricIdx = computed(() => {
  if (!lyrics.value.length) return -1
  let idx = 0
  for (let i = 0; i < lyrics.value.length; i++) {
    if (lyrics.value[i].time <= currentPosition.value) idx = i
    else break
  }
  return idx
})

// When room state changes, sync audio
watch(() => store.room?.playback, async (pb) => {
  if (!pb?.song) {
    audioSrc.value = ''
    isPlaying.value = false
    clearInterval(positionTimer.value)
    currentPosition.value = 0
    return
  }
  // Load new song if changed
  if (pb.song.id !== lastSongId.value) {
    lastSongId.value = pb.song.id
    lyrics.value = []
    await loadSong(pb.song.id)
  }
  isPlaying.value = pb.isPlaying
  if (pb.isPlaying) {
    if (shouldPlayAudio.value && audioEl.value) {
      const elapsed = (Date.now() - pb.startedAt) / 1000
      const seekPos = pb.position + elapsed
      if (Math.abs(audioEl.value.currentTime - seekPos) > 2) {
        audioEl.value.currentTime = seekPos
      }
      audioEl.value.play().catch(() => {})
    }
    startPositionTimer()
  } else {
    if (audioEl.value) audioEl.value.pause()
    currentPosition.value = pb.position
    clearInterval(positionTimer.value)
  }
}, { deep: true })

const lastSongId = ref(null)

async function loadSong(id) {
  try {
    const url = await store.getMusicURL(id)
    audioSrc.value = url
    const lrc = await store.getLyric(id)
    lyrics.value = lrc || []
  } catch (e) {
    console.error('load song error', e)
  }
}

function startPositionTimer() {
  clearInterval(positionTimer.value)
  positionTimer.value = setInterval(() => {
    if (shouldPlayAudio.value && audioEl.value) {
      // Track from actual audio element when playing locally
      currentPosition.value = audioEl.value.currentTime
    } else {
      // Calculate position from server's startedAt timestamp (no local audio)
      const pb = store.room?.playback
      if (pb?.isPlaying && pb.startedAt) {
        const elapsed = (Date.now() - pb.startedAt) / 1000
        const pos = (pb.position || 0) + elapsed
        currentPosition.value = Math.min(pos, store.currentSong?.duration || pos)
      }
    }
  }, 300)
}

function onTimeUpdate() {
  if (shouldPlayAudio.value && audioEl.value) {
    currentPosition.value = audioEl.value.currentTime
  }
}

function onCanPlay() {
  if (!shouldPlayAudio.value) return
  const pb = store.room?.playback
  if (pb?.isPlaying) {
    const elapsed = (Date.now() - pb.startedAt) / 1000
    const seekPos = pb.position + elapsed
    if (audioEl.value && Math.abs(audioEl.value.currentTime - seekPos) > 1) {
      audioEl.value.currentTime = Math.min(seekPos, audioEl.value.duration || seekPos)
    }
    audioEl.value?.play().catch(() => {})
    startPositionTimer()
  }
}

function onError() {}

function onEnded() {
  if (store.isHost) store.sendNextSong()
}

function togglePlay() {
  if (!audioEl.value) return
  if (isPlaying.value) {
    audioEl.value.pause()
    isPlaying.value = false
    store.sendPlaybackSync(false, audioEl.value.currentTime)
    clearInterval(positionTimer.value)
  } else {
    audioEl.value.play().then(() => {
      isPlaying.value = true
      store.sendPlaybackSync(true, audioEl.value.currentTime)
      startPositionTimer()
    }).catch(() => {})
  }
}

function seekTo(e) {
  if (!store.isHost || !store.currentSong || !audioEl.value) return
  const rect = e.currentTarget.getBoundingClientRect()
  const pct = (e.clientX - rect.left) / rect.width
  const pos = pct * store.currentSong.duration
  audioEl.value.currentTime = pos
  currentPosition.value = pos
  store.sendPlaybackSync(isPlaying.value, pos)
}

function voteUp(qid, amount) {
  const me = store.me
  if (!me || me.coins < amount) return
  store.sendVoteUp(qid, amount)
}

async function sendChat() {
  if (!chatInput.value.trim()) return
  store.sendChat(chatInput.value.trim())
  chatInput.value = ''
  await nextTick()
  if (chatEl.value) chatEl.value.scrollTop = chatEl.value.scrollHeight
}

watch(() => store.room?.messages?.length, async (newLen, oldLen) => {
  await nextTick()
  if (chatEl.value) chatEl.value.scrollTop = chatEl.value.scrollHeight
  if (newLen && newLen > (oldLen || 0)) {
    const msg = store.room.messages[store.room.messages.length - 1]
    if (msg) addDanmaku(msg)
  }
})

async function doSearch() {
  if (!searchQuery.value.trim()) return
  searching.value = true
  searchError.value = ''
  searchResults.value = []
  try {
    searchResults.value = await store.searchMusic(searchQuery.value.trim())
    if (!searchResults.value.length) searchError.value = ''
  } catch (e) {
    searchError.value = '搜索失败：网络错误，请检查后端服务是否启动（go run main.go）'
  } finally {
    searching.value = false
  }
}

function addToQueue(song) {
  store.sendQueueAdd({
    id: song.id,
    name: song.name,
    artist: song.artist,
    album: song.album,
    cover: song.cover,
    duration: song.duration,
  })
  showSearch.value = false
  mobileTab.value = 'now'
}

function memberColor(nickname) {
  const colors = ['#7C5CFA','#2563EB','#059669','#DC2626','#D97706','#DB2777']
  let h = 0
  for (const ch of nickname) h = h * 31 + ch.charCodeAt(0)
  return colors[Math.abs(h) % colors.length]
}

function formatTime(secs) {
  if (!secs) return '0:00'
  const m = Math.floor(secs / 60)
  const s = Math.floor(secs % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

// Keyboard shortcut: space = play/pause
function onKeydown(e) {
  if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return
  if (e.code === 'Space' && store.isHost) {
    e.preventDefault()
    togglePlay()
  }
}

onMounted(async () => {
  const roomId = route.params.id
  const storedMemberId = localStorage.getItem('bfm_memberId')
  const storedRoomId = localStorage.getItem('bfm_roomId')

  if (storedRoomId !== roomId || !storedMemberId) {
    router.push('/')
    return
  }

  store.roomId = roomId
  store.memberId = storedMemberId
  store.connect()

  // Watch until room loads, then set localAudio based on role
  const unwatch = watch(() => store.me, (me) => {
    if (me) {
      localAudio.value = me.isHost
      unwatch()
    }
  }, { immediate: true })

  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
  clearInterval(positionTimer.value)
})
</script>

<style scoped>
.room-layout {
  display: grid;
  grid-template-columns: 280px 1fr 280px;
  height: 100vh;
  overflow: hidden;
  background: var(--bg0);
  position: relative;
}

/* Reaction float layer */
.reaction-layer {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 100;
}
.float-reaction {
  position: absolute;
  bottom: 80px;
  font-size: 28px;
  animation: floatUp 3.2s ease-out forwards;
  pointer-events: none;
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.4));
}

.float-enter-active { animation: floatUp 3.2s ease-out forwards; }
.float-leave-active { animation: none; }

/* Sidebars */
.sidebar {
  display: flex;
  flex-direction: column;
  border-right: 1px solid var(--border);
  height: 100vh;
  overflow: hidden;
  background: var(--bg1);
}
.right-sidebar {
  border-right: none;
  border-left: 1px solid var(--border);
}

.room-header {
  padding: 16px;
  flex-shrink: 0;
}
.room-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}
.room-name {
  font-size: 15px;
  font-weight: 700;
  color: var(--text0);
}
.room-code {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.1em;
  transition: opacity 0.15s;
}
.room-code:hover { opacity: 0.8; }
.room-code-wrap { display: flex; align-items: center; gap: 6px; position: relative; }
.copy-tip {
  position: absolute;
  right: -50px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 11px;
  color: var(--green);
  white-space: nowrap;
  animation: fadeSlideUp 0.2s ease;
}

.vip-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 7px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 600;
  border: 1px solid rgba(255,209,102,0.25);
  background: rgba(255,209,102,0.08);
  color: var(--text3);
  cursor: pointer;
  transition: all 0.15s;
  letter-spacing: 0.02em;
}
.vip-btn:hover { background: rgba(255,209,102,0.15); color: var(--text2); }
.vip-btn.active { color: var(--yellow); border-color: rgba(255,209,102,0.4); background: rgba(255,209,102,0.12); }

/* VIP modal */
.vip-tabs {
  display: flex;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.vip-tab {
  flex: 1;
  padding: 10px;
  font-size: 13px;
  font-weight: 500;
  color: var(--text2);
  background: none;
  border: none;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
  transition: all 0.15s;
  font-family: inherit;
}
.vip-tab.active { color: var(--text0); border-bottom-color: var(--accent); }
.vip-tab:hover { color: var(--text1); }

.vip-modal-body { padding: 20px; display: flex; flex-direction: column; gap: 16px; }
.vip-tip { font-size: 12px; color: var(--text2); line-height: 1.8; }
.vip-tip code { background: var(--bg4); padding: 1px 5px; border-radius: 3px; font-size: 11px; color: var(--accent-light); font-family: monospace; }

/* QR Code */
.vip-qr-wrap { display: flex; justify-content: center; padding: 4px 0; }
.qr-img-wrap {
  position: relative;
  width: 164px; height: 164px;
  padding: 2px;
  background: #fff;
  border-radius: var(--radius-sm);
  display: flex; align-items: center; justify-content: center;
}
.qr-img-wrap.expired img { filter: blur(4px) brightness(0.5); }
.qr-expired-overlay {
  position: absolute; inset: 0;
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  gap: 8px; cursor: pointer;
  color: #fff; font-size: 12px; font-weight: 500;
  background: rgba(0,0,0,0.5);
  border-radius: var(--radius-sm);
}
.qr-loading {
  width: 164px; height: 164px;
  display: flex; align-items: center; justify-content: center;
  background: var(--bg3); border-radius: var(--radius-sm);
}

.vip-qr-hint {
  display: flex; align-items: center; justify-content: center;
  gap: 8px; font-size: 13px; color: var(--text2);
}
.qr-status-dot {
  width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0;
}
.dot-pulse { background: var(--accent); animation: pulse 1.5s ease-in-out infinite; }
.dot-yellow { background: var(--yellow); }
.dot-green-solid { background: var(--green); }
.dot-red { background: var(--red); }

.vip-qr-success {
  display: flex; flex-direction: column; align-items: center;
  gap: 12px; padding: 16px 0 8px;
}
.qr-success-title { font-size: 16px; font-weight: 700; color: var(--green); }
.qr-success-sub { font-size: 13px; color: var(--text2); }

.cookie-steps { display: flex; flex-direction: column; gap: 9px; }
.cookie-step {
  display: flex; align-items: flex-start; gap: 10px;
  font-size: 12px; color: var(--text2); line-height: 1.6;
}
.cookie-step code {
  background: var(--bg4); padding: 0 5px; border-radius: 3px;
  font-size: 11px; color: var(--text1); font-family: monospace;
}
.cookie-step strong { color: var(--text1); font-weight: 500; }
.step-num {
  width: 18px; height: 18px; border-radius: 50%;
  background: var(--bg4); border: 1px solid var(--border-active);
  display: flex; align-items: center; justify-content: center;
  font-size: 10px; font-weight: 700; color: var(--accent-light);
  flex-shrink: 0; margin-top: 2px;
}
.link { color: var(--accent-light); text-decoration: none; }
.link:hover { text-decoration: underline; }
.cookie-textarea {
  height: auto; padding: 10px 12px; resize: none;
  font-family: monospace; font-size: 11px; line-height: 1.6;
}

.vip-input-row { display: flex; }
.vip-actions { display: flex; align-items: center; gap: 8px; }
.vip-saved-tip {
  display: flex; align-items: center; gap: 6px;
  font-size: 12px; color: var(--green);
  background: rgba(34,217,122,0.06);
  border: 1px solid rgba(34,217,122,0.15);
  border-radius: var(--radius-sm);
  padding: 8px 12px;
}
.room-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text2);
}
.dot-green {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--green);
  animation: pulse 2s infinite;
  display: inline-block;
}
.ml-2 { margin-left: 4px; }

.section-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  font-size: 11px;
  font-weight: 600;
  color: var(--text2);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  flex-shrink: 0;
}
.flex-row { display: flex; align-items: center; gap: 6px; }
.text-muted { color: var(--text2); }

/* Members */
.member-list {
  flex-shrink: 0;
  max-height: 160px;
  overflow-y: auto;
  padding: 0 8px 8px;
}
.member-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  border-radius: var(--radius-sm);
  transition: background 0.1s;
}
.member-item:hover { background: var(--bg3); }
.member-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}
.member-info { flex: 1; min-width: 0; }
.member-name { font-size: 12px; font-weight: 500; display: flex; align-items: center; gap: 4px; }
.member-persona { font-size: 10px; color: var(--text3); }
.member-coins { font-size: 11px; color: var(--yellow); display: flex; align-items: center; gap: 3px; flex-shrink: 0; }

/* Queue */
.queue-list { flex: 1; padding: 0 8px; overflow-y: auto; overflow-x: hidden; }
.queue-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  border-radius: var(--radius-sm);
  transition: background 0.1s;
}
.queue-item:hover { background: var(--bg3); }
.queue-rank { font-size: 11px; color: var(--text3); width: 16px; text-align: center; flex-shrink: 0; }
.queue-cover { width: 36px; height: 36px; border-radius: 4px; overflow: hidden; flex-shrink: 0; }
.queue-cover img { width: 100%; height: 100%; object-fit: cover; }
.queue-info { flex: 1; min-width: 0; }
.queue-song { font-size: 12px; font-weight: 500; line-height: 1.4; }
.queue-meta { display: flex; align-items: center; gap: 6px; margin-top: 2px; }
.queue-meta span { font-size: 11px; }
.queue-adder { color: var(--accent-light); font-size: 10px; font-weight: 500; }
.queue-actions { display: flex; flex-direction: column; align-items: flex-end; gap: 4px; flex-shrink: 0; }
.vote-score { font-size: 11px; color: var(--yellow); display: flex; align-items: center; gap: 2px; }
.vote-btns { display: none; gap: 2px; align-items: center; }
.queue-item:hover .vote-btns { display: flex; }

.empty-state { text-align: center; padding: 32px 16px; color: var(--text2); font-size: 12px; line-height: 1.8; }

/* Center */
.center-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100vh;
  overflow: hidden;
  background: var(--bg0);
  position: relative;
  background-image:
    radial-gradient(ellipse 70% 50% at 50% 30%, rgba(124,92,250,0.08) 0%, transparent 70%);
}

.now-playing {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
  width: 100%;
  padding: 24px 32px 8px;
  overflow: hidden;
}

/* Cover */
.cover-wrap {
  position: relative;
  margin-bottom: 20px;
  flex-shrink: 0;
}
.cover-glow {
  position: absolute;
  inset: -20px;
  background-size: cover;
  background-position: center;
  filter: blur(40px) saturate(1.5);
  opacity: 0.3;
  border-radius: 50%;
  z-index: 0;
}
.cover-img-wrap {
  width: 160px;
  height: 160px;
  border-radius: 50%;
  overflow: hidden;
  position: relative;
  z-index: 1;
  border: 3px solid rgba(255,255,255,0.1);
  box-shadow: 0 20px 60px rgba(0,0,0,0.5);
}
.cover-img-wrap.spinning {
  animation: spin 20s linear infinite;
}
.cover-img { width: 100%; height: 100%; object-fit: cover; }

.skip-overlay {
  position: absolute;
  inset: 0;
  z-index: 2;
  background: rgba(0,0,0,0.6);
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.skip-count { font-size: 24px; font-weight: 700; color: var(--red); }
.skip-label { font-size: 11px; color: var(--text2); }

/* Song info */
.song-info { text-align: center; margin-bottom: 12px; flex-shrink: 0; }
.song-name { font-size: 18px; font-weight: 700; color: var(--text0); letter-spacing: -0.3px; }
.song-artist { font-size: 14px; color: var(--text2); margin-top: 4px; }
.song-album { font-size: 12px; margin-top: 2px; }

/* Progress */
.progress-section {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  max-width: 360px;
  flex-shrink: 0;
  margin-bottom: 16px;
}
.time-text { font-size: 11px; color: var(--text3); width: 34px; text-align: center; font-variant-numeric: tabular-nums; }
.progress-bar {
  flex: 1;
  height: 4px;
  background: var(--bg4);
  border-radius: 2px;
  position: relative;
  cursor: pointer;
}
.progress-bar:hover .progress-thumb { opacity: 1; }
.progress-fill {
  height: 100%;
  background: var(--accent);
  border-radius: 2px;
  transition: width 0.5s linear;
}
.progress-thumb {
  position: absolute;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 12px;
  height: 12px;
  background: #fff;
  border-radius: 50%;
  opacity: 0;
  transition: opacity 0.2s;
  box-shadow: 0 1px 4px rgba(0,0,0,0.5);
}

/* Controls */
.controls {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-shrink: 0;
  margin-bottom: 16px;
}
.ctrl-btn {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: var(--bg3);
  border: 1px solid var(--border-active);
  color: var(--text1);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s;
}
.ctrl-btn.playing { background: var(--accent); border-color: transparent; color: #fff; }
.ctrl-btn-indicator { cursor: default; opacity: 0.5; }
.ctrl-btn:not(.ctrl-btn-indicator):hover { transform: scale(1.05); }

.local-audio-btn { color: var(--text3); }
.local-audio-btn.active { color: var(--green); border-color: rgba(34,217,122,0.3); background: rgba(34,217,122,0.06); }
.local-audio-btn:hover { color: var(--text1); }

/* Lyrics */
.lyrics-section {
  flex: 1;
  width: 100%;
  max-width: 400px;
  text-align: center;
  padding-bottom: 12px;
}
.lyric-line {
  padding: 5px 0;
  font-size: 13px;
  color: var(--text3);
  transition: all 0.3s ease;
  line-height: 1.8;
}
.lyric-line.active {
  font-size: 15px;
  font-weight: 600;
  color: var(--text0);
  transform: scale(1.02);
}
.lyric-line.past { color: var(--text3); }

/* Empty playing */
.empty-playing {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  flex: 1;
  padding: 24px;
}
.empty-disc {
  position: relative;
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 4px;
}
.empty-disc-ring {
  position: absolute;
  border-radius: 50%;
  border: 1px solid var(--border-active);
  animation: pulse 2.5s ease-in-out infinite;
}
.empty-disc-ring.outer { width: 100px; height: 100px; animation-delay: 0s; }
.empty-disc-ring.inner { width: 66px; height: 66px; animation-delay: 0.5s; border-color: var(--accent); opacity: 0.3; }
.empty-disc-center {
  width: 16px; height: 16px;
  border-radius: 50%;
  background: var(--bg4);
  border: 2px solid var(--border-active);
}
.empty-title { font-size: 18px; font-weight: 600; color: var(--text1); }
.empty-sub { font-size: 13px; }

/* Reactions */
.reaction-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 0;
  flex-shrink: 0;
}
.emoji-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--bg2);
  border: 1px solid var(--border);
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
  line-height: 1;
}
.emoji-btn:hover {
  background: var(--bg4);
  border-color: var(--border-active);
  transform: scale(1.15);
}
.emoji-btn:active { transform: scale(0.95); }

/* Chat */
.chat-messages {
  flex: 1;
  padding: 8px 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.chat-msg { display: flex; flex-direction: column; gap: 3px; animation: slideIn 0.2s ease; }
.chat-nick { font-size: 11px; color: var(--text3); padding-left: 4px; }
.chat-bubble {
  display: inline-block;
  max-width: 85%;
  background: var(--bg3);
  padding: 7px 11px;
  border-radius: 12px 12px 12px 4px;
  font-size: 13px;
  color: var(--text1);
  word-break: break-word;
  line-height: 1.5;
}
.chat-msg.mine { align-items: flex-end; }
.chat-msg.mine .chat-bubble {
  background: var(--accent-bg);
  border-color: rgba(124,92,250,0.2);
  border-radius: 12px 12px 4px 12px;
  color: var(--text0);
}
.chat-input-row {
  display: flex;
  gap: 8px;
  padding: 12px;
  flex-shrink: 0;
  border-top: 1px solid var(--border);
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.7);
  backdrop-filter: blur(8px);
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}
.modal-box {
  background: var(--bg2);
  border: 1px solid var(--border-active);
  border-radius: var(--radius-xl);
  width: 100%;
  max-width: 520px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.modal-title { font-size: 15px; font-weight: 700; }
.modal-search {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  flex-shrink: 0;
  border-bottom: 1px solid var(--border);
}
.search-results { flex: 1; padding: 8px; }
.search-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.1s;
}
.search-item:hover { background: var(--bg3); }
.search-cover { width: 40px; height: 40px; border-radius: 4px; object-fit: cover; flex-shrink: 0; }
.search-info { flex: 1; min-width: 0; }
.search-name { font-size: 13px; font-weight: 500; }
.search-meta { font-size: 11px; margin-top: 2px; }
.search-duration { font-size: 12px; flex-shrink: 0; }

.modal-enter-active, .modal-leave-active { transition: opacity 0.2s, transform 0.2s; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
.modal-enter-from .modal-box, .modal-leave-to .modal-box { transform: scale(0.95); }

.loading-screen {
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: var(--bg0);
}
.spinner-lg {
  width: 36px;
  height: 36px;
  border: 3px solid var(--bg4);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

/* Danmaku */
.danmaku-layer {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 50;
  overflow: hidden;
}
.danmaku-item {
  position: absolute;
  left: 0;
  white-space: nowrap;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 0.01em;
  text-shadow: 0 1px 4px rgba(0,0,0,0.95), 0 0 10px rgba(0,0,0,0.8);
  animation: danmakuFlow linear forwards;
}
@keyframes danmakuFlow {
  from { transform: translateX(100vw); }
  to { transform: translateX(-100%); }
}

/* ====== Mobile Responsive ====== */
@media (max-width: 768px) {
  .room-layout {
    display: flex;
    flex-direction: column;
    height: 100vh;
    height: 100dvh;
  }

  /* Sidebars hidden by default, shown based on mobileTab */
  .sidebar {
    display: none;
    height: auto;
    flex: 1;
    overflow-y: auto;
    border-right: none;
    border-left: none;
  }
  .sidebar.left-sidebar {
    display: none;
  }
  .sidebar.right-sidebar {
    display: none;
  }

  /* Show sidebar when its tab is active */
  .room-layout.show-queue .sidebar.left-sidebar,
  .room-layout.show-chat .sidebar.right-sidebar {
    display: flex;
  }

  /* Center panel */
  .center-panel {
    flex: 1;
    height: auto;
    overflow-y: auto;
    min-height: 0;
  }

  /* Now playing compact */
  .now-playing {
    padding: 12px 16px 8px;
    gap: 8px;
  }

  .cover-img-wrap {
    width: 120px;
    height: 120px;
  }
  .cover-glow {
    inset: -12px;
  }

  .song-name {
    font-size: 16px;
  }
  .song-artist {
    font-size: 13px;
  }

  .progress-section {
    max-width: 100%;
    margin-bottom: 12px;
  }

  .controls {
    gap: 12px;
    margin-bottom: 12px;
  }
  .ctrl-btn {
    width: 44px;
    height: 44px;
  }
  .ctrl-btn svg {
    width: 18px;
    height: 18px;
  }

  .lyrics-section {
    max-width: 100%;
    padding-bottom: 40px;
  }
  .lyric-line {
    font-size: 12px;
    padding: 4px 0;
  }
  .lyric-line.active {
    font-size: 14px;
  }

  /* Reactions bar */
  .reaction-bar {
    padding: 8px 0;
    gap: 6px;
  }
  .emoji-btn {
    width: 36px;
    height: 36px;
    font-size: 16px;
  }

  /* Members: compact */
  .member-list {
    max-height: 120px;
  }

  /* Queue: compact */
  .queue-item {
    padding: 8px;
  }
  .queue-cover {
    width: 32px;
    height: 32px;
  }
  .vote-btns {
    display: flex !important;
  }

  /* Chat: compact */
  .chat-input-row {
    padding: 8px;
  }

  /* Room header compact */
  .room-header {
    padding: 10px 12px;
  }
  .room-name {
    font-size: 14px;
  }

  .section-header {
    padding: 8px 12px;
  }

  /* Modals: full screen */
  .modal-box {
    max-width: 100%;
    max-height: 90vh;
    border-radius: var(--radius-lg) var(--radius-lg) 0 0;
    margin-top: auto;
  }
  .modal-overlay {
    align-items: flex-end;
    padding: 0;
  }

  /* Empty playing */
  .empty-playing {
    padding: 16px;
    gap: 12px;
  }
  .empty-disc {
    width: 80px;
    height: 80px;
  }
  .empty-disc-ring.outer {
    width: 80px;
    height: 80px;
  }
  .empty-disc-ring.inner {
    width: 50px;
    height: 50px;
  }
  .empty-title {
    font-size: 16px;
  }
  .empty-sub {
    font-size: 12px;
  }

  /* Search modal */
  .search-item {
    padding: 10px;
  }
  .search-cover {
    width: 36px;
    height: 36px;
  }
  .search-name {
    font-size: 12px;
  }
  .search-meta {
    font-size: 10px;
  }
  .modal-search {
    padding: 10px 12px;
  }

  /* VIP modal */
  .vip-modal-body {
    padding: 14px;
  }
  .cookie-step {
    font-size: 11px;
  }
}

/* Mobile bottom nav */
.mobile-nav {
  display: none;
}

@media (max-width: 768px) {
  .mobile-nav {
    display: flex;
    align-items: center;
    justify-content: space-around;
    background: var(--bg1);
    border-top: 1px solid var(--border);
    padding: 0;
    flex-shrink: 0;
    padding-bottom: env(safe-area-inset-bottom, 0);
  }
  .mobile-nav-btn {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    padding: 8px 0;
    flex: 1;
    background: none;
    border: none;
    color: var(--text3);
    cursor: pointer;
    font-family: inherit;
    font-size: 10px;
    transition: color 0.15s;
    -webkit-tap-highlight-color: transparent;
  }
  .mobile-nav-btn.active {
    color: var(--accent-light);
  }
  .mobile-nav-btn svg {
    width: 20px;
    height: 20px;
  }

  /* Panel visibility by mobile tab */
  .left-sidebar { display: none !important; }
  .right-sidebar { display: none !important; }
  .center-panel { display: none !important; }

  .room-layout.tab-queue .left-sidebar { display: flex !important; }
  .room-layout.tab-now .center-panel { display: flex !important; }
  .room-layout.tab-chat .right-sidebar { display: flex !important; }
}
</style>
