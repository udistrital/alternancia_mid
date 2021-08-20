package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/alternancia_mid/helpers"
)

// AccesoController operations for Acceso
type AccesoController struct {
	beego.Controller
}

// URLMapping ...
func (c *AccesoController) URLMapping() {
	c.Mapping("GetAutorizacion", c.GetAutorizacion)
	c.Mapping("GetIngreso", c.GetIngreso)
}

// GetAutorizacion ...
// @Title GetAutorizacion
// @Description Consulta autorización de acceso para quien muestra el QR
// @Param	idQr		path 	int			true			"Id de terceros de quien solicita acceso"
// @Param	idScan		path 	int			true			"Id de terceros de quien escanea"
// @Param	tipo		path 	string		true			"Tipo de escaneo (in/out)"
// @Param	sede		query 	int			true			"Id de la sede a consultar"
// @Param	edificio	query 	int			""				"Id del edificio a consultar"
// @Param	aula		query 	string		""				"Id del aula a consultar"
// @Success 200 {object} models.Persona
// @Failure 404	No found resource
// @router /:idQr/:idScan/:tipo [get]
func (c *AccesoController) GetAutorizacion() {
	idQrStr := c.GetString(":idQr")
	idScanStr := c.GetString(":idScan")
	tipo := c.GetString(":tipo")
	sede := c.GetString("sede")
	edificio := c.GetString("edificio")
	aula := c.GetString("aula")

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/GetAutorizacion/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	_, err := strconv.Atoi(idQrStr)
	if err != nil {
		panic(map[string]interface{}{"funcion": "GetAutorizacion", "err": "Error parametro de ingreso \"idQr\"", "status": "400", "log": err})
	}
	_, err = strconv.Atoi(idScanStr)
	if err != nil {
		panic(map[string]interface{}{"funcion": "GetAutorizacion", "err": "Error parametro de ingreso \"idScan\"", "status": "400", "log": err})
	}
	if tipo != "in" && tipo != "out" {
		panic(map[string]interface{}{"funcion": "GetAutorizacion", "err": "Error parametro de ingreso \"tipo\"", "status": "400", "log": err})
	}
	_, err = strconv.Atoi(sede)
	if err != nil {
		panic(map[string]interface{}{"funcion": "GetAutorizacion", "err": "Error parametro de ingreso \"sede\"", "status": "400", "log": err})
	}
	if c.GetString("edificio") != "" {
		_, err = strconv.Atoi(edificio)
		if err != nil {
			panic(map[string]interface{}{"funcion": "GetAutorizacion", "err": "Error parametro de ingreso \"edificio\"", "status": "400", "log": err})
		}
	}

	if respuesta, err := helpers.Autorizacion(idQrStr, idScanStr, aula, edificio, sede, tipo); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}

	c.ServeJSON()
}

// GetAcceso ...
// @Title GetAcceso
// @Description Da el acceso al estudiante
// @Param	idEspacio		path 	int			true			"Id del salon al que se accede"
// @Param	tipoQR			path 	string		true			"Determina si el qr escaneado es de entrada o salida (in/out)"
// @Success 200
// @Failure 404	No found resource
// @router /idEspacio/tipoQR [get]
func (c *AccesoController) GetIngreso() {

	idEspacio := c.GetString("idEspacio")
	tipoQr := c.GetString("tipoQR")

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/GetAutorizacion/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	if _, err := strconv.Atoi(idEspacio); err != nil {
		panic(map[string]interface{}{"funcion": "GetAutorizacion", "err": "Error parametro de ingreso \"idEspacio\"", "status": "400", "log": err})
	}
	if tipoQr != "out" && tipoQr != "in" {
		panic(map[string]interface{}{"funcion": "GetAutorizacion", "err": "Error parametro de ingreso \"tipoQR\"", "status": "400", "log": "El parámetro no cumple con las condiciones"})
	}

	if respuesta, err := helpers.ActualizarAforo(idEspacio, tipoQr); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}

	c.ServeJSON()
}
