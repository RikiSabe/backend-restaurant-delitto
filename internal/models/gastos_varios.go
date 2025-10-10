package models

import "time"

type GastoVario struct {
	ID           uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre       string `gorm:"column:nombre;not null" json:"nombre"`
	UnidadMedida string `gorm:"column:unidad_medida" json:"unidad_medida"`

	IDProveedor uint      `gorm:"column:id_proveedor" json:"id_proveedor"`
	Proveedor   Proveedor `gorm:"foreignKey:IDProveedor;reference:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (GastoVario) TableName() string {
	return "gastos_varios"
}
