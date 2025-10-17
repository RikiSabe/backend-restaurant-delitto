package querys

var Proovedor = `
	select 
		pro.id, 
		pro.tipo,
		pro.nombre, 
		pro.telefono,
		pro.correo,
		pro.direccion,
		pro.ci,
		pro.nombre_encargado,
		CASE 
			WHEN pro.estado THEN 'Activo'
			ELSE 'Inactivo'
		END AS estado
	from proveedores as pro
	where pro.id = ?
	limit 1;`

var Proovedores = `
	select 
		pro.id, 
		pro.tipo,
		pro.nombre, 
		pro.telefono,
		pro.correo,
		pro.direccion,
		pro.ci,
		pro.nombre_encargado,
		CASE 
			WHEN pro.estado THEN 'Activo'
			ELSE 'Inactivo'
		END AS estado
	from proveedores as pro
	ORDER BY pro.id asc;`
