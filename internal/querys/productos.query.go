package querys

var Productos = `
	select p.id, p.nombre, p.descripcion, p.precio, p.imagen,
		case when p.estado
		then 'Activo'
		else 'Inactivo'
		end as estado, cp.nombre as categoria
	from productos p 
	left join categoria_producto as cp on cp.id = p.id_categoria
	order by p.id asc;`

var Producto = `
	select p.id, p.nombre, p.descripcion, p.precio, p.imagen,
		case when p.estado
		then 'Activo'
		else 'Inactivo'
		end as estado, cp.nombre as categoria
	from productos p 
	left join categoria_producto as cp on cp.id = p.id_categoria
	where p.id = ?
	limit 1;`

var ProductoPorNombre = `
	select id from productos
	where nombre = ?
	limit 1;
`
