package models

<<<<<<< HEAD
import "gorm.io/gorm"

type GroupInfo struct {
	GroupID int `json:"group_id" gorm:"primary_key"`
	gorm.Model
}
=======
import (
	"common/snowflake"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type GroupInfo struct {
	ID          int ` gorm:"primarykey"`
	UID         string
	GroupName   string
	GroupDesc   string
	GroupLeader string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	gorm.Model
}

// CreateGroup 创建群组并添加成员，使用嵌套事务优化粒度
// 返回值：群组UID、错误信息
func CreateGroup(db *gorm.DB, groupName string, groupDesc string, groupLeader string, groupMembers []string, groupConfig GroupSettings) (string, error) {
	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return "", tx.Error
	}
	defer func() {
		// 发生panic时回滚事务
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 创建群组信息
	groupInfo := GroupInfo{
		GroupName:   groupName,
		GroupDesc:   groupDesc,
		GroupLeader: groupLeader,
		CreatedAt:   time.Now(),
		UID:         strconv.FormatInt(snowflake.GetID(), 10), //生成唯一的账号号, // 自动生成ID
	}
	if err := tx.Create(&groupInfo).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	groupID := groupInfo.ID // 获取自动生成的群组ID

	// 2. 创建群组配置
	if err := tx.Create(&groupConfig).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	// 3. 添加群主（单独处理，确保群主一定能加入）
	leaderMember := GroupMember{
		GroupID:    groupID,
		AccountNum: groupLeader,
		Purview:    "leader", // 群主权限
	}
	if err := tx.Create(&leaderMember).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	LastSavePoint := "create_members"
	err := tx.SavePoint(LastSavePoint).Error
	if err != nil {
		tx.Rollback()
		return "", err
	}

	faultMember := make([]string, 0)

	// 4. 添加其他成员
	for _, memberAccount := range groupMembers {
		// 跳过群主（避免重复添加）
		if memberAccount == groupLeader {
			continue
		}

		member := GroupMember{
			GroupID:    groupID,
			AccountNum: memberAccount,
		}

		// 尝试添加成员
		if err = tx.Create(&member).Error; err != nil {
			// 成员添加失败，回滚到上一个保存点
			faultMember = append(faultMember, memberAccount)
			if err = tx.RollbackTo(LastSavePoint).Error; err != nil {
				tx.Rollback()
				return "", err
			}
		}
		LastSavePoint = memberAccount
		err = tx.SavePoint(LastSavePoint).Error
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}
	if len(faultMember) > 0 {
		//TODO: 记录没有成功拉群的成员日志，但不回滚事务，可以异步重试
	}

	// 5. 所有操作成功，提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return "", err
	}

	return groupInfo.UID, nil
}
>>>>>>> 2f4b59e ("latest update")
