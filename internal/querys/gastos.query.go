package querys

var Gastos = `
	select 
		gv.nombre,
		gv.unidad_medida,
		gv.id_proveedor,
		p.nombre
	from gastos_varios gv
	join proveedor p on gv.id_proveedor = p.id`

var Gasto = `
	select
		nombre, 
		unidad_medida,
		id_proveedor
	from gastos_varios
	where id = ?`
