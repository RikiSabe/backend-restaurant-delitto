package models

import (
	"time"
)

type Producto struct {
	ID          uint    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre      string  `gorm:"column:nombre;not null" json:"nombre"`
	Descripcion string  `gorm:"column:descripcion" json:"descripcion"`
	Precio      float64 `gorm:"not null" json:"precio" format:"%.2f"`
	Imagen      string  `gorm:"column:imagen" json:"imagen"`
	Estado      bool    `gorm:"column:estado" json:"estado"`
	IDCategoria uint    `gorm:"column:id_categoria" json:"id_categoria"`

	Categoria CategoriaProducto `gorm:"foreignKey:IDCategoria;reference:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Producto) TableName() string {
	return "productos"
}
