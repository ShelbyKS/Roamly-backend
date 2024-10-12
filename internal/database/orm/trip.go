package orm

type Trip struct {
	ID        int    `gorm:"type:int;primaryKey"`
	StartTime string `gorm:"type:TIME"`
	EndTime   string `gorm:"type:TIME"`
	AreaID    string
	Users     []*User  `gorm:"many2many:trip_user"`
	Places    []*Place `gorm:"many2many:trip_place"`
}
