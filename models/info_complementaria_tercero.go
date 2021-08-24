package models

type InfoComplementariaTercero struct {
	Id                       int
	TerceroId                *Tercero
	InfoComplementariaId     *InfoComplementaria
	Dato                     string
	Activo                   bool
	InfoCompleTerceroPadreId *InfoComplementariaTercero
}
