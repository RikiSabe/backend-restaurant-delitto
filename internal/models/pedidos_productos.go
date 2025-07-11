package models

type PedidosProductos struct {
	IDPedido   uint `gorm:"id_pedido;not null" json:"id_pedido"`
	IDProducto uint `gorm:"id_producto;not null" json:"id_producto"`

	Pedido   Pedido   `gorm:"foreignKey:IDPedido;reference:ID" json:"-"`
	Producto Producto `gorm:"foreignKey:IDProducto;reference:ID" json:"-"`
}

func (PedidosProductos) TableName() string {
	return "pedidos_productos"
}
