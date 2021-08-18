package models

type EspacioFisico struct {
	Id            int
	Nombre        string
	TipoEspacioId *TipoEspacio `json:"TipoEspacio"`
}
