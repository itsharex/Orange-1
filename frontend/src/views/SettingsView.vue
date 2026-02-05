<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authApi } from '@/api/auth'
import api from '@/api'

import UserManagement from '@/components/settings/UserManagement.vue'
import DataSyncPanel from '@/components/settings/DataSyncPanel.vue'
import TokenManagement from '@/components/settings/TokenManagement.vue'
import NotificationManagement from '@/components/settings/NotificationManagement.vue'
import DictionaryManagement from '@/components/settings/DictionaryManagement.vue'
import GlassCard from '@/components/common/GlassCard.vue'
import { useConfirm } from '@/composables/useConfirm'
import { useToast } from '@/composables/useToast'
import { Browser } from '@wailsio/runtime'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import pkg from '../../package.json'


const route = useRoute()
const router = useRouter()
const { confirm } = useConfirm()
const toast = useToast()

const activeTab = ref(route.query.tab as string || 'profile')

// Watch route tab change
watch(() => route.query.tab, (newTab) => {
  if (newTab) {
    activeTab.value = newTab as string
  }
})

// Update route when tab changes
watch(activeTab, (newTab) => {
  router.replace({ query: { ...route.query, tab: newTab } })
})

const authStore = useAuthStore()

// 简单的权限判断
const isAdmin = computed(() => authStore.user?.role === 'admin')

// 计算设置导航菜单 (根据权限动态显示)
const settingsNav = computed(() => {
  const items = [
    { key: 'profile', icon: 'ri-user-line', label: '个人信息' },
    // Admin only
    ...(isAdmin.value ? [{ key: 'users', icon: 'ri-team-line', label: '用户管理' }] : []),
    { key: 'security', icon: 'ri-lock-line', label: '安全设置' },
    { key: 'data-sync', icon: 'ri-cloud-line', label: '数据同步' },
    { key: 'appearance', icon: 'ri-palette-line', label: '外观设置' },
    { key: 'notification', icon: 'ri-notification-3-line', label: '通知设置' },
    { key: 'developer', icon: 'ri-terminal-box-line', label: '开发设置' },
    { key: 'about', icon: 'ri-information-line', label: '关于' },
  ]
  
  if (isAdmin.value) {
    // 管理员专属菜单 (Additional items if needed)
    const dictIndex = items.findIndex(i => i.key === 'security')
    items.splice(dictIndex, 0, { key: 'dictionary', icon: 'ri-book-2-line', label: '字典管理' })
  }
  
  return items
})

// ============ Profile Logic ============
const profile = ref({
  name: '',
  position: '',
  email: '',
  phone: '',
  department: '',
})

// Store original profile for comparison
const originalProfile = ref({ ...profile.value }) // Initialize empty

// Fetch current user
const fetchProfile = async () => {
  try {
    const res = await authApi.getCurrentUser()
    if (res.data.code === 0 && res.data.data) {
      const user = res.data.data
      profile.value = {
        name: user.name || '',
        position: user.position || '',
        email: user.email || '',
        phone: user.phone || '',
        department: user.department || '',
      }
      // Update original profile
      originalProfile.value = { ...profile.value }
    }
  } catch (error) {
    console.error('Failed to fetch profile:', error)
  }
}

// 保存个人信息
const saveProfile = async () => {
  // 检查是否有变更
  const isModified =
    profile.value.name !== originalProfile.value.name ||
    profile.value.position !== originalProfile.value.position ||
    profile.value.email !== originalProfile.value.email ||
    profile.value.phone !== originalProfile.value.phone ||
    profile.value.department !== originalProfile.value.department

  if (!isModified) {
    toast.info('未做任何修改')
    return
  }

  // 校验邮箱格式（如果有输入）
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (profile.value.email && !emailRegex.test(profile.value.email)) {
    toast.warning('邮箱格式不正确')
    return
  }

  // 校验手机号格式（如果有输入）
  const phoneRegex = /^1[3-9]\d{9}$/
  if (profile.value.phone && !phoneRegex.test(profile.value.phone)) {
    toast.warning('手机号格式不正确')
    return
  }

  try {
    const res = await authApi.updateProfile({
      name: profile.value.name,
      position: profile.value.position,
      email: profile.value.email, 
      phone: profile.value.phone,
      department: profile.value.department,
    })
    if (res.data.code === 0) {
      toast.success('保存成功')
      // 更新原始快照
      originalProfile.value = { ...profile.value }
    } else {
      toast.error(`保存失败: ${res.data.message}`)
    }
  } catch (error) {
    console.error('Failed to update profile:', error)
    toast.error('保存失败')
  }
}







// ============ Appearance Logic ============
const themeStore = useThemeStore()

const themes = [
  { value: 'auto', label: '跟随系统', icon: 'ri-computer-line' },
  { value: 'light', label: '浅色模式', icon: 'ri-sun-line' },
  { value: 'dark', label: '深色模式', icon: 'ri-moon-line' },
]

// ============ About Page Logic ============
const openGitHub = () => {
  Browser.OpenURL('https://github.com/FruitsAI/Orange')
}

// ============ Security Logic ============
const securityForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const handlePasswordChange = async () => {
  const { oldPassword, newPassword, confirmPassword } = securityForm.value

  if (!oldPassword || !newPassword || !confirmPassword) {
    toast.warning('请填写所有密码字段')
    return
  }

  if (newPassword !== confirmPassword) {
    toast.warning('两次输入的新密码不一致')
    return
  }

  if (newPassword.length < 6) {
    toast.warning('新密码长度不能少于6位')
    return
  }

  const success = await authStore.changePassword(oldPassword, newPassword)

  if (success) {
    toast.success('密码修改成功')
    // Reset form
    securityForm.value = {
      oldPassword: '',
      newPassword: '',
      confirmPassword: '',
    }
  } else {
    toast.error(authStore.error || '密码修改失败')
  }
}


// ============ About Logic ============
const checkingUpdate = ref(false)

// 检查软件更新
const checkUpdate = async () => {
  checkingUpdate.value = true
  try {
    // 调用后端代理接口，避免跨域 (CORS) 问题
    const { data } = await api.get('/system/updates/check')
    
    const releaseInfo = data.data
    const latestVersionTag = releaseInfo.tag_name 
    const htmlUrl = releaseInfo.html_url
    
    // 简单的版本号比较逻辑 (去掉 'v' 前缀)
    const cleanLatest = latestVersionTag.replace(/^v/, '')
    const cleanCurrent = pkg.version.replace(/^v/, '')
    
    if (cleanLatest === cleanCurrent) {
       toast.success(`当前已是最新版本 (v${cleanLatest})`)
    } else {
       const confirmed = await confirm({ 
         title: '发现新版本', 
         message: `检测到新版本 ${latestVersionTag}，当前版本 v${cleanCurrent}。是否前往下载？` 
       })
       
       if (confirmed) {
         Browser.OpenURL(htmlUrl)
       }
    }
  } catch (error) {
    console.error('Check update failed:', error)
    toast.error('检查更新失败，请稍后重试')
  } finally {
    checkingUpdate.value = false
  }
}

onMounted(() => {
  fetchProfile()

})
</script>

<template>
  <div class="settings-view grid gap-lg">
    <!-- Settings Navigation -->
    <GlassCard class="p-0 h-fit nav-card">
      <div class="nav-header">设置</div>
      <div class="nav-list">
        <a
          v-for="item in settingsNav"
          :key="item.key"
          href="#"
          class="nav-item-settings"
          :class="{ active: activeTab === item.key }"
          @click.prevent="activeTab = item.key"
        >
          <i :class="[item.icon, 'nav-icon']"></i> <span class="nav-label">{{ item.label }}</span>
        </a>
      </div>
    </GlassCard>

    <!-- Main Content Area -->
    <GlassCard v-if="activeTab === 'profile'" class="profile-card" no-padding>
      <div class="profile-panel">
        <!-- 头部区域 -->
        <div class="dev-header">
          <div class="dev-header-content">
            <div class="dev-title-section">
              <div class="dev-icon-wrapper">
                <i class="ri-user-3-line"></i>
              </div>
              <div class="dev-title-info">
                <h2 class="dev-title">个人信息</h2>
                <p class="dev-subtitle">管理您的个人资料和联系方式</p>
              </div>
            </div>
            <!-- 右上角按钮 -->
            <button class="dev-create-btn" @click="saveProfile">
              <i class="ri-save-line"></i>
              <span>保存更改</span>
            </button>
          </div>
        </div>

        <!-- 表单区域 -->
        <div class="dev-content">
          <div class="profile-form-grid">
            <div class="form-group">
              <label class="form-label">姓名</label>
              <input type="text" v-model="profile.name" class="form-input" spellcheck="false" autocomplete="off" autocorrect="off" autocapitalize="off" />
            </div>
            <div class="form-group">
              <label class="form-label">职位</label>
              <input type="text" v-model="profile.position" class="form-input" spellcheck="false" autocomplete="off" autocorrect="off" autocapitalize="off" />
            </div>
            <div class="form-group">
              <label class="form-label">邮箱</label>
              <input type="email" v-model="profile.email" class="form-input" spellcheck="false" autocomplete="off" autocorrect="off" autocapitalize="off" />
            </div>
            <div class="form-group">
              <label class="form-label">手机</label>
              <input type="tel" v-model="profile.phone" class="form-input" spellcheck="false" autocomplete="off" autocorrect="off" autocapitalize="off" />
            </div>
            <div class="form-group form-group-full">
              <label class="form-label">部门</label>
              <input type="text" v-model="profile.department" class="form-input" spellcheck="false" autocomplete="off" autocorrect="off" autocapitalize="off" />
            </div>
          </div>
        </div>
      </div>
    </GlassCard>

    <!-- User Management (Admin) -->
    <GlassCard
      v-else-if="activeTab === 'users' && isAdmin"
      class="h-fit flex flex-col p-0 overflow-hidden"
    >
      <UserManagement />
    </GlassCard>

    <!-- Dictionary Management -->
    <!-- Dictionary Management -->
    <GlassCard
      v-else-if="activeTab === 'dictionary'"
      class="flex flex-col p-0"
    >
      <DictionaryManagement />
    </GlassCard>

    <!-- Data Sync Panel -->
    <GlassCard
      v-else-if="activeTab === 'data-sync'"
      class="h-full"
      no-padding
    >
      <DataSyncPanel />
    </GlassCard>

    <!-- Developer Settings (Token Management) -->
    <GlassCard
      v-else-if="activeTab === 'developer'"
      class="h-full"
      no-padding
    >
      <TokenManagement />
    </GlassCard>



    <!-- Security Settings -->
    <GlassCard v-else-if="activeTab === 'security'" class="security-card" no-padding>
      <div class="security-panel">
        <!-- 头部区域 -->
        <div class="dev-header">
          <div class="dev-header-content">
            <div class="dev-title-section">
              <div class="dev-icon-wrapper">
                <i class="ri-lock-line"></i>
              </div>
              <div class="dev-title-info">
                <h2 class="dev-title">安全设置</h2>
                <p class="dev-subtitle">管理账户安全，保护个人信息</p>
              </div>
            </div>
            <!-- 右上角按钮 -->
            <button class="dev-create-btn" @click="handlePasswordChange">
              <i class="ri-lock-password-line"></i>
              <span>修改密码</span>
            </button>
          </div>
        </div>

        <!-- 内容区域 - 三列布局 -->
        <div class="dev-content">
          <div class="security-form-grid">
            <div class="form-group">
              <label class="form-label">当前密码</label>
              <input
                type="password"
                v-model="securityForm.oldPassword"
                class="form-input"
                placeholder="请输入当前密码"
                spellcheck="false"
                autocomplete="off"
                autocorrect="off"
                autocapitalize="off"
              />
            </div>
            <div class="form-group">
              <label class="form-label">新密码</label>
              <input
                type="password"
                v-model="securityForm.newPassword"
                class="form-input"
                placeholder="请输入新密码（至少6位）"
                spellcheck="false"
                autocomplete="off"
                autocorrect="off"
                autocapitalize="off"
              />
            </div>
            <div class="form-group">
              <label class="form-label">确认新密码</label>
              <input
                type="password"
                v-model="securityForm.confirmPassword"
                class="form-input"
                placeholder="请再次输入新密码"
                spellcheck="false"
                autocomplete="off"
                autocorrect="off"
                autocapitalize="off"
              />
            </div>
          </div>
        </div>
      </div>
    </GlassCard>

    <!-- Notification Settings -->
    <GlassCard v-else-if="activeTab === 'notification'" class="h-full" no-padding>
      <NotificationManagement :is-admin="isAdmin" />
    </GlassCard>

    <!-- Appearance Settings -->
    <GlassCard v-else-if="activeTab === 'appearance'" class="appearance-card" no-padding>
      <div class="appearance-panel">
        <!-- 头部区域 -->
        <div class="appearance-header">
          <div class="appearance-header-main">
            <div class="appearance-title-wrapper">
              <div class="appearance-icon">
                <i class="ri-palette-line"></i>
              </div>
              <div class="appearance-title-content">
                <h2 class="appearance-title">外观设置</h2>
                <p class="appearance-subtitle">自定义界面主题，打造专属的使用体验</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 主题选择内容 -->
        <div class="appearance-content">
          <div class="grid grid-cols-3 gap-4">
            <div
              v-for="t in themes"
              :key="t.value"
              class="theme-card"
              :class="{ active: themeStore.theme === t.value }"
              @click="themeStore.setTheme(t.value)"
            >
              <div class="theme-icon">
                 <i :class="t.icon"></i>
              </div>
              <span class="theme-label">{{ t.label }}</span>
              <div class="theme-check" v-if="themeStore.theme === t.value">
                <i class="ri-check-line"></i>
              </div>
            </div>
          </div>
        </div>
      </div>
    </GlassCard>



    <!-- About Page - Liquid Glass Premium Edition -->
    <GlassCard v-else-if="activeTab === 'about'" class="about-page h-auto min-h-full">
      <div class="about-container">
        <!-- Hero Section -->
        <div class="about-hero">
          <div class="logo-wrapper">
            <div class="logo-glow"></div>
            <img src="/orange.png" alt="Orange Logo" class="about-logo" />
          </div>
          <h1 class="about-title">Orange</h1>
          <div class="about-subtitle">
            <span class="version-badge">v{{ pkg.version }}</span>
            <span class="tagline">小旭姐专属记账工具</span>
          </div>
        </div>

        <!-- Floating Info Cards -->
        <div class="info-cards-grid">
          <!-- Author Card -->
          <div class="info-card info-card-author">
            <div class="card-glass">
              <div class="card-icon-wrapper">
                <i class="ri-user-smile-line"></i>
              </div>
              <div class="card-label">作者</div>
              <div class="card-value">willxue</div>
            </div>
          </div>
          
          <!-- WeChat Card -->
          <div class="info-card info-card-wechat">
            <div class="card-glass">
              <div class="card-icon-wrapper">
                <i class="ri-wechat-line"></i>
              </div>
              <div class="card-label">微信公众号</div>
              <div class="card-value">为学书院</div>
            </div>
          </div>

          <!-- GitHub Card -->
          <div class="info-card info-card-github" @click="openGitHub">
            <div class="card-glass">
              <div class="card-icon-wrapper">
                <i class="ri-github-line"></i>
              </div>
              <div class="card-label">开源地址</div>
              <div class="card-value">
                FruitsAI/Orange
                <i class="ri-arrow-right-up-line external-icon"></i>
              </div>
            </div>
          </div>
        </div>

        <!-- Tech Stack Pills -->
        <div class="tech-stack">
          <div class="tech-pill" data-tech="wails">
            <span class="tech-dot"></span>
            <span class="tech-name">Wails v3</span>
          </div>
          <div class="tech-pill" data-tech="vue">
            <span class="tech-dot"></span>
            <span class="tech-name">Vue 3</span>
          </div>
          <div class="tech-pill" data-tech="ts">
            <span class="tech-dot"></span>
            <span class="tech-name">TypeScript</span>
          </div>
          <div class="tech-pill" data-tech="go">
            <span class="tech-dot"></span>
            <span class="tech-name">Go</span>
          </div>
        </div>

        <!-- Update Button -->
        <button 
          class="update-btn" 
          @click="checkUpdate" 
          :disabled="checkingUpdate"
          :class="{ 'updating': checkingUpdate }"
        >
          <span class="btn-glow"></span>
          <span class="btn-content">
            <i class="ri-loop-left-line btn-icon" :class="{ 'spinning': checkingUpdate }"></i>
            <span class="btn-text">{{ checkingUpdate ? '正在检测更新...' : '检测更新' }}</span>
          </span>
        </button>

        <!-- Copyright -->
        <div class="copyright">
          <span class="copyright-text">© {{ new Date().getFullYear() }} FruitsAI</span>
          <span class="copyright-divider">·</span>
          <span class="copyright-rights">All rights reserved</span>
        </div>
      </div>
    </GlassCard>



  </div>
</template>

<style scoped>
.settings-view {
  display: grid;
  grid-template-columns: 250px 1fr;
  min-height: 100%;
}

@media (max-width: 768px) {
  .settings-view {
    grid-template-columns: 1fr;
  }
}

.border-color-border {
  border-color: var(--separator-color);
}

/* ===== 个人信息面板 ===== */
.profile-card {
  /* 与左侧导航栏对齐 */
  align-self: stretch;
}

.profile-panel {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  padding: 1.5rem;
  min-height: 100%;
  box-sizing: border-box;
}

/* 个人信息表单网格 */
.profile-form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.25rem;
}

.form-group-full {
  grid-column: span 2;
}

/* ===== 外观设置面板 ===== */
.appearance-card {
  /* 与左侧导航栏对齐 */
  align-self: stretch;
}

.appearance-panel {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  padding: 1.5rem;
  min-height: 100%;
  box-sizing: border-box;
}

.appearance-header {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.appearance-header-main {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
}

.appearance-title-wrapper {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.appearance-icon {
  width: 3rem;
  height: 3rem;
  border-radius: 0.75rem;
  background: linear-gradient(135deg, rgba(255, 159, 10, 0.2) 0%, rgba(255, 159, 10, 0.05) 100%);
  border: 1px solid rgba(255, 159, 10, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  color: #FF9F0A;
  box-shadow: 0 4px 12px rgba(255, 159, 10, 0.15);
}

.appearance-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: -0.02em;
}

.appearance-subtitle {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0.25rem 0 0 0;
}

.appearance-content {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

/* ===== 安全设置面板 ===== */
.security-card {
  /* 与左侧导航栏对齐 */
  align-self: stretch;
}

.security-panel {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  padding: 1.5rem;
  min-height: 100%;
  box-sizing: border-box;
}

/* 复用通知管理的头部样式 */
.dev-header {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.dev-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
}

.dev-title-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.dev-icon-wrapper {
  width: 3rem;
  height: 3rem;
  border-radius: 0.75rem;
  background: linear-gradient(135deg, rgba(255, 159, 10, 0.15) 0%, rgba(255, 159, 10, 0.05) 100%);
  border: 1px solid rgba(255, 159, 10, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  color: #FF9F0A;
  backdrop-filter: blur(8px);
}


.dev-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: -0.02em;
}

.dev-subtitle {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0.25rem 0 0 0;
}

/* 复用统计卡片样式 */
.dev-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
}

.dev-stat-card {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  padding: 1rem 1.25rem;
  background: var(--bg-elevated);
  border: 1px solid var(--border-color);
  border-radius: 0.75rem;
  transition: all 0.2s ease;
}

.dev-stat-card:hover {
  border-color: rgba(255, 159, 10, 0.3);
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.dev-stat-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.625rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
}

.dev-stat-icon.blue {
  background: rgba(59, 130, 246, 0.12);
  color: #3B82F6;
}

.dev-stat-icon.orange {
  background: rgba(245, 158, 11, 0.12);
  color: #F59E0B;
}

.dev-stat-icon.green {
  background: rgba(34, 197, 94, 0.12);
  color: #22C55E;
}

.dev-stat-info {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.dev-stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
}

.dev-stat-label {
  font-size: 0.75rem;
  color: var(--text-secondary);
  font-weight: 500;
}


/* 创建按钮样式 - 与通知管理保持一致 */
.dev-create-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1.25rem;
  background: #FF9F0A;
  color: white;
  border: none;
  border-radius: 0.625rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 4px 14px rgba(255, 159, 10, 0.3);
  position: relative;
  z-index: 10;
}

.dev-create-btn:hover {
  background: #E58909;
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(255, 159, 10, 0.4);
}

/* 内容区域 */
.dev-content {
  flex: 1;
}

.security-form-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.25rem;
}

/* 安全设置提交按钮 */
.security-submit-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.875rem 2rem;
  background: #FF9F0A;
  color: white;
  border: none;
  border-radius: 0.625rem;
  font-size: 0.9375rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 4px 14px rgba(255, 159, 10, 0.35);
  margin-top: 0.5rem;
}

.security-submit-btn:hover {
  background: #E58909;
  box-shadow: 0 6px 20px rgba(255, 159, 10, 0.45);
  transform: translateY(-1px);
}

.theme-card {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all 0.2s;
  gap: 12px;
}

.theme-card:hover {
  background: var(--bg-hover);
  border-color: var(--color-primary);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.theme-card.active {
  background: rgba(var(--color-primary-rgb), 0.05);
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.theme-icon {
  font-size: 24px;
  margin-bottom: 4px;
}

.theme-label {
  font-size: 14px;
  font-weight: 500;
}

.theme-check {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--color-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}


.nav-header {
  padding: var(--spacing-md);
  font-weight: 600;
  border-bottom: 1px solid var(--separator-color);
}

.nav-list {
  display: flex;
  flex-direction: column;
}

.nav-item-settings {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  color: var(--text-secondary);
  transition: all 0.2s;
  border-radius: 0;
  text-decoration: none;
  border-left: 3px solid transparent;
  white-space: nowrap;
}

.nav-icon {
  margin-right: var(--spacing-sm);
}

.nav-item-settings:hover {
  background: rgba(var(--text-primary-rgb), 0.03);
  color: var(--text-primary);
}

[data-theme='dark'] .nav-item-settings:hover {
  background: rgba(255, 255, 255, 0.05);
}

.nav-item-settings.active {
  background: rgba(var(--color-primary-rgb), 0.1);
  color: var(--color-primary);
  border-left-color: var(--color-primary);
}

.form-label {
  display: block;
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: var(--spacing-sm);
}

.form-input,
.form-select {
  width: 100%;
  padding: 8px 12px;
  background: var(--bg-base);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  outline: none;
  transition: border-color 0.2s;
  font-size: 14px;
}

.form-input:focus,
.form-select:focus {
  border-color: var(--color-primary);
}

.form-select {
  appearance: none;
  padding-right: 28px;
}

.input-wrapper {
  position: relative;
}

.input-wrapper .select-arrow {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 16px;
  color: var(--text-tertiary);
  pointer-events: none;
}

[data-theme='dark'] .form-input,
[data-theme='dark'] .form-select {
  background: rgba(0, 0, 0, 0.2);
  border-color: rgba(255, 255, 255, 0.1);
}

[data-theme='dark'] .form-input:focus,
[data-theme='dark'] .form-select:focus {
  border-color: var(--color-primary);
  background: rgba(0, 0, 0, 0.4);
}

/* 表单网格布局 */
.profile-form {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-lg);
}

.profile-form .col-span-2 {
  grid-column: span 2;
}

/* 小屏幕响应式 */
@media (max-width: 768px) {
  .profile-form {
    grid-template-columns: 1fr;
  }

  .profile-form .col-span-2 {
    grid-column: span 1;
  }

  /* 导航卡片改为水平滚动 */
  .nav-card {
    order: -1;
  }

  .nav-header {
    display: none;
  }

  .nav-list {
    flex-direction: row;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  .nav-list::-webkit-scrollbar {
    display: none;
  }

  .nav-item-settings {
    border-left: none;
    border-bottom: 2px solid transparent;
    padding: var(--spacing-sm) var(--spacing-md);
    flex-shrink: 0;
  }

  .nav-item-settings.active {
    border-bottom-color: var(--color-primary);
    border-left-color: transparent;
  }

  .nav-icon {
    margin-right: var(--spacing-xs);
  }

  /* 外观设置响应式 */
  .appearance-panel {
    padding: 1rem;
  }

  .appearance-header-main {
    flex-direction: column;
    align-items: stretch;
  }

  /* 安全设置响应式 */
  .security-panel {
    padding: 1rem;
  }

  .dev-header-content {
    flex-direction: column;
    align-items: stretch;
  }

  .security-form-grid {
    grid-template-columns: 1fr;
  }

  .dev-create-btn {
    width: 100%;
    justify-content: center;
  }
}

/* Dictionary Management Styles */
.dict-management {
  display: flex;
  flex-direction: column;
}

.dict-management .dev-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.dict-management .dev-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.dict-management .dev-title-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.dict-management .dev-icon-wrapper {
  width: 3rem;
  height: 3rem;
  border-radius: 0.75rem;
  background: linear-gradient(135deg, rgba(255, 159, 10, 0.15) 0%, rgba(255, 159, 10, 0.05) 100%);
  border: 1px solid rgba(255, 159, 10, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.375rem;
  color: #FF9F0A;
}

.dict-management .dev-title-info {
  display: flex;
  flex-direction: column;
}

.dict-management .dev-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.dict-management .dev-subtitle {
  font-size: 0.8125rem;
  color: var(--text-secondary);
  margin: 0;
}

.dict-management .dev-create-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1.25rem;
  background: #FF9F0A;
  border: none;
  border-radius: 0.625rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4px 12px rgba(255, 159, 10, 0.3);
}

.dict-management .dev-create-btn:hover {
  background: #F59300;
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(255, 159, 10, 0.4);
}

.dict-management .dev-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
}

.dict-management .dev-stat-card {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: var(--bg-elevated);
  border: 1px solid var(--border-color);
  border-radius: 0.75rem;
}

.dict-management .dev-stat-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.625rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.125rem;
}

.dict-management .dev-stat-icon.total {
  background: linear-gradient(135deg, rgba(255, 159, 10, 0.15) 0%, rgba(255, 159, 10, 0.05) 100%);
  color: #FF9F0A;
}

.dict-management .dev-stat-icon.items {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.15) 0%, rgba(59, 130, 246, 0.05) 100%);
  color: #3B82F6;
}

.dict-management .dev-stat-icon.active {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.15) 0%, rgba(34, 197, 94, 0.05) 100%);
  color: #22C55E;
}

.dict-management .dev-stat-info {
  display: flex;
  flex-direction: column;
}

.dict-management .dev-stat-value {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--text-primary);
}

.dict-management .dev-stat-value.dict-code {
  font-size: 0.875rem;
  font-weight: 600;
  font-family: monospace;
}

.dict-management .dev-stat-label {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.dict-layout {
  display: flex;
}

.dict-sidebar {
  width: 180px;
  border-right: 1px solid var(--border-color);
  background: transparent;
  overflow-y: auto;
  padding: 0.75rem 0;
}

.dict-sidebar-title {
  padding: 0.5rem 1rem;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.dict-nav-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  cursor: pointer;
  font-size: 0.875rem;
  color: var(--text-secondary);
  border-left: 3px solid transparent;
  transition: all 0.2s;
}

.dict-nav-item i {
  font-size: 1rem;
  opacity: 0.7;
}

.dict-nav-item:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.dict-nav-item.active {
  background: rgba(255, 159, 10, 0.08);
  color: #FF9F0A;
  border-left-color: #FF9F0A;
  font-weight: 500;
}

.dict-nav-item.active i {
  opacity: 1;
}

.dict-content {
  flex: 1;
  padding: 1rem;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.dict-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  text-align: center;
}

.dict-empty-icon {
  width: 4rem;
  height: 4rem;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(255, 159, 10, 0.1) 0%, rgba(255, 159, 10, 0.05) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  color: #FF9F0A;
  margin-bottom: 1rem;
}

.dict-empty-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

.dict-empty-desc {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

.dict-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.dict-item-card {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  background: var(--bg-elevated);
  border: 1px solid var(--border-color);
  border-radius: 0.625rem;
  transition: all 0.2s;
}

.dict-item-card:hover {
  border-color: rgba(255, 159, 10, 0.3);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
}

.dict-item-icon {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 0.5rem;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.15) 0%, rgba(59, 130, 246, 0.05) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  color: #3B82F6;
  flex-shrink: 0;
}



/* Reusing ConfirmModal styles for consistency */
.confirm-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition: all 0.5s cubic-bezier(0.19, 1, 0.22, 1);
}

[data-theme='dark'] .confirm-overlay {
  background: rgba(0, 0, 0, 0.5);
}

.confirm-modal {
  /* Liquid Glass Base */
  background: var(--bg-elevated);
  backdrop-filter: 
    blur(var(--glass-blur)) 
    saturate(var(--glass-saturation))
    brightness(var(--glass-brightness));
  -webkit-backdrop-filter: 
    blur(var(--glass-blur)) 
    saturate(var(--glass-saturation))
    brightness(var(--glass-brightness));
    
  border: var(--glass-border);
  border-radius: var(--radius-xl);
  box-shadow: 
    var(--glass-shadow-outer),
    0 20px 40px rgba(0,0,0,0.1);
    
  padding: 32px;
  width: 90%;
  max-width: 400px;
  text-align: center;
  position: relative;
  overflow: hidden;
  isolation: isolate;
}

/* Specular Highlight */
.confirm-modal::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--glass-specular);
  pointer-events: none;
  z-index: 0;
  border-radius: inherit;
  opacity: 0.6;
  mix-blend-mode: overlay;
}

/* Inner Glow */
.confirm-modal::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  box-shadow:
    inset 0 0 0 1px rgba(255, 255, 255, 0.15),
    inset 0 1px 2px rgba(255, 255, 255, 0.2);
  pointer-events: none;
  z-index: 1;
}

[data-theme='dark'] .confirm-modal::after {
  box-shadow:
    inset 0 0 0 1px rgba(255, 255, 255, 0.08),
    inset 0 1px 2px rgba(255, 255, 255, 0.05);
}

.confirm-title, .form-group, .confirm-actions {
  position: relative;
  z-index: 2;
}

.confirm-title {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 24px;
  color: var(--text-primary);
  letter-spacing: -0.01em;
}

.confirm-actions {
  display: flex;
  gap: 16px;
  justify-content: flex-end;
}
/* Transitions - Apple Spring */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.4s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-active .confirm-modal,
.fade-leave-active .confirm-modal {
  transition: all 0.5s cubic-bezier(0.19, 1, 0.22, 1);
}

.fade-enter-from .confirm-modal,
.fade-leave-to .confirm-modal {
  transform: scale(0.9);
  opacity: 0;
}

/* ============================================
   About Page - Liquid Glass Premium Edition
   ============================================ */

.about-page {
  --about-orange: #FF9F0A;
  --about-orange-soft: rgba(255, 159, 10, 0.15);
  --about-orange-glow: rgba(255, 159, 10, 0.4);
  --about-card-bg: rgba(255, 255, 255, 0.6);
  --about-card-border: rgba(255, 255, 255, 0.8);
  --about-text-muted: rgba(60, 60, 67, 0.6);
  
  container-type: inline-size;
}

[data-theme='dark'] .about-page {
  --about-card-bg: rgba(30, 30, 30, 0.5);
  --about-card-border: rgba(255, 255, 255, 0.12);
  --about-text-muted: rgba(255, 255, 255, 0.5);
}

.about-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 2rem;
  gap: 2.5rem;
  min-height: 100%;
  animation: about-fade-in 0.8s cubic-bezier(0.19, 1, 0.22, 1) forwards;
}

@keyframes about-fade-in {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Hero Section */
.about-hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  animation: hero-entrance 0.9s cubic-bezier(0.19, 1, 0.22, 1) 0.1s both;
}

@keyframes hero-entrance {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.logo-wrapper {
  position: relative;
  margin-bottom: 1.25rem;
}

.logo-glow {
  position: absolute;
  inset: -20%;
  background: radial-gradient(circle at center, var(--about-orange-glow) 0%, transparent 70%);
  opacity: 0.6;
  filter: blur(20px);
  animation: logo-pulse 4s ease-in-out infinite;
}

@keyframes logo-pulse {
  0%, 100% { transform: scale(1); opacity: 0.6; }
  50% { transform: scale(1.1); opacity: 0.8; }
}

.about-logo {
  position: relative;
  width: 6rem;
  height: 6rem;
  object-fit: contain;
  filter: drop-shadow(0 8px 24px rgba(255, 159, 10, 0.35));
  transition: transform 0.5s cubic-bezier(0.19, 1, 0.22, 1);
}

.about-logo:hover {
  transform: scale(1.08) rotate(-3deg);
}

.about-title {
  font-size: 2.5rem;
  font-weight: 800;
  background: linear-gradient(135deg, var(--about-orange) 0%, #FFB340 50%, #FF8C00 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.03em;
  margin-bottom: 0.75rem;
  text-shadow: 0 2px 40px rgba(255, 159, 10, 0.2);
}

.about-subtitle {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  flex-wrap: wrap;
  justify-content: center;
}

.version-badge {
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0.35em 0.75em;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.12) 0%, rgba(59, 130, 246, 0.05) 100%);
  color: #3B82F6;
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 100px;
  letter-spacing: 0.02em;
  box-shadow: 
    0 2px 8px rgba(59, 130, 246, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

[data-theme='dark'] .version-badge {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.2) 0%, rgba(59, 130, 246, 0.1) 100%);
  border-color: rgba(59, 130, 246, 0.3);
}

.tagline {
  font-size: 0.9375rem;
  color: var(--about-text-muted);
  font-weight: 500;
  letter-spacing: 0.01em;
}

/* Floating Info Cards */
.info-cards-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
  width: 100%;
  max-width: 720px;
  animation: cards-entrance 0.9s cubic-bezier(0.19, 1, 0.22, 1) 0.2s both;
}

@keyframes cards-entrance {
  from {
    opacity: 0;
    transform: translateY(40px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.info-card {
  position: relative;
  perspective: 1000px;
}

.card-glass {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 1.75rem 1.25rem;
  min-height: 180px;
  background: var(--about-card-bg);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid var(--about-card-border);
  border-radius: 20px;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.04),
    0 10px 20px -4px rgba(0, 0, 0, 0.06),
    inset 0 1px 0 rgba(255, 255, 255, 0.6);
  transition: all 0.4s cubic-bezier(0.19, 1, 0.22, 1);
  overflow: hidden;
}

[data-theme='dark'] .card-glass {
  box-shadow: 
    0 4px 6px -1px rgba(0, 0, 0, 0.2),
    0 10px 20px -4px rgba(0, 0, 0, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.08);
}

.card-glass::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.4) 0%,
    rgba(255, 255, 255, 0.1) 50%,
    transparent 100%
  );
  border-radius: inherit;
  pointer-events: none;
}

.info-card:hover .card-glass {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 
    0 20px 40px -8px rgba(0, 0, 0, 0.1),
    0 8px 16px -4px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
}

.info-card-author:hover .card-glass {
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 
    0 20px 40px -8px rgba(59, 130, 246, 0.2),
    0 8px 16px -4px rgba(59, 130, 246, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
}

.info-card-wechat:hover .card-glass {
  border-color: rgba(34, 197, 94, 0.4);
  box-shadow: 
    0 20px 40px -8px rgba(34, 197, 94, 0.2),
    0 8px 16px -4px rgba(34, 197, 94, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
}

.info-card-github {
  cursor: pointer;
}

.info-card-github:hover .card-glass {
  border-color: rgba(100, 100, 100, 0.4);
  box-shadow: 
    0 20px 40px -8px rgba(100, 100, 100, 0.2),
    0 8px 16px -4px rgba(100, 100, 100, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
}

.card-icon-wrapper {
  position: relative;
  width: 3.5rem;
  height: 3.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  margin-bottom: 1rem;
  font-size: 1.5rem;
  transition: all 0.4s cubic-bezier(0.19, 1, 0.22, 1);
}

.info-card-author .card-icon-wrapper {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.12) 0%, rgba(59, 130, 246, 0.05) 100%);
  color: #3B82F6;
  box-shadow: 
    0 4px 12px rgba(59, 130, 246, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

.info-card-wechat .card-icon-wrapper {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.12) 0%, rgba(34, 197, 94, 0.05) 100%);
  color: #22C55E;
  box-shadow: 
    0 4px 12px rgba(34, 197, 94, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

.info-card-github .card-icon-wrapper {
  background: linear-gradient(135deg, rgba(100, 100, 100, 0.12) 0%, rgba(100, 100, 100, 0.05) 100%);
  color: #666;
  box-shadow: 
    0 4px 12px rgba(100, 100, 100, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

[data-theme='dark'] .info-card-github .card-icon-wrapper {
  color: #999;
}

.info-card:hover .card-icon-wrapper {
  transform: scale(1.1) rotate(-5deg);
}

.card-label {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--about-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin-bottom: 0.375rem;
}

.card-value {
  font-size: 1.0625rem;
  font-weight: 700;
  color: var(--about-orange);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.25rem;
  white-space: nowrap;
}

.external-icon {
  font-size: 0.875rem;
  opacity: 0.6;
  transition: all 0.3s ease;
}

.info-card-github:hover .external-icon {
  opacity: 1;
  transform: translate(2px, -2px);
}

/* Tech Stack */
.tech-stack {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.5rem;
  animation: tech-entrance 0.9s cubic-bezier(0.19, 1, 0.22, 1) 0.3s both;
}

@keyframes tech-entrance {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.tech-pill {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.875rem;
  background: var(--about-card-bg);
  backdrop-filter: blur(12px);
  border: 1px solid var(--about-card-border);
  border-radius: 100px;
  font-size: 0.75rem;
  font-weight: 600;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', monospace;
  color: var(--text-secondary);
  transition: all 0.3s cubic-bezier(0.19, 1, 0.22, 1);
  cursor: default;
  box-shadow:
    0 2px 4px rgba(0, 0, 0, 0.02),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

.tech-pill:hover {
  transform: translateY(-2px);
  box-shadow: 
    0 8px 16px rgba(0, 0, 0, 0.06),
    inset 0 1px 0 rgba(255, 255, 255, 0.6);
}

.tech-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  animation: tech-dot-pulse 2s ease-in-out infinite;
}

.tech-pill[data-tech="wails"] .tech-dot {
  background: #E53E3E;
  box-shadow: 0 0 8px rgba(229, 62, 62, 0.5);
}

.tech-pill[data-tech="vue"] .tech-dot {
  background: #42B883;
  box-shadow: 0 0 8px rgba(66, 184, 131, 0.5);
  animation-delay: 0.3s;
}

.tech-pill[data-tech="ts"] .tech-dot {
  background: #3178C6;
  box-shadow: 0 0 8px rgba(49, 120, 198, 0.5);
  animation-delay: 0.6s;
}

.tech-pill[data-tech="go"] .tech-dot {
  background: #00ADD8;
  box-shadow: 0 0 8px rgba(0, 173, 216, 0.5);
  animation-delay: 0.9s;
}

@keyframes tech-dot-pulse {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.3); opacity: 0.7; }
}

.tech-pill:hover .tech-dot {
  animation-duration: 0.8s;
}

/* Update Button */
.update-btn {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.625rem 1.5rem;
  background: linear-gradient(135deg, var(--about-orange) 0%, #FFB340 100%);
  border: none;
  border-radius: 10px;
  font-size: 0.875rem;
  font-weight: 600;
  color: white;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.19, 1, 0.22, 1);
  box-shadow:
    0 4px 12px rgba(255, 159, 10, 0.3),
    0 6px 20px rgba(255, 159, 10, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
  animation: btn-entrance 0.9s cubic-bezier(0.19, 1, 0.22, 1) 0.4s both;
}

@keyframes btn-entrance {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.btn-glow {
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at center, rgba(255, 255, 255, 0.3) 0%, transparent 70%);
  opacity: 0;
  transition: opacity 0.4s ease;
}

.update-btn:hover {
  transform: translateY(-3px) scale(1.02);
  box-shadow: 
    0 8px 24px rgba(255, 159, 10, 0.45),
    0 16px 48px rgba(255, 159, 10, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.4);
}

.update-btn:hover .btn-glow {
  opacity: 1;
}

.update-btn:active:not(:disabled) {
  transform: translateY(-1px) scale(0.98);
  box-shadow: 
    0 4px 12px rgba(255, 159, 10, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.update-btn:disabled {
  cursor: not-allowed;
  opacity: 0.8;
}

.update-btn.updating {
  background: linear-gradient(135deg, #9CA3AF 0%, #6B7280 100%);
  box-shadow: 
    0 4px 16px rgba(0, 0, 0, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.btn-content {
  position: relative;
  display: flex;
  align-items: center;
  gap: 0.625rem;
  z-index: 1;
}

.btn-icon {
  font-size: 1.125rem;
  transition: transform 0.3s ease;
}

.btn-icon.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.btn-text {
  letter-spacing: 0.02em;
}

/* Copyright */
.copyright {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.8125rem;
  color: var(--about-text-muted);
  font-weight: 500;
  animation: copyright-entrance 0.9s cubic-bezier(0.19, 1, 0.22, 1) 0.5s both;
}

@keyframes copyright-entrance {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.copyright-divider {
  opacity: 0.4;
}

.copyright-rights {
  opacity: 0.7;
}

/* Responsive Adjustments */
@media (max-width: 640px) {
  .info-cards-grid {
    grid-template-columns: 1fr;
    max-width: 320px;
  }
  
  .about-subtitle {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .tech-stack {
    gap: 0.5rem;
  }
  
  .tech-pill {
    padding: 0.375rem 0.75rem;
    font-size: 0.75rem;
  }
  
  .update-btn {
    padding: 0.875rem 2rem;
    width: 100%;
    max-width: 280px;
  }
}

@media (min-width: 641px) and (max-width: 1024px) {
  .info-cards-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: 1rem;
  }
  
  .card-glass {
    padding: 1.5rem 1rem;
  }
  
  .card-icon-wrapper {
    width: 3rem;
    height: 3rem;
    font-size: 1.25rem;
  }
}

/* Reduced Motion */
@media (prefers-reduced-motion: reduce) {
  .about-container,
  .about-hero,
  .info-cards-grid,
  .tech-stack,
  .update-btn,
  .copyright {
    animation: none;
    opacity: 1;
    transform: none;
  }
  
  .logo-glow {
    animation: none;
  }
  
  .tech-dot {
    animation: none;
  }
}
</style>