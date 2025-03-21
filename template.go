package pgtmpl

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"
	"text/template"
)

func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

type (
	Template struct {
		mu      sync.RWMutex
		text    *template.Template
		funcmap FuncMap
	}
	FuncMap = template.FuncMap
)

const (
	FuncHold = "hold"
)

var holderPool = [4096]string{}

func init() {
	for i := range holderPool {
		holderPool[i] = "$" + strconv.Itoa(i+1)
	}
}

func New(name string) (tpl *Template) {
	tpl = &Template{}
	tpl.funcmap = template.FuncMap{
		FuncHold: func(v any) any {
			return v
		},
	}
	tpl.text = template.New(name).Funcs(tpl.funcmap)
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
		return nil, fmt.Errorf("template: no files named in call to ParseFiles")
	}
	rootname := filepath.Base(filenames[0])
	t = New(rootname)
	_, err = t.ParseFiles(filenames...)
	return t, err
}
func (tpl *Template) ParseGlob(pattern string) (*Template, error) {
	_, err := tpl.text.ParseGlob(pattern)
	return tpl, err
}
func ParseGlob(pattern string) (t *Template, err error) {
	var rootname string
	{
		var filenames []string
		filenames, err = filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		if len(filenames) == 0 {
			return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
		}
		rootname = filepath.Base(filenames[0])
	}
	t = New(rootname)
	_, err = t.ParseGlob(pattern)
	return t, err
}
func (tpl *Template) ParseFS(fsys fs.FS, patterns ...string) (*Template, error) {
	_, err := tpl.text.ParseFS(fsys, patterns...)
	return tpl, err
}
func ParseFS(fsys fs.FS, patterns ...string) (t *Template, err error) {
	if len(patterns) == 0 {
		return nil, fmt.Errorf("template: no patterns in call to ParseFS")
	}
	var rootname string
	{
		var filenames []string
		filenames, err = fs.Glob(fsys, patterns[0])
		if err != nil {
			return nil, err
		}
		if len(filenames) == 0 {
			return nil, fmt.Errorf("template: pattern matches no files: %#q", patterns[0])
		}
		rootname = filepath.Base(filenames[0])
	}
	t = New(rootname)
	_, err = t.ParseFS(fsys, patterns...)
	return t, err
}
func (tpl *Template) Name() string {
	return tpl.text.Name()
}
func (tpl *Template) Lookup(name string) *Template {
	ntext := tpl.text.Lookup(name)
	if ntext == nil {
		return nil
	}
	tpl.mu.RLock()
	defer tpl.mu.RUnlock()
	return &Template{
		funcmap: tpl.funcmap,
		text:    ntext,
	}
}
func (tpl *Template) Clone() (*Template, error) {
	var err error
	tpl.mu.RLock()
	tplcp := &Template{funcmap: tpl.funcmap}
	tpl.mu.RUnlock()
	tplcp.text, err = tpl.text.Clone()
	return tplcp, err
}
func (tpl *Template) Funcs(funcmap FuncMap) *Template {
	tpl.mu.Lock()
	for k, v := range funcmap {
		tpl.funcmap[k] = v
	}
	tpl.mu.Unlock()
	tpl.text.Funcs(funcmap)
	return tpl
}
func (tpl *Template) Execute(q *Query, data interface{}) (err error) {
	var selftpl *Template
	if selftpl, err = tpl.Clone(); err != nil {
		return err
	}
	if q.Cap() == 0 {
		*q = *NewQuery()
	} else if q.Len() != 0 {
		_, _ = q.WriteString(";")
	}
	q.name = tpl.Name()
	if err = selftpl.text.Funcs(tpl.framefunc(q)).Execute(q, data); err != nil {
		return fmt.Errorf("%q: %w", tpl.Name(), err)
	}
	return err
}
func (tpl *Template) ExecuteTemplate(q *Query, name string, data interface{}) error {
	if t := tpl.Lookup(name); t != nil {
		return t.Execute(q, data)
	}
	return fmt.Errorf("template: no template %q associated with template %q", name, tpl.Name())
}
func (tpl *Template) framefunc(q *Query) FuncMap {
	tpl.mu.RLock()
	holdfunc := tpl.funcmap[FuncHold]
	tpl.mu.RUnlock()

	return FuncMap{
		FuncHold: func(v any) (holder any, err error) {
			if holdfunc != nil {
				if v, err = tpl.callfunc(holdfunc, v); err != nil {
					return "", err
				}
			}
			q.args = append(q.args, v)
			var i = len(q.args)
			if i <= len(holderPool) {
				return holderPool[i-1], nil
			}
			return "$" + strconv.Itoa(len(q.args)), nil
		},
	}
}
func (tpl *Template) callfunc(fn any, in ...any) (out any, err error) {
	var ok bool
	var fnin = make([]reflect.Value, len(in))
	for i := range in {
		fnin[i] = reflect.ValueOf(in[i])
	}
	var fnout = reflect.ValueOf(fn).Call(fnin)

	switch len(fnout) {
	case 0:
		return out, err
	case 1:
		out = fnout[0].Interface()
	case 2: //nolint:gomnd // binary return
		out = fnout[0].Interface()
		if err, ok = fnout[1].Interface().(error); !ok {
			err = fmt.Errorf("template: 2nd returned value must be an error type, got %+v", fnout[1].Interface())
		}
	default:
		err = fmt.Errorf("template: bad number of returned values")
	}
	if err != nil {
		return nil, err
	}
	return out, nil
}
