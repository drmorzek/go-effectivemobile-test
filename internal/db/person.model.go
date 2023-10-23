package db

type Person struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"type:varchar(100);not null"`
	Surname     string `json:"surname" gorm:"type:varchar(100);not null"`
	Patronymic  string `json:"patronymic,omitempty" gorm:"type:varchar(100);not null"`
	Age         int    `gorm:"type:int" json:"age"`
	Gender      string `gorm:"type:varchar(10)" json:"gender"`
	Nationality string `gorm:"type:varchar(100)" json:"nationality"`
}
