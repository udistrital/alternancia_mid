package models

type Vinculacion struct {
	Id                 int
	TerceroPrincipalId *Tercero
	TipoVinculacionId  int
	Activo             bool
	Alternancia        bool
}
