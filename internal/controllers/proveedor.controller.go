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

type ProveedorDAO struct {
	ID        uint   `json:"id"`
	Nombre    string `json:"nombre"`
	Telefono  string `json:"telefono"`
	Correo    string `json:"correo"`
	Direccion string `json:"direccion"`
	Estado    string `json:"estado"`
}

type ProovedorModificado struct {
	Nombre    string `json:"nombre"`
	Telefono  string `json:"telefono"`
	Correo    string `json:"correo"`
	Direccion string `json:"direccion"`
	Estado    string `json:"estado"`
}

func ObtenerProveedores(w http.ResponseWriter, r *http.Request) {

	var proveedores []ProveedorDAO

	err := db.GDB.Raw(querys.Proovedores).Scan(&proveedores).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proveedores)
}

func ObtenerProveedor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var proveedor ProveedorDAO

	err := db.GDB.Raw(querys.Proovedor, id).First(&proveedor).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proveedor)
}

func AgregarProveedor(w http.ResponseWriter, r *http.Request) {
	var proveedor models.Proveedor

	if err := json.NewDecoder(r.Body).Decode(&proveedor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&proveedor).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Proveedor", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&proveedor)
}

func ModificarProveedor(w http.ResponseWriter, r *http.Request) {
	id_proveedor := mux.Vars(r)["id"]
	var proveedorExistente models.Proveedor

	err := db.GDB.Where("id = ?", id_proveedor).First(&proveedorExistente).Error
	if err != nil {
		http.Error(w, "Proveedor no encontrado", http.StatusNotFound)
		return
	}

	var proveedorActualizado ProovedorModificado
	if err := json.NewDecoder(r.Body).Decode(&proveedorActualizado); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	proveedorExistente.Nombre = proveedorActualizado.Nombre
	proveedorExistente.Telefono = proveedorActualizado.Telefono
	proveedorExistente.Correo = proveedorActualizado.Correo
	proveedorExistente.Direccion = proveedorActualizado.Direccion
	nuevoEstado, err := functions.ActualizarEstado(proveedorActualizado.Estado)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	proveedorExistente.Estado = nuevoEstado

	if err := db.GDB.Save(&proveedorExistente).Error; err != nil {
		http.Error(w, "Error al actualizar usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&proveedorExistente)
}
