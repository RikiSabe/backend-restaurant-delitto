package main

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/routers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	var err error

	// Cargar el archivo .env
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error al conectarse a la base de datos: %v", err)
	}

	port := "5000"

	err = db.Connection()
	if err != nil {
		log.Printf("Error al conectar a la base de datos: %v", err)
		return
	}

	if err := db.GDB.AutoMigrate(
		/* migraciones */
		&models.Categoria{},
		&models.Proveedor{},
		&models.Mesa{},
		&models.Usuario{},

		&models.Producto{},
		&models.Insumo{},
		&models.Pedido{},

		&models.PedidosProductos{},
		&models.UsuariosPedidos{},
		&models.PedidosInsumos{},
	); err != nil {
		log.Fatal("Error al migrar los modelos de la db:", err)
	}

	var count int64
	if err := db.GDB.Model(&models.Usuario{}).Count(&count).Error; err != nil {
		log.Printf("Error al contar usuarios: %v", err)
	} else if count == 0 {
		usuario := models.Usuario{Nombre: "user", Contra: "admin"}
		if err := db.GDB.Create(&usuario).Error; err != nil {
			log.Printf("Error al crear usuario por defecto: %v", err)
		} else {
			log.Println("Usuario por defecto creado: user / admin")
		}
	} else {
		log.Printf("Ya hay usuarios registrados")
	}

	r := mux.NewRouter()
	routers.InitEndPoints(r)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	// Iniciar el servidor
	fmt.Printf("Servidor corriendo en puerto: %s\n", port)
	if err := http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(r)); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
