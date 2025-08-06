package utils

// 错误码定义（按业务模块分类）
const (
	// 通用错误
	ErrCodeSuccess         = 0     // 成功
	ErrCodeParamInvalid    = 10001 // 参数验证失败
	ErrCodeParamBind       = 10002 // 参数绑定失败
	ErrCodeParamTypeError  = 10003 // 参数类型错误
	ErrCodeParamOutOfRange = 10004 // 参数值超出合法范围
	ErrCodeDataFormatError = 10005 // 数据格式错误（如 JSON/XML 格式解析失败）

	// 用户/权限相关
	ErrCodePermissionDenied = 20001 // 权限不足（无访问该资源的权限）
	ErrCodeAccessForbidden  = 20002 // 禁止访问（如 IP 被拉黑、资源被限制）
	ErrCodeRoleInvalid      = 20003 // 角色无效（关联的角色已删除或未启用）
	ErrCodeUserNotFound     = 20004 // 用户不存在

	// 资源相关
	ErrCodeResourceNotFound      = 30001 // 资源不存在（如查询的用户 ID / 订单 ID 不存在）
	ErrCodeResourceExists        = 30002 // 资源已存在（如创建时唯一键冲突，如 “用户名已注册”）
	ErrCodeResourceExpired       = 30003 // 资源已过期（如临时链接失效、活动已结束）
	ErrCodeResourceLocked        = 30004 // 资源被锁定（如数据正被其他操作占用，无法修改）
	ErrCodeResourceNotAllowed    = 30005 // 资源操作不允许（如删除系统保护资源）
	ErrCodeResourceQuotaExceeded = 30006 // 资源配额超限（如存储容量已满、创建数量达上限）

	// 认证相关
	ErrCodeAuthFailed       = 40001 // 认证失败（如账号密码错误、Token 无效）
	ErrCodeAuthTokenExpired = 40002 // 认证令牌过期（如 JWT 过期）
	ErrCodeAuthTokenInvalid = 40003 // 令牌格式无效（如 Token 被篡改、格式错误）
	ErrCodeAuthRequired     = 40004 // 需要先认证（未登录时访问需登录的资源）
	ErrCodeGetUserIDFailed  = 40005 // 验证码错误（如登录时验证码不匹配）

	// 服务器/系统相关
	ErrCodeServerInternalError = 50001 // 服务器内部错误（如代码异常、未捕获的异常）
	ErrCodeServerOverload      = 50002 // 服务器过载（如请求量超过承载上限）
	ErrCodeServiceUnavailable  = 50003 // 服务暂不可用（如维护中）

	// 限流/频率控制
	ErrCodeRateLimitExceeded = 60001 // 请求频率超限（触发限流策略）
	ErrCodeIpLimitExceeded   = 60002 // IP 请求次数超限

	// 业务逻辑相关
	ErrCodeBusinessLogicError = 70001 // 业务逻辑错误（如订单状态不允许操作）
	ErrCodeDataConflict       = 70002 // 数据冲突（如更新时数据已被修改）

	// 依赖/外部服务相关
	ErrCodeDependencyServiceError = 80001 // 依赖服务调用失败（如调用第三方 API 超时）
	ErrCodeDatabaseError          = 80002 // 数据库操作失败（如查询、写入异常）
)
