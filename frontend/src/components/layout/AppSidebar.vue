<script setup lang="ts">
import { useLayoutStore } from '@/stores/layout'

const layoutStore = useLayoutStore()

const menuItems = [
  { name: '工作台', icon: 'ri-dashboard-3-line', path: '/dashboard' },
  { name: '项目管理', icon: 'ri-folder-line', path: '/projects' },
  { name: '收款日历', icon: 'ri-calendar-line', path: '/calendar' },
  { name: '数据分析', icon: 'ri-bar-chart-grouped-line', path: '/analytics' },
]

const systemItems = [
  { name: '系统设置', icon: 'ri-settings-4-line', path: '/settings' },
]
</script>

<template>
  <aside class="sidebar liquid-glass" :class="{ collapsed: layoutStore.sidebarCollapsed }" id="sidebar">
    <div class="sidebar-header">
      <div class="sidebar-logo">
        <div class="sidebar-logo-icon">
          <img src="/orange.png" alt="Orange Logo" />
        </div>
        <span class="sidebar-logo-text">Orange</span>
      </div>
    </div>

    <nav class="nav-menu">
      <div class="nav-section">
        <div class="nav-section-title">主菜单</div>
        <router-link
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          active-class="active"
        >
          <span class="nav-item-icon"><i :class="item.icon"></i></span>
          <span class="nav-item-text">{{ item.name }}</span>
        </router-link>
      </div>

      <div class="nav-section">
        <div class="nav-section-title">系统</div>
        <router-link
          v-for="item in systemItems"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          active-class="active"
        >
          <span class="nav-item-icon"><i :class="item.icon"></i></span>
          <span class="nav-item-text">{{ item.name }}</span>
        </router-link>
      </div>
    </nav>
  </aside>
</template>

<style scoped>
.sidebar-logo-text {
  margin-left: 12px;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 20px;
  color: var(--text-primary);
}
</style>

