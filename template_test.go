package pgtmpl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typomaker/pgtmpl"
)

func TestTemplate_new(t *testing.T) {
	q := pgtmpl.Query{}
	tpl := pgtmpl.Must(pgtmpl.New("T1").Parse(`T1 {{template "T2"}}`))
	pgtmpl.Must(tpl.New("T2").Parse("T2"))
	err := tpl.Execute(&q, nil)
	assert.NoError(t, err)
	assert.Equal(t, "T1 T2", q.String())
}
func TestTemplate_name(t *testing.T) {
	tpl := pgtmpl.New("T2")
	assert.Equal(t, "T2", tpl.Name())
}
func TestTemplate_concatQuery(t *testing.T) {
	q := pgtmpl.Query{}
	tpl := pgtmpl.New("")
	pgtmpl.Must(tpl.New("Q1").Parse("SELECT 1, {{hold .}}"))
	pgtmpl.Must(tpl.New("Q2").Parse("SELECT 2, {{hold .}}"))

	err := tpl.ExecuteTemplate(&q, "Q1", 1)
	assert.NoError(t, err)
	err = tpl.ExecuteTemplate(&q, "Q2", 2)
	assert.NoError(t, err)

	assert.Equal(t, "SELECT 1, $1;SELECT 2, $2", q.String())
	assert.Len(t, q.Args(), 2)
	assert.Equal(t, q.Args(), []interface{}{1, 2})
}
func TestTemplate_resetQuery(t *testing.T) {
	q := pgtmpl.Query{}
	tpl := pgtmpl.New("")
	pgtmpl.Must(tpl.New("Q1").Parse("SELECT 1, {{hold .}}"))
	pgtmpl.Must(tpl.New("Q2").Parse("SELECT 2, {{hold .}}"))

	err := tpl.ExecuteTemplate(&q, "Q1", 1)
	assert.NoError(t, err)
	q.Reset()
	err = tpl.ExecuteTemplate(&q, "Q2", 2)
	assert.NoError(t, err)

	assert.Equal(t, "SELECT 2, $1", q.String())
	assert.Len(t, q.Args(), 1)
	assert.Equal(t, q.Args(), []interface{}{2})
}
