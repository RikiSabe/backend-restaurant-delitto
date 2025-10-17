package models

import "time"

type Proveedor struct {
	ID              uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Tipo            string `gorm:"column:tipo;not null;default:'persona'" json:"tipo"` // persona | empresa
	Nombre          string `gorm:"column:nombre;not null" json:"nombre"`
	Telefono        string `gorm:"column:telefono" json:"telefono"`
	Correo          string `gorm:"column:correo" json:"correo"`
	Direccion       string `gorm:"column:direccion" json:"direccion"`
	CI              string `gorm:"column:ci" json:"ci"`
	NIT             string `gorm:"column:nit" json:"nit"`
	NombreEncargado string `gorm:"column:nombre_encargado" json:"nombre_encargado"`
	Estado          bool   `gorm:"column:estado" json:"estado"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Proveedor) TableName() string {
	return "proveedores"
}
