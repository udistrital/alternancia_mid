package models

type EspacioFisicoPadre struct {
	Id      int
	PadreId *EspacioFisico
	HijoId  *EspacioFisico
}
