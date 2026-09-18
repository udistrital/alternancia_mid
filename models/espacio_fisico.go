package models

type EspacioFisico struct {
	Id                  int
	Estado              string
	Nombre              string
	TipoEspacioFisicoId *TipoEspacio
	Codigo              string
}
