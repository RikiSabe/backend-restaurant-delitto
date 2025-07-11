package models

import "time"

type Categoria struct {
	ID     uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre string `gorm:"column:nombre; not null" json:"nombre"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Categoria) TableName() string {
	return "categorias"
}
