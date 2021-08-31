package models

type EspacioFisico struct {
	Id          int
	Estado      string
	Nombre      string
	TipoEspacio *TipoEspacio
	Codigo      string
}
