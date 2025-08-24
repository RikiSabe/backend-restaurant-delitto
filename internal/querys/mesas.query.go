package querys

var Mesas = `
	select
		id, nombre, estado, capacidad
	from mesas
	order by id asc;`

var Mesa = `
	select
		nombre, estado, capacidad
	from mesas
	where id = ?
	limit 1;`
