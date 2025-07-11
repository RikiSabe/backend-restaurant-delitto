package models

import "time"

type Mesa struct {
	ID        uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Estado    string `gorm:"column:estado;not null" json:"estado"`
	Capacidad uint   `gorm:"column:capacidad;not null" json:"capacidad"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Mesa) TableName() string {
	return "mesas"
}
