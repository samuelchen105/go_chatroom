package common

import (
	"fmt"
	"html/template"
	"net/http"
)

func GenerateHTML(w http.ResponseWriter, data interface{}, files ...string) {
	var tmp []string
	for _, f := range files {
		tmp = append(tmp, fmt.Sprintf("template/%s.html", f))
	}

	tmpl := template.Must(template.ParseFiles(tmp...))
	tmpl.ExecuteTemplate(w, "layout", data)
}

func Redirect(w http.ResponseWriter, target string) {
	w.Header().Set("Location", target)
	w.WriteHeader(http.StatusFound)
}
