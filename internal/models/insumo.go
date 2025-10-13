package models

import "time"

type Insumo struct {
	ID          uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre      string `gorm:"column:nombre;not null" json:"nombre"`
	StockActual uint   `gorm:"column:stock_actual;not null" json:"stock_actual"`
	StockMinimo uint   `gorm:"column:stock_minimo" json:"stock_minimo"`
	// Unidad de Medida ya no es necesaria
	UnidadMedida string `gorm:"column:unidad_medida" json:"unidad_medida"`

	IDProveedor uint `gorm:"column:id_proveedor" json:"id_proveedor"`
	IDCategoria uint `gorm:"column:id_categoria" json:"id_categoria"`

	Categoria CategoriaInsumo `gorm:"foreignKey:IDCategoria;reference:ID" json:"-"`
	Proveedor Proveedor       `gorm:"foreignKey:IDProveedor;reference:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Insumo) TableName() string {
	return "insumos"
}
