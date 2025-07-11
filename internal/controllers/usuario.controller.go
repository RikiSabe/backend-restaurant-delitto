package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type TypeUsuario struct {
	ID     uint   `json:"id"`
	Rol    string `json:"rol"`
	Nombre string `json:"nombre"`
}

type TypeModificadoUsuario struct {
	Rol    string `json:"rol"`
	Nombre string `json:"nombre"`
}

func ObtenerUsuarios(w http.ResponseWriter, r *http.Request) {

	var usuarios []TypeUsuario

	query := `select u.id, u.rol, u.nombre from usuarios as u;`

	err := db.GDB.Raw(query).Scan(&usuarios).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(usuarios); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var usuario TypeUsuario

	query := `select u.id, u.rol, u.nombre from usuarios as u where u.id = ?`

	err := db.GDB.Raw(query, id).Scan(&usuario).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(usuario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AgregarUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario models.Usuario

	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	usuario.Contra = usuario.Rol

	tx := db.GDB.Begin()
	if err := tx.Create(&usuario).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Usuario", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ModificarUsuario(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]
	var usuarioExistente models.Usuario

	err := db.GDB.Where("id = ?", id_usuario).First(&usuarioExistente).Error
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	var usuarioActualizado TypeModificadoUsuario
	if err := json.NewDecoder(r.Body).Decode(&usuarioActualizado); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	usuarioExistente.Rol = usuarioActualizado.Rol
	usuarioExistente.Nombre = usuarioActualizado.Nombre

	if err := db.GDB.Save(&usuarioExistente).Error; err != nil {
		http.Error(w, "Error al actualizar usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&usuarioExistente); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
