package models

import (
	"time"
)

type Usuario struct {
	ID     uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Rol    string `gorm:"column:rol;not null" json:"rol"`
	Nombre string `gorm:"column:nombre;size:255;not null" json:"nombre"`
	Contra string `gorm:"colum:contra;size:255" json:"contra"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Usuario) TableName() string {
	return "usuarios"
}
