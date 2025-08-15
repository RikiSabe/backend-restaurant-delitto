package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"encoding/json"
	"net/http"
)

type auth struct{}

var Auth auth

func (auth) AuthLoginWeb(w http.ResponseWriter, r *http.Request) {
	var usuario struct {
		Usuario string `json:"usuario"`
		Contra  string `json:"contra"`
	}
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var usuarioExistente models.Usuario
	if err := db.GDB.Where("usuario = ?", usuario.Usuario).First(&usuarioExistente).Error; err != nil {
		http.Error(w, "Credenciales incorrectas", http.StatusInternalServerError)
		return
	}

	if usuarioExistente.Contra != usuario.Contra {
		http.Error(w, "Credenciales incorrectas", http.StatusInternalServerError)
		return
	}

	type RespuestaUsuario struct {
		ID  uint   `json:"id"`
		Rol string `json:"rol"`
	}

	var rol models.Rol
	if err := db.GDB.Table("roles").Where("id = ?", usuarioExistente.IDRol).Scan(&rol).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respuesta := RespuestaUsuario{
		ID:  usuarioExistente.ID,
		Rol: rol.Nombre,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&respuesta)
}
