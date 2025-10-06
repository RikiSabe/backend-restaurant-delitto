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

func FacturaPedido(w http.ResponseWriter, r *http.Request) {
	id_pedido := mux.Vars(r)["id"]
	m, err := makePDFRecibo(id_pedido)
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
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_recibo.pdf\"")

	if _, err := w.Write(doc.GetBytes()); err != nil {
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
		return
	}
}

var QueryPedido = `
	SELECT 
		p.id,
		p.fecha,
		p.estado,
		p.origen,
		COALESCE(m.nombre, 'N/A') as mesa_nombre
	FROM pedidos p
	LEFT JOIN mesas m ON p.id_mesa = m.id
	WHERE p.id = ?
	LIMIT 1
`

var QueryDetallesPedido = `
	SELECT 
		dp.cantidad,
		pr.nombre,
		pr.precio,
		dp.subtotal
	FROM detalle_pedido dp
	INNER JOIN productos pr ON dp.id_producto = pr.id
	WHERE dp.id_pedido = ?`

type PedidoData struct {
	ID         uint      `json:"id"`
	Fecha      time.Time `json:"fecha"`
	Estado     string    `json:"estado"`
	Origen     string    `json:"origen"`
	MesaNombre string    `json:"mesa_nombre"`
}

type DetalleProducto struct {
	Nombre   string  `json:"nombre"`
	Cantidad uint    `json:"cantidad"`
	Precio   float64 `json:"precio"`
	Subtotal float64 `json:"subtotal"`
}

func makePDFRecibo(id_pedido string) (core.Maroto, error) {
	// Obtener datos del pedido
	var pedido PedidoData
	err := db.GDB.Raw(QueryPedido, id_pedido).Scan(&pedido).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener pedido: %v", err)
	}

	// Obtener detalles del pedido
	var detalles []DetalleProducto
	err = db.GDB.Raw(QueryDetallesPedido, id_pedido).Scan(&detalles).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener detalles del pedido: %v", err)
	}

	if len(detalles) == 0 {
		return nil, fmt.Errorf("no se encontraron productos para el pedido")
	}

	cfg := config.NewBuilder().
		WithLeftMargin(5).
		WithDimensions(120, 220).
		WithCompression(false).
		Build()

	m := maroto.New(cfg)

	// Encabezado
	m.AddRows(text.NewRow(10, "RESTAURANT DELITTO", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Size:  11,
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
		text.NewCol(12, "RECIBO ELECTRÓNICO", props.Text{
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

	// Información del pedido
	fecha := pedido.Fecha.Format("02/01/2006")
	hora := pedido.Fecha.Format("15:04:05")

	m.AddRow(8,
		text.NewCol(12, fmt.Sprintf("PEDIDO Nº: %d", pedido.ID), props.Text{
			Left:  6.0,
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
	)

	m.AddRow(8,
		text.NewCol(12, fmt.Sprintf("FECHA: %s  HORA: %s", fecha, hora), props.Text{
			Left:  6.0,
			Top:   1.0,
			Align: align.Left,
		}),
	)

	m.AddRow(8,
		text.NewCol(12, fmt.Sprintf("MESA: %s", pedido.MesaNombre), props.Text{
			Left:  6.0,
			Top:   1.0,
			Align: align.Left,
		}),
	)

	m.AddRows(line.NewRow(1))

	// Encabezado de tabla
	m.AddRow(6,
		text.NewCol(4, "CANT.", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(4, "PRECIO", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(4, "TOTAL", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
	)

	m.AddRows(line.NewRow(1))

	// Productos
	total := 0.0
	for _, detalle := range detalles {
		m.AddRow(6,
			text.NewCol(12, detalle.Nombre, props.Text{
				Left:  6.0,
				Align: align.Left,
			}),
		)
		m.AddRow(5,
			text.NewCol(4, fmt.Sprintf("%d", detalle.Cantidad), props.Text{
				Bottom: 1.0,
				Align:  align.Center,
			}),
			text.NewCol(4, fmt.Sprintf("%.2f", detalle.Precio), props.Text{
				Bottom: 1.0,
				Align:  align.Center,
			}),
			text.NewCol(4, fmt.Sprintf("%.2f", detalle.Subtotal), props.Text{
				Bottom: 1.0,
				Align:  align.Center,
			}),
		)

		total += detalle.Subtotal
	}

	m.AddRows(line.NewRow(1))

	// Total final
	m.AddRow(10,
		text.NewCol(8, "TOTAL: ", props.Text{
			Left:  6.0,
			Style: fontstyle.Bold,
			Align: align.Left,
		}),
		text.NewCol(4, fmt.Sprintf("%.2f Bs", total), props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
	)

	m.AddRow(6,
		text.NewCol(12, "SON: "+strings.ToUpper(utils.ConvertES(int(total)))+" BOLIVIANOS", props.Text{
			Left:  6.0,
			Align: align.Left,
		}),
	)

	centavos := int((total - float64(int(total))) * 100)
	m.AddRow(6,
		text.NewCol(12, fmt.Sprintf("CON %02d/100", centavos), props.Text{
			Left:  6.0,
			Align: align.Left,
		}),
	)

	m.AddRows(
		line.NewRow(1, props.Line{
			Style: linestyle.Dashed,
		}),
	)

	m.AddRow(8,
		text.NewCol(12, "¡GRACIAS POR SU PREFERENCIA!", props.Text{
			Top:   1.0,
			Align: align.Center,
			Style: fontstyle.Bold,
		}),
	)

	m.AddRows(
		line.NewRow(1, props.Line{
			Style: linestyle.Dashed,
		}),
	)

	return m, nil
}
