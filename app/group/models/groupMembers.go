package models

<<<<<<< HEAD
type GroupMembers struct {
	ID        int      `json:"id"`
	GroupID   int      `json:"group_id"`
	MembersID []string `json:"members"`
=======
import "time"

type GroupMember struct {
	ID            int ` gorm:"primarykey"`
	GroupID       int ` gorm:"foreignKey:GroupID;references:GroupInfo.ID;onDelete:CASCADE"`
	AccountNum    string
	Purview       string
	GroupJoinTime time.Time
>>>>>>> 2f4b59e ("latest update")
}
