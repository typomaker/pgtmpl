package pgtemplate

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"text/template"
)

func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

type Template struct {
	text *template.Template
}

func New(name string) (tpl *Template) {
	tpl = &Template{}
	tpl.text = template.New(name).Funcs(template.FuncMap{
		"hold": func(v interface{}) string {
			return "$?"
		},
	})
	return tpl
}
func (tpl *Template) New(name string) (t *Template) {
	return &Template{
		text: tpl.text.New(name),
	}
}
func (tpl *Template) Parse(text string) (*Template, error) {
	_, err := tpl.text.Parse(text)
	return tpl, err
}
func (tpl *Template) ParseFiles(filenames ...string) (*Template, error) {
	_, err := tpl.text.ParseFiles(filenames...)
	return tpl, err
}
func ParseFiles(filenames ...string) (t *Template, err error) {
	if len(filenames) == 0 {
		return nil, fmt.Errorf("pgtemplate: no files named in call to ParseFiles")
	}
	rootname := filepath.Base(filenames[0])
	t = New(rootname)
	_, err = t.ParseFiles(filenames...)
	return
}
func (tpl *Template) ParseGlob(pattern string) (*Template, error) {
	_, err := tpl.text.ParseGlob(pattern)
	return tpl, err
}
func ParseGlob(pattern string) (t *Template, err error) {
	var rootname string
	{
		filenames, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		if len(filenames) == 0 {
			return nil, fmt.Errorf("pgtemplate: pattern matches no files: %#q", pattern)
		}
		rootname = filepath.Base(filenames[0])
	}
	t = New(rootname)
	_, err = t.ParseGlob(pattern)
	return
}
func (tpl *Template) ParseFS(fsys fs.FS, patterns ...string) (*Template, error) {
	_, err := tpl.text.ParseFS(fsys, patterns...)
	return tpl, err
}
func ParseFS(fsys fs.FS, patterns ...string) (t *Template, err error) {
	if len(patterns) == 0 {
		return nil, fmt.Errorf("pgtemplate: no patterns in call to ParseFS")
	}
	var rootname string
	{
		filenames, err := fs.Glob(fsys, patterns[0])
		if err != nil {
			return nil, err
		}
		if len(filenames) == 0 {
			return nil, fmt.Errorf("pgtemplate: pattern matches no files: %#q", patterns[0])
		}
		rootname = filepath.Base(filenames[0])
	}
	t = New(rootname)
	_, err = t.ParseFS(fsys, patterns...)
	return
}
func (tpl *Template) Name() string {
	return tpl.text.Name()
}
func (tpl *Template) Lookup(name string) *Template {
	ntext := tpl.text.Lookup(name)
	if ntext == nil {
		return nil
	}
	return &Template{
		text: ntext,
	}
}
func (tpl *Template) Clone() (*Template, error) {
	var err error
	tplcp := &Template{}
	tplcp.text, err = tpl.text.Clone()
	return tplcp, err
}
func (tpl *Template) Funcs(funcMap template.FuncMap) *Template {
	tpl.text.Funcs(funcMap)
	return tpl
}
func (tpl *Template) Execute(q *Query, data interface{}) (err error) {
	if tpl, err = tpl.Clone(); err != nil {
		return
	}
	if q.Len() != 0 {
		q.WriteString(";")
	}
	tpl.Funcs(template.FuncMap{
		"hold": func(v interface{}) string {
			q.args = append(q.args, v)
			return "$" + strconv.Itoa(len(q.args))
		},
	})
	if err = tpl.text.Execute(q, data); err != nil {
		return fmt.Errorf("%q: %w", tpl.Name(), err)
	}
	return err
}
func (tpl *Template) ExecuteTemplate(q *Query, name string, data interface{}) error {
	if t := tpl.Lookup(name); t != nil {
		return t.Execute(q, data)
	} else {
		return fmt.Errorf("pgtemplate: no template %q associated with template %q", name, tpl.Name())
	}
}
