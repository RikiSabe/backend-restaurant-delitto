package querys

var PedidosJSONB = `
	SELECT 
		p.id,
		p.fecha,
		p.estado,
		p.origen,
		p.comentario,
		jsonb_agg(
			jsonb_build_object(
				'nombre', pr.nombre,
				'precio', pr.precio,
				'cantidad', dp.cantidad,
				'subtotal', dp.subtotal,
				'categoria', cp.nombre
			)
		) AS productos
	FROM pedidos p
	JOIN detalle_pedido dp ON dp.id_pedido = p.id
	JOIN productos pr ON pr.id = dp.id_producto
	LEFT JOIN categoria_producto cp ON cp.id = pr.id_categoria
	WHERE DATE(p.fecha) = CURRENT_DATE
	GROUP BY p.id
	ORDER BY p.fecha DESC;`
