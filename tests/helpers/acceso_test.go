package test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/udistrital/alternancia_mid/helpers"
	"github.com/udistrital/alternancia_mid/models"
	_ "github.com/udistrital/alternancia_mid/routers"

	"github.com/astaxie/beego"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

func TestConsultarAforo(t *testing.T) {
	_, err := helpers.ConsultarAforo("1096")
	if err != nil {
		t.Error("La consulta de aforo ha dado error")
		t.Error(err)
		t.Fail()
	} else {
		t.Log("TestConsultarAforo finalizado correctamente")
	}
}

func TestConsultarCupo(t *testing.T) {
	cupo := models.EspacioFisicoCampo{
		Id:    0,
		Valor: "",
	}
	err := helpers.ConsultarCupo("1096", &cupo)
	if err != nil {
		t.Error("La consulta de cupo ha dado error")
		t.Error(err)
		t.Fail()
	} else {
		t.Log("TestConsultarCupo finalizado correctamente")
	}
}

func TestConsultarSintomas(t *testing.T) {
	if _, err := helpers.ConsultarSintomas("9851"); err != nil {
		t.Error("La consulta de sintomas ha dado error")
		t.Error(err)
		t.Fail()
	} else {
		t.Log("TestConsultarSintomas finalizado correctamente")
	}
}

func TestConsultarVacuna(t *testing.T) {
	if _, err := helpers.ConsultarVacuna("9851"); err != nil {
		t.Error("La consulta de vacunacion ha dado error")
		t.Error(err)
		t.Fail()
	} else {
		t.Log("TestConsultarVacuna finalizado correctamente")
	}
}

func TestConsultarComorbilidades(t *testing.T) {
	if _, err := helpers.ConsultarComorbilidades("9851"); err != nil {
		t.Error("La consulta de comorbilidades ha dado error")
		t.Error(err)
		t.Fail()
	} else {
		t.Log("TestConsultarComorbilidades finalizado correctamente")
	}
}
