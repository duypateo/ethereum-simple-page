package renderer

import (
	"net/http"
	"strconv"
	"text/template"
	"time"
)

// packpagePath to this package
// const packpagePath = "src/ethereum-simple-page/"
const packpagePath = "./"

// templatePath to template
const templatePath = packpagePath + "template/"

// baseLayoutTmpl - layout template
const baseLayoutTmpl = templatePath + "base.html"

var functionMap = template.FuncMap{
	"convertToDatetime": func(timestamp string) time.Time {
		if timestampInt, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
			return time.Unix(timestampInt, 0)
		}
		return time.Now()
	},
}

// RenderTemplate - Render ui template by name
func RenderTemplate(w http.ResponseWriter, tmplName string, pageData map[string]interface{}) {
	_, ok := pageData["Title"]
	if !ok {
		pageData["Title"] = "ETHEREUM SIMPLE GATE"
	}

	tmpl, err := template.New("").Funcs(functionMap).ParseFiles(templatePath+tmplName+".html", baseLayoutTmpl)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = tmpl.ExecuteTemplate(w, "base", pageData)
}
