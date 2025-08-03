package models

import "time"

type Factura struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	FechaEmision  time.Time `gorm:"column:fecha_emision" json:"fecha_emision"`
	NitCliente    string    `gorm:"column:nit_cliente" json:"nit_cliente"`
	Total         float64   `gorm:"column:total" json:"total"`
	CodigoControl string    `gorm:"column:codigo_control" json:"codigo_control"`
	IDPedido      uint      `gorm:"column:id_pedido" json:"id_pedido,omitempty"`

	Pedido Pedido `gorm:"foreignKey:IDPedido;reference:ID" json:"-"`
}

func (Factura) TableName() string {
	return "facturas"
}
