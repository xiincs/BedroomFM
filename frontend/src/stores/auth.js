import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

const API = `${import.meta.env.VITE_API_BASE}/api`

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('bfm_token') || '')
  const user = ref(JSON.parse(localStorage.getItem('bfm_user') || 'null'))
  const tier = ref(null)
  const xpProgress = ref(0)
  const xpToNext = ref(100)
  const allTiers = ref([])

  const isLoggedIn = computed(() => !!token.value && !!user.value)

  if (token.value) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  function _persist(tokenVal, userData) {
    token.value = tokenVal
    user.value = userData
    localStorage.setItem('bfm_token', tokenVal)
    localStorage.setItem('bfm_user', JSON.stringify(userData))
    axios.defaults.headers.common['Authorization'] = `Bearer ${tokenVal}`
  }

  async function register(username, password, nickname) {
    const { data } = await axios.post(`${API}/auth/register`, { username, password, nickname })
    _persist(data.token, data.user)
    tier.value = data.tier
    xpProgress.value = data.xpProgress ?? 0
    xpToNext.value = data.xpToNext ?? 100
    allTiers.value = data.allTiers ?? []
    return data
  }

  async function login(username, password) {
    const { data } = await axios.post(`${API}/auth/login`, { username, password })
    _persist(data.token, data.user)
    tier.value = data.tier
    xpProgress.value = data.xpProgress ?? 0
    xpToNext.value = data.xpToNext ?? 100
    allTiers.value = data.allTiers ?? []
    return data
  }

  function logout() {
    token.value = ''
    user.value = null
    tier.value = null
    localStorage.removeItem('bfm_token')
    localStorage.removeItem('bfm_user')
    delete axios.defaults.headers.common['Authorization']
  }

  async function fetchMe() {
    if (!token.value) return
    try {
      const { data } = await axios.get(`${API}/auth/me`)
      user.value = data.user
      tier.value = data.tier
      xpProgress.value = data.xpProgress ?? 0
      xpToNext.value = data.xpToNext ?? 100
      allTiers.value = data.allTiers ?? []
      localStorage.setItem('bfm_user', JSON.stringify(data.user))
    } catch (e) {
      if (e?.response?.status === 401) logout()
    }
  }

  async function gainXP(action) {
    if (!token.value) return null
    try {
      const { data } = await axios.post(`${API}/user/xp`, { action })
      if (user.value) {
        user.value = { ...user.value, xp: data.xp, level: data.level }
        localStorage.setItem('bfm_user', JSON.stringify(user.value))
      }
      if (data.tier) tier.value = data.tier
      xpProgress.value = data.xpProgress ?? xpProgress.value
      xpToNext.value = data.xpToNext ?? xpToNext.value
      return data
    } catch {
      return null
    }
  }

  return {
    token, user, tier, xpProgress, xpToNext, allTiers,
    isLoggedIn,
    register, login, logout, fetchMe, gainXP,
  }
})
