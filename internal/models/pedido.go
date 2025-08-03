package models

import (
	"time"
)

type Pedido struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Fecha     time.Time `gorm:"column:fecha;not null" json:"fecha"`
	Estado    string    `gorm:"column:estado;size:100;not null" json:"estado"`
	Origen    string    `gorm:"column:origen" json:"origen"`
	IDMesa    uint      `gorm:"column:id_mesa" json:"id_mesa,omitempty"`
	IDUsuario uint      `gorm:"column:id_usuario" json:"id_usuario,omitempty"`

	Mesa    Mesa    `gorm:"foreignkey:IDMesa;reference:ID" json:"-"`
	Usuario Usuario `gorm:"foreignKey:IDUsuario;reference:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Pedido) TableName() string {
	return "pedidos"
}
