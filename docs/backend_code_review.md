# 后端代码深度审查报告

**审查日期**: 2026-01-08
**审查范围**: `internal/` 所有模块及 `main.go`
**状态**: 已完成

本报告汇总了在深度审查中发现的所有待解决问题。请按照优先级规划后续修复工作。

## 1. 安全性 (Security)

### 1.1 CORS 策略过于宽泛 [中]

- **位置**: `internal/router/router.go` (lines 130)
- **问题**: `Access-Control-Allow-Origin: *` 允许任意来源的跨域请求。在生产环境中（尤其是结合浏览器使用时），这可能导致 CSRF 或信息泄露风险。虽然 Wails 本地应用通常较为安全，但建议收紧。
- **建议**: 将其限制为特定协议或域名，如 `wails://` (Wails 前端协议) 和 `http://localhost` (开发环境)。

### 1.2 缺乏密码复杂度校验 [低]

- **位置**: `internal/service/auth.go` (`Register` 方法)
- **问题**: 注册时仅检查了空值，未对密码长度、字符组合（大小写、数字、符号）进行强制要求。
- **建议**: 引入密码强度校验库或正则表达式，强制密码最小长度（如 8 位）。

### 1.3 注册接口开放性 [低]

- **位置**: `internal/service/auth.go` (`Register` 方法)
- **问题**: 任何人都可以调用注册接口创建一个 `user` 角色的账户。对于企业内部工具，通常应禁止自行注册，改由管理员创建或邀请。
- **建议**: 根据业务需求，考虑关闭公开注册接口，或添加审批流程。

## 2. 代码质量与架构 (Code Quality & Architecture)

### 2.1 全局数据库实例 (Global State) [低]

- **位置**: `internal/database/database.go`, `internal/repository/*.go`
- **问题**: Repository 层直接调用 `database.GetDB()` 获取单例。这导致代码与具体数据库实现紧密耦合，极难编写单元测试（无法 Mock 数据库连接）。
- **建议**: 采用依赖注入（Dependency Injection）。在 `main.go` 中初始化 `*gorm.DB`，并层层传递给 Repository -> Service -> Handler。

### 2.2 紧密耦合的依赖初始化 [中]

- **位置**: 所有 Service 和 Handler 的 `New` 方法
- **问题**: 例如 `NewAuthService` 内部直接调用 `repository.NewUserRepository()`。这同样阻碍了单元测试（无法注入 Mock Repository）。
- **建议**: 修改构造函数签名，接受接口类型的依赖。例如 `func NewAuthService(userRepo repository.UserRepositoryInterface) *AuthService`。

### 2.3 错误处理不够精细 [低]

- **位置**: 各个 Handler 的 CRUD 操作
- **问题**: 许多 Handler 直接捕获 Service 层返回的 `error` 并返回 `500 Internal Error`。实际上，很多错误可能是 `RecordNotFound` (应返 404) 或由于输入导致的数据库约束冲突 (应返 400)。
- **建议**: 在 Service 层定义统一的业务错误类型（如 `ErrNotFound`, `ErrDuplicate`），Handler 层根据错误类型返回正确的 HTTP 状态码。

### 2.4 硬编码字符串 [低]

- **位置**: `internal/models/models.go`, `internal/service/auth.go` 等
- **问题**: 角色名称（"user", "admin"）、状态值（"active", "pending"）以字符串字面量散落在代码中。
- **建议**: 建立 `internal/consts` 或在 Model 中定义常量（如 `RoleAdmin`, `RoleUser`, `StatusActive`）以统一管理。

## 3. 数据库规范 (Database Naming)

### 3.1 表名命名不一致 [低]

- **位置**: `internal/models/models.go`
- **问题**: `DictionaryItem` 模型的表名被指定为 `dictionary_item` (单数)，而 User/Project/Notification 等均为复数 (`users`, `projects`)。
- **建议**: 将 `DictionaryItem` 的表名修改为 `dictionary_items` 以保持风格一致。

## 4. 之前已修复的重大问题 (已解决)

- [高] **支付确认逻辑漏洞**: 已增加事务与幂等性检查。
- [高] **项目删除逻辑漏洞**: 已增加级联删除款项的逻辑。
- [中] **统计性能优化**: 已将内存计算改为数据库聚合查询。
- [中] **字典管理权限**: 已增加 Admin 角色校验。

---

建议在下一阶段优先处理 **CORS 策略** 和 **表名一致性** 问题，架构重构可视长期维护计划决定。
