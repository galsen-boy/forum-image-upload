package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func IsvalidEmail(email string) bool {

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func Handle404Error(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	DataExec := map[string]interface{}{
		"ErrNum":  http.StatusNotFound,
		"TextErr": "Page Not Found"}
	RenderTemplate(w, "errorPage", DataExec)
}

func Handle405Error(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	DataExec := map[string]interface{}{
		"ErrNum":  http.StatusMethodNotAllowed,
		"TextErr": "Method Not Allowed"}
	RenderTemplate(w, "errorPage", DataExec)
}

func Handle400Error(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	DataExec := map[string]interface{}{
		"ErrNum":  http.StatusBadRequest,
		"TextErr": "Bad Request"}
	RenderTemplate(w, "errorPage", DataExec)
}

func Handle500Error(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	DataExec := map[string]interface{}{
		"ErrNum":  http.StatusInternalServerError,
		"TextErr": "Internal Server Error"}
	RenderTemplate(w, "errorPage", DataExec)
}

func RenderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		Handle500Error(w)
	}
}

func FormatDate(dateStr string) string {
	dateSt := strings.Split(dateStr, " ")[0]
	date, err := time.Parse("2006-01-02", dateSt)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		os.Exit(1)
	}
	formatedDate := date.Format("Jan 2, 2006")

	return fmt.Sprintf(formatedDate + " at " + strings.Split(dateStr, " ")[1][:5])
}
