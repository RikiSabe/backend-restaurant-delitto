package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/querys"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type TypeInsumo struct {
	ID              uint   `json:"id"`
	Nombre          string `json:"nombre"`
	StockActual     uint   `json:"stock_actual"`
	StockMinimo     uint   `json:"stock_minimo"`
	UnidadMedida    string `json:"unidad_medida"`
	IDProveedor     uint   `json:"id_proveedor"`
	NombreProveedor string `json:"nombre_proveedor"`
	IDCategoria     uint   `json:"id_categoria"`
	NombreCategoria string `json:"nombre_categoria"`
}

type TypeModificadoInsumo struct {
	Nombre       string `json:"nombre"`
	StockActual  uint   `json:"stock_actual"`
	StockMinimo  uint   `json:"stock_minimo"`
	UnidadMedida string `json:"unidad_medida"`
	IDProveedor  uint   `json:"id_proveedor"`
	IDCategoria  uint   `json:"id_categoria"`
}

func ObtenerInsumos(w http.ResponseWriter, r *http.Request) {
	var insumos []TypeInsumo

	err := db.GDB.Raw(querys.Insumos).Scan(&insumos).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insumos)
}

func ObtenerInsumo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var insumo TypeInsumo

	err := db.GDB.Raw(querys.Insumo, id).Scan(&insumo).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insumo)
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
	json.NewEncoder(w).Encode(insumo)
}

func ModificarInsumo(w http.ResponseWriter, r *http.Request) {
	id_insumo := mux.Vars(r)["id"]
	var insumoExistente models.Insumo

	err := db.GDB.Where("id = ?", id_insumo).First(&insumoExistente).Error
	if err != nil {
		http.Error(w, "Insumo no encontrado", http.StatusNotFound)
		return
	}

	var insumoActualizado TypeModificadoInsumo
	if err := json.NewDecoder(r.Body).Decode(&insumoActualizado); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	insumoExistente.Nombre = insumoActualizado.Nombre
	insumoExistente.StockActual = insumoActualizado.StockActual
	insumoExistente.StockMinimo = insumoActualizado.StockMinimo
	insumoExistente.UnidadMedida = insumoActualizado.UnidadMedida
	insumoExistente.IDProveedor = insumoActualizado.IDProveedor
	insumoExistente.IDCategoria = insumoActualizado.IDCategoria

	if err := db.GDB.Save(&insumoExistente).Error; err != nil {
		http.Error(w, "Error al actualizar Insumo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&insumoExistente)
}
