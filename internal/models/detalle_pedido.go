package models

type DetallePedido struct {
	ID         uint    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SubTotal   float64 `gorm:"column:subtotal" json:"subtotal" format:"%.2f"`
	Cantidad   uint    `gorm:"column:cantidad" json:"cantidad"`
	IDPedido   uint    `gorm:"id_pedido;not null" json:"id_pedido"`
	IDProducto uint    `gorm:"id_producto;not null" json:"id_producto"`

	Pedido   Pedido   `gorm:"foreignKey:IDPedido;reference:ID" json:"-"`
	Producto Producto `gorm:"foreignKey:IDProducto;reference:ID" json:"-"`
}

func (DetallePedido) TableName() string {
	return "detalle_pedido"
}
