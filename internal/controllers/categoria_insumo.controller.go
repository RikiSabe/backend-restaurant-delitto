package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/functions"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/querys"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CategoriaInsumoDAO struct {
	ID          uint   `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Estado      string `json:"estado"`
}

type CategoriaInsumoMOD struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Estado      string `json:"estado"`
}

func ObtenerCategoriasInsumos(w http.ResponseWriter, r *http.Request) {
	var categorias []CategoriaInsumoDAO

	err := db.GDB.Raw(querys.CategoriasInsumos).Scan(&categorias).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categorias)
}

func ObtenerCategoriaInsumo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var categoria CategoriaInsumoMOD

	err := db.GDB.Raw(querys.CategoriaInsumos, id).Scan(&categoria).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoria)
}

func AgregarCategoriaInsumos(w http.ResponseWriter, r *http.Request) {
	var categoria CategoriaInsumoMOD

	if err := json.NewDecoder(r.Body).Decode(&categoria); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nuevoEstado, err := functions.ActualizarEstado(categoria.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}

	nuevaCategoria := models.CategoriaInsumo{
		Nombre:      categoria.Nombre,
		Descripcion: categoria.Descripcion,
		Estado:      nuevoEstado,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevaCategoria).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Categoria", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoria)
}

func ModificarCategoriaInsumo(w http.ResponseWriter, r *http.Request) {
	id_categoria := mux.Vars(r)["id"]
	var categoriaExistente models.CategoriaInsumo

	err := db.GDB.Where("id = ?", id_categoria).First(&categoriaExistente).Error
	if err != nil {
		http.Error(w, "Categoria no encontrada", http.StatusNotFound)
		return
	}

	var categoriaActualizado CategoriaInsumoMOD
	if err := json.NewDecoder(r.Body).Decode(&categoriaActualizado); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	categoriaExistente.Nombre = categoriaActualizado.Nombre
	categoriaExistente.Descripcion = categoriaActualizado.Descripcion
	nuevoEstado, err := functions.ActualizarEstado(categoriaActualizado.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}
	categoriaExistente.Estado = nuevoEstado

	if err := db.GDB.Save(&categoriaExistente).Error; err != nil {
		http.Error(w, "Error al actualizar Categoria", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&categoriaExistente)
}
