package main

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/functions"
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
		&models.Rol{},
		&models.CategoriaProducto{},
		&models.CategoriaInsumo{},
		&models.Proveedor{},
		&models.Mesa{},
		&models.Usuario{},
		&models.Producto{},
		&models.Insumo{},
		&models.Pedido{},
		&models.UsuariosPedidos{},
		&models.DetallePedido{},
		&models.Factura{},

		&models.GastoVario{},
		&models.Movimiento{},
	); err != nil {
		log.Fatal("Error al migrar los modelos de la db:", err)
	}

	if err := functions.CreacionInicial(); err != nil {
		log.Fatal("Error al iniciar los datos predeterminados")
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
