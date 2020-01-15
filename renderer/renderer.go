package renderer

import (
	"net/http"
	"os"
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
	"convertToDatetime": func(timestamp string) string {
		if timestampInt, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
			unixTime := time.Unix(timestampInt, 0)
			return unixTime.Format("2006-01-02 15:04:05")
		}
		return "failed"
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

	pageData["Timezone"] = "Asia/Ho_Chi_Minh"

	configTz := os.Getenv("TZ")
	if configTz != "" {
		pageData["Timezone"] = configTz
	}

	_ = tmpl.ExecuteTemplate(w, "base", pageData)
}
