package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type TypeMesa struct {
	ID        uint   `json:"id"`
	Estado    string `json:"estado"`
	Capacidad uint   `json:"capacidad"`
}

type TypeModificadoMesa struct {
	Estado    string `json:"estado"`
	Capacidad uint   `json:"capacidad"`
}

func ObtenerMesas(w http.ResponseWriter, r *http.Request) {
	var mesas []TypeMesa

	query := `select me.id, me.estado, me.capacidad from mesas as me`

	err := db.GDB.Raw(query).Scan(&mesas).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mesas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerMesa(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var mesa TypeMesa

	query := `select me.id, me.estado, me.capacidad from mesas as me where me.id = ?`

	err := db.GDB.Raw(query, id).Scan(&mesa).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mesa); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mesa); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ModificarMesa(w http.ResponseWriter, r *http.Request) {
	id_mesa := mux.Vars(r)["id"]
	var mesaExistente models.Mesa

	err := db.GDB.Where("id = ?", id_mesa).First(&mesaExistente).Error
	if err != nil {
		http.Error(w, "Mesa no encontrada", http.StatusNotFound)
		return
	}

	var mesaActualizada TypeModificadoMesa
	if err := json.NewDecoder(r.Body).Decode(&mesaActualizada); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	mesaExistente.Estado = mesaActualizada.Estado
	mesaExistente.Capacidad = mesaActualizada.Capacidad

	if err := db.GDB.Save(&mesaExistente).Error; err != nil {
		http.Error(w, "Error al actualizar Mesa", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&mesaExistente); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
