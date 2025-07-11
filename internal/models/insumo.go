package models

import "time"

type Insumo struct {
	ID          uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre      string `gorm:"column:nombre;not null" json:"nombre"`
	Cantidad    uint   `gorm:"column:cantidad;not null" json:"cantidad"`
	IDProveedor uint   `gorm:"column:id_proveedor;not null" json:"id_proveedor"`

	Proveedor Proveedor `gorm:"foreignKey:IDProveedor;reference:ID" json:"-"`
	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Insumo) TableName() string {
	return "insumos"
}
