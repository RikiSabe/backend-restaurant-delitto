package routers

import (
	c "backend-restaurant-delitto/internal/controllers"
	r "backend-restaurant-delitto/internal/reports"
	"net/http"

	"github.com/gorilla/mux"
)

func InitEndPoints(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()
	endPointsAPI(api)
	reports(api)
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

	v1CategoriaProducto := v1.PathPrefix("/categorias-productos").Subrouter()
	v1CategoriaProducto.HandleFunc("/{id}", c.ObtenerCategoria).Methods(http.MethodGet)
	v1CategoriaProducto.HandleFunc("/{id}", c.ModificarCategoria).Methods(http.MethodPut)
	v1CategoriaProducto.HandleFunc("", c.ObtenerCategorias).Methods(http.MethodGet)
	v1CategoriaProducto.HandleFunc("", c.AgregarCategoria).Methods(http.MethodPost)

	v1CategoriaInsumos := v1.PathPrefix("/categorias-insumos").Subrouter()
	v1CategoriaInsumos.HandleFunc("/{id}", c.ObtenerCategoriaInsumo).Methods(http.MethodGet)
	v1CategoriaInsumos.HandleFunc("/{id}", c.ModificarCategoriaInsumo).Methods(http.MethodPut)
	v1CategoriaInsumos.HandleFunc("", c.ObtenerCategoriasInsumos).Methods(http.MethodGet)
	v1CategoriaInsumos.HandleFunc("", c.AgregarCategoriaInsumos).Methods(http.MethodPost)

	v1Mesa := v1.PathPrefix("/mesas").Subrouter()
	v1Mesa.HandleFunc("/{id}", c.ObtenerMesa).Methods(http.MethodGet)
	v1Mesa.HandleFunc("/{id}", c.ModificarMesa).Methods(http.MethodPut)
	v1Mesa.HandleFunc("", c.ObtenerMesas).Methods(http.MethodGet)
	v1Mesa.HandleFunc("", c.AgregarMesa).Methods(http.MethodPost)
	v1Mesa.HandleFunc("/liberar/{id}", c.LiberarMesa).Methods(http.MethodPut)

	v1Pedido := v1.PathPrefix("/pedidos").Subrouter()
	v1Pedido.HandleFunc("", c.RegistrarPedido).Methods(http.MethodPost)
	v1Pedido.HandleFunc("", c.ObtenerPedidos).Methods(http.MethodGet)
}

func reports(api *mux.Router) {
	v1 := api.PathPrefix("/v1").Subrouter()

	v1Reporte := v1.PathPrefix("/reportes").Subrouter()
	v1Reporte.HandleFunc("/factura/{id}", r.FacturaPedido).Methods(http.MethodGet)
	v1Reporte.HandleFunc("/productos", r.ReporteProductos).Methods(http.MethodGet)
	v1Reporte.HandleFunc("/insumos", r.ReporteInsumos).Methods(http.MethodGet)
	v1Reporte.HandleFunc("/proveedores", r.ReporteProveedores).Methods(http.MethodGet)
	v1Reporte.HandleFunc("/diario", r.ReportePedidosDiarios).Methods(http.MethodGet)
	v1Reporte.HandleFunc("/individual/{id}", r.ReportePedidoIndividual).Methods(http.MethodGet)
}
