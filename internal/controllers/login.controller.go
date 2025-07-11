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
		Nombre string `json:"nombre"`
		Contra string `json:"contra"`
	}
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var usuarioExistente models.Usuario
	if err := db.GDB.Where("nombre = ?", usuario.Nombre).First(&usuarioExistente).Error; err != nil {
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

	respuesta := RespuestaUsuario{
		ID:  usuarioExistente.ID,
		Rol: usuarioExistente.Rol,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&respuesta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
