package reports

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/linestyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

var QueryPedidoIndividual = `
	SELECT 
		p.id,
		p.fecha,
		p.estado,
		p.origen,
		COALESCE(m.nombre, 'Sin Mesa') as mesa_nombre,
		COALESCE(m.capacidad, 0) as mesa_capacidad
	FROM pedidos p
	LEFT JOIN mesas m ON p.id_mesa = m.id
	WHERE p.id = ?
	LIMIT 1
`

var QueryDetallesPedidoIndividual = `
	SELECT 
		pr.nombre,
		dp.cantidad,
		pr.precio,
		dp.subtotal
	FROM detalle_pedido dp
	INNER JOIN productos pr ON dp.id_producto = pr.id
	WHERE dp.id_pedido = ?
	ORDER BY pr.nombre
`

var QueryUsuariosPedido = `
	SELECT 
		u.nombre,
		u.apellido
	FROM usuario_pedidos up
	INNER JOIN usuarios u ON up.id_usuario = u.id
	WHERE up.id_pedido = ?
`

type PedidoIndividual struct {
	ID            uint      `json:"id"`
	Fecha         time.Time `json:"fecha"`
	Estado        string    `json:"estado"`
	Origen        string    `json:"origen"`
	MesaNombre    string    `json:"mesa_nombre"`
	MesaCapacidad uint      `json:"mesa_capacidad"`
}

type DetalleProductoIndividual struct {
	Nombre   string  `json:"nombre"`
	Cantidad uint    `json:"cantidad"`
	Precio   float64 `json:"precio"`
	Subtotal float64 `json:"subtotal"`
}

type UsuarioPedido struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

func ReportePedidoIndividual(w http.ResponseWriter, r *http.Request) {
	id_pedido := mux.Vars(r)["id"]

	m, err := makePDFPedidoIndividual(id_pedido)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	doc, err := m.Generate()
	if err != nil {
		http.Error(w, "Error al generar pdf: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"pedido_%s.pdf\"", id_pedido))

	if _, err := w.Write(doc.GetBytes()); err != nil {
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
		return
	}
}

func makePDFPedidoIndividual(id_pedido string) (core.Maroto, error) {
	var pedido PedidoIndividual
	err := db.GDB.Raw(QueryPedidoIndividual, id_pedido).Scan(&pedido).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener pedido: %v", err)
	}

	if pedido.ID == 0 {
		return nil, fmt.Errorf("pedido no encontrado")
	}

	var detalles []DetalleProductoIndividual
	err = db.GDB.Raw(QueryDetallesPedidoIndividual, id_pedido).Scan(&detalles).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener detalles del pedido: %v", err)
	}

	if len(detalles) == 0 {
		return nil, fmt.Errorf("no se encontraron productos para el pedido")
	}

	var usuarios []UsuarioPedido
	db.GDB.Raw(QueryUsuariosPedido, id_pedido).Scan(&usuarios)

	cfg := config.NewBuilder().
		WithLeftMargin(15).
		WithRightMargin(15).
		WithTopMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	m.AddRows(text.NewRow(12, "RESTAURANT DELITTO", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Size:  18,
	}))

	m.AddRow(8,
		text.NewCol(12, "TARIJA - BOLIVIA", props.Text{
			Top:   1.0,
			Align: align.Center,
		}),
	)

	m.AddRow(8,
		text.NewCol(12, "Correo: restaurant.delitto.tja@gmail.com", props.Text{
			Top:   1.0,
			Align: align.Center,
		}),
	)

	m.AddRows(
		line.NewRow(1, props.Line{
			Style: linestyle.Dashed,
		}),
	)

	m.AddRow(10,
		text.NewCol(12, "DETALLE DE PEDIDO", props.Text{
			Top:   1.0,
			Align: align.Center,
			Style: fontstyle.Bold,
			Size:  16,
		}),
	)

	m.AddRows(
		line.NewRow(1, props.Line{
			Style: linestyle.Dashed,
		}),
	)

	fecha := pedido.Fecha.Format("02/01/2006")
	hora := pedido.Fecha.Format("15:04:05")

	m.AddRow(8,
		text.NewCol(6, "PEDIDO Nº:", props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(6, fmt.Sprintf("%d", pedido.ID), props.Text{
			Top:   1.0,
			Align: align.Left,
		}),
	)

	m.AddRow(7,
		text.NewCol(6, "FECHA:", props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(6, fecha, props.Text{
			Top:   1.0,
			Align: align.Left,
		}),
	)

	m.AddRow(7,
		text.NewCol(6, "HORA:", props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(6, hora, props.Text{
			Top:   1.0,
			Align: align.Left,
		}),
	)

	m.AddRow(7,
		text.NewCol(6, "ESTADO:", props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(6, strings.ToUpper(pedido.Estado), props.Text{
			Top:   1.0,
			Align: align.Left,
		}),
	)

	m.AddRow(7,
		text.NewCol(6, "ORIGEN:", props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(6, pedido.Origen, props.Text{
			Top:   1.0,
			Align: align.Left,
		}),
	)

	mesaInfo := pedido.MesaNombre
	if pedido.MesaCapacidad > 0 {
		mesaInfo = fmt.Sprintf("%s (Cap: %d)", pedido.MesaNombre, pedido.MesaCapacidad)
	}

	m.AddRow(7,
		text.NewCol(6, "MESA:", props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(6, mesaInfo, props.Text{
			Top:   1.0,
			Align: align.Left,
		}),
	)

	if len(usuarios) > 0 {
		nombresUsuarios := []string{}
		for _, u := range usuarios {
			nombresUsuarios = append(nombresUsuarios, fmt.Sprintf("%s %s", u.Nombre, u.Apellido))
		}

		m.AddRow(7,
			text.NewCol(6, "ATENDIDO POR:", props.Text{
				Top:   1.0,
				Align: align.Left,
				Style: fontstyle.Bold,
			}),
			text.NewCol(6, strings.Join(nombresUsuarios, ", "), props.Text{
				Top:   1.0,
				Align: align.Left,
			}),
		)
	}

	m.AddRows(line.NewRow(2))

	// Encabezado de tabla de productos
	m.AddRow(8,
		text.NewCol(5, "PRODUCTO", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(2, "CANTIDAD", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(2, "PRECIO", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(3, "SUBTOTAL", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
	)

	m.AddRows(line.NewRow(1))

	// Productos
	total := 0.0
	for _, detalle := range detalles {
		m.AddRow(7,
			text.NewCol(5, detalle.Nombre, props.Text{
				Top:   1.0,
				Align: align.Left,
			}),
			text.NewCol(2, fmt.Sprintf("%d", detalle.Cantidad), props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
			text.NewCol(2, fmt.Sprintf("%.2f", detalle.Precio), props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
			text.NewCol(3, fmt.Sprintf("%.2f", detalle.Subtotal), props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
		)

		total += detalle.Subtotal
	}

	m.AddRows(line.NewRow(2))

	// Total
	m.AddRow(10,
		text.NewCol(9, "TOTAL:", props.Text{
			Top:   1.0,
			Style: fontstyle.Bold,
			Align: align.Right,
			Size:  14,
		}),
		text.NewCol(3, fmt.Sprintf("%.2f Bs", total), props.Text{
			Top:   1.0,
			Style: fontstyle.Bold,
			Align: align.Center,
			Size:  14,
		}),
	)

	m.AddRows(line.NewRow(1))

	m.AddRow(7,
		text.NewCol(12, "SON: "+strings.ToUpper(utils.ConvertES(int(total)))+" BOLIVIANOS", props.Text{
			Top:   1.0,
			Align: align.Center,
			Style: fontstyle.Bold,
		}),
	)

	centavos := int((total - float64(int(total))) * 100)
	m.AddRow(6,
		text.NewCol(12, fmt.Sprintf("CON %02d/100", centavos), props.Text{
			Top:   1.0,
			Align: align.Center,
		}),
	)

	m.AddRows(
		line.NewRow(2, props.Line{
			Style: linestyle.Dashed,
		}),
	)

	horaGeneracion := time.Now().Format("02/01/2006 15:04:05")
	m.AddRow(6,
		text.NewCol(12, "Generado: "+horaGeneracion, props.Text{
			Top:   1.0,
			Align: align.Center,
			Size:  8,
		}),
	)

	return m, nil
}
