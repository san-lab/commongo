package templates

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer() *Renderer {
	r := &Renderer{}
	r.LoadTemplates()
	return r
}

//Taken out of the constructor with the idea of forced template reloading
func (r *Renderer) LoadTemplates() {
	var allFiles []string
	files, err := ioutil.ReadDir("./templates")
	if err != nil {
		log.Println(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".htemplate") {
			allFiles = append(allFiles, "./templates/"+filename)
		}
	}
	r.templates, err = template.ParseFiles(allFiles...) //parses all .tmpl files in the 'templates' folder
	if err != nil {
		log.Println(err)
	}
}

// TODO Load uri->template association from json
func (r *Renderer) LoadTemplateMapping() {

}

func (r *Renderer) RenderResponse(w io.Writer, data RenderData) error {
	if len(data.TemplateName) < 1 {
		data.TemplateName = "home"
	}
	err := r.templates.ExecuteTemplate(w, data.TemplateName, data)
	if err != nil {
		log.Println(err)
	}
	return err

}

//This is a try to bring some uniformity to passing data to the templates
//The "RenderData" container is a wrapper for the header/body/footer containers
type RenderData struct {
	Error        string
	TemplateName string
	User         string
	InTEE        bool
	HeaderData   interface{}
	BodyData     interface{}
	FooterData   interface{}
}
