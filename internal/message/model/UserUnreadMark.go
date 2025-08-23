package model

// UserUnreadMark 对应user_unread_marks表的数据模型，用于记录用户在各消息组中的最后已读消息ID
type UserUnreadMark struct {
	ID            int    `json:"id" gorm:"primaryKey;column:id"`
	UserID        int    `json:"user_id" gorm:"column:user_id"`
	MsgGroupID    int    `json:"msg_group_id" gorm:"column:msg_group_id"`
	LastReadMsgID int    `json:"last_read_msg_id" gorm:"column:last_read_msg_id"`
	IsDeleted     string `json:"is_deleted" gorm:"column:is_deleted;default:N"`
	CreateTime    string `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime    string `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
	CreateUser    int    `json:"create_user" gorm:"column:create_user"`
	UpdateUser    int    `json:"update_user" gorm:"column:update_user"`
}

// TableName 设置表名
func (*UserUnreadMark) TableName() string {
	return "user_unread_marks" // 表名指定为user_unread_marks
}
