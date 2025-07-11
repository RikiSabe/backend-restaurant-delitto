package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type TypeProveedor struct {
	ID            uint   `json:"id"`
	NombreEmpresa string `json:"nombre_empresa"`
}

type TypeModificadoProveedor struct {
	NombreEmpresa string `json:"nombre_empresa"`
}

func ObtenerProveedores(w http.ResponseWriter, r *http.Request) {

	var proveedores []TypeProveedor

	query := `select pro.id, pro.nombre_empresa from proveedores as pro`

	err := db.GDB.Raw(query).Scan(&proveedores).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(proveedores); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerProveedor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var proveedor TypeProveedor

	query := `select pro.id, pro.nombre_empresa from proveedores as pro where pro.id = ?`

	err := db.GDB.Raw(query, id).Scan(&proveedor).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(proveedor); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&proveedor); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ModificarProveedor(w http.ResponseWriter, r *http.Request) {
	id_proveedor := mux.Vars(r)["id"]
	var proveedorExistente models.Proveedor

	err := db.GDB.Where("id = ?", id_proveedor).First(&proveedorExistente).Error
	if err != nil {
		http.Error(w, "Proveedor no encontrado", http.StatusNotFound)
		return
	}

	var proveedorActualizado TypeModificadoProveedor
	if err := json.NewDecoder(r.Body).Decode(&proveedorActualizado); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	proveedorExistente.NombreEmpresa = proveedorActualizado.NombreEmpresa

	if err := db.GDB.Save(&proveedorExistente).Error; err != nil {
		http.Error(w, "Error al actualizar usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&proveedorExistente); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
