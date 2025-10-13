package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/querys"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type GetMovimientos struct {
	NombreGasto  string    `json:"nombre_gasto"`
	Cantidad     float64   `json:"cantidad"`
	Precio       uint      `json:"precio"`
	Fecha        time.Time `json:"fecha"`
	UnidadMedida string    `json:"unidad_medida"`
}

func ObtenerMovimientos(w http.ResponseWriter, r *http.Request) {
	id_gasto := mux.Vars(r)["id"]
	var movimientos []GetMovimientos

	err := db.GDB.Raw(querys.Movimientos, id_gasto).Scan(&movimientos).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movimientos)
}

type NewMovimiento struct {
	Cantidad float64 `json:"cantidad"`
	Precio   uint    `json:"precio"`
}

func AgregarMovimiento(w http.ResponseWriter, r *http.Request) {
	id_gasto := mux.Vars(r)["id"]
	var movimiento NewMovimiento

	if err := json.NewDecoder(r.Body).Decode(&movimiento); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var idGastoUint uint
	if _, err := fmt.Sscanf(id_gasto, "%d", &idGastoUint); err != nil {
		http.Error(w, "ID de gasto inválido", http.StatusBadRequest)
		return
	}

	data := models.Movimiento{
		Cantidad: movimiento.Cantidad,
		Precio:   movimiento.Precio,
		Fecha:    time.Now(),
		IDGasto:  idGastoUint,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Movimiento", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
