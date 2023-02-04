package main

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"text/template"

	"github.com/asmsh/lazy"
)

// helloTmpl will panic if it fails to load the waned value,
// hence no need to check its Err() return.
var helloTmpl = lazy.NewValue(func() (*template.Template, error) {
	tmplTxt := `Hello {{.}}!`
	tmpl := template.Must(template.New("hello").Parse(tmplTxt))
	return tmpl, nil
})

func main() {
	users := []string{
		"UserA",
		"UserB",
		"UserC",
	}

	wg := sync.WaitGroup{}
	wg.Add(len(users))

	for _, v := range users {
		user := v

		// greet each user concurrently, using the same template
		go func() {
			defer wg.Done()

			b := bytes.Buffer{}
			err := helloTmpl.Val().Execute(&b, user)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(b.String())
		}()
	}

	wg.Wait()
}
