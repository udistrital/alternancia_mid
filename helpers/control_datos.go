package helpers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/alternancia_mid/models"
)

var idEstudiante string

func ConsultarTraza(idPersona string) (traza models.TrazaEstudiante, outputError map[string]interface{}) {
	var persona models.Tercero
	idEstudiante = idPersona

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetTraza", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	if status, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"tercero/"+idEstudiante, &persona); status != 200 || err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetTraza/GetTercero", "err": err, "responseStatus": status, "status": "502"}
		return models.TrazaEstudiante{}, outputError
	}
	regEntrada, err := GetSeguimiento("I")
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetTraza/GetSeguimiento", "err": err, "status": "502"}
		return models.TrazaEstudiante{}, outputError
	}
	regSalida, err := GetSeguimiento("S")
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetTraza/GetSeguimiento", "err": err, "status": "502"}
		return models.TrazaEstudiante{}, outputError
	}
	logs.Debug(regEntrada)
	logs.Debug(regSalida)
	for index, reg := range regEntrada {
		var espacio models.EspacioFisico
		if status, err := getJsonTest(beego.AppConfig.String("UrlCrudOikos")+"espacio_fisico/"+strconv.Itoa(reg.EspacioId), &espacio); status != 200 || err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/GetTraza/GetOikos/id:" + strconv.Itoa(reg.EspacioId), "err": err, "responseStatus": status, "status": "502"}
			return models.TrazaEstudiante{}, outputError
		}
		traza.Estudiante = persona.NombreCompleto
		traza.Espacios = append(traza.Espacios, struct {
			Nombre       string
			FechaEntrada string
			FechaSalida  string
		}{
			Nombre:       espacio.Nombre,
			FechaEntrada: regEntrada[index].FechaCreacion,
			FechaSalida:  regSalida[index].FechaCreacion,
		})
	}

	return
}

func GetSeguimiento(tipoReg string) (res_seguimiento []models.RegistroTraza, outputError map[string]interface{}) {
	var res map[string]interface{}
	if status, err := getJsonTest(beego.AppConfig.String("UrlCrudSeguimiento")+"seguimiento/?limit=0&query=tercero_id:"+idEstudiante+",tipo_registro:"+tipoReg, &res); status != 200 || err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetTraza/GetSeguimiento", "err": err, "responseStatus": status, "status": "502"}
		return nil, outputError
	}
	LimpiezaRespuestaRefactor(res, &res_seguimiento)
	return
}
