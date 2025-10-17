package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/functions"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/querys"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type UsuarioDAO struct {
	ID       uint   `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	CI       string `json:"ci"`
	Celular  string `json:"celular"`
	Usuario  string `json:"usuario"`
	Contra   string `json:"contra"`
	Estado   string `json:"estado"`
	Rol      string `json:"rol"`
}

type UsuarioModificado struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	CI       string `json:"ci"`
	Celular  string `json:"celular"`
	Usuario  string `json:"usuario"`
	Contra   string `json:"contra"`
	Estado   string `json:"estado"`
	Rol      string `json:"rol"`
}

// validatePassword checks if the password meets the criteria
func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	if !regexp.MustCompile(`[a-zA-Z]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false
	}
	return true
}

func ObtenerUsuarios(w http.ResponseWriter, r *http.Request) {

	var usuarios []UsuarioDAO

	err := db.GDB.Raw(querys.Usuarios).Scan(&usuarios).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

func ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var usuario UsuarioDAO

	err := db.GDB.Raw(querys.Usuario, id).Scan(&usuario).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

func AgregarUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario UsuarioDAO

	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate password
	if !validatePassword(usuario.Contra) {
		http.Error(w, "La contraseña debe tener al menos 8 caracteres, incluyendo letras y números.", http.StatusBadRequest)
		return
	}

	nuevoEstado, err := functions.ActualizarEstado(usuario.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}
	nuevoRol, err := functions.ActualizarRol(usuario.Rol)
	if err != nil {
		http.Error(w, "Rol no valido:"+err.Error(), http.StatusBadRequest)
		return
	}

	nuevoUsuario := models.Usuario{
		Nombre:   usuario.Nombre,
		Apellido: usuario.Apellido,
		CI:       usuario.CI,
		Celular:  usuario.Celular,
		Usuario:  usuario.Usuario,
		Contra:   usuario.Contra,
		Estado:   nuevoEstado,
		IDRol:    nuevoRol,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevoUsuario).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Usuario", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&nuevoUsuario)
}

func ModificarUsuario(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]
	var usuarioExistente models.Usuario

	err := db.GDB.Where("id = ?", id_usuario).First(&usuarioExistente).Error
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	var usuarioActualizado UsuarioModificado
	if err := json.NewDecoder(r.Body).Decode(&usuarioActualizado); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate password if provided
	if usuarioActualizado.Contra != "" {
		if !validatePassword(usuarioActualizado.Contra) {
			http.Error(w, "La contraseña debe tener al menos 8 caracteres, incluyendo letras y números.", http.StatusBadRequest)
			return
		}
		usuarioExistente.Contra = usuarioActualizado.Contra
	}

	// Cambios
	usuarioExistente.Nombre = usuarioActualizado.Nombre
	usuarioExistente.Apellido = usuarioActualizado.Apellido
	usuarioExistente.CI = usuarioActualizado.CI
	usuarioExistente.Celular = usuarioActualizado.Celular
	usuarioExistente.Usuario = usuarioActualizado.Usuario
	nuevoEstado, err := functions.ActualizarEstado(usuarioActualizado.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}
	usuarioExistente.Estado = nuevoEstado
	nuevoRol, err := functions.ActualizarRol(usuarioActualizado.Rol)
	if err != nil {
		http.Error(w, "Rol no valido", http.StatusBadRequest)
		return
	}
	usuarioExistente.IDRol = nuevoRol

	if err := db.GDB.Save(&usuarioExistente).Error; err != nil {
		http.Error(w, "Error al actualizar usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&usuarioExistente)
}
