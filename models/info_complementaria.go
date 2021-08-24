package models

type InfoComplementaria struct {
	Id                        int
	Nombre                    string
	CodigoAbreviacion         string
	Activo                    bool
	TipoDeDato                string
	GrupoInfoComplementariaId *GrupoInfoComplementaria
}
