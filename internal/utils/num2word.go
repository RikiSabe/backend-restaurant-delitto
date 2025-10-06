package utils

var unidades = []string{"", "uno", "dos", "tres", "cuatro", "cinco", "seis", "siete", "ocho", "nueve"}
var especiales = []string{"diez", "once", "doce", "trece", "catorce", "quince", "dieciséis", "diecisiete", "dieciocho", "diecinueve"}
var decenas = []string{"", "", "veinte", "treinta", "cuarenta", "cincuenta", "sesenta", "setenta", "ochenta", "noventa"}
var centenas = []string{"", "cien", "doscientos", "trescientos", "cuatrocientos", "quinientos", "seiscientos", "setecientos", "ochocientos", "novecientos"}

// Función principal
func ConvertES(n int) string {
	if n == 0 {
		return "cero"
	}
	return convertirNumero(n)
}

// Convierte números en diferentes rangos
func convertirNumero(n int) string {
	if n < 10 {
		return unidades[n]
	} else if n < 20 {
		return especiales[n-10]
	} else if n < 100 {
		if n%10 == 0 {
			return decenas[n/10]
		}
		return decenas[n/10] + " y " + unidades[n%10]
	} else if n < 1000 {
		if n == 100 {
			return "cien"
		}
		if n%100 == 0 {
			return centenas[n/100]
		}
		return centenas[n/100] + " " + convertirNumero(n%100)
	} else if n < 1000000 {
		if n == 1000 {
			return "mil"
		}
		if n < 2000 {
			return "mil " + convertirNumero(n%1000)
		}
		if n%1000 == 0 {
			return convertirNumero(n/1000) + " mil"
		}
		return convertirNumero(n/1000) + " mil " + convertirNumero(n%1000)
	} else if n < 1000000000 {
		if n == 1000000 {
			return "un millón"
		}
		if n < 2000000 {
			return "un millón " + convertirNumero(n%1000000)
		}
		if n%1000000 == 0 {
			return convertirNumero(n/1000000) + " millones"
		}
		return convertirNumero(n/1000000) + " millones " + convertirNumero(n%1000000)
	}
	return "número fuera de rango"
}
