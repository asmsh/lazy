# Lazy

A lazy value loader, it can be used to load expensive computation only when it is needed 
and cache the result for future use.

## Installation

```bash
go get github.com/asmsh/lazy
```

## Example

```go
package main

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"

	"github.com/asmsh/lazy"
)

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

		go func() {
			defer wg.Done()

			b := bytes.Buffer{}
			helloTmpl.Val().Execute(&b, user)
			fmt.Println(b.String())
		}()
	}

	wg.Wait()
}
```
