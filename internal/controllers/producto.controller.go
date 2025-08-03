package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/functions"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/querys"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProductoDAO struct {
	ID          uint    `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio" format:"%.2f"`
	Imagen      string  `json:"imagen"`
	Estado      string  `json:"estado"`
	Categoria   string  `json:"categoria"`
}

type ProductoMOD struct {
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio" format:"%.2f"`
	Imagen      string  `json:"imagen"`
	Estado      string  `json:"estado"`
	Categoria   string  `json:"categoria"`
}

func encodeImageToBase64(path string) (string, error) {
	imgBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	mimeType := "image/jpeg"
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".jpg":
		mimeType = "image/jpg"
	case ".jpeg":
		mimeType = "image/jpeg"
	default:
		mimeType = "application/octet-stream"
	}

	encoded := base64.StdEncoding.EncodeToString(imgBytes)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

func ObtenerProductos(w http.ResponseWriter, r *http.Request) {
	var productos []ProductoDAO

	err := db.GDB.Raw(querys.Productos).Scan(&productos).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	for i, producto := range productos {
		if producto.Imagen != "N/A" {
			encodeImagen, err := encodeImageToBase64(producto.Imagen)
			if err == nil {
				productos[i].Imagen = encodeImagen
			} else {
				productos[i].Imagen = "N/A"
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)
}

func ObtenerProducto(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var producto ProductoMOD

	err := db.GDB.Raw(querys.Producto, id).Scan(&producto).Error
	if err != nil {
		http.Error(w, "No existe el producto", http.StatusInternalServerError)
		return
	}

	if producto.Imagen != "N/A" {
		encodeImagen, err := encodeImageToBase64(producto.Imagen)
		if err == nil {
			producto.Imagen = encodeImagen
		} else {
			producto.Imagen = "N/A"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(producto)
}

func AgregarProducto(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusInternalServerError)
		return
	}

	nuevoPrecio, err := strconv.ParseFloat(r.FormValue("precio"), 64)
	if err != nil {
		http.Error(w, "Precio no valido", http.StatusBadRequest)
		return
	}
	nuevoEstado, err := functions.ActualizarEstado(r.FormValue("estado"))
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}

	direccionImagen := "N/A"
	file, handler, err := r.FormFile("imagen")
	if err == nil {
		defer file.Close()

		nombreImagen := fmt.Sprintf("producto-%s%s", uuid.New().String(), filepath.Ext(handler.Filename))
		rutaImagen := "internal/images/productos/" + nombreImagen

		outFile, err := os.Create(rutaImagen)
		if err != nil {
			http.Error(w, "Error al guardar la foto", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Error al escribir la foto", http.StatusInternalServerError)
			return
		}

		direccionImagen = rutaImagen
	} else if err != http.ErrMissingFile {
		http.Error(w, "Error al obtener la foto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	nuevaCategoria, err := functions.ActualizarCategoria(r.FormValue("categoria"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nuevoProducto := models.Producto{
		Nombre:      r.FormValue("nombre"),
		Descripcion: r.FormValue("descripcion"),
		Precio:      nuevoPrecio,
		Imagen:      direccionImagen,
		Estado:      nuevoEstado,
		IDCategoria: nuevaCategoria,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevoProducto).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al guardar el producto: "+err.Error(), http.StatusInternalServerError)
		_ = os.Remove(direccionImagen)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nuevoProducto)
}

func ModificarProducto(w http.ResponseWriter, r *http.Request) {
	id_producto := mux.Vars(r)["id"]
	var productoExistente models.Producto

	err := db.GDB.Where("id = ?", id_producto).First(&productoExistente).Error
	if err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("imagen")
	if err == nil {
		defer file.Close()

		if productoExistente.Imagen != "" && productoExistente.Imagen != "N/A" {
			_ = os.Remove(productoExistente.Imagen) // ignorar error si no existe
		}

		nombreImagen := fmt.Sprintf("producto-%s%s", uuid.New().String(), filepath.Ext(handler.Filename))
		rutaImagen := "internal/images/productos/" + nombreImagen

		outFile, err := os.Create(rutaImagen)
		if err != nil {
			http.Error(w, "Error al guardar la nueva foto", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, file); err != nil {
			http.Error(w, "Error al escribir la nueva foto", http.StatusInternalServerError)
			return
		}

		productoExistente.Imagen = rutaImagen
	} else if err != http.ErrMissingFile {
		http.Error(w, "Error al procesar la foto: "+err.Error(), http.StatusBadRequest)
		return
	}

	productoExistente.Nombre = r.FormValue("nombre")
	productoExistente.Descripcion = r.FormValue("descripcion")
	nuevoPrecio, err := strconv.ParseFloat(r.FormValue("precio"), 64)
	if err != nil {
		http.Error(w, "Precio invalido", http.StatusBadRequest)
		return
	}
	productoExistente.Precio = nuevoPrecio
	nuevoEstado, err := functions.ActualizarEstado(r.FormValue("estado"))
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}
	productoExistente.Estado = nuevoEstado
	nuevaCategoria, err := functions.ActualizarCategoria(r.FormValue("categoria"))
	if err != nil {
		http.Error(w, "Categoria no valida", http.StatusBadRequest)
		return
	}
	productoExistente.IDCategoria = nuevaCategoria

	if err := db.GDB.Save(&productoExistente).Error; err != nil {
		http.Error(w, "Error al actualizar producto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productoExistente)
}
