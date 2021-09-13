package models

import "time"

type RegistroTraza struct {
	Id            string
	TerceroId     int
	EspacioId     int
	TipoEspacioId int
	TipoEscaneo   string
	FechaReg      time.Time
}
