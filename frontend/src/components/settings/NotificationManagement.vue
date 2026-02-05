<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { notificationApi, type Notification, type UserBrief } from '@/api/notification'
import { useToast } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'
import { Events } from '@wailsio/runtime'
import NotificationDetailModal from '@/components/notification/NotificationDetailModal.vue'

const props = defineProps<{
  isAdmin: boolean
}>()

const toast = useToast()
const { confirm } = useConfirm()

// 通知列表
const notifications = ref<Notification[]>([])
const notificationTotal = ref(0)
const targetUsers = ref<UserBrief[]>([])
const notificationLoading = ref(false)

// 分页状态
const notificationCurrentPage = ref(1)
const notificationPageSize = ref(5)

// 创建/编辑弹窗
const showCreateNotificationModal = ref(false)
const creatingNotification = ref(false)
const newNotification = ref({
  title: '',
  content: '',
  type: 'system',
  target_user_id: 0,
})

// 详情/编辑弹窗状态
const showNotificationDetailModal = ref(false)
const selectedNotification = ref<Notification | null>(null)
const isEditingNotification = ref(false)

// 统计计算
const unreadNotifications = computed(() => notifications.value.filter(n => !n.is_read).length)
const systemNotifications = computed(() => notifications.value.filter(n => n.type === 1).length)

// 分页计算
const notificationTotalPages = computed(() => Math.ceil(notificationTotal.value / notificationPageSize.value))

const notificationPaginationInfo = computed(() => {
  const total = notificationTotal.value
  if (total === 0) return '暂无数据'
  const start = (notificationCurrentPage.value - 1) * notificationPageSize.value + 1
  const end = Math.min(notificationCurrentPage.value * notificationPageSize.value, total)
  return `显示 ${start}-${end} 条，共 ${total} 条`
})

// 分页操作
const notificationPrevPage = () => {
  if (notificationCurrentPage.value > 1) {
    notificationCurrentPage.value--
    loadNotifications()
  }
}

const notificationNextPage = () => {
  if (notificationCurrentPage.value < notificationTotalPages.value) {
    notificationCurrentPage.value++
    loadNotifications()
  }
}

const notificationGoToPage = (page: number) => {
  if (page !== notificationCurrentPage.value) {
    notificationCurrentPage.value = page
    loadNotifications()
  }
}

const visibleNotificationPages = computed(() => {
  const total = notificationTotalPages.value
  const current = notificationCurrentPage.value
  const delta = 2 // Show 2 pages before and after current
  
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }
  
  const pages: (number | string)[] = [1]
  
  // Calculate range
  let start = current - delta
  let end = current + delta
  
  // Adjust window if close to start
  if (start <= 2) {
    start = 2
    end = Math.min(6, total - 1)
  }
  // Adjust window if close to end
  else if (end >= total - 1) {
    end = total - 1
    start = Math.max(total - 5, 2)
  }
  
  if (start > 2) {
    pages.push('...')
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  
  if (end < total - 1) {
    pages.push('...')
  }
  
  if (total > 1) {
    pages.push(total)
  }
  
  return pages
})

// 监听 pageSize 变化重置页码
watch(notificationPageSize, () => {
  notificationCurrentPage.value = 1
  loadNotifications()
})

// 加载通知列表
const loadNotifications = async () => {
  notificationLoading.value = true
  try {
    const res = await notificationApi.list(notificationCurrentPage.value, notificationPageSize.value)
    if (res.data.code === 0) {
      notifications.value = res.data.data.list
      notificationTotal.value = res.data.data.total
    }
  } catch (error) {
    console.error('Failed to load notifications:', error)
  } finally {
    notificationLoading.value = false
  }
}

// 加载目标用户列表
const loadTargetUsers = async () => {
  try {
    const res = await notificationApi.getUsers()
    if (res.data.code === 0) {
      targetUsers.value = res.data.data
    }
  } catch (error) {
    console.error('Failed to load users:', error)
  }
}

// 创建通知
const handleCreateNotification = async () => {
  if (!newNotification.value.title || !newNotification.value.content) {
    toast.error('请填写标题和内容')
    return
  }
  creatingNotification.value = true
  try {
    const res = await notificationApi.create(newNotification.value)
    if (res.data.code === 0) {
      toast.success('发送成功')
      showCreateNotificationModal.value = false
      newNotification.value = { title: '', content: '', type: 'system', target_user_id: 0 }
      loadNotifications()
    } else {
      toast.error(res.data.message || '发送失败')
    }
  } catch (error) {
    console.error('Failed to create notification:', error)
    toast.error('发送失败')
  } finally {
    creatingNotification.value = false
  }
}

// 删除通知
const handleDeleteNotification = async (id: number) => {
  const confirmed = await confirm({ title: '确认删除', message: '确定要删除这条通知吗？' })
  if (confirmed) {
    try {
      const res = await notificationApi.delete(id)
      if (res.data.code === 0) {
        toast.success('删除成功')
        loadNotifications()
      }
    } catch (error) {
      console.error('Failed to delete notification:', error)
      toast.error('删除失败')
    }
  }
}

// 查看详情（自动标记已读）
const viewNotificationDetail = async (notification: Notification) => {
  selectedNotification.value = notification
  isEditingNotification.value = false
  showNotificationDetailModal.value = true
  
  if (!notification.is_read) {
    try {
      await notificationApi.markAsRead(notification.id)
      notification.is_read = true
      Events.Emit('notification_updated')
    } catch (error) {
      console.error('Mark as read failed', error)
    }
  }
}

// 编辑通知
const editNotification = (notification: Notification) => {
  let typeStr = 'system'
  if (notification.type === 2) typeStr = 'activity'
  
  newNotification.value = {
    title: notification.title,
    content: notification.content,
    type: typeStr,
    target_user_id: notification.is_global === 1 ? 0 : 0
  }
  selectedNotification.value = notification
  isEditingNotification.value = true
  showCreateNotificationModal.value = true
}

// 更新通知
const handleUpdateNotification = async () => {
  if (!selectedNotification.value) return

  if (!newNotification.value.title || !newNotification.value.content) {
    toast.error('请填写标题和内容')
    return
  }
  
  creatingNotification.value = true
  try {
    const res = await notificationApi.update(selectedNotification.value.id, newNotification.value)
    if (res.data.code === 0) {
      toast.success('更新成功')
      showCreateNotificationModal.value = false
      loadNotifications()
    } else {
      toast.error(res.data.message || '更新失败')
    }
  } catch (error) {
    console.error('Failed to update notification:', error)
    toast.error('更新失败')
  } finally {
    creatingNotification.value = false
  }
}

// 打开创建弹窗
const openCreateModal = () => {
  isEditingNotification.value = false
  newNotification.value = { title: '', content: '', type: 'system', target_user_id: 0 }
  showCreateNotificationModal.value = true
}

onMounted(() => {
  loadNotifications()
  if (props.isAdmin) {
    loadTargetUsers()
  }
})
</script>

<template>
  <div class="notification-management">
    <!-- 头部区域 -->
    <div class="dev-header">
      <div class="dev-header-content">
        <div class="dev-title-section">
          <div class="dev-icon-wrapper notif-icon">
            <i class="ri-notification-3-line"></i>
          </div>
          <div class="dev-title-info">
            <h2 class="dev-title">通知管理</h2>
            <p class="dev-subtitle">查看系统消息{{ isAdmin ? '，管理员可发送通知' : '' }}</p>
          </div>
        </div>
        <button v-if="isAdmin" class="dev-create-btn" @click="openCreateModal">
          <i class="ri-add-line"></i>
          <span>发送通知</span>
        </button>
      </div>

      <!-- 统计卡片 -->
      <div class="dev-stats">
        <div class="dev-stat-card">
          <div class="dev-stat-icon total">
            <i class="ri-mail-line"></i>
          </div>
          <div class="dev-stat-info">
            <span class="dev-stat-value">{{ notificationTotal }}</span>
            <span class="dev-stat-label">总通知</span>
          </div>
        </div>
        <div class="dev-stat-card">
          <div class="dev-stat-icon unread">
            <i class="ri-mail-unread-line"></i>
          </div>
          <div class="dev-stat-info">
            <span class="dev-stat-value">{{ unreadNotifications }}</span>
            <span class="dev-stat-label">未读</span>
          </div>
        </div>
        <div class="dev-stat-card">
          <div class="dev-stat-icon system">
            <i class="ri-computer-line"></i>
          </div>
          <div class="dev-stat-info">
            <span class="dev-stat-value">{{ systemNotifications }}</span>
            <span class="dev-stat-label">系统</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="dev-content">
      <!-- 加载状态 -->
      <div v-if="notificationLoading" class="dev-loading">
        <div class="dev-loading-spinner">
          <i class="ri-loader-4-line animate-spin"></i>
        </div>
        <span>正在加载通知列表...</span>
      </div>

      <!-- 空状态 -->
      <div v-else-if="notifications.length === 0" class="dev-empty">
        <div class="dev-empty-icon notif-empty-icon">
          <i class="ri-notification-off-line"></i>
        </div>
        <h3 class="dev-empty-title">暂无通知</h3>
        <p class="dev-empty-desc">目前没有收到任何系统消息</p>
      </div>

      <!-- 通知列表 -->
      <div v-else class="notification-list">
        <div
          v-for="notification in notifications"
          :key="notification.id"
          class="notification-card"
          :class="{ 
            'unread': !notification.is_read,
            'type-system': notification.type === 1,
            'type-activity': notification.type === 2,
            'type-private': notification.type === 3
          }"
          @click="viewNotificationDetail(notification)"
        >
          <div class="notification-inner">
            <!-- 类型图标 -->
            <div class="notification-type-icon" :class="'type-' + notification.type">
              <i :class="notification.type === 2 ? 'ri-calendar-event-line' : (notification.type === 3 ? 'ri-mail-send-line' : 'ri-settings-3-line')"></i>
            </div>
            
            <!-- 内容区域 -->
            <div class="notification-body">
              <div class="notification-header-row">
                <div class="notification-tags">
                  <span class="notification-tag" :class="'tag-' + notification.type">
                    {{ notification.type === 2 ? '活动' : (notification.type === 3 ? '私信' : '系统') }}
                  </span>
                  <span class="notification-scope">{{ notification.is_global === 1 ? '全员' : '私信' }}</span>
                </div>
              </div>
              
              <h4 class="notification-title">
                <span v-if="!notification.is_read" class="unread-dot"></span>
                {{ notification.title }}
              </h4>
              
              <p class="notification-desc">{{ notification.content }}</p>
              
              <div v-if="notification.sender" class="notification-sender">
                <i class="ri-user-line"></i>
                <span>{{ notification.sender.name }}</span>
              </div>
            </div>
            
            <!-- 右侧区域：操作按钮 + 时间 -->
            <div class="notification-right">
              <div class="notification-actions" v-if="isAdmin">
                <button
                  class="action-btn"
                  @click.stop="editNotification(notification)"
                  title="编辑"
                >
                  <i class="ri-edit-line"></i>
                </button>
                <button
                  class="action-btn delete"
                  @click.stop="handleDeleteNotification(notification.id)"
                  title="删除"
                >
                  <i class="ri-delete-bin-line"></i>
                </button>
              </div>
              <span class="notification-time">{{ new Date(notification.create_time).toLocaleString() }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="notifications.length > 0" class="notification-pagination">
        <div class="pagination-inner">
          <span class="pagination-info">{{ notificationPaginationInfo }}</span>
          
          <div class="pagination-controls">
            <button 
              class="page-btn" 
              :disabled="notificationCurrentPage === 1" 
              @click="notificationPrevPage"
            >
              <i class="ri-arrow-left-s-line"></i>
            </button>
            
            <div class="page-numbers">
              <button
                v-for="(page, index) in visibleNotificationPages"
                :key="index"
                class="page-number"
                :class="{ active: notificationCurrentPage === page, 'cursor-default': page === '...' }"
                :disabled="notificationCurrentPage === page || page === '...'"
                @click="typeof page === 'number' && notificationGoToPage(page)"
              >
                {{ page }}
              </button>
            </div>
            
            <button 
              class="page-btn" 
              :disabled="notificationCurrentPage === notificationTotalPages" 
              @click="notificationNextPage"
            >
              <i class="ri-arrow-right-s-line"></i>
            </button>
          </div>
          
          <div class="page-size">
            <select v-model="notificationPageSize" class="page-select">
              <option :value="5">5条/页</option>
              <option :value="10">10条/页</option>
            </select>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑通知弹窗 -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showCreateNotificationModal" class="modal-overlay open" @click.self="showCreateNotificationModal = false">
          <div class="modal">
            <div class="modal-header">
              <h3 class="modal-title">{{ isEditingNotification ? '编辑通知' : '发送通知' }}</h3>
              <button class="modal-close" @click="showCreateNotificationModal = false">
                <i class="ri-close-line"></i>
              </button>
            </div>
            <div class="modal-body">
              <div class="form-group">
                <label class="form-label">通知标题</label>
                <input
                  type="text"
                  v-model="newNotification.title"
                  class="form-input"
                  placeholder="请输入通知标题"
                />
              </div>
              <div class="form-group">
                <label class="form-label">通知内容</label>
                <textarea
                  v-model="newNotification.content"
                  class="form-input"
                  rows="4"
                  placeholder="请输入通知内容"
                ></textarea>
              </div>
              <div class="form-group">
                <label class="form-label">通知类型</label>
                <div class="input-wrapper">
                  <select v-model="newNotification.type" class="form-select">
                    <option value="system">系统通知</option>
                    <option value="activity">活动通知</option>
                  </select>
                  <i class="ri-arrow-down-s-line select-arrow"></i>
                </div>
              </div>
              <div class="form-group" v-if="isAdmin">
                <label class="form-label">发送对象</label>
                <div class="input-wrapper">
                  <select v-model="newNotification.target_user_id" class="form-select">
                    <option :value="0">全员通知</option>
                    <option v-for="user in targetUsers" :key="user.id" :value="user.id">
                      {{ user.name }} ({{ user.username }})
                    </option>
                  </select>
                  <i class="ri-arrow-down-s-line select-arrow"></i>
                </div>
              </div>
            </div>
            <div class="modal-footer">
              <button class="btn btn-secondary" @click="showCreateNotificationModal = false">取消</button>
              <button 
                class="btn btn-primary" 
                @click="isEditingNotification ? handleUpdateNotification() : handleCreateNotification()" 
                :disabled="creatingNotification"
              >
                {{ creatingNotification ? '提交中...' : (isEditingNotification ? '更新' : '发送') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- 通知详情弹窗 -->
    <NotificationDetailModal
      v-model="showNotificationDetailModal"
      :notification="selectedNotification"
    />
  </div>
</template>

<style scoped>
/* ===== 通知管理组件 - 与开发设置风格统一 ===== */
.notification-management {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  padding: 1.5rem;
  min-height: 100%;
}

/* 头部区域 */
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

/* 统计卡片 */
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

.dev-stat-icon.total {
  background: rgba(59, 130, 246, 0.12);
  color: #3B82F6;
}

.dev-stat-icon.unread {
  background: rgba(245, 158, 11, 0.12);
  color: #F59E0B;
}

.dev-stat-icon.system {
  background: rgba(139, 92, 246, 0.12);
  color: #8B5CF6;
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

/* 内容区域 */
.dev-content {
  flex: 1;
}

/* 加载状态 */
.dev-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 4rem 0;
  color: var(--text-secondary);
}

.dev-loading-spinner {
  width: 3rem;
  height: 3rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  color: #FF9F0A;
}

/* 空状态 */
.dev-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 0;
  text-align: center;
}

.dev-empty-icon {
  width: 5rem;
  height: 5rem;
  border-radius: 1rem;
  background: linear-gradient(135deg, rgba(255, 159, 10, 0.1) 0%, rgba(255, 159, 10, 0.02) 100%);
  border: 1px solid rgba(255, 159, 10, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2.5rem;
  color: #FF9F0A;
  margin-bottom: 1.5rem;
}

.dev-empty-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

.dev-empty-desc {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

/* 通知列表 */
.notification-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

/* 通知卡片 - 使用 !important 覆盖 main.css */
.notification-card {
  position: relative !important;
  display: block !important;
  background: var(--bg-elevated) !important;
  border: 1px solid var(--border-color) !important;
  border-radius: 0.875rem !important;
  overflow: visible !important;
  transition: all 0.2s ease;
  cursor: pointer;
  padding: 0 !important;
}

.notification-card:hover {
  border-color: rgba(255, 159, 10, 0.3) !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
  transform: translateY(-2px) !important;
}

.notification-card.unread {
  border-left: 3px solid #F59E0B;
}

.notification-inner {
  display: flex !important;
  align-items: center !important;
  gap: 1rem;
  padding: 1.25rem;
  padding-right: 5rem;
  position: relative;
}

/* 类型图标 */
.notification-type-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.625rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
  flex-shrink: 0;
}

.notification-type-icon.type-1 {
  background: rgba(100, 116, 139, 0.1);
  color: #64748B;
}

.notification-type-icon.type-2 {
  background: rgba(239, 68, 68, 0.1);
  color: #EF4444;
}

.notification-type-icon.type-3 {
  background: rgba(139, 92, 246, 0.1);
  color: #8B5CF6;
}

/* 通知内容 */
.notification-body {
  flex: 1;
  min-width: 0;
}

.notification-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.notification-tags {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.notification-tag {
  padding: 0.25rem 0.625rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
}

.notification-tag.tag-1 {
  background: #F1F5F9;
  color: #64748B;
  border: 1px solid #CBD5E1;
}

.notification-tag.tag-2 {
  background: #FEF2F2;
  color: #EF4444;
  border: 1px solid #FECACA;
}

.notification-tag.tag-3 {
  background: #F5F3FF;
  color: #8B5CF6;
  border: 1px solid #DDD6FE;
}

.notification-scope {
  font-size: 0.75rem;
  color: var(--text-tertiary);
}

.notification-time {
  font-size: 0.8125rem;
  color: var(--text-tertiary);
}

.notification-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.375rem 0;
  line-height: 1.4;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.unread-dot {
  width: 0.5rem;
  height: 0.5rem;
  border-radius: 9999px;
  background: #F59E0B;
  flex-shrink: 0;
}

.notification-desc {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0 0 0.5rem 0;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notification-sender {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  font-size: 0.8125rem;
  color: var(--text-tertiary);
}

/* 右侧区域 - 包含操作按钮和时间 */
.notification-right {
  position: absolute !important;
  right: 1.25rem !important;
  top: 50% !important;
  transform: translateY(-50%) !important;
  display: flex !important;
  flex-direction: column !important;
  align-items: flex-end !important;
  gap: 0.5rem !important;
  min-height: 4rem;
  justify-content: space-between !important;
  z-index: 10;
}

.notification-right .notification-time {
  font-size: 0.75rem;
  color: var(--text-tertiary);
  white-space: nowrap;
}

/* 操作按钮 */
.notification-actions {
  display: flex !important;
  flex-direction: column !important;
  gap: 0.5rem !important;
  opacity: 0;
  transition: opacity 0.2s;
  align-items: center !important;
  pointer-events: auto;
}

.notification-card:hover .notification-actions {
  opacity: 1;
}

.action-btn {
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-base);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  font-size: 1rem;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background: rgba(255, 159, 10, 0.08);
  color: #FF9F0A;
  border-color: rgba(255, 159, 10, 0.3);
}

.action-btn.delete {
  color: #EF4444;
}

.action-btn.delete:hover {
  background: rgba(239, 68, 68, 0.08);
  border-color: rgba(239, 68, 68, 0.3);
}

/* 分页 - 使用 !important 覆盖 main.css 全局样式 */
.notification-pagination {
  display: block !important;
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
}

.pagination-inner {
  display: grid !important;
  grid-template-columns: 1fr auto 1fr !important;
  align-items: center !important;
  gap: 1rem;
}

.pagination-info {
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-weight: 500;
  justify-self: start !important;
}

.pagination-controls {
  display: flex !important;
  align-items: center !important;
  gap: 0.375rem !important;
  justify-self: center !important;
}

.page-btn {
  width: 2.25rem;
  height: 2.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-base);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  font-size: 1.125rem;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  background: var(--bg-elevated);
  color: var(--text-primary);
  border-color: var(--text-tertiary);
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-numbers {
  display: flex;
  gap: 0.25rem;
}

.page-number {
  min-width: 2.25rem;
  height: 2.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  font-size: 0.875rem;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
  padding: 0 0.625rem;
}

.page-number:hover {
  background: var(--bg-elevated);
  color: var(--text-primary);
}

.page-number.active {
  background: #FF9F0A;
  border-color: transparent;
  color: white;
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(255, 159, 10, 0.3);
}

.page-size {
  display: flex !important;
  justify-content: flex-end !important;
  justify-self: end !important;
}

.page-size .page-select {
  padding: 0.5rem 2rem 0.5rem 0.875rem;
  background: var(--bg-base);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  font-size: 0.875rem;
  color: var(--text-primary);
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%239CA3AF' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 0.5rem center;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  padding: 1rem;
}

.modal {
  background: var(--bg-elevated);
  border: 1px solid var(--border-color);
  border-radius: 1rem;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem;
  border-bottom: 1px solid var(--border-color);
}

.modal-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.modal-close {
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 1.25rem;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
  border-radius: 0.5rem;
}

.modal-close:hover {
  background: var(--bg-base);
  color: var(--text-primary);
}

.modal-body {
  padding: 1.25rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1.25rem;
  border-top: 1px solid var(--border-color);
}

/* 表单样式 */
.form-group {
  margin-bottom: 1rem;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.form-input {
  width: 100%;
  padding: 0.625rem 0.875rem;
  background: var(--bg-base);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  font-size: 0.9375rem;
  color: var(--text-primary);
  transition: all 0.2s;
  outline: none;
}

.form-input:focus {
  border-color: #FF9F0A;
  background: var(--bg-elevated);
  box-shadow: 0 0 0 3px rgba(255, 159, 10, 0.1);
}

.form-input::placeholder {
  color: var(--text-tertiary);
}

textarea.form-input {
  resize: vertical;
  min-height: 100px;
}

.form-select {
  width: 100%;
  padding: 0.625rem 2rem 0.625rem 0.875rem;
  background: var(--bg-base);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  font-size: 0.9375rem;
  color: var(--text-primary);
  cursor: pointer;
  appearance: none;
  outline: none;
}

.form-select:focus {
  border-color: #FF9F0A;
  background: var(--bg-elevated);
  box-shadow: 0 0 0 3px rgba(255, 159, 10, 0.1);
}

.input-wrapper {
  position: relative;
}

.select-arrow {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  font-size: 1.25rem;
  color: var(--text-tertiary);
  pointer-events: none;
}

/* 按钮样式 */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.625rem 1.25rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-primary {
  background: #FF9F0A;
  color: white;
  box-shadow: 0 4px 14px rgba(255, 159, 10, 0.3);
}

.btn-primary:hover:not(:disabled) {
  background: #E58909;
  box-shadow: 0 6px 20px rgba(255, 159, 10, 0.4);
  transform: translateY(-1px);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: var(--bg-base);
  color: var(--text-primary);
  border-color: var(--text-tertiary);
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-active .modal,
.fade-leave-active .modal {
  transition: all 0.3s cubic-bezier(0.19, 1, 0.22, 1);
}

.fade-enter-from .modal,
.fade-leave-to .modal {
  transform: scale(0.95);
  opacity: 0;
}

/* 响应式适配 */
@media (max-width: 1024px) {
  .dev-stats {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .notification-management {
    padding: 1rem;
    gap: 1rem;
  }
  
  .dev-header-content {
    flex-direction: column;
    align-items: stretch;
  }
  
  .dev-create-btn {
    width: 100%;
    justify-content: center;
  }
  
  .dev-stats {
    grid-template-columns: 1fr;
    gap: 0.75rem;
  }
  
  .notification-inner {
    flex-direction: column;
    gap: 1rem;
  }
  
  .notification-actions {
    flex-direction: row;
    opacity: 1;
    width: 100%;
    justify-content: flex-end;
  }
  
  .notification-header-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .pagination-inner {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .pagination-controls {
    justify-content: center;
  }
  
  .page-size {
    display: flex;
    justify-content: center;
  }
}
</style>

<!-- 弹窗样式 - 非 scoped，因为 Teleport 会将弹窗传送到 body 外部 -->
<style>
/* 模态框样式 - 必须非 scoped 才能应用到 Teleport 传送的元素 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  padding: 1rem;
}

.modal {
  background: var(--bg-elevated);
  border: 1px solid var(--border-color);
  border-radius: 1rem;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem;
  border-bottom: 1px solid var(--border-color);
}

.modal-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.modal-close {
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 1.25rem;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
  border-radius: 0.5rem;
}

.modal-close:hover {
  background: var(--bg-base);
  color: var(--text-primary);
}

.modal-body {
  padding: 1.25rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1.25rem;
  border-top: 1px solid var(--border-color);
}

/* 表单样式 */
.modal .form-group {
  margin-bottom: 1rem;
}

.modal .form-group:last-child {
  margin-bottom: 0;
}

.modal .form-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.modal .form-input {
  width: 100%;
  padding: 0.625rem 0.875rem;
  background: var(--bg-base);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  font-size: 0.9375rem;
  color: var(--text-primary);
  transition: all 0.2s;
  outline: none;
}

.modal .form-input:focus {
  border-color: #FF9F0A;
  background: var(--bg-elevated);
  box-shadow: 0 0 0 3px rgba(255, 159, 10, 0.1);
}

.modal .form-input::placeholder {
  color: var(--text-tertiary);
}

.modal textarea.form-input {
  resize: vertical;
  min-height: 100px;
}

.modal .form-select {
  width: 100%;
  padding: 0.625rem 2rem 0.625rem 0.875rem;
  background: var(--bg-base);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  font-size: 0.9375rem;
  color: var(--text-primary);
  cursor: pointer;
  appearance: none;
  outline: none;
}

.modal .form-select:focus {
  border-color: #FF9F0A;
  background: var(--bg-elevated);
  box-shadow: 0 0 0 3px rgba(255, 159, 10, 0.1);
}

.modal .input-wrapper {
  position: relative;
}

.modal .select-arrow {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  font-size: 1.25rem;
  color: var(--text-tertiary);
  pointer-events: none;
}

/* 按钮样式 */
.modal .btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.625rem 1.25rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.modal .btn-primary {
  background: #FF9F0A;
  color: white;
  box-shadow: 0 4px 14px rgba(255, 159, 10, 0.3);
}

.modal .btn-primary:hover:not(:disabled) {
  background: #E58909;
  box-shadow: 0 6px 20px rgba(255, 159, 10, 0.4);
  transform: translateY(-1px);
}

.modal .btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.modal .btn-secondary {
  background: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.modal .btn-secondary:hover {
  background: var(--bg-base);
  color: var(--text-primary);
  border-color: var(--text-tertiary);
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-active .modal,
.fade-leave-active .modal {
  transition: all 0.3s cubic-bezier(0.19, 1, 0.22, 1);
}

.fade-enter-from .modal,
.fade-leave-to .modal {
  transform: scale(0.95);
  opacity: 0;
}
</style>
