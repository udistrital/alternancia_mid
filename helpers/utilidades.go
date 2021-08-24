package helpers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/alternancia_mid/models"
	"github.com/udistrital/utils_oas/request"
)

//-------------------------Utilidades generales---------------------------

func getJsonTest(url string, target interface{}) (status int, err error) {
	r, err := http.Get(url)
	if err != nil {
		return r.StatusCode, err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()
	return r.StatusCode, json.NewDecoder(r.Body).Decode(target)
}

func putJson(url string, id string, body models.EspacioFisicoCampo) (outputError map[string]interface{}) {
	var res map[string]interface{}
	var env map[string]interface{}

	e, err := json.Marshal(body)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/PutJson", "err": err, "status": "502"}
		return outputError
	}

	json.Unmarshal(e, &env)
	if err := request.SendJson(url+"/"+strconv.Itoa(body.Id), "PUT", &res, env); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/PutJson", "err": err, "status": "502"}
		return outputError
	}
	return
}

func LimpiezaRespuestaRefactor(respuesta map[string]interface{}, v interface{}) {
	b, err := json.Marshal(respuesta["Data"])
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, &v)
}
