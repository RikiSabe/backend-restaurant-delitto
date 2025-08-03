package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type TypeInsumo struct {
	ID          uint   `json:"id"`
	Nombre      string `json:"nombre"`
	Cantidad    uint   `json:"cantidad"`
	IDProveedor uint   `json:"id_proveedor"`
}

type TypeModificadoInsumo struct {
	Nombre      string `json:"nombre"`
	Cantidad    uint   `json:"cantidad"`
	IDProveedor uint   `json:"id_proveedor"`
}

func ObtenerInsumos(w http.ResponseWriter, r *http.Request) {
	var insumos []TypeInsumo

	query := `select ins.id, ins.nombre, ins.cantidad, ins.id_proveedor from insumos as ins`

	err := db.GDB.Raw(query).Scan(&insumos).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(insumos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerInsumo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var insumo TypeInsumo

	query := `select ins.id, ins.nombre, ins.cantidad, ins.id_proveedor from insumos as ins where ins.id = ?`

	err := db.GDB.Raw(query, id).Scan(&insumo).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(insumo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AgregarInsumo(w http.ResponseWriter, r *http.Request) {
	var insumo models.Insumo

	if err := json.NewDecoder(r.Body).Decode(&insumo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&insumo).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Insumo", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(insumo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// func ModificarInsumo(w http.ResponseWriter, r *http.Request) {
// 	id_insumo := mux.Vars(r)["id"]
// 	var insumoExistente models.Insumo

// 	err := db.GDB.Where("id = ?", id_insumo).First(&insumoExistente).Error
// 	if err != nil {
// 		http.Error(w, "Insumo no encontrado", http.StatusNotFound)
// 		return
// 	}

// 	var insumoActualizado TypeModificadoInsumo
// 	if err := json.NewDecoder(r.Body).Decode(&insumoActualizado); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Cambios
// 	insumoExistente.Nombre = insumoActualizado.Nombre
// 	insumoExistente.Cantidad = insumoActualizado.Cantidad
// 	insumoExistente.IDProveedor = insumoActualizado.IDProveedor

// 	if err := db.GDB.Save(&insumoExistente).Error; err != nil {
// 		http.Error(w, "Error al actualizar Insumo", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	if err := json.NewEncoder(w).Encode(&insumoExistente); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }
