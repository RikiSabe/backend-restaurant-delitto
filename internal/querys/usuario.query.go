package querys

var Usuarios = `
	select u.id, 
	concat(u.nombre, ' ', u.apellido) as nombre, 
	u.ci, u.celular, u.usuario,
	case 
		when u.estado then 'Activo'
		else 'Inactivo'
	end as estado, r.nombre as rol
	from usuarios as u
	join roles as r on r.id = u.id_rol
	order by u.id asc;`

var Usuario = `
	select u.id, u.nombre, u.apellido, u.ci, u.celular, u.usuario, u.contra,
		CASE 
			WHEN u.estado THEN 'Activo'
			ELSE 'Inactivo'
		end AS estado, r.nombre as rol
	FROM usuarios AS u
	LEFT JOIN roles AS r ON r.id = u.id_rol
	WHERE u.id = ?
	LIMIT 1;`