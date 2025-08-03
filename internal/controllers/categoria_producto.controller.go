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

type CategoriaDAO struct {
	ID          uint   `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Estado      string `json:"estado"`
}

type CategoriaMOD struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Estado      string `json:"estado"`
}

func ObtenerCategorias(w http.ResponseWriter, r *http.Request) {
	var categorias []CategoriaDAO

	err := db.GDB.Raw(querys.CategoriasProductos).Scan(&categorias).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categorias)
}

func ObtenerCategoria(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var categoria CategoriaMOD

	err := db.GDB.Raw(querys.CategoriaProductos, id).Scan(&categoria).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoria)
}

func AgregarCategoria(w http.ResponseWriter, r *http.Request) {
	var categoria CategoriaMOD

	if err := json.NewDecoder(r.Body).Decode(&categoria); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nuevoEstado, err := functions.ActualizarEstado(categoria.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}

	nuevaCategoria := models.CategoriaProducto{
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

func ModificarCategoria(w http.ResponseWriter, r *http.Request) {
	id_categoria := mux.Vars(r)["id"]
	var categoriaExistente models.CategoriaProducto

	err := db.GDB.Where("id = ?", id_categoria).First(&categoriaExistente).Error
	if err != nil {
		http.Error(w, "Categoria no encontrada", http.StatusNotFound)
		return
	}

	var categoriaActualizado CategoriaMOD
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
