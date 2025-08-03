package models

import "time"

type Proveedor struct {
	ID        uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre    string `gorm:"column:nombre;not null" json:"nombre"`
	Telefono  string `gorm:"column:telefono" json:"telefono"`
	Correo    string `gorm:"column:correo" json:"correo"`
	Direccion string `gorm:"direccion" json:"direccion"`
	Estado    bool   `gorm:"column:estado" json:"estado"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Proveedor) TableName() string {
	return "proveedores"
}
