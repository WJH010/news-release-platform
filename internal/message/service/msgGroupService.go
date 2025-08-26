package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"news-release/internal/message/dto"
	"news-release/internal/message/model"
	"news-release/internal/message/repository"
	"news-release/internal/utils"
)

type MsgGroupService interface {
	// AddUserToGroup 用户入群
	AddUserToGroup(ctx context.Context, msgGroupID int, userIDs []int, operateUser int) error
	// CreateMsgGroup 创建消息群组
	CreateMsgGroup(ctx context.Context, msgGroup *model.UserMessageGroup, userIDs []int) error
	// DeleteUserFromGroup 用户退群
	DeleteUserFromGroup(ctx context.Context, msgGroupID int, userIDs []int, operateUser int) error
	// UpdateMsgGroup 更新消息群组
	UpdateMsgGroup(ctx context.Context, msgGroupID int, request dto.UpdateMsgGroupRequest, userID int) error
	// DeleteMsgGroup 删除消息群组
	DeleteMsgGroup(ctx context.Context, msgGroupID int, userID int) error
	// ListMsgGroups 列表查询消息群组
	ListMsgGroups(ctx context.Context, page int, pageSize int, groupName string, eventID int, queryScope string) ([]dto.ListMsgGroupResponse, int64, error)
	// ListGroupsUsers 获取指定群组内用户
	ListGroupsUsers(ctx context.Context, page int, pageSize int, msgGroupID int) ([]dto.ListGroupsUsersResponse, int64, error)
	// ListNotInGroupUsers 查询不在指定组内的用户
	ListNotInGroupUsers(ctx context.Context, page int, pageSize int, msgGroupID int, req dto.ListNotInGroupUsersRequest) ([]dto.ListGroupsUsersResponse, int64, error)
}

// MsgGroupServiceImpl 实现接口的具体结构体，持有数据访问层接口 Repository 的实例
type MsgGroupServiceImpl struct {
	msgGroupRepo repository.MsgGroupRepository
	msgRepo      repository.MessageRepository
}

// NewMsgGroupService 创建服务实例
func NewMsgGroupService(msgGroupRepo repository.MsgGroupRepository, msgRepo repository.MessageRepository) MsgGroupService {
	return &MsgGroupServiceImpl{msgGroupRepo: msgGroupRepo, msgRepo: msgRepo}
}

// 工具函数：分类用户（需要新增/需要恢复/无需操作）
func classifyUsers(
	userIDs []int,
	existingMappings map[int]model.UserMsgGroupMapping, // key: userID, value: 已有记录
	groupID, latestMsgID, operateUser int,
) (needCreate []model.UserMsgGroupMapping, needRecover []model.UserMsgGroupMapping) {
	for _, userID := range userIDs {
		if mapping, exists := existingMappings[userID]; exists {
			// 记录存在：若is_deleted=Y则需要恢复，否则无需操作
			if mapping.IsDeleted == "Y" {
				needRecover = append(needRecover, mapping)
			}
		} else {
			// 记录不存在：需要新增
			needCreate = append(needCreate, model.UserMsgGroupMapping{
				MsgGroupID:    groupID,
				UserID:        userID,
				JoinMsgID:     latestMsgID, // 入群时的最新消息ID（用于过滤历史消息）
				LastReadMsgID: latestMsgID,
				CreateUser:    operateUser,
				UpdateUser:    operateUser,
				IsDeleted:     "N",
			})
		}
	}
	return
}

// 工具函数：从需要恢复的记录中提取用户ID
func extractUserIDs(mappings []model.UserMsgGroupMapping) []int {
	var ids []int
	for _, m := range mappings {
		ids = append(ids, m.UserID)
	}
	return ids
}

// AddUserToGroup 用户入群
func (svc *MsgGroupServiceImpl) AddUserToGroup(ctx context.Context, msgGroupID int, userIDs []int, operateUser int) error {
	// 检查群组是否存在
	group, err := svc.msgGroupRepo.GetMsgGroupByID(ctx, msgGroupID)
	if err != nil {
		return err
	}
	if group == nil {
		return utils.NewBusinessError(utils.ErrCodeResourceNotFound, "数据异常，消息群组不存在")
	}
	// 获取当前群组最新消息ID
	latestMsgID, err := svc.msgRepo.GetLatestMsgIDInGroup(ctx, msgGroupID)
	if err != nil {
		return err
	}

	// 查询已存在的用户-群组关联记录（过滤无效用户ID）
	existingMappings, err := svc.msgGroupRepo.GetExistingMappings(ctx, msgGroupID, userIDs)
	if err != nil {
		return fmt.Errorf("查询已有记录失败: %w", err)
	}

	// 分类处理：区分需要新增、需要恢复（is_deleted=Y）、无需操作（is_deleted=N）的用户
	needCreate, needRecover := classifyUsers(userIDs, existingMappings, msgGroupID, latestMsgID, operateUser)

	// 使用 GORM 函数式事务执行
	err = svc.msgGroupRepo.ExecTransaction(ctx, func(tx *gorm.DB) error {
		// 事务内的业务操作：所有数据库操作必须使用 tx 作为 DB 实例
		if len(needCreate) > 0 {
			// 批量新增用户-群组关联记录
			if err := svc.msgGroupRepo.CreateUserMsgGroupMappings(ctx, tx, needCreate); err != nil {
				return err // 返回错误，GORM 自动回滚
			}
		}

		if len(needRecover) > 0 {
			// 批量恢复用户-群组关联记录
			userIDsToRecover := extractUserIDs(needRecover)
			if err := svc.msgGroupRepo.RecoverUserMsgGroupMappings(ctx, tx, msgGroupID, userIDsToRecover, latestMsgID, operateUser); err != nil {
				return err // 返回错误，GORM 自动回滚
			}
		}

		return nil // 返回 nil，GORM 自动提交
	})

	// 处理事务执行结果
	if err != nil {
		return utils.NewSystemError(fmt.Errorf("事务执行失败: %w", err))
	}

	return nil
}

// CreateMsgGroup 创建消息群组
// 没有进行事务控制，允许群组创建成功但用户添加失败
func (svc *MsgGroupServiceImpl) CreateMsgGroup(ctx context.Context, msgGroup *model.UserMessageGroup, userIDs []int) error {
	// 创建消息群组
	err := svc.msgGroupRepo.CreateMsgGroup(ctx, msgGroup)
	if err != nil {
		return err
	}
	// 如果群组创建成功且包含用户，则添加用户到群组
	if msgGroup.IncludeAllUser == "N" && len(userIDs) > 0 {
		err = svc.AddUserToGroup(ctx, msgGroup.ID, userIDs, msgGroup.CreateUser)
		if err != nil {
			logrus.Errorf("添加用户到群组失败" + err.Error())
			return utils.NewBusinessError(utils.ErrCodeServerInternalError, "添加用户到群组失败，请手动添加")
		}
	}
	return nil
}

// DeleteUserFromGroup 用户退群
func (svc *MsgGroupServiceImpl) DeleteUserFromGroup(ctx context.Context, msgGroupID int, userIDs []int, operateUser int) error {
	// 检查群组是否存在
	group, err := svc.msgGroupRepo.GetMsgGroupByID(ctx, msgGroupID)
	if err != nil {
		return err
	}
	if group == nil {
		return utils.NewBusinessError(utils.ErrCodeResourceNotFound, "数据异常，消息群组不存在")
	}
	// 删除用户-群组关联记录（软删除）
	err = svc.msgGroupRepo.DeleteUserMsgGroupMappings(ctx, msgGroupID, userIDs, operateUser)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMsgGroup 更新消息群组
func (svc *MsgGroupServiceImpl) UpdateMsgGroup(ctx context.Context, msgGroupID int, request dto.UpdateMsgGroupRequest, userID int) error {
	// 检查群组是否存在
	group, err := svc.msgGroupRepo.GetMsgGroupByID(ctx, msgGroupID)
	if err != nil {
		return err
	}
	if group == nil {
		return utils.NewBusinessError(utils.ErrCodeResourceNotFound, "数据异常，消息群组不存在")
	}
	// 构建更新字段
	updateField := make(map[string]interface{})
	if request.GroupName != nil {
		updateField["group_name"] = request.GroupName
	}
	if request.Desc != nil {
		updateField["desc"] = request.Desc
	}

	updateField["update_user"] = userID

	// 更新消息群组信息
	err = svc.msgGroupRepo.UpdateMsgGroup(ctx, msgGroupID, updateField)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMsgGroup 删除消息群组
func (svc *MsgGroupServiceImpl) DeleteMsgGroup(ctx context.Context, msgGroupID int, userID int) error {
	// 检查群组是否存在
	group, err := svc.msgGroupRepo.GetMsgGroupByID(ctx, msgGroupID)
	if err != nil {
		return err
	}
	if group == nil {
		return utils.NewBusinessError(utils.ErrCodeResourceNotFound, "数据异常，消息群组不存在")
	}

	// 软删除消息群组，复用 UpdateMsgGroup 方法
	updateField := map[string]interface{}{
		"is_deleted":  "Y",
		"update_user": userID,
	}
	if err = svc.msgGroupRepo.UpdateMsgGroup(ctx, msgGroupID, updateField); err != nil {
		return err
	}

	return nil
}

// ListMsgGroups 列表查询消息群组
func (svc *MsgGroupServiceImpl) ListMsgGroups(ctx context.Context, page int, pageSize int, groupName string, eventID int, queryScope string) ([]dto.ListMsgGroupResponse, int64, error) {
	return svc.msgGroupRepo.ListMsgGroups(ctx, page, pageSize, groupName, eventID, queryScope)
}

// ListGroupsUsers 获取指定群组内用户
func (svc *MsgGroupServiceImpl) ListGroupsUsers(ctx context.Context, page int, pageSize int, msgGroupID int) ([]dto.ListGroupsUsersResponse, int64, error) {
	return svc.msgGroupRepo.ListGroupsUsers(ctx, page, pageSize, msgGroupID)
}

// ListNotInGroupUsers 查询不在指定组内的用户
func (svc *MsgGroupServiceImpl) ListNotInGroupUsers(ctx context.Context, page int, pageSize int, msgGroupID int, req dto.ListNotInGroupUsersRequest) ([]dto.ListGroupsUsersResponse, int64, error) {
	return svc.msgGroupRepo.ListNotInGroupUsers(ctx, page, pageSize, msgGroupID, req)
}
