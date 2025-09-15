package models
<<<<<<< HEAD
=======

type GroupSettings struct {
	GroupID       int ` gorm:"foreignKey:GroupID;references:GroupInfo.ID;onDelete:CASCADE"`
	MaxNumOfUsers int
	GroupType     string
}
>>>>>>> 2f4b59e ("latest update")
