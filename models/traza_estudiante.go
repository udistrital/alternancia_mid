package models

type TrazaEstudiante struct {
	Estudiante string
	Espacios   []struct {
		Nombre       string
		FechaEntrada string
		FechaSalida  string
	}
}
