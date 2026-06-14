import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

const API = 'http://localhost:8080/api'
const WS_BASE = 'ws://localhost:8080/ws'

export const useRoomStore = defineStore('room', () => {
  const roomId = ref(null)
  const memberId = ref(null)
  const room = ref(null)
  const ws = ref(null)
  const connected = ref(false)
  const reactions = ref([]) // floating reactions

  const me = computed(() => {
    if (!room.value || !memberId.value) return null
    return room.value.members?.find(m => m.id === memberId.value)
  })

  const isHost = computed(() => me.value?.isHost === true)

  const currentSong = computed(() => room.value?.playback?.song || null)

  function applyState(state) {
    room.value = state
  }

  async function createRoom(roomName, nickname) {
    const { data } = await axios.post(`${API}/room/create`, { roomName, nickname })
    roomId.value = data.roomId
    memberId.value = data.memberId
    localStorage.setItem('bfm_roomId', data.roomId)
    localStorage.setItem('bfm_memberId', data.memberId)
    return data
  }

  async function joinRoom(code, nickname) {
    const { data } = await axios.post(`${API}/room/join`, { code, nickname })
    roomId.value = data.roomId
    memberId.value = data.memberId
    localStorage.setItem('bfm_roomId', data.roomId)
    localStorage.setItem('bfm_memberId', data.memberId)
    return data
  }

  function connect() {
    if (ws.value) ws.value.close()
    const url = `${WS_BASE}/${roomId.value}?memberId=${memberId.value}`
    const socket = new WebSocket(url)

    socket.onopen = () => { connected.value = true }
    socket.onclose = () => {
      connected.value = false
      // auto reconnect after 3s
      setTimeout(() => { if (roomId.value) connect() }, 3000)
    }
    socket.onmessage = (e) => {
      try {
        const msg = JSON.parse(e.data)
        handleMessage(msg)
      } catch {}
    }
    ws.value = socket
  }

  function handleMessage(msg) {
    switch (msg.type) {
      case 'room_state': {
        const state = typeof msg.payload === 'string' ? JSON.parse(msg.payload) : msg.payload
        applyState(state)
        break
      }
      case 'chat': {
        const m = typeof msg.payload === 'string' ? JSON.parse(msg.payload) : msg.payload
        if (room.value) {
          if (!room.value.messages) room.value.messages = []
          room.value.messages.push(m)
          if (room.value.messages.length > 200) room.value.messages.splice(0, 50)
        }
        break
      }
      case 'reaction': {
        const r = typeof msg.payload === 'string' ? JSON.parse(msg.payload) : msg.payload
        const id = Date.now() + Math.random()
        reactions.value.push({ ...r, id, x: 20 + Math.random() * 60 })
        setTimeout(() => {
          reactions.value = reactions.value.filter(rx => rx.id !== id)
        }, 3500)
        break
      }
      case 'playback_sync': {
        const p = typeof msg.payload === 'string' ? JSON.parse(msg.payload) : msg.payload
        if (room.value?.playback) {
          room.value.playback.isPlaying = p.isPlaying
          room.value.playback.position = p.position
          room.value.playback.startedAt = p.startedAt
        }
        break
      }
      case 'skip_vote_update': {
        const p = typeof msg.payload === 'string' ? JSON.parse(msg.payload) : msg.payload
        if (room.value) room.value._skipVoteStatus = p
        break
      }
      case 'member_join':
      case 'member_leave':
        // room_state handles full sync
        break
    }
  }

  function send(type, payload) {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify({ type, payload }))
    }
  }

  function sendChat(content) { send('chat', { content }) }
  function sendReaction(emoji) { send('reaction', { emoji }) }
  function sendQueueAdd(song) { send('queue_add', { song }) }
  function sendQueueRemove(qid) { send('queue_remove', { qid }) }
  function sendVoteUp(qid, amount) { send('vote_up', { qid, amount }) }
  function sendVoteSkip() { send('vote_skip', {}) }
  function sendPlaybackSync(isPlaying, position) { send('playback_sync', { isPlaying, position }) }
  function sendNextSong() { send('next_song', {}) }

  async function searchMusic(q) {
    const { data } = await axios.get(`${API}/music/search`, { params: { q } })
    return data
  }

  function getCookie() {
    return localStorage.getItem('bfm_netease_cookie') || ''
  }
  function saveCookie(val) {
    if (val) localStorage.setItem('bfm_netease_cookie', val)
    else localStorage.removeItem('bfm_netease_cookie')
  }

  async function getMusicURL(id) {
    const params = { id }
    const cookie = getCookie()
    if (cookie) params.cookie = cookie
    const { data } = await axios.get(`${API}/music/url`, { params })
    return data.url
  }

  async function getLyric(id) {
    const { data } = await axios.get(`${API}/music/lyric`, { params: { id } })
    return data
  }

  function disconnect() {
    ws.value?.close()
    ws.value = null
    roomId.value = null
    memberId.value = null
    room.value = null
    connected.value = false
    localStorage.removeItem('bfm_roomId')
    localStorage.removeItem('bfm_memberId')
  }

  return {
    roomId, memberId, room, ws, connected, reactions,
    me, isHost, currentSong,
    createRoom, joinRoom, connect, disconnect,
    sendChat, sendReaction, sendQueueAdd, sendQueueRemove,
    sendVoteUp, sendVoteSkip, sendPlaybackSync, sendNextSong,
    searchMusic, getMusicURL, getLyric,
    getCookie, saveCookie,
  }
})
