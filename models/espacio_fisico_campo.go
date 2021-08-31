package models

type EspacioFisicoCampo struct {
	Id            int
	Valor         string
	EspacioFisico *EspacioFisico
	Campo         *Campo
}
