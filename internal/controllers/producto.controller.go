package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type TypeProducto struct {
	ID        uint    `json:"id"`
	Nombre    string  `json:"nombre"`
	Precio    float64 `json:"precio"`
	Categoria string  `json:"categoria"`
}

type TypeModificadoProducto struct {
	Nombre      string  `json:"nombre"`
	Precio      float64 `json:"precio"`
	IDCategoria uint    `json:"id_categoria"`
}

func ObtenerProductos(w http.ResponseWriter, r *http.Request) {
	var productos []TypeProducto

	query := `
		select 
			p.id, p.nombre, p.precio, ct.nombre as categoria 
		from productos as p 
		join categorias as ct on ct.id = p.id_categoria`

	err := db.GDB.Raw(query).Scan(&productos).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(productos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerProducto(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var producto TypeProducto

	query := `
		select p.id, p.nombre, p.precio, ct.nombre as categoria 
		from productos as p 
		join categorias as ct on ct.id = p.id_categoria 
		where p.id = ?`

	err := db.GDB.Raw(query, id).Scan(&producto).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(producto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AgregarProducto(w http.ResponseWriter, r *http.Request) {
	var producto models.Producto

	if err := json.NewDecoder(r.Body).Decode(&producto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&producto).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Producto", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(producto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ModificarProducto(w http.ResponseWriter, r *http.Request) {
	id_producto := mux.Vars(r)["id"]
	var productoExistente models.Producto

	err := db.GDB.Where("id = ?", id_producto).First(&productoExistente).Error
	if err != nil {
		http.Error(w, "Producot no encontrado", http.StatusNotFound)
		return
	}

	var productoActualizado TypeModificadoProducto
	if err := json.NewDecoder(r.Body).Decode(&productoActualizado); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	productoExistente.Nombre = productoActualizado.Nombre
	productoExistente.Precio = productoActualizado.Precio
	productoExistente.IDCategoria = productoActualizado.IDCategoria

	if err := db.GDB.Save(&productoExistente).Error; err != nil {
		http.Error(w, "Error al actualizar producto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&productoExistente); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
