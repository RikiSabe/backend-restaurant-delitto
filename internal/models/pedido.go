package models

import (
	"time"
)

type Pedido struct {
	ID     uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Estado string    `gorm:"column:estado;size:100;not null" json:"estado"`
	Fecha  time.Time `gorm:"column:fecha;not null" json:"fecha"`
	IDMesa uint      `gorm:"column:id_mesa" json:"id_mesa,omitempty"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Pedido) TableName() string {
	return "pedidos"
}
