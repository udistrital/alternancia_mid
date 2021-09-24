package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/alternancia_mid/models"
	"github.com/udistrital/utils_oas/time_bogota"
)

func Autorizacion(idQr string, idScan string, salon string, idEdificio string, idSede string, tipoScan string) (persona models.Persona, outputError map[string]interface{}) {

	var respuesta_peticion []models.InfoComplementariaTercero

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/Autorizacion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"info_complementaria_tercero/?limit=-1&query=tercero_id:"+idQr, &respuesta_peticion); (err == nil) && (response == 200) {
		if len(respuesta_peticion) != 0 {

			//Declaracion de variables a usar
			var comorbilidad bool
			var vacuna bool
			var sintomas bool
			var cupoDisponible bool
			var permiso bool
			var clase bool
			var aforo int
			var cupo int
			id := idSede
			var idRol int
			//var materiasDia []models.CargaAcademica
			tercero := respuesta_peticion[0].TerceroId
			persona.Nombre = tercero.NombreCompleto
			persona.Fecha = time_bogota.Tiempo_bogota().Format("2006-01-02 15:04")
			persona.Acceso = "No autorizado"

			//Consulta de aforo
			aforo, err1 := ConsultarAforo(id)
			if err1 != nil {
				logs.Error(err1)
				outputError = map[string]interface{}{"funcion": "/ConsultarAforo", "err": err1, "status": "502"}
				return models.Persona{}, outputError
			}

			//Consulta de cupo
			cupo, err1 = ConsultarCupo(id)
			if err1 != nil {
				logs.Error(err1)
				outputError = map[string]interface{}{"funcion": "/ConsultarCupo", "err": err1, "status": "502"}
				return models.Persona{}, outputError
			}
			persona.Cupo = aforo - cupo

			//Consulta de comorbilidades
			comorbilidad, msg, err1 := ConsultarComorbilidades(strconv.Itoa(tercero.Id))
			if err1 != nil {
				logs.Error(err1)
				return models.Persona{}, err1
			}
			if msg != "" {
				persona.Causa = msg
				return
			}

			//Consulta de vacunacion
			vacuna, msg, err1 = ConsultarVacuna(strconv.Itoa(tercero.Id))
			if err1 != nil {
				logs.Error(err1)
				return models.Persona{}, err1
			}
			if msg != "" {
				persona.Causa = msg
				return
			}

			//Consulta de sintomas
			sintomas, msg, err1 = ConsultarSintomas(strconv.Itoa(tercero.Id))
			if err1 != nil {
				logs.Error(err1)
				return models.Persona{}, err1
			}
			if msg != "" {
				persona.Causa = msg
				return
			}

			if idEdificio != "" {
				id = idEdificio
			}
			if salon != "" {
				//Busqueda de id de espacio
				var respuesta_peticion_salones []models.EspacioFisicoPadre
				salon = strings.ToUpper(salon)
				coincidencias := 0
				if response, err := getJsonTest(beego.AppConfig.String("UrlCrudOikos")+"espacio_fisico_padre/?limit=-1&query=Padre.Id:"+idEdificio, &respuesta_peticion_salones); (err == nil) && (response == 200) {
					if len(respuesta_peticion_salones) > 0 {
						for _, espacio := range respuesta_peticion_salones {
							hijo := espacio.Hijo
							if strings.Contains(hijo.Nombre, salon) {
								id = strconv.Itoa(hijo.Id)
								coincidencias++
							}
						}
						if coincidencias > 1 {
							logs.Error("Se encontró más de un salon")
							outputError = map[string]interface{}{"funcion": "/ConsultaIdSalon", "err": "Se encontró más de un salon que coincide con la busqueda", "status": "502"}
							return models.Persona{}, outputError
						} else if coincidencias == 0 {
							logs.Error("No se encontró ningun salon")
							outputError = map[string]interface{}{"funcion": "/ConsultaIdSalon", "err": "No se encontró ningun salon", "status": "502"}
							return models.Persona{}, outputError
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "/ConsultaIdSalon", "err": "Se encontró más de un salon que coincide con la busqueda", "status": "502"}
						return models.Persona{}, outputError
					}
				} else {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/ConsultaIdSalon", "err": err, "status": "502"}
					return models.Persona{}, outputError
				}
			}

			if err != nil {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/Autorizacion", "err": err, "status": "502"}
				return models.Persona{}, outputError
			}

			if tipoScan == "in" {
				cupoDisponible = persona.Cupo > 0
			} else if tipoScan == "out" {
				cupoDisponible = persona.Cupo < aforo
			}

			//Consulta de roles para conceder permisos
			var respuesta_peticion_permisos []models.Vinculacion
			if response, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"vinculacion/?query=Activo:true,TerceroPrincipalId.Id:"+idScan, &respuesta_peticion_permisos); (err == nil) && (response == 200) {
				for _, vinculacion := range respuesta_peticion_permisos {
					idRol = vinculacion.TipoVinculacionId
					if idRol == 377 || idRol == 292 || idRol == 294 || (idRol >= 296 && idRol <= 299) {
						permiso = true
					}
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/Autorizacion/ConsultarPermisos", "err": err, "status": "502"}
				return models.Persona{}, outputError
			}

			//Consulta de horario
			/*respuesta := append(materiasDia,
				models.CargaAcademica{
					CODIGO_DIA: "7",
					HORA:       "15",
					SALON:      "AULA 201",
				}, models.CargaAcademica{
					CODIGO_DIA: "4",
					HORA:       "10",
					SALON:      "AULA 501",
				})
			var minutosAHora float64 = float64(time.Now().Local().Minute()) / 60
			var horaActual float64 = math.Round(float64(time.Now().Local().Hour()) + minutosAHora)
			for _, materia := range respuesta {
				if hora, err := strconv.ParseFloat(materia.HORA, 64); hora == horaActual && err == nil {
					clase = true
				}
			}*/
			clase = true

			espacioFisico, err := ConsultarEspacio(id)
			if err != nil {
				logs.Error(err)
				return models.Persona{}, err
			}

			//Rellenado de modelo persona
			validacion, err := validarFlujo(idQr, espacioFisico, tipoScan)
			if err != nil {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/Autorizacion", "err": err, "status": "502"}
				return models.Persona{}, outputError
			}
			if (vacuna || !comorbilidad) && !sintomas && cupoDisponible && clase && validacion {
				persona.Acceso = "Autorizado"
				if permiso {
					err := registrarFlujo(idQr, espacioFisico, tipoScan)
					if err != nil {
						logs.Error(err)
						return models.Persona{}, err
					}
					if tipoScan == "in" {
						persona.Cupo = persona.Cupo - 1
					} else if tipoScan == "out" {
						persona.Cupo = persona.Cupo + 1
					}
				}
			} else {
				if comorbilidad && !vacuna {
					persona.Causa = "Presenta comorbilidad(es) y no tiene vacunacion"
				} else if sintomas {
					persona.Causa = "Presenta sintomas"
				} else if !cupoDisponible {
					persona.Causa = "No hay cupo disponible en el espacio"
				} else if !clase {
					persona.Causa = "No tiene clase en este momento"
				} else if !validacion {
					var res map[string]interface{}
					var seguimiento []models.RegistroTraza
					if status, err := getJsonTest(beego.AppConfig.String("UrlCrudSeguimiento")+"seguimiento/?limit=1&order=desc&sortby=fecha_creacion", &res); status != 200 || err != nil {
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "/ValidarFlujo/GetSeguimiento", "err": err, "responseStatus": status, "status": "502"}
						return models.Persona{}, outputError
					}
					LimpiezaRespuestaRefactor(res, &seguimiento)
					ultimoReg := seguimiento[0]
					ultimoEspacio, err := ConsultarEspacio(strconv.Itoa(ultimoReg.EspacioId))
					if err != nil {
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "/ValidarFlujo/ConsultarEspacio", "err": err, "status": "502"}
						return models.Persona{}, outputError
					}
					if ultimoReg.TipoEspacioId == 3 && ultimoReg.TipoEscaneo == "I" && ((tipoScan == "out" && espacioFisico.TipoEspacio.Id == 1) || (tipoScan == "in" && espacioFisico.TipoEspacio.Id > 2)) {
						err := registrarFlujo(idQr, ultimoEspacio, "out")
						if err != nil {
							logs.Error(err)
							return models.Persona{}, err
						}
						err = registrarFlujo(idQr, espacioFisico, tipoScan)
						if err != nil {
							logs.Error(err)
							return models.Persona{}, err
						}
						persona.Acceso = "Autorizado"
						if tipoScan == "in" {
							persona.Cupo = persona.Cupo - 1
						} else if tipoScan == "out" {
							persona.Cupo = persona.Cupo + 1
						}
					} else {
						persona.Causa = "El registro del espacio es invalido, por favor asegurese de haber registrado todas las entradas y salidas"
					}
				}
			}
		} else {
			logs.Error("No hay datos de caracterización registrados para el usuario")
			outputError = map[string]interface{}{"funcion": "/Autorizacion", "err": "No hay datos de caracterización registrados para el usuario", "status": "502"}
			return models.Persona{}, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/Autorizacion", "err": err, "status": "502"}
		return models.Persona{}, outputError
	}

	return
}

func ActualizarAforo(idPersona string, idEspacio string, tipoQr string) (persona models.Persona, outputError map[string]interface{}) {
	cupo, _ := ConsultarCupo(idEspacio)
	var respuesta_peticion []models.InfoComplementariaTercero
	espacioFisico, err := ConsultarEspacio(idEspacio)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ActualizarAforo", "err": err, "status": "502"}
		return models.Persona{}, outputError
	}
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"info_complementaria_tercero/?limit=-1&query=tercero_id:"+idPersona, &respuesta_peticion); (err == nil) && (response == 200) {
		if len(respuesta_peticion) != 0 {
			persona.Nombre = respuesta_peticion[0].TerceroId.NombreCompleto
			persona.Fecha = time_bogota.Tiempo_bogota().Format("2006-01-02 15:04")
		} else {
			logs.Error("No hay datos de caracterización registrados para el usuario")
			outputError = map[string]interface{}{"funcion": "/ActualizarAforo", "err": "No hay datos de caracterización registrados para el usuario", "status": "502"}
			return models.Persona{}, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ActualizarAforo", "err": err, "status": "502"}
		return models.Persona{}, outputError
	}
	aforo, err := ConsultarAforo(idEspacio)
	if err != nil {
		logs.Error(err)
		return models.Persona{}, err
	}
	persona.Cupo = aforo - cupo
	persona.Acceso = "No autorizado"
	if tipoQr == "in" {
		comorbilidades, msg, err := ConsultarComorbilidades(idPersona)
		if err != nil {
			logs.Error(err)
			return models.Persona{}, err
		}
		if msg != "" {
			persona.Causa = msg
			return
		}
		vacuna, msg, err := ConsultarVacuna(idPersona)
		if err != nil {
			logs.Error(err)
			return models.Persona{}, err
		}
		if msg != "" {
			persona.Causa = msg
			return
		}
		sintomas, msg, err := ConsultarSintomas(idPersona)
		if err != nil {
			logs.Error(err)
			return models.Persona{}, err
		}
		if msg != "" {
			persona.Causa = msg
			return
		}
		if cupo < aforo && (vacuna || !comorbilidades) && !sintomas {
			//Registro de salida automático
			if val, err := validarFlujo(idPersona, espacioFisico, tipoQr); !val && err == nil {
				//Get hacia el último registro de tipo entrada hacia un aula
				var res map[string]interface{}
				var seguimiento []models.RegistroTraza
				if status, err := getJsonTest(beego.AppConfig.String("UrlCrudSeguimiento")+"seguimiento/?limit=1&order=desc&sortby=fecha_creacion", &res); status != 200 || err != nil {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/ValidarFlujo/GetSeguimiento", "err": err, "responseStatus": status, "status": "502"}
					return models.Persona{}, outputError
				}
				LimpiezaRespuestaRefactor(res, &seguimiento)
				ultimoReg := seguimiento[0]
				ultimoEspacio, err := ConsultarEspacio(strconv.Itoa(ultimoReg.EspacioId))
				if err != nil {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/ValidarFlujo/ConsultarEspacio", "err": err, "status": "502"}
					return models.Persona{}, outputError
				}
				if ultimoReg.TipoEspacioId == 3 && ultimoReg.TipoEscaneo == "I" && espacioFisico.TipoEspacio.Id > 2 {
					err = registrarFlujo(idPersona, ultimoEspacio, "out")
					if err != nil {
						logs.Error(err)
						return models.Persona{}, err
					}
					err = registrarFlujo(idPersona, espacioFisico, tipoQr)
					if err != nil {
						logs.Error(err)
						return models.Persona{}, err
					}
					persona.Acceso = "Autorizado"
					persona.Cupo--
				} else {
					persona.Causa = "Registro invalido, por favor asegurese de haber registrado todos los QR de entrada y salida"
				}
			} else if val {
				err := registrarFlujo(idPersona, espacioFisico, tipoQr)
				if err != nil {
					logs.Error(err)
					return models.Persona{}, err
				}
				persona.Acceso = "Autorizado"
				persona.Cupo--
			} else {
				logs.Error(err)
				return models.Persona{}, err
			}
		} else {
			if cupo == aforo {
				persona.Causa = "No hay cupo disponible en el espacio"
			} else if comorbilidades {
				persona.Causa = "Presenta comorbilidades"
			} else if sintomas {
				persona.Causa = "Presenta sintomas"
			}
		}
	} else if tipoQr == "out" {
		if cupo > 0 {
			if val, err := validarFlujo(idPersona, espacioFisico, tipoQr); val && err == nil {
				err1 := registrarFlujo(idPersona, espacioFisico, tipoQr)
				if err1 != nil {
					logs.Error(err1)
					return models.Persona{}, err1
				}
				persona.Acceso = "Autorizado"
				persona.Cupo++
			} else if !val {
				persona.Causa = "Registro invalido, por favor asegurese de haber escaneado todos los QR de entrada y salida"
				var res map[string]interface{}
				var seguimiento []models.RegistroTraza
				if status, err := getJsonTest(beego.AppConfig.String("UrlCrudSeguimiento")+"seguimiento/?limit=1&order=desc&sortby=fecha_creacion", &res); status != 200 || err != nil {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/ValidarFlujo/GetSeguimiento", "err": err, "responseStatus": status, "status": "502"}
					return models.Persona{}, outputError
				}
				LimpiezaRespuestaRefactor(res, &seguimiento)
				ultimoReg := seguimiento[0]
				ultimoEspacio, err := ConsultarEspacio(strconv.Itoa(ultimoReg.EspacioId))
				if err != nil {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/ValidarFlujo/ConsultarEspacio", "err": err, "status": "502"}
					return models.Persona{}, outputError
				}
				if ultimoReg.TipoEspacioId == 3 && ultimoReg.TipoEscaneo == "I" && espacioFisico.TipoEspacio.Id == 1 {
					err = registrarFlujo(idPersona, ultimoEspacio, "out")
					if err != nil {
						logs.Error(err)
						return models.Persona{}, err
					}
					err = registrarFlujo(idPersona, espacioFisico, tipoQr)
					if err != nil {
						logs.Error(err)
						return models.Persona{}, err
					}
					persona.Acceso = "Autorizado"
					persona.Causa = ""
					persona.Cupo++
				}
			} else {
				logs.Error(err)
				return models.Persona{}, err
			}

		} else {
			persona.Causa = "El cupo del espacio no concuerda"
		}
	}
	return
}

func ConsultarAforo(id string) (aforo int, outputError map[string]interface{}) {
	var respuesta_peticion_aforo []models.EspacioFisicoCampo
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudOikos")+"espacio_fisico_campo/?query=Campo.Id:5,EspacioFisico.Id:"+id, &respuesta_peticion_aforo); (err == nil) && (response == 200) {
		if len(respuesta_peticion_aforo) != 0 {
			aforoStr := respuesta_peticion_aforo[0].Valor
			aforo, _ = strconv.Atoi(aforoStr)
		} else {
			outputError = map[string]interface{}{"funcion": "/ConsultarAforo", "err": "No hay aforo registrado para el espacio", "status": "502"}
			return 0, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ConsultarAforo", "err": err, "status": "502"}
		return 0, outputError
	}
	return
}

func ConsultarCupo(id string) (cupo int, outputError map[string]interface{}) {

	entradas, err := contarGet(beego.AppConfig.String("UrlCrudSeguimiento") + "seguimiento/?query=oikos_id:" + id + ",tipo_registro:I")
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ConsultarCupo/contarEntradas", "err": err, "status": "502"}
		return 0, outputError
	}
	salidas, err := contarGet(beego.AppConfig.String("UrlCrudSeguimiento") + "seguimiento/?query=oikos_id:" + id + ",tipo_registro:S")
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ConsultarAforo/contarSalidas", "err": err, "status": "502"}
		return 0, outputError
	}
	cupo = entradas - salidas
	return
}

func ConsultarSintomas(idQr string) (sintomas bool, msg string, outputError map[string]interface{}) {
	var respuesta_peticion_sintomas map[string]interface{}
	var sintoma []models.Sintomas
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudSintomas")+"sintomas?limit=1&order=desc&sortby=fecha_creacion&query=terceroId:"+idQr, &respuesta_peticion_sintomas); (err == nil) && (response == 200) {
		LimpiezaRespuestaRefactor(respuesta_peticion_sintomas, &sintoma)
		if len(sintoma) != 0 && strings.Contains(sintoma[0].FechaCreacion, time.Now().UTC().Format("2006-01-02T")) {
			sintomasRegistrados := sintoma[0].InfoSalud
			sintomas = (sintomasRegistrados.Agotamiento ||
				sintomasRegistrados.CongestionNasal ||
				sintomasRegistrados.ContactoCovid ||
				sintomasRegistrados.DificultadRespiratoria ||
				sintomasRegistrados.EstadoEmbarazo ||
				sintomasRegistrados.Fiebre ||
				sintomasRegistrados.MalestarGeneral)
		} else {
			msg = "No se encontraron registros de sintomas para hoy"
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ConsultarSintomas", "err": err, "status": "502"}
		return false, "", outputError
	}
	return
}

func ConsultarVacuna(idQr string) (vacuna bool, msg string, outputError map[string]interface{}) {
	var respuesta_peticion_vacuna []models.InfoComplementariaTercero
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"info_complementaria_tercero/?query=tercero_id:"+idQr+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:V&order=asc&sortby=InfoComplementariaId", &respuesta_peticion_vacuna); (err == nil) && (response == 200) {
		if len(respuesta_peticion_vacuna) != 0 {
			layout := "2006-01-02T15:04:05.000Z"
			var dato map[string]interface{}
			json.Unmarshal([]byte(respuesta_peticion_vacuna[1].Dato), &dato)
			var fecha = dato["dato"]
			data, _ := json.Marshal(fecha)
			if string(data) != "\"\"" {
				str := fmt.Sprintf("%v", fecha)
				t, err := time.Parse(layout, str)

				if err != nil {
					outputError = map[string]interface{}{"funcion": "/ConsultarVacuna", "err": err, "status": "502"}
					return false, msg, outputError
				}
				duracion := time.Since(t)
				dias := int(duracion.Hours() / 24)
				if dias > 14 {
					vacuna = true
				}
			}
		} else {
			msg = "No se encontraron datos de vacunación registrados"
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ConsultarVacuna", "err": err, "status": "502"}
		return false, "", outputError
	}
	return
}

func ConsultarEspacio(idEspacio string) (espacioFisico models.EspacioFisico, outputError map[string]interface{}) {
	var list []models.EspacioFisico
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudOikos")+"espacio_fisico/?limit=1&query=Id:"+idEspacio, &list); (err != nil) || (response != 200) {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ConsultaEspacio", "err": err, "status": "502"}
		return models.EspacioFisico{}, outputError
	}

	return list[0], outputError
}

func ConsultarComorbilidades(idQr string) (comorbilidad bool, msg string, outputError map[string]interface{}) {
	var respuesta_peticion_comorbilidades []models.InfoComplementariaTercero
	var dato map[string]interface{}
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"info_complementaria_tercero/?limit=-1&query=tercero_id:"+idQr+",InfoComplementariaId.GrupoInfoComplementariaId.Id:47", &respuesta_peticion_comorbilidades); (err == nil) && (response == 200) {
		if len(respuesta_peticion_comorbilidades) > 0 {
			for _, info := range respuesta_peticion_comorbilidades {
				json.Unmarshal([]byte(info.Dato), &dato)
				if dato["dato"] == true {
					comorbilidad = true
				}
			}
		} else {
			msg = "No se encontró información de comorbilidades registrada"
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ConsultarComorbilidades", "err": err, "status": "502"}
		return false, "", outputError
	}
	return
}

func registrarFlujo(idTercero string, espacio models.EspacioFisico, tipo string) (outputError map[string]interface{}) {
	idTInt, _ := strconv.Atoi(idTercero)
	var tipoS string
	if tipo == "in" {
		tipoS = "I"
	} else {
		tipoS = "S"
	}
	var tipoEsp int
	if espacio.TipoEspacio.Id == 1 || espacio.TipoEspacio.Id == 2 {
		tipoEsp = espacio.TipoEspacio.Id
	} else {
		tipoEsp = 3
	}
	var body = models.RegistroTraza{
		TerceroId:         idTInt,
		EspacioId:         espacio.Id,
		TipoEspacioId:     tipoEsp,
		TipoEscaneo:       tipoS,
		Activo:            true,
		FechaCreacion:     time.Now().Format("2006-01-02T15:04:05.000Z"),
		FechaModificacion: time.Now().Format("2006-01-02T15:04:05.000Z"),
	}
	var res map[string]interface{}
	if err := SendJson(beego.AppConfig.String("UrlCrudSeguimiento")+"seguimiento", "POST", &res, body); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/registrarFlujo/PostSeguimiento", "resStatus": res["status"], "err": err, "status": "502"}
		return outputError
	}
	return
}

func validarFlujo(idTercero string, espacio models.EspacioFisico, tipo string) (validacion bool, outputError map[string]interface{}) {
	//get del nuevo api que traiga el último registro del tipo de espacio
	var res map[string]interface{}
	var seguimiento []models.RegistroTraza
	if status, err := getJsonTest(beego.AppConfig.String("UrlCrudSeguimiento")+"seguimiento/?limit=1&query=tercero_id:"+idTercero+"&order=desc&sortby=fecha_creacion", &res); status != 200 || err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ValidarFlujo/GetSeguimiento", "err": err, "responseStatus": status, "status": "502"}
		return false, outputError
	}
	LimpiezaRespuestaRefactor(res, &seguimiento)
	if len(seguimiento) != 0 {
		regSeguimiento := seguimiento[0]
		tipoReg := regSeguimiento.TipoEscaneo
		idEspReg := regSeguimiento.EspacioId
		var tipoEsp int
		if espacio.TipoEspacio.Id == 1 || espacio.TipoEspacio.Id == 2 {
			tipoEsp = espacio.TipoEspacio.Id
		} else {
			tipoEsp = 3
		}
		if tipo == "in" {
			return (tipoReg == "S" && seguimiento[0].TipoEspacioId == tipoEsp) || (tipoEsp > seguimiento[0].TipoEspacioId && tipoReg == "I"), nil
		}
		return (tipoReg == "I" && idEspReg == espacio.Id) || (tipoEsp < seguimiento[0].TipoEspacioId && tipoReg == "S"), nil
	} else {
		return true, nil
	}
}
