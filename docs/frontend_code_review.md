# 前端代码深度审查报告

**审查日期**: 2026-01-08
**审查范围**: `frontend/src` 核心模块
**状态**: 已完成

本报告汇总了在深度审查中发现的前端待解决问题。

## 1. 安全性 (Security)

### 1.1 Token 存储机制 [中]

- **位置**: `frontend/src/api/index.ts` (interceptors), `frontend/src/stores/auth.ts`
- **问题**: 身份验证 Token 存储在 `localStorage` 中。这使得应用易受 **XSS (跨站脚本攻击)** 影响。如果攻击者能在页面注入脚本，即可轻松窃取 Token。
- **建议**:
  1.  **Web 环境**: 使用 `HttpOnly` Cookie 存储 Token。
  2.  **Wails 环境**: 可使用 Wails 的 Runtime Store 或更为安全的本地存储机制（如系统钥匙串，如果可用）。
  3.  **缓解措施**: 确保严格的 CSP (内容安全策略) 以防止 XSS。

### 1.2 敏感错误信息暴露 [低]

- **位置**: `api/index.ts`
- **问题**: 接口响应拦截器直接将后端返回的 `message` 字段展示给用户（或用于 Reject）。虽然方便调试，但可能暴露后端详细错误。
- **建议**: 生产环境中统一错误提示，或仅透出这一白名单内的错误码消息。

## 2. 逻辑与正确性 (Logic & Correctness)

### 2.1 状态管理不一致 [低]

- **位置**: `frontend/src/stores/auth.ts`
- **问题**: `isLoggedIn` 和 `isAuthenticated` 计算属性逻辑完全相同，存在冗余。且 Token 存在性检查较为简单 (`!!token.value`)。
- **建议**: 统一使用 `isAuthenticated`，并增加对 Token 格式/过期时间的初步本地校验（如 decode JWT 查看 exp）。

### 2.2 表单/输入校验缺失 [中]

- **位置**: `src/views/ProjectsView.vue`, `src/views/SettingsView.vue` (推测)
- **问题**: 主要依赖后端返回的 400 错误进行反馈。前端缺乏即时的输入校验（如邮箱格式、必填项、最大长度），不仅用户体验差，也增加了不必要的网络请求。
- **建议**: 引入 `vee-validate` 或 `zod` 进行前端表单验证。

## 3. 性能优化 (Performance)

### 3.1 长列表渲染 [低]

- **位置**: `ProjectsView.vue`
- **问题**: 虽然当前分页大小 (`pageSize`) 默认为 5，掩盖了性能问题。但如果用户选择每页 50 或 100 条，表格渲染可能会变慢。
- **建议**: 保持分页是合理的。如果未来需要无限滚动 (Infinite Scroll)，必须引入虚拟滚动 (Virtual Scroller)。

### 3.2 频繁的 localStorage 读写 [低]

- **位置**: `stores/auth.ts`
- **问题**: 每次更新用户信息都写入 `localStorage`。
- **建议**: 考虑使用 Pinia 的持久化插件 (`pinia-plugin-persistedstate`) 来统一管理，而非手动 `setItem`，代码更简洁且性能更好（有防抖机制）。

## 4. 代码质量与可维护性 (Code Quality)

### 4.1 硬编码文本与 I18n [低]

- **位置**: `src/views/ProjectsView.vue` (getStatusLabel) 等
- **问题**: "进行中", "已完成" 等中文文案硬编码在组件中。不利于多语言扩展。
- **建议**: 即使目前只需要中文，也建议提取到独立的 `locales/zh-CN.json` 或 `constants.ts` 中管理。

### 4.2 TypeScript 类型定义不完善 [低]

- **位置**: `src/views/ProjectsView.vue` (line 40)
- **问题**: 使用了 `// eslint-disable-next-line @typescript-eslint/no-explicit-any` 绕过类型检查。
- **建议**: 为查询参数定义明确的 `interface ProjectQueryParams`，避免使用 `any`。

### 4.3 样式重复 [低]

- **位置**: `AppSidebar.vue`, `ProjectsView.vue`
- **问题**: 依然存在部分硬编码的颜色值（如 `rgba(0, 0, 0, 0.05)`），而非全部使用 CSS 变量。这使得主题切换（深色模式）维护困难。
- **建议**: 将所有颜色提取到 CSS 变量（Token）中。

---

建议下一阶段重点关注 **Token 存储安全优化** 和 **引入前端表单验证框架**。
