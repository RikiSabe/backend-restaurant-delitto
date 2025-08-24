package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/querys"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Mesas struct {
	ID        uint   `json:"id"`
	Nombre    string `json:"nombre"`
	Estado    string `json:"estado"`
	Capacidad uint   `json:"capacidad"`
}

type Mesa struct {
	Nombre    string `json:"nombre"`
	Estado    string `json:"estado"`
	Capacidad uint   `json:"capacidad"`
}

func ObtenerMesas(w http.ResponseWriter, r *http.Request) {
	var mesas []Mesas

	err := db.GDB.Raw(querys.Mesas).Scan(&mesas).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mesas)
}

func ObtenerMesa(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var mesa Mesa

	err := db.GDB.Raw(querys.Mesa, id).Scan(&mesa).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mesa)
}

func AgregarMesa(w http.ResponseWriter, r *http.Request) {
	var mesa models.Mesa

	if err := json.NewDecoder(r.Body).Decode(&mesa); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&mesa).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Mesa", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mesa)
}

func ModificarMesa(w http.ResponseWriter, r *http.Request) {
	id_mesa := mux.Vars(r)["id"]
	var mesaExistente models.Mesa

	err := db.GDB.Where("id = ?", id_mesa).First(&mesaExistente).Error
	if err != nil {
		http.Error(w, "Mesa no encontrada", http.StatusNotFound)
		return
	}

	var mesaActualizada Mesa
	if err := json.NewDecoder(r.Body).Decode(&mesaActualizada); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	mesaExistente.Nombre = mesaActualizada.Nombre
	mesaExistente.Estado = mesaActualizada.Estado
	mesaExistente.Capacidad = mesaActualizada.Capacidad

	if err := db.GDB.Save(&mesaExistente).Error; err != nil {
		http.Error(w, "Error al actualizar Mesa", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&mesaExistente)
}

func LiberarMesa(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var mesa models.Mesa

	err := db.GDB.Where("id = ?", id).First(&mesa).Error
	if err != nil {
		http.Error(w, "Mesa no encontrada", http.StatusNotFound)
		return
	}
	mesa.Estado = "Disponible"

	if err := db.GDB.Save(&mesa).Error; err != nil {
		http.Error(w, "Error al guardar estado de la mesa", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&mesa)
}
