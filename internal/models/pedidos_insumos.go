package models

type PedidosInsumos struct {
	IDPedido uint `gorm:"column:id_pedido;not null" json:"id_pedido"`
	IDInsumo uint `gorm:"column:id_insumo;not null" json:"id_insumo"`

	Pedido Pedido `gorm:"foreignKey:IDPedido;references:ID" json:"-"`
	Insumo Insumo `gorm:"foreignKey:IDInsumo;references:ID" json:"-"`
}

func (PedidosInsumos) TableName() string {
	return "pedidos_insumos"
}
