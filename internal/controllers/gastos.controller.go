package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/querys"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type GetGastos struct {
	ID           uint   `json:"id"`
	Nombre       string `json:"nombre"`
	UnidadMedida string `json:"unidad_medida"`

	IDProveedor     uint   `json:"id_proveedor"`
	NombreProveedor string `json:"nombre_proveedor"`
}

func ObtenerGastos(w http.ResponseWriter, r *http.Request) {
	var gastos []GetGastos

	err := db.GDB.Raw(querys.Gastos).Scan(&gastos).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gastos)
}

type GetGasto struct {
	ID           uint   `json:"id"`
	Nombre       string `json:"nombre"`
	UnidadMedida string `json:"unidad_medida"`

	IDProveedor uint `json:"id_proveedor"`
}

func ObtenerGasto(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var gasto GetGasto

	err := db.GDB.Raw(querys.Gasto, id).Scan(&gasto).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gasto)
}

func AgregarGasto(w http.ResponseWriter, r *http.Request) {
	var gasto models.GastoVario

	if err := json.NewDecoder(r.Body).Decode(&gasto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&gasto).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Gasto", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gasto)
}

func ModificarGasto(w http.ResponseWriter, r *http.Request) {
	id_gasto := mux.Vars(r)["id"]
	var gastoExistente models.GastoVario

	err := db.GDB.Where("id = ?", id_gasto).First(&gastoExistente).Error
	if err != nil {
		http.Error(w, "Gasto no encontrado", http.StatusNotFound)
		return
	}

	var gastoActualizado GetGasto
	if err := json.NewDecoder(r.Body).Decode(&gastoActualizado); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	gastoExistente.Nombre = gastoActualizado.Nombre
	gastoExistente.UnidadMedida = gastoActualizado.UnidadMedida
	gastoExistente.IDProveedor = gastoActualizado.IDProveedor

	if err := db.GDB.Save(&gastoExistente).Error; err != nil {
		http.Error(w, "Error al actualizar Gasto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&gastoExistente)
}
