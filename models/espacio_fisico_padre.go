package models

type EspacioFisicoPadre struct {
	Id    int
	Padre *EspacioFisico
	Hijo  *EspacioFisico
}
