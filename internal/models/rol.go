package models

import "time"

type Rol struct {
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre   string `gorm:"column:nombre" json:"nombre"`
	Permisos string `gorm:"column:permisos" json:"permisos"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
}

func (Rol) TableName() string {
	return "roles"
}
