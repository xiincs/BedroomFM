<template>
  <span class="vip-badge" :class="`vip-tier-${tierId}`" :style="badgeStyle">
    <span class="vip-lv">Lv.{{ level }}</span>
    <span class="vip-name">{{ tierName }}</span>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  level: { type: Number, default: 1 },
})

const TIERS = [
  { id: 0, name: '听众',    min: 1,  color: '#6b7280', g: ['#6b7280','#9ca3af'] },
  { id: 1, name: '小歌迷',  min: 5,  color: '#3b82f6', g: ['#3b82f6','#60a5fa'] },
  { id: 2, name: '节奏达人',min: 10, color: '#06b6d4', g: ['#06b6d4','#10b981'] },
  { id: 3, name: '黄金DJ',  min: 20, color: '#f59e0b', g: ['#f59e0b','#ef4444'] },
  { id: 4, name: '白金制作人',min:35, color: '#8b5cf6', g: ['#8b5cf6','#ec4899'] },
  { id: 5, name: '钻石传说',min: 50, color: '#06b6d4', g: ['#06b6d4','#8b5cf6','#ec4899'] },
  { id: 6, name: '传奇之声',min: 75, color: '#a855f7', g: ['#f59e0b','#8b5cf6','#06b6d4'] },
  { id: 7, name: '音乐之神',min: 100,color: '#fbbf24', g: ['#fbbf24','#ef4444','#8b5cf6'] },
]

const currentTier = computed(() => {
  for (let i = TIERS.length - 1; i >= 0; i--) {
    if (props.level >= TIERS[i].min) return TIERS[i]
  }
  return TIERS[0]
})

const tierId = computed(() => currentTier.value.id)
const tierName = computed(() => currentTier.value.name)

const badgeStyle = computed(() => {
  const g = currentTier.value.g
  const gradient = g.length > 1
    ? `linear-gradient(90deg, ${g.join(', ')})`
    : g[0]
  return {
    '--badge-grad': gradient,
    '--badge-color': currentTier.value.color,
  }
})
</script>

<style scoped>
.vip-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 7px;
  border-radius: 20px;
  font-size: 10px;
  font-weight: 600;
  line-height: 1.4;
  background: var(--badge-grad);
  color: #fff;
  white-space: nowrap;
  text-shadow: 0 1px 2px rgba(0,0,0,0.4);
  box-shadow: 0 0 8px color-mix(in srgb, var(--badge-color) 40%, transparent);
  position: relative;
  overflow: hidden;
}

.vip-lv {
  opacity: 0.85;
  font-size: 9px;
  letter-spacing: 0.02em;
  position: relative;
}

.vip-name {
  font-size: 10px;
  position: relative;
}

/* shimmer sweep for high tiers */
.vip-tier-5::after,
.vip-tier-6::after,
.vip-tier-7::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, transparent 0%, rgba(255,255,255,0.35) 50%, transparent 100%);
  animation: badge-shimmer 2s linear infinite;
}

@keyframes badge-shimmer {
  from { transform: translateX(-100%); }
  to   { transform: translateX(100%); }
}
</style>
