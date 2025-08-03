package functions

import (
	"backend-restaurant-delitto/internal/db"
	"fmt"
)

func ActualizarEstado(estado string) (bool, error) {
	switch estado {
	case "Activo":
		return true, nil
	case "Inactivo":
		return false, nil
	default:
		return false, fmt.Errorf("error al actualizar estado")
	}
}

func ActualizarRol(rol string) (uint, error) {
	var id_rol uint
	err := db.GDB.Raw("SELECT id FROM roles WHERE nombre = ? LIMIT 1", rol).Scan(&id_rol).Error
	if err != nil {
		return 0, fmt.Errorf("error al actualizar rol: %w", err)
	}
	if id_rol == 0 {
		return 0, fmt.Errorf("rol no encontrado")
	}
	return id_rol, nil
}

func ActualizarCategoria(categoria string) (uint, error) {
	var id_categoria uint
	err := db.GDB.Raw("SELECT id FROM categoria_producto where nombre = ? LIMIT 1", categoria).Scan(&id_categoria).Error
	if err != nil {
		return 0, fmt.Errorf("error al actualizar categoria: %w", err)
	}
	if id_categoria == 0 {
		return 0, fmt.Errorf("categoria no encontrada")
	}
	return id_categoria, nil
}
