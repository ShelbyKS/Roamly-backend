package model

type User struct {
	ID   int  `gorm:"type:int;primaryKey"`
	Name string `gorm:"type:TEXT;unique"`
	Password string `gorm:"type:TEXT"`
}
