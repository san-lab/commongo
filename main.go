package main

import (
	"net/http"
	"san-lab/commongo/gohttpservice"
	"san-lab/commongo/gohttpservice/templates"
)

func main() {
	gohttpservice.Startserver(testHandler)
}

var testHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		data := new(templates.RenderData)
		data.SessionID, _ = r.Cookie(gohttpservice.SessionIdName)

		//path := r.URL.Path[1:]
		data.User, _, _ = r.BasicAuth()
		r.ParseForm()

		data.TemplateName = "test"

		renderer.RenderResponse(w, *data)

	})

var renderer *templates.Renderer

func init() {
	renderer = new(templates.Renderer)
	renderer.LoadTemplates("./gohttpservice/templates")
}
