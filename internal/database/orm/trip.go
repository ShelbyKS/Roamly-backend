package orm

type Trip struct {
	ID        int     `gorm:"type:int;primaryKey"`
	StartTime string  `gorm:"type:TIME"`
	EndTime   string  `gorm:"type:TIME"`
	Region    string  `gorm:"type:TEXT"`
	Users     []*User `gorm:"many2many:trip_users"`
}
