package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/asmsh/lazy"
)

// helloTmpl will panic if it fails to load the waned value,
// hence no need to check its Err() return.
var helloTmpl = lazy.NewValue(func() (*template.Template, error) {
	tmplTxt := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Greetings</title>
</head>
<body>
{{if . -}} Hello {{.}}! {{- else -}} Hello! {{- end}}
</body>
</html>
`
	tmpl := template.Must(template.New("hello").Parse(tmplTxt))
	return tmpl, nil
})

func main() {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		name := request.FormValue("name")

		// handle failing to load the template for any reason
		if helloTmpl.Err() != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// the template is loaded successfully, execute it
		err := helloTmpl.Val().Execute(writer, name)

		// handle failing to execute the template for any reason
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	log.Println("listening on http://localhost:8888")
	err := http.ListenAndServe(":8888", handler)
	if err != nil {
		log.Fatal(err)
	}
}
