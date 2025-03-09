package email

import (
	"bytes"
	"html/template"
	"os"
)

func parseTemplate(name string) (*template.Template, error) {
	path := "email/templates/" + name + ".html"

	content, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(name).Parse(string(content))

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func renderTemplate(tmpl *template.Template, data map[string]string) (string, error) {
	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
