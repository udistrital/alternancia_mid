package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/alternancia_mid/models"
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
	if err := SendJson(url+"/"+strconv.Itoa(body.Id), "PUT", &res, env); err != nil || res["status"] != 200 {
		logs.Error(err)
		logs.Error(res["status"])
		logs.Error(res["message"])
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

func SendJson(urlp string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		json.NewEncoder(b).Encode(datajson)
	}
	//proxyUrl, err := url.Parse("http://10.20.4.15:3128")
	//http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

	client := &http.Client{}
	req, err := http.NewRequest(trequest, urlp, b)
	if err != nil {
		logs.Error(err)
		return err
	}

	//Se intenta acceder a cabecera, si no existe, se realiza peticion normal.
	defer func() {
		//Catch
		if r := recover(); r != nil {

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				beego.Error("Error reading response. ", err)
			}

			defer resp.Body.Close()
			respuesta := map[string]interface{}{"message": resp.Body, "status": resp.StatusCode}
			e, err := json.Marshal(respuesta)
			if err != nil {
				logs.Error(err)
			}

			json.Unmarshal(e, &target)
		}
	}()

	//try
	req.Header.Set("Authorization", "")

	resp, err := client.Do(req)
	if err != nil {
		beego.Error("Error reading response. ", err)
	}
	respuesta := map[string]interface{}{"message": resp.Body, "status": resp.StatusCode}
	e, err := json.Marshal(respuesta)
	if err != nil {
		logs.Error(err)
	}

	defer resp.Body.Close()
	return json.Unmarshal(e, &target)
}
