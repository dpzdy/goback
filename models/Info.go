package models

type Info struct {
	//ID            int    `gorm:"primaryKey"`
	UserEmail     string `gorm:"type:varchar(255)"`
	UserName      string `gorm:"type:varchar(255)"`
	UserAuthority string `gorm:"type:varchar(255)"`
}
