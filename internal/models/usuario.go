package models

import (
	"time"
)

type Usuario struct {
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre   string `gorm:"column:nombre;size:255" json:"nombre"`
	Apellido string `gorm:"column:apellido;size:255" json:"apellido"`
	CI       string `gorm:"column:ci" json:"ci"`
	Usuario  string `gorm:"column:usuario;not null" json:"usuario"`
	Contra   string `gorm:"column:contra;size:255" json:"contra"`
	Estado   bool   `gorm:"column:estado" json:"estado"`
	IDRol    uint   `gorm:"column:id_rol" json:"id_rol"`

	Rol Rol `gorm:"foreignKey:IDRol;reference:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Usuario) TableName() string {
	return "usuarios"
}
