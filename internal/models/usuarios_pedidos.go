package models

type UsuariosPedidos struct {
	IDUsuario uint `gorm:"column:id_usuario;not null" json:"id_usuario"`
	IDPedido  uint `gorm:"column:id_pedido;not null" json:"id_pedido"`

	Usuario Usuario `gorm:"foreignKey:IDUsuario;references:ID" json:"-"`
	Pedido  Pedido  `gorm:"foreignKey:IDPedido;references:ID" json:"-"`
}

func (UsuariosPedidos) TableName() string {
	return "usuario_pedidos"
}
