package controllers

import (
	"time"

	"gorm.io/datatypes"
)

type TypePedido struct {
	IDUsuario   uint   `json:"id_usuario"`
	IDMesa      uint   `json:"id_mesa"`
	IDProductos []uint `json:"id_productos"`
	IDInsumos   []uint `json:"id_insumos"`
}

type UsuarioEnlazado struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

type ProductoEnlazado struct {
	ID     uint    `json:"id"`
	Nombre string  `json:"nombre"`
	Precio float64 `json:"precio"`
}

type InsumoEnlazado struct {
	ID       uint   `json:"id"`
	Nombre   string `json:"nombre"`
	Cantidad uint   `json:"cantidad"`
}

type MesaEnlazada struct {
	ID        uint   `json:"id"`
	Estado    string `json:"estado"`
	Capacidad uint   `json:"capacidad"`
}

type Pedido struct {
	ID        uint                                   `json:"id"`
	Estado    string                                 `json:"estado"`
	Fecha     time.Time                              `json:"fecha"`
	IDMesa    uint                                   `json:"id_mesa"`
	Mesa      datatypes.JSONType[MesaEnlazada]       `json:"mesa"`
	Usuarios  datatypes.JSONType[[]UsuarioEnlazado]  `json:"usuarios"`
	Productos datatypes.JSONType[[]ProductoEnlazado] `json:"productos"`
	Insumos   datatypes.JSONType[[]InsumoEnlazado]   `json:"insumos"`
}

// func RegistrarPedido(w http.ResponseWriter, r *http.Request) {
// 	var pedidoBody TypePedido

// 	if err := json.NewDecoder(r.Body).Decode(&pedidoBody); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	var nuevo_pedido = models.Pedido{
// 		Estado: "Activo",
// 		Fecha:  time.Now(),
// 		IDMesa: pedidoBody.IDMesa,
// 	}

// 	tx := db.GDB.Begin()
// 	if err := tx.Create(&nuevo_pedido).Error; err != nil {
// 		tx.Rollback()
// 		http.Error(w, "Error al registrar pedido", http.StatusInternalServerError)
// 		return
// 	}
// 	if err := tx.Create(&models.UsuariosPedidos{
// 		IDUsuario: pedidoBody.IDUsuario,
// 		IDPedido:  nuevo_pedido.ID,
// 	}).Error; err != nil {
// 		tx.Rollback()
// 		http.Error(w, "Error en relacion Usuarios - Pedidos", http.StatusInternalServerError)
// 		return
// 	}
// 	for _, idProducto := range pedidoBody.IDProductos {
// 		if err := tx.Create(&models.PedidosProductos{
// 			IDPedido:   nuevo_pedido.ID,
// 			IDProducto: idProducto,
// 		}).Error; err != nil {
// 			tx.Rollback()
// 			http.Error(w, "Error en relacion Pedido - Producto", http.StatusInternalServerError)
// 			return
// 		}
// 	}
// 	for _, idInsumo := range pedidoBody.IDInsumos {
// 		if err := tx.Create(&models.PedidosInsumos{
// 			IDPedido: nuevo_pedido.ID,
// 			IDInsumo: idInsumo,
// 		}).Error; err != nil {
// 			tx.Rollback()
// 			http.Error(w, "Error en relacion Pedido - Insumo", http.StatusInternalServerError)
// 			return
// 		}
// 	}
// 	tx.Commit()

// 	w.Header().Set("Content-Type", "applicaction/json")
// 	w.WriteHeader(http.StatusOK)

// 	if err := json.NewEncoder(w).Encode(nuevo_pedido); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// func ObtenerPedidos(w http.ResponseWriter, r *http.Request) {
// 	var pedidos []Pedido

// 	query := `
// 		SELECT
// 			p.id,
// 			p.estado,
// 			p.fecha,
// 			p.id_mesa,
// 			jsonb_build_object(
// 				'id', m.id, 'estado', m.estado, 'capacidad', m.capacidad
// 			) AS mesa,
// 			json_agg(
// 				DISTINCT jsonb_build_object(
// 					'id', u.id, 'nombre', u.nombre, 'rol', u.rol
// 				)
// 			) FILTER (WHERE u.id IS NOT NULL) AS usuarios,
// 			json_agg(
// 				DISTINCT jsonb_build_object(
// 					'id', pr.id, 'nombre', pr.nombre, 'precio', pr.precio
// 				)
// 			) FILTER (WHERE pr.id IS NOT NULL) AS productos,
// 			json_agg(
// 				DISTINCT jsonb_build_object(
// 					'id', i.id, 'nombre', i.nombre, 'cantidad', i.cantidad
// 				)
// 			) FILTER (WHERE i.id IS NOT NULL) AS insumos
// 		FROM pedidos AS p
// 		LEFT JOIN mesas AS m ON p.id_mesa = m.id
// 		LEFT JOIN usuario_pedidos AS up ON p.id = up.id_pedido
// 		LEFT JOIN usuarios AS u ON up.id_usuario = u.id
// 		LEFT JOIN pedidos_productos AS pp ON p.id = pp.id_pedido
// 		LEFT JOIN productos AS pr ON pp.id_producto = pr.id
// 		LEFT JOIN pedidos_insumos AS pi ON p.id = pi.id_pedido
// 		LEFT JOIN insumos AS i ON pi.id_insumo = i.id
// 		GROUP BY
// 			p.id, p.estado, p.fecha, p.id_mesa,
// 			m.id, m.estado, m.capacidad
// 		ORDER BY p.id DESC;`

// 	err := db.GDB.Raw(query).Scan(&pedidos).Error
// 	if err != nil {
// 		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	if err := json.NewEncoder(w).Encode(pedidos); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }
