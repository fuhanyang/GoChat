package models
<<<<<<< HEAD
<<<<<<< HEAD
=======

type GroupSettings struct {
	GroupID       int ` gorm:"foreignKey:GroupID;references:GroupInfo.ID;onDelete:CASCADE"`
	MaxNumOfUsers int
	GroupType     string
}
>>>>>>> 2f4b59e ("latest update")
=======
>>>>>>> f07d8cb712af6cc475a8a13376802bdb57e3b5b5
