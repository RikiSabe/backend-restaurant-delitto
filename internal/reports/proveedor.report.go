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

var QueryProveedores = `
	SELECT 
		p.nombre,
		p.telefono,
		p.correo,
		p.direccion,
		p.estado,
		COUNT(i.id) as total_insumos
	FROM proveedores p
	LEFT JOIN insumos i ON p.id = i.id_proveedor
	GROUP BY p.id, p.nombre, p.telefono, p.correo, p.direccion, p.estado
	ORDER BY p.nombre`

type ProveedorReporte struct {
	Nombre       string `json:"nombre"`
	Telefono     string `json:"telefono"`
	Correo       string `json:"correo"`
	Direccion    string `json:"direccion"`
	Estado       bool   `json:"estado"`
	TotalInsumos int    `json:"total_insumos"`
}

func ReporteProveedores(w http.ResponseWriter, r *http.Request) {
	m, err := makePDFProveedores()
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
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_proveedores.pdf\"")

	if _, err := w.Write(doc.GetBytes()); err != nil {
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
		return
	}
}

func makePDFProveedores() (core.Maroto, error) {
	var proveedores []ProveedorReporte
	err := db.GDB.Raw(QueryProveedores).Scan(&proveedores).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener proveedores: %v", err)
	}

	if len(proveedores) == 0 {
		return nil, fmt.Errorf("no se encontraron proveedores")
	}

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
		text.NewCol(12, "REPORTE DE PROVEEDORES", props.Text{
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

	activos := 0
	inactivos := 0
	for _, prov := range proveedores {
		if prov.Estado {
			activos++
		} else {
			inactivos++
		}
	}

	m.AddRow(7,
		text.NewCol(4, fmt.Sprintf("Proveedores Activos: %d", activos), props.Text{
			Top:   1.0,
			Align: align.Left,
			Style: fontstyle.Bold,
		}),
		text.NewCol(4, fmt.Sprintf("Proveedores Inactivos: %d", inactivos), props.Text{
			Top:   1.0,
			Align: align.Center,
			Style: fontstyle.Bold,
		}),
		text.NewCol(4, fmt.Sprintf("Total: %d", len(proveedores)), props.Text{
			Top:   1.0,
			Align: align.Right,
			Style: fontstyle.Bold,
		}),
	)

	m.AddRows(line.NewRow(2))

	for i, proveedor := range proveedores {
		if i > 0 {
			m.AddRows(line.NewRow(1))
		}

		estadoTexto := "ACTIVO"
		if !proveedor.Estado {
			estadoTexto = "INACTIVO"
		}

		m.AddRow(8,
			text.NewCol(9, proveedor.Nombre, props.Text{
				Top:   1.0,
				Align: align.Left,
				Style: fontstyle.Bold,
				Size:  11,
			}),
			text.NewCol(3, estadoTexto, props.Text{
				Top:   1.0,
				Align: align.Right,
				Style: func() fontstyle.Type {
					if proveedor.Estado {
						return fontstyle.Bold
					}
					return fontstyle.Normal
				}(),
			}),
		)

		telefono := proveedor.Telefono
		if telefono == "" {
			telefono = "No registrado"
		}

		correo := proveedor.Correo
		if correo == "" {
			correo = "No registrado"
		}

		m.AddRow(6,
			text.NewCol(2, "Teléfono:", props.Text{
				Top:   1.0,
				Align: align.Left,
				Style: fontstyle.Bold,
			}),
			text.NewCol(4, telefono, props.Text{
				Top:   1.0,
				Align: align.Left,
			}),
			text.NewCol(2, "Insumos:", props.Text{
				Top:   1.0,
				Align: align.Left,
				Style: fontstyle.Bold,
			}),
			text.NewCol(4, fmt.Sprintf("%d", proveedor.TotalInsumos), props.Text{
				Top:   1.0,
				Align: align.Left,
			}),
		)

		m.AddRow(6,
			text.NewCol(2, "Correo:", props.Text{
				Top:   1.0,
				Align: align.Left,
				Style: fontstyle.Bold,
			}),
			text.NewCol(10, correo, props.Text{
				Top:   1.0,
				Align: align.Left,
			}),
		)

		direccion := proveedor.Direccion
		if direccion == "" {
			direccion = "No registrada"
		}

		m.AddRow(6,
			text.NewCol(2, "Dirección:", props.Text{
				Top:   1.0,
				Align: align.Left,
				Style: fontstyle.Bold,
			}),
			text.NewCol(10, direccion, props.Text{
				Top:   1.0,
				Align: align.Left,
			}),
		)
	}

	m.AddRows(line.NewRow(2))

	m.AddRow(6,
		text.NewCol(12, "Este reporte contiene información de contacto de todos los proveedores registrados", props.Text{
			Top:   1.0,
			Align: align.Center,
			Size:  8,
		}),
	)

	return m, nil
}
