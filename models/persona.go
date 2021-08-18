package models

import "time"

type Persona struct {
	Nombre string
	Fecha  time.Time
	Acceso string
	Causa  string
	Cupo   int
}
