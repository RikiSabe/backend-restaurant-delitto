package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type TypeCategoria struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
}

type TypeModificadoCategoria struct {
	Nombre string `json:"nombre"`
}

func ObtenerCategorias(w http.ResponseWriter, r *http.Request) {
	var categorias []TypeCategoria

	query := `select ct.id, ct.nombre from categorias as ct`

	err := db.GDB.Raw(query).Scan(&categorias).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(categorias); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerCategoria(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var categoria TypeCategoria

	query := `select ct.id, ct.nombre from categorias as ct where ct.id = ?`

	err := db.GDB.Raw(query, id).Scan(&categoria).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(categoria); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AgregarCategoria(w http.ResponseWriter, r *http.Request) {
	var categoria models.Categoria

	if err := json.NewDecoder(r.Body).Decode(&categoria); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&categoria).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Categoria", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(categoria); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ModificarCategoria(w http.ResponseWriter, r *http.Request) {
	id_categoria := mux.Vars(r)["id"]
	var categoriaExistente models.Categoria

	err := db.GDB.Where("id = ?", id_categoria).First(&categoriaExistente).Error
	if err != nil {
		http.Error(w, "Categoria no encontrada", http.StatusNotFound)
		return
	}

	var categoriaActualizado TypeModificadoCategoria
	if err := json.NewDecoder(r.Body).Decode(&categoriaActualizado); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	categoriaExistente.Nombre = categoriaActualizado.Nombre

	if err := db.GDB.Save(&categoriaExistente).Error; err != nil {
		http.Error(w, "Error al actualizar Categoria", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&categoriaExistente); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
