package querys

var CategoriasProductos = `
	select cp.id, cp.nombre, cp.descripcion,
		case when cp.estado 
		then 'Activo'
		else 'Inactivo'
		end as estado
	from categoria_producto cp
	order by cp.id asc;`

var CategoriaProductos = `
	select cp.nombre, cp.descripcion,
		case when cp.estado 
		then 'Activo'
		else 'Inactivo'
		end as estado
	from categoria_producto cp
	where cp.id = ?
	limit 1;`
