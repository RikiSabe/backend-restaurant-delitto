package querys

var CategoriasInsumos = `
	select id, nombre, descripcion,
		case when estado 
		then 'Activo'
		else 'Inactivo'
		end as estado
	from categoria_insumo
	order by id asc;`

var CategoriaInsumos = `
	select
		id, nombre, descripcion, 
		case
			when estado then 'Activo'
			else 'Inactivo'
		end as estado
	from categoria_insumo
	where id = ?
	limit 1;`
