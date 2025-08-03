package functions

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
)

func CreacionInicial() error {
	if err := Roles(); err != nil {
		return err
	}
	if err := PrimerUsuario(); err != nil {
		return err
	}
	return nil
}

func Roles() error {
	var count int64
	if err := db.GDB.Model(&models.Rol{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		roles := []models.Rol{
			{Nombre: "admin", Permisos: "todos los permisos"},
			{Nombre: "cajero", Permisos: "ventas"},
			{Nombre: "cliente", Permisos: "productos"},
		}

		if err := db.GDB.Create(&roles).Error; err != nil {
			return err
		}
	}

	return nil
}

func PrimerUsuario() error {

	var count int64
	if err := db.GDB.Model(&models.Usuario{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		var id_admin uint
		query := `select id from roles where nombre = 'admin' limit 1`
		if err := db.GDB.Raw(query).Scan(&id_admin).Error; err != nil {
			return err
		}
		usuario := models.Usuario{
			Usuario: "user",
			Contra:  "admin",
			IDRol:   id_admin,
			Estado:  true,
		}
		if err := db.GDB.Create(&usuario).Error; err != nil {
			return err
		}
	}
	return nil
}
