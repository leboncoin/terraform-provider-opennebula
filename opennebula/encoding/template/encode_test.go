package template

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Template struct {
	Attr1 string `template:"ATTR1"`
	Bttr2 string `template:"BTTR2"`
	Cttr3 int    `template:"CTTR3"`
	Dttr4 string
	Ettr5 bool
}

func (d *Template) marshalTemplate() string {
	template, err := MarshalWrap(d)
	if err != nil {
		fmt.Printf("failed to solve problem: %s\n", err)
	}
	return template
}
func (d *Template) marshalTemplateWithoutWrap() string {
	template, err := Marshal(d)
	if err != nil {
		fmt.Printf("failed to solve problem: %s\n", err)
	}
	return template
}
func TestEncodeTemplate(t *testing.T) {
	tpl := &Template{
		Attr1: "test1",
		Bttr2: "test2",
		Cttr3: 1,
	}
	attempt := `TEMPLATE = [ATTR1=test1, BTTR2=test2, CTTR3=1]`
	result := tpl.marshalTemplate()

	assert.Equal(t, result, attempt, "The two string should be the same.")
}

func TestEncodeTemplateWithoutAttr(t *testing.T) {
	tpl := &Template{
		Dttr4: "test4",
	}
	attempt := "TEMPLATE = [DTTR4=test4]"
	result := tpl.marshalTemplate()

	assert.Equal(t, result, attempt, "The two string should be the same.")
}

func TestEncodeTemplateWithEmptyValue(t *testing.T) {
	tpl := &Template{
		Attr1: "",
		Cttr3: 0,
	}
	attempt := "TEMPLATE = []"
	result := tpl.marshalTemplate()

	assert.Equal(t, result, attempt, "The two string should be the same.")
}

func TestEncodeTemplateWithUnavailableType(t *testing.T) {
	tpl := &Template{
		Ettr5: true,
	}
	attempt := "TEMPLATE = []"
	result := tpl.marshalTemplate()

	assert.Equal(t, result, attempt, "The two string should be the same.")
}

func TestEncodeTemplateWithoutWrapper(t *testing.T) {
	tpl := &Template{
		Attr1: "test1",
		Bttr2: "test2",
		Cttr3: 1,
	}
	attempt := "ATTR1=test1, BTTR2=test2, CTTR3=1"
	result := tpl.marshalTemplateWithoutWrap()

	assert.Equal(t, result, attempt, "The two string should be the same.")
}
