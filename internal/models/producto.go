package models

import (
	"time"
)

type Producto struct {
	ID          uint    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre      string  `gorm:"column:nombre;not null" json:"nombre"`
	Precio      float64 `gorm:"not null" json:"precio" format:"%.2f"`
	IDCategoria uint    `gorm:"column:id_categoria;not null" json:"id_categoria"`

	Categoria Categoria `gorm:"foreignKey:IDCategoria;references:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Producto) TableName() string {
	return "productos"
}
