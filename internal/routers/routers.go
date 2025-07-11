package routers

import (
	c "backend-restaurant-delitto/internal/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func InitEndPoints(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()
	endPointsAPI(api)
}

func endPointsAPI(api *mux.Router) {
	v1 := api.PathPrefix("/v1").Subrouter()

	v1.HandleFunc("/login", c.Auth.AuthLoginWeb).Methods(http.MethodPost)

	v1Usuarios := v1.PathPrefix("/usuarios").Subrouter()

	v1Usuarios.HandleFunc("/{id}", c.ObtenerUsuario).Methods(http.MethodGet)
	v1Usuarios.HandleFunc("/{id}", c.ModificarUsuario).Methods(http.MethodPut)
	v1Usuarios.HandleFunc("", c.ObtenerUsuarios).Methods(http.MethodGet)
	v1Usuarios.HandleFunc("", c.AgregarUsuario).Methods(http.MethodPost)

	v1Producto := v1.PathPrefix("/productos").Subrouter()

	v1Producto.HandleFunc("/{id}", c.ObtenerProducto).Methods(http.MethodGet)
	v1Producto.HandleFunc("/{id}", c.ModificarProducto).Methods(http.MethodPut)
	v1Producto.HandleFunc("", c.ObtenerProductos).Methods(http.MethodGet)
	v1Producto.HandleFunc("", c.AgregarProducto).Methods(http.MethodPost)

	v1Proveedor := v1.PathPrefix("/proveedores").Subrouter()

	v1Proveedor.HandleFunc("/{id}", c.ObtenerProveedor).Methods(http.MethodGet)
	v1Proveedor.HandleFunc("/{id}", c.ModificarProveedor).Methods(http.MethodPut)
	v1Proveedor.HandleFunc("", c.ObtenerProveedores).Methods(http.MethodGet)
	v1Proveedor.HandleFunc("", c.AgregarProveedor).Methods(http.MethodPost)

	v1Insumo := v1.PathPrefix("/insumos").Subrouter()

	v1Insumo.HandleFunc("/{id}", c.ObtenerInsumo).Methods(http.MethodGet)
	v1Insumo.HandleFunc("/{id}", c.ModificarInsumo).Methods(http.MethodPut)
	v1Insumo.HandleFunc("", c.ObtenerInsumos).Methods(http.MethodGet)
	v1Insumo.HandleFunc("", c.AgregarInsumo).Methods(http.MethodPost)

	v1Categoria := v1.PathPrefix("/categorias").Subrouter()

	v1Categoria.HandleFunc("/{id}", c.ObtenerCategoria).Methods(http.MethodGet)
	v1Categoria.HandleFunc("/{id}", c.ModificarCategoria).Methods(http.MethodPut)
	v1Categoria.HandleFunc("", c.ObtenerCategorias).Methods(http.MethodGet)
	v1Categoria.HandleFunc("", c.AgregarCategoria).Methods(http.MethodPost)

	v1Mesa := v1.PathPrefix("/mesas").Subrouter()

	v1Mesa.HandleFunc("/{id}", c.ObtenerMesa).Methods(http.MethodGet)
	v1Mesa.HandleFunc("/{id}", c.ModificarMesa).Methods(http.MethodPut)
	v1Mesa.HandleFunc("", c.ObtenerMesas).Methods(http.MethodGet)
	v1Mesa.HandleFunc("", c.AgregarMesa).Methods(http.MethodPost)

	v1Pedido := v1.PathPrefix("/pedidos").Subrouter()

	v1Pedido.HandleFunc("", c.RegistrarPedido).Methods(http.MethodPost)
	v1Pedido.HandleFunc("", c.ObtenerPedidos).Methods(http.MethodGet)
}
