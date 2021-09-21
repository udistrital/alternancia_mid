package helpers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/alternancia_mid/models"
)

var idEstudiante string
var persona models.Tercero

func initPersona(idPersona string) (outputError map[string]interface{}) {
	idEstudiante = idPersona
	if status, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"tercero/"+idEstudiante, &persona); status != 200 || err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/initPersona", "err": err, "responseStatus": status, "status": "502"}
		return outputError
	}
	return
}

func ConsultarTraza(idPersona string) (trazaRes models.TrazaEstudiante, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetTraza", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	err := initPersona(idPersona)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetTraza", "err": err, "status": "502"}
		return models.TrazaEstudiante{}, outputError
	}

	traza, err := RegistrarTraza()
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetTraza/RegistrarTraza", "err": err, "status": "502"}
		return models.TrazaEstudiante{}, outputError
	}

	return traza, nil
}

func GenerarInformeTraza(idPersona string) (outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GenerarInformeTraza", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	err := initPersona(idPersona)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GenerarInformeTraza", "err": err, "status": "502"}
		return outputError
	}

	return
}

func RegistrarTraza() (traza models.TrazaEstudiante, outputError map[string]interface{}) {
	var res map[string]interface{}
	var res_seguimiento []models.RegistroTraza
	if status, err := getJsonTest(beego.AppConfig.String("UrlCrudSeguimiento")+"seguimiento/?limit=0&query=tercero_id:"+idEstudiante+"&sortby=fecha_creacion&order=asc", &res); status != 200 || err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetTraza/GetSeguimiento", "err": err, "responseStatus": status, "status": "502"}
		return models.TrazaEstudiante{}, outputError
	}
	LimpiezaRespuestaRefactor(res, &res_seguimiento)
	traza.Estudiante = persona.NombreCompleto
	var temp int
	var tempAula int
	for _, reg := range res_seguimiento {
		var list []models.EspacioFisico
		if status, err := getJsonTest(beego.AppConfig.String("UrlCrudOikos")+"espacio_fisico/?limit=1&query=Id:"+strconv.Itoa(reg.EspacioId), &list); status != 200 || err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/GetTraza/GetOikos/id:" + strconv.Itoa(reg.EspacioId), "err": err, "responseStatus": status, "status": "502"}
			return models.TrazaEstudiante{}, outputError
		}
		espacio := list[0]
		if reg.TipoEspacioId == 1 {
			if reg.TipoEscaneo == "I" {
				traza.Sedes = append(traza.Sedes, struct {
					Nombre       string
					FechaEntrada string
					FechaSalida  string
					Aulas        []struct {
						Nombre       string
						FechaEntrada string
						FechaSalida  string
					}
				}{
					Nombre:       espacio.Nombre,
					FechaEntrada: reg.FechaCreacion,
				})
				temp = len(traza.Sedes) - 1
			} else if reg.TipoEscaneo == "S" {
				traza.Sedes[temp].FechaSalida = reg.FechaCreacion
				temp = -1
			}
		} else if reg.TipoEspacioId == 3 {
			if reg.TipoEscaneo == "I" {
				traza.Sedes[temp].Aulas = append(traza.Sedes[temp].Aulas, struct {
					Nombre       string
					FechaEntrada string
					FechaSalida  string
				}{
					Nombre:       espacio.Nombre,
					FechaEntrada: reg.FechaCreacion,
				})
				tempAula = len(traza.Sedes[temp].Aulas) - 1
			} else if reg.TipoEscaneo == "S" {
				traza.Sedes[temp].Aulas[tempAula].FechaSalida = reg.FechaCreacion
				tempAula = -1
			}

		}
	}
	return
}

/*
func GenerarPDF(traza models.TrazaEstudiante) (outputError map[string]interface{}) {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitPoint, gofpdf.PageSizeLetter, "")

	pdf.AddPage()

	// ::Cosas de texto basicas::
	pdf.MoveTo(0, 0) // opcional
	pdf.SetFont("times", "B", 24)
	//pdf.Cell(40, 10, "Hola mundo")
	_, lineHeight := pdf.GetFontSize()
	pdf.SetTextColor(255, 0, 0) // rgb
	pdf.Text(0, lineHeight, "Contenido de prueba")

	pdf.MoveTo(0, lineHeight*2.0)
	pdf.SetFont("arial", "", 18)
	pdf.SetTextColor(100, 100, 100)
	_, lineHeight = pdf.GetFontSize()

	pdf.MultiCell(0, lineHeight*1.5, "Sed eu mauris nulla. In vitae laoreet nisi, eget maximus orci. Sed vel aliquam nisl, id suscipit tellus. Ut id tempor orci, non gravida odio.\n Nullam in scelerisque felis, id aliquet urna. Nam dapibus maximus ante, eu tempor neque lobortis eget.\n Suspendisse ut enim tincidunt, condimentum arcu id, ultrices urna. Maecenas sollicitudin fringilla aliquet.", gofpdf.BorderFull, gofpdf.AlignRight, false)

	// ::Cosas de formas basicas::
	pdf.SetFillColor(0, 255, 0) //F
	pdf.SetDrawColor(0, 0, 255) //D
	//-> Dibujar rectangulo:
	pdf.Rect(10, 100, 100, 100, "FD") //<-(SetFillColor + SetDrawColor)

	pdf.SetFillColor(100, 200, 200) //F
	//-> Dibujar poligono:
	pdf.Polygon([]gofpdf.PointType{{110, 250}, {160, 300}, {110, 350}, {60, 300}}, "F")

	//-> Crear cuadricula:
	drawGrid(pdf)

	//-> Agregar imagen:
	pdf.ImageOptions("imagenes/train.png", 275, 275, 92, 0, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")

	err := pdf.OutputFileAndClose("pf1.pdf")
	if err != nil {
		panic(err)
	}
	return
}

func drawGrid(pdf *gofpdf.Fpdf) {
	w, h := pdf.GetPageSize()
	pdf.SetFont("courier", "", 12)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetDrawColor(200, 200, 200)
	for x := 0.0; x < w; x = x + (w / 20.0) { //del Letter size
		pdf.Line(x, 0, x, h)
		_, lineHeight := pdf.GetFontSize()
		pdf.Text(x, lineHeight, fmt.Sprintf("%d", int(x)))
	}

	for y := 0.0; y < h; y = y + (w / 20.0) { //del Letter size
		pdf.Line(0, y, w, y)
		//_, lineHeight := pdf.GetFontSize()
		pdf.Text(0, y, fmt.Sprintf("%d", int(y)))
	}

}
*/
