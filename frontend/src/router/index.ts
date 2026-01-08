import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/login',
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { title: '登录', requiresAuth: false },
    },
    {
      path: '/',
      component: AppLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue'),
          meta: { title: '工作台', requiresAuth: true },
        },
        {
          path: 'projects',
          name: 'projects',
          component: () => import('@/views/ProjectsView.vue'),
          meta: { title: '项目管理', requiresAuth: true },
        },
        {
          path: 'projects/create',
          name: 'project-create',
          component: () => import('@/views/ProjectCreateView.vue'),
          meta: { title: '新建项目', requiresAuth: true },
        },
        {
          path: 'projects/edit/:id',
          name: 'project-edit',
          component: () => import('@/views/ProjectCreateView.vue'),
          meta: { title: '编辑项目', requiresAuth: true },
        },
        {
          path: 'projects/:id',
          name: 'project-detail',
          component: () => import('@/views/ProjectDetailView.vue'),
          meta: { title: '项目详情', requiresAuth: true },
        },
        {
          path: 'projects/:id/payment/create',
          name: 'payment-create',
          component: () => import('@/views/PaymentCreateView.vue'),
          meta: { title: '添加收款', requiresAuth: true },
        },
        {
          path: 'payment/create',
          name: 'payment-create-global',
          component: () => import('@/views/PaymentCreateView.vue'),
          meta: { title: '添加收款', requiresAuth: true },
        },
        {
          path: 'calendar',
          name: 'calendar',
          component: () => import('@/views/CalendarView.vue'),
          meta: { title: '收款日历', requiresAuth: true },
        },
        {
          path: 'analytics',
          name: 'analytics',
          component: () => import('@/views/AnalyticsView.vue'),
          meta: { title: '数据分析', requiresAuth: true },
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/SettingsView.vue'),
          meta: { title: '系统设置', requiresAuth: true },
        },
      ],
    },
  ],
})

// Navigation Guard for Auth and Title
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  // 设置页面标题
  const title = to.meta.title as string
  if (title) {
    document.title = `${title} - Orange`
  } else {
    document.title = 'Orange - 项目收款管理系统'
  }

  // 权限检查
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    // 未登录，重定向到登录页
    next({ name: 'login' })
  } else if (to.name === 'login' && authStore.isLoggedIn) {
    // 已登录，访问登录页时重定向到工作台
    next({ name: 'dashboard' })
  } else {
    next()
  }
})

export default router
