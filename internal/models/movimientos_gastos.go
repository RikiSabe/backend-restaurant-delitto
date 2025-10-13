package models

import "time"

type Movimiento struct {
	ID       uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Cantidad float64   `gorm:"column:cantidad;not null" json:"cantidad"`
	Precio   uint      `gorm:"column:precio;not null" json:"precio"`
	Fecha    time.Time `gorm:"column:fecha;not null" json:"fecha"`

	IDGasto uint       `gorm:"column:id_gasto; not null" json:"id_gasto"`
	Gasto   GastoVario `gorm:"foreignKey:IDGasto;reference:ID" json:"-"`
}

func (Movimiento) TableName() string {
	return "movimientos"
}
