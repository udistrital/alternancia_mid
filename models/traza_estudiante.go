package models

type TrazaEstudiante struct {
	Estudiante string
	Sedes      []struct {
		Nombre       string
		FechaEntrada string
		FechaSalida  string
		Aulas        []struct {
			Nombre       string
			FechaEntrada string
			FechaSalida  string
		}
	}
}
