package templates

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/sirupsen/logrus"
)

type TemplateName string

const (
	TemplateNameQuotes TemplateName = "quotes"
)

func CompileTemplate(
	tmpl string,
	params interface{},
) (string, error) {
	templateInstance, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = templateInstance.Execute(&buf, params)
	if err != nil {
		err = fmt.Errorf("CompileTemplate: %v", err)
		logrus.Error(err)
		return "", err
	}

	return buf.String(), nil
}
