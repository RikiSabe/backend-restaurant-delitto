package querys

var Gastos = `
	select 
		gv.id,
		gv.nombre,
		gv.unidad_medida,
		gv.id_proveedor,
		p.nombre as nombre_proveedor
	from gastos_varios gv
	join proveedores p on gv.id_proveedor = p.id`

var Gasto = `
	select
		id,
		nombre, 
		unidad_medida,
		id_proveedor
	from gastos_varios
	where id = ?`
