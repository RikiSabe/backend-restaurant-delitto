package querys

var Insumos = `
	select
		ins.id, ins.nombre,	ins.stock_actual, ins.stock_minimo,	ins.unidad_medida,
		cat.id as id_categoria, cat.nombre as nombre_categoria,
		pro.id as id_proveedor, pro.nombre as nombre_proveedor
	from insumos as ins
	left join categoria_insumo cat on ins.id_categoria = cat.id 
	left join proveedores pro on ins.id_proveedor = pro.id
	order by ins.id asc;`

var Insumo = `
	select
		ins.id, ins.nombre,	ins.stock_actual, ins.stock_minimo,	ins.unidad_medida,
		cat.id as id_categoria, cat.nombre as nombre_categoria,
		pro.id as id_proveedor, pro.nombre as nombre_proveedor
	from insumos as ins
	left join categoria_insumo cat on ins.id_categoria = cat.id 
	left join proveedores pro on ins.id_proveedor = pro.id
	where ins.id = ?
	limit 1;`
