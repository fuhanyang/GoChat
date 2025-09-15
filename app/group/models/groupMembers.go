package models

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> f07d8cb712af6cc475a8a13376802bdb57e3b5b5
type GroupMembers struct {
	ID        int      `json:"id"`
	GroupID   int      `json:"group_id"`
	MembersID []string `json:"members"`
<<<<<<< HEAD
=======
import "time"

type GroupMember struct {
	ID            int ` gorm:"primarykey"`
	GroupID       int ` gorm:"foreignKey:GroupID;references:GroupInfo.ID;onDelete:CASCADE"`
	AccountNum    string
	Purview       string
	GroupJoinTime time.Time
>>>>>>> 2f4b59e ("latest update")
=======
>>>>>>> f07d8cb712af6cc475a8a13376802bdb57e3b5b5
}
