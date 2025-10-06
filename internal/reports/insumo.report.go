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

var QueryInsumos = `
	SELECT 
		i.nombre,
		i.stock_actual,
		i.stock_minimo,
		i.unidad_medida,
		COALESCE(ci.nombre, 'Sin Categoría') as categoria,
		COALESCE(p.nombre, 'Sin Proveedor') as proveedor
	FROM insumos i
	LEFT JOIN categoria_insumo ci ON i.id_categoria = ci.id
	LEFT JOIN proveedores p ON i.id_proveedor = p.id
	ORDER BY ci.nombre, i.nombre
`

type InsumoReporte struct {
	Nombre       string `json:"nombre"`
	StockActual  uint   `json:"stock_actual"`
	StockMinimo  uint   `json:"stock_minimo"`
	UnidadMedida string `json:"unidad_medida"`
	Categoria    string `json:"categoria"`
	Proveedor    string `json:"proveedor"`
}

func ReporteInsumos(w http.ResponseWriter, r *http.Request) {
	m, err := makePDFInsumos()
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
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_insumos.pdf\"")

	if _, err := w.Write(doc.GetBytes()); err != nil {
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
		return
	}
}

func makePDFInsumos() (core.Maroto, error) {
	// Obtener insumos
	var insumos []InsumoReporte
	err := db.GDB.Raw(QueryInsumos).Scan(&insumos).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener insumos: %v", err)
	}

	if len(insumos) == 0 {
		return nil, fmt.Errorf("no se encontraron insumos")
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
		text.NewCol(12, "REPORTE DE INSUMOS", props.Text{
			Top:   1.0,
			Align: align.Center,
			Style: fontstyle.Bold,
			Size:  14,
		}),
	)

	fechaActual := time.Now().Format("02/01/2006 15:04:05")
	m.AddRow(8,
		text.NewCol(12, "Fecha de generación: "+fechaActual, props.Text{
			Top:   1.0,
			Align: align.Center,
		}),
	)

	m.AddRows(line.NewRow(2))

	// Encabezado de tabla
	m.AddRow(8,
		text.NewCol(3, "NOMBRE", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(2, "STOCK ACTUAL", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(2, "STOCK MÍNIMO", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(2, "CATEGORÍA", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(3, "PROVEEDOR", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
	)

	m.AddRows(line.NewRow(1))

	// Listado de insumos con alertas
	categoriaActual := ""
	insumosAlerta := 0

	for _, insumo := range insumos {
		// Separador por categoría
		if categoriaActual != insumo.Categoria {
			if categoriaActual != "" {
				m.AddRows(line.NewRow(1))
			}
			categoriaActual = insumo.Categoria
		}

		// Verificar si está en alerta (stock actual <= stock mínimo)
		enAlerta := insumo.StockActual <= insumo.StockMinimo
		if enAlerta {
			insumosAlerta++
		}

		// Stock con unidad de medida
		stockActual := fmt.Sprintf("%d %s", insumo.StockActual, insumo.UnidadMedida)
		stockMinimo := fmt.Sprintf("%d %s", insumo.StockMinimo, insumo.UnidadMedida)

		// Nombre con indicador de alerta
		nombre := insumo.Nombre
		if enAlerta {
			nombre = "⚠ " + nombre
		}

		m.AddRow(7,
			text.NewCol(3, nombre, props.Text{
				Top:   1.0,
				Align: align.Center,
				Style: func() fontstyle.Type {
					if enAlerta {
						return fontstyle.Bold
					}
					return fontstyle.Normal
				}(),
			}),
			text.NewCol(2, stockActual, props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
			text.NewCol(2, stockMinimo, props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
			text.NewCol(2, insumo.Categoria, props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
			text.NewCol(3, insumo.Proveedor, props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
		)
	}

	m.AddRows(line.NewRow(2))

	// Resumen final
	m.AddRow(8,
		text.NewCol(6, fmt.Sprintf("Total de insumos: %d", len(insumos)), props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(6, fmt.Sprintf("Insumos en alerta: %d", insumosAlerta), props.Text{
			Top:   1.0,
			Align: align.Right,
			Style: fontstyle.Bold,
		}),
	)

	// Nota sobre alertas
	if insumosAlerta > 0 {
		m.AddRow(6,
			text.NewCol(12, "Los Insumos en negritas indican insumos con stock actual igual o menor al stock mínimo", props.Text{
				Top:   1.0,
				Align: align.Center,
				Size:  8,
			}),
		)
	}

	return m, nil
}
