// Package utils 用于定义全局常量
package utils

const (
	DeletedFlagYes    = "Y"       // 软删除标志，表示已删除
	DeletedFlagNo     = "N"       // 软删除标志，表示未删除
	TypeEvent         = "EVENT"   // 活动类型常量
	TypeArticle       = "ARTICLE" // 新闻类型常量
	TypeGroup         = "GROUP"   // 群组类型常量
	TypeSystem        = "SYSTEM"  // 系统消息类型常量
	QueryScopeAll     = "ALL"     // 查询范围常量，表示查询全部
	QueryScopeDeleted = "DELETED" // 查询范围常量，表示查询
	// 角色常量
	RoleUser       = 1 // 普通用户角色
	RoleAdmin      = 2 // 管理员角色
	RoleSuperAdmin = 3 // 超级管理员角色
)

var QueryScopeList = []string{
	QueryScopeAll,
	QueryScopeDeleted,
}

var UserGroupMessageTypeList = []string{
	TypeGroup,
	TypeSystem,
}
