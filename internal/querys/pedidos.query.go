package querys

var PedidosJSONB = `
	SELECT 
		p.id,
		p.fecha,
		p.estado,
		p.origen,
		jsonb_agg(
			jsonb_build_object(
				'nombre', pr.nombre,
				'precio', pr.precio,
				'cantidad', dp.cantidad,
				'subtotal', dp.subtotal
			)
		) AS productos
	FROM pedidos p
	JOIN detalle_pedido dp ON dp.id_pedido = p.id
	JOIN productos pr ON pr.id = dp.id_producto
	WHERE DATE(p.fecha) = CURRENT_DATE
	GROUP BY p.id
	ORDER BY p.fecha DESC;`
