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

var QueryProductos = `
	SELECT 
		p.nombre,
		p.precio,
		p.descripcion,
		COALESCE(cp.nombre, 'Sin Categoría') as categoria
	FROM productos p
	LEFT JOIN categoria_producto cp ON p.id_categoria = cp.id
	WHERE p.estado = true
	ORDER BY cp.nombre, p.nombre`

type ProductoReporte struct {
	Nombre      string  `json:"nombre"`
	Precio      float64 `json:"precio"`
	Descripcion string  `json:"descripcion"`
	Categoria   string  `json:"categoria"`
}

func ReporteProductos(w http.ResponseWriter, r *http.Request) {
	m, err := makePDFProductos()
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
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_productos.pdf\"")

	if _, err := w.Write(doc.GetBytes()); err != nil {
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
		return
	}
}

func makePDFProductos() (core.Maroto, error) {
	// Obtener productos
	var productos []ProductoReporte
	err := db.GDB.Raw(QueryProductos).Scan(&productos).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener productos: %v", err)
	}

	if len(productos) == 0 {
		return nil, fmt.Errorf("no se encontraron productos")
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
		text.NewCol(12, "REPORTE DE PRODUCTOS", props.Text{
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
		text.NewCol(4, "NOMBRE", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(2, "PRECIO (Bs)", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(3, "CATEGORÍA", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(3, "DESCRIPCIÓN", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
	)

	m.AddRows(line.NewRow(1))

	// Listado de productos
	categoriaActual := ""
	for _, producto := range productos {
		// Separador por categoría
		if categoriaActual != producto.Categoria {
			if categoriaActual != "" {
				m.AddRows(line.NewRow(1))
			}
			categoriaActual = producto.Categoria
		}

		// Descripción truncada si es muy larga
		descripcion := producto.Descripcion
		if len(descripcion) > 50 {
			descripcion = descripcion[:47] + "..."
		}

		m.AddRow(7,
			text.NewCol(4, producto.Nombre, props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
			text.NewCol(2, fmt.Sprintf("%.2f", producto.Precio), props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
			text.NewCol(3, producto.Categoria, props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
			text.NewCol(3, descripcion, props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
		)
	}

	m.AddRows(line.NewRow(2))

	// Pie de página con total de productos
	m.AddRow(8,
		text.NewCol(12, fmt.Sprintf("Total de productos: %d", len(productos)), props.Text{
			Top:   1.0,
			Align: align.Right,
			Style: fontstyle.Bold,
		}),
	)

	return m, nil
}
