package orm

type User struct {
	ID       int     `gorm:"type:int;primaryKey"`
	Login    string  `gorm:"type:TEXT"`
	Password string  `gorm:"type:TEXT"`
	Trips    []*Trip `gorm:"many2many:trip_users"`
}
