package reports

import (
	"backend-restaurant-delitto/internal/db"
	"fmt"
	"net/http"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

var QueryPedidosDiarios = `
	SELECT 
		p.id,
		p.fecha,
		p.estado,
		p.origen,
		COALESCE(m.nombre, 'Sin Mesa') as mesa_nombre,
		COALESCE(SUM(dp.subtotal), 0) as total
	FROM pedidos p
	LEFT JOIN mesas m ON p.id_mesa = m.id
	LEFT JOIN detalle_pedido dp ON p.id = dp.id_pedido
	WHERE DATE(p.fecha) = ?
	GROUP BY p.id, p.fecha, p.estado, p.origen, m.nombre
	ORDER BY p.fecha DESC`

var QueryProductoMasVendido = `
	SELECT 
		pr.nombre,
		SUM(dp.cantidad) as total_vendido,
		SUM(dp.subtotal) as total_ganancia
	FROM detalle_pedido dp
	INNER JOIN productos pr ON dp.id_producto = pr.id
	INNER JOIN pedidos p ON dp.id_pedido = p.id
	WHERE DATE(p.fecha) = ?
	GROUP BY pr.id, pr.nombre
	ORDER BY total_vendido DESC
	LIMIT 1`

var QueryEstadisticasDiarias = `
	SELECT 
		COUNT(DISTINCT p.id) as total_pedidos,
		COALESCE(SUM(dp.subtotal), 0) as ganancia_total,
		SUM(CASE WHEN p.origen = 'caja' THEN 1 ELSE 0 END) as pedidos_caja,
		SUM(CASE WHEN p.origen = 'panel-cliente' THEN 1 ELSE 0 END) as pedidos_panel
	FROM pedidos p
	LEFT JOIN detalle_pedido dp ON p.id = dp.id_pedido
	WHERE DATE(p.fecha) = ?`

type PedidoDiario struct {
	ID         uint      `json:"id"`
	Fecha      time.Time `json:"fecha"`
	Estado     string    `json:"estado"`
	Origen     string    `json:"origen"`
	MesaNombre string    `json:"mesa_nombre"`
	Total      float64   `json:"total"`
}

type ProductoMasVendido struct {
	Nombre        string  `json:"nombre"`
	TotalVendido  int     `json:"total_vendido"`
	TotalGanancia float64 `json:"total_ganancia"`
}

type EstadisticasDiarias struct {
	TotalPedidos  int     `json:"total_pedidos"`
	GananciaTotal float64 `json:"ganancia_total"`
	PedidosCaja   int     `json:"pedidos_caja"`
	PedidosPanel  int     `json:"pedidos_panel"`
}

func ReportePedidosDiarios(w http.ResponseWriter, r *http.Request) {
	// Obtener fecha del query parameter, si no existe usar fecha actual
	fechaStr := r.URL.Query().Get("fecha")
	var fecha time.Time
	var err error

	if fechaStr == "" {
		fecha = time.Now()
	} else {
		fecha, err = time.Parse("2006-01-02", fechaStr)
		if err != nil {
			http.Error(w, "Formato de fecha inválido. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	m, err := makePDFPedidosDiarios(fecha)
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
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"reporte_pedidos_%s.pdf\"", fecha.Format("2006-01-02")))

	if _, err := w.Write(doc.GetBytes()); err != nil {
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
		return
	}
}

func makePDFPedidosDiarios(fecha time.Time) (core.Maroto, error) {
	fechaStr := fecha.Format("2006-01-02")

	// Obtener estadísticas diarias
	var stats EstadisticasDiarias
	err := db.GDB.Raw(QueryEstadisticasDiarias, fechaStr).Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas: %v", err)
	}

	// Obtener producto más vendido
	var productoTop ProductoMasVendido
	err = db.GDB.Raw(QueryProductoMasVendido, fechaStr).Scan(&productoTop).Error
	if err != nil {
		productoTop.Nombre = "Sin ventas"
		productoTop.TotalVendido = 0
	}

	// Obtener todos los pedidos del día
	var pedidos []PedidoDiario
	err = db.GDB.Raw(QueryPedidosDiarios, fechaStr).Scan(&pedidos).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener pedidos: %v", err)
	}

	cfg := config.NewBuilder().
		WithLeftMargin(15).
		WithRightMargin(15).
		WithTopMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	// Encabezado
	m.AddRows(text.NewRow(12, "RESTAURANT DELITTO", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Size:  18,
	}))

	m.AddRow(8,
		text.NewCol(12, "REPORTE DIARIO DE PEDIDOS", props.Text{
			Top:   1.0,
			Align: align.Center,
			Style: fontstyle.Bold,
			Size:  14,
		}),
	)

	m.AddRow(8,
		text.NewCol(12, "Fecha: "+fecha.Format("02/01/2006"), props.Text{
			Top:   1.0,
			Align: align.Center,
			Style: fontstyle.Bold,
			Size:  12,
		}),
	)

	m.AddRows(line.NewRow(2))

	// SECCIÓN DE ESTADÍSTICAS
	m.AddRow(10,
		text.NewCol(12, "RESUMEN DEL DÍA", props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
			Size:  12,
		}),
	)

	m.AddRows(line.NewRow(1))

	// Estadísticas en 2 filas
	m.AddRow(8,
		text.NewCol(4, "Total de Pedidos:", props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(2, fmt.Sprintf("%d", stats.TotalPedidos), props.Text{
			Top:   1.0,
			Align: align.Left,
		}),
		text.NewCol(4, "Ganancia Total:", props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(2, fmt.Sprintf("%.2f Bs", stats.GananciaTotal), props.Text{
			Top:   1.0,
			Align: align.Left,
		}),
	)

	m.AddRow(8,
		text.NewCol(4, "Pedidos desde Caja:", props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(2, fmt.Sprintf("%d", stats.PedidosCaja), props.Text{
			Top:   1.0,
			Align: align.Left,
		}),
		text.NewCol(4, "Pedidos Panel Cliente:", props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(2, fmt.Sprintf("%d", stats.PedidosPanel), props.Text{
			Top:   1.0,
			Align: align.Left,
		}),
	)

	m.AddRows(line.NewRow(1))

	// Producto más vendido
	m.AddRow(8,
		text.NewCol(4, "Producto Más Vendido:", props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(8, productoTop.Nombre, props.Text{
			Top:   1.0,
			Align: align.Left,
		}),
	)

	if productoTop.TotalVendido > 0 {
		m.AddRow(6,
			text.NewCol(4, "", props.Text{}),
			text.NewCol(8, fmt.Sprintf("Cantidad vendida: %d unidades - Ganancia: %.2f Bs",
				productoTop.TotalVendido, productoTop.TotalGanancia), props.Text{
				Top:   1.0,
				Align: align.Left,
				Size:  9,
			}),
		)
	}

	m.AddRows(line.NewRow(2))

	// SECCIÓN DE PEDIDOS DETALLADOS
	if len(pedidos) > 0 {
		m.AddRow(10,
			text.NewCol(12, "DETALLE DE PEDIDOS", props.Text{
				Top:   1.0,
				Align: align.Left,
				Style: fontstyle.Bold,
				Size:  12,
			}),
		)

		m.AddRows(line.NewRow(1))

		// Encabezado de tabla
		m.AddRow(8,
			text.NewCol(1, "ID", props.Text{
				Style: fontstyle.Bold,
				Align: align.Center,
			}),
			text.NewCol(2, "HORA", props.Text{
				Style: fontstyle.Bold,
				Align: align.Center,
			}),
			text.NewCol(2, "MESA", props.Text{
				Style: fontstyle.Bold,
				Align: align.Center,
			}),
			text.NewCol(2, "ORIGEN", props.Text{
				Style: fontstyle.Bold,
				Align: align.Center,
			}),
			text.NewCol(3, "ESTADO", props.Text{
				Style: fontstyle.Bold,
				Align: align.Center,
			}),
			text.NewCol(2, "TOTAL (Bs)", props.Text{
				Style: fontstyle.Bold,
				Align: align.Center,
			}),
		)

		m.AddRows(line.NewRow(1))

		// Listado de pedidos
		for _, pedido := range pedidos {
			hora := pedido.Fecha.Format("15:04:05")

			m.AddRow(7,
				text.NewCol(1, fmt.Sprintf("%d", pedido.ID), props.Text{
					Top:   1.0,
					Align: align.Center,
				}),
				text.NewCol(2, hora, props.Text{
					Top:   1.0,
					Align: align.Center,
				}),
				text.NewCol(2, pedido.MesaNombre, props.Text{
					Top:   1.0,
					Align: align.Center,
				}),
				text.NewCol(2, pedido.Origen, props.Text{
					Top:   1.0,
					Align: align.Center,
				}),
				text.NewCol(3, pedido.Estado, props.Text{
					Top:   1.0,
					Align: align.Center,
				}),
				text.NewCol(2, fmt.Sprintf("%.2f", pedido.Total), props.Text{
					Top:   1.0,
					Align: align.Center,
				}),
			)
		}
	} else {
		m.AddRow(10,
			text.NewCol(12, "No se registraron pedidos en esta fecha", props.Text{
				Top:   1.0,
				Align: align.Center,
				Style: fontstyle.Bold,
			}),
		)
	}

	m.AddRows(line.NewRow(2))

	// Pie de página
	horaGeneracion := time.Now().Format("15:04:05")
	m.AddRow(6,
		text.NewCol(12, fmt.Sprintf("Reporte generado el %s a las %s",
			time.Now().Format("02/01/2006"), horaGeneracion), props.Text{
			Top:   1.0,
			Align: align.Center,
			Size:  8,
		}),
	)

	return m, nil
}
