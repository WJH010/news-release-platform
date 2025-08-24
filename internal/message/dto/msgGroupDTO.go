package dto

type MsgGroupIDRequest struct {
	MsgGroupID int `uri:"msg_group_id" binding:"required,numeric"` // 消息组ID，必须为数字
}

type AddUserToGroupRequest struct {
	UserIDs []int `json:"user_ids" binding:"required,dive,numeric"` // 用户ID列表，必须为数字
}

type CreateMsgGroupRequest struct {
	GroupName      string `json:"group_name" binding:"required,max=255"`          // 群组名称，必填，最大长度255
	Desc           string `json:"desc" binding:"omitempty"`                       // 群组描述，选填
	IncludeAllUser string `json:"include_all_user" binding:"omitempty,oneof=Y N"` // 是否包含所有用户，选填，默认N
	UserIDs        []int  `json:"user_ids" binding:"omitempty,dive,numeric"`      // 初始用户ID列表，选填，必须为数字
}
