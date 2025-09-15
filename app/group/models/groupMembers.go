package models

type GroupMembers struct {
	ID        int      `json:"id"`
	GroupID   int      `json:"group_id"`
	MembersID []string `json:"members"`
}
