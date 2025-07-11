package models

import "time"

type Proveedor struct {
	ID            uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NombreEmpresa string `gorm:"column:nombre_empresa;not null" json:"nombre_empresa"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Proveedor) TableName() string {
	return "proveedores"
}
