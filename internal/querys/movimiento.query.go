package querys

var Movimientos = `
	select
		gv.nombre_gasto,
		gv.unidad_medida,
		m.cantidad,
		m.precio,
		m.fecha
	from movimientos m, gastos_varios gv
	where m.id_gasto = gv.id and gv.id = ?`
