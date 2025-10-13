package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/querys"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/datatypes"
)

type Productos struct {
	ID       uint    `json:"id"`
	Nombre   string  `json:"nombre"`
	Precio   float64 `json:"precio"`
	Cantidad uint    `json:"cantidad"`
}

type PedidoEntrante struct {
	Origen    string                          `json:"origen"`
	Comentario    string                          `json:"comentario"`
	Productos datatypes.JSONType[[]Productos] `json:"productos"`
	Estado    string                          `json:"estado"`
	IDMesa    uint                            `json:"id_mesa"`
}

func RegistrarPedido(w http.ResponseWriter, r *http.Request) {
	var pedidoEntrante PedidoEntrante

	if err := json.NewDecoder(r.Body).Decode(&pedidoEntrante); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pedido, err := json.MarshalIndent(pedidoEntrante, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(pedido))

	var nuevo_pedido = models.Pedido{
		Estado: pedidoEntrante.Estado,
		Fecha:  time.Now(),
		Origen: pedidoEntrante.Origen,
		Comentario: pedidoEntrante.Comentario,
		IDMesa: pedidoEntrante.IDMesa,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevo_pedido).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al registrar pedido", http.StatusInternalServerError)
		return
	}
	productosEntrada := pedidoEntrante.Productos.Data()
	for _, producto := range productosEntrada {
		if err := tx.Create(&models.DetallePedido{
			IDProducto: producto.ID,
			SubTotal:   producto.Precio * float64(producto.Cantidad),
			Cantidad:   producto.Cantidad,
			IDPedido:   nuevo_pedido.ID,
		}).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error en relacion Pedido - Producto", http.StatusInternalServerError)
			return
		}
	}
	var mesa models.Mesa
	if err := tx.Where("id = ?", pedidoEntrante.IDMesa).First(&mesa).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Mesa no encontrada", http.StatusNotFound)
		return
	}
	mesa.Estado = "Ocupado"
	if err := tx.Save(&mesa).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar estado de la mesa", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nuevo_pedido)
}

type PedidoJSON struct {
	ID        uint            `json:"id"`
	Fecha     string          `json:"fecha"`
	Estado    string          `json:"estado"`
	Origen    string          `json:"origen"`
	Productos json.RawMessage `json:"productos"`
}

func ObtenerPedidos(w http.ResponseWriter, r *http.Request) {
	var pedidos []PedidoJSON

	err := db.GDB.Raw(querys.PedidosJSONB).Scan(&pedidos).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidos)
}
