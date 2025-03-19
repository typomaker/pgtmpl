// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pgtmpl_test

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/typomaker/pgtmpl"
)

func ExampleTemplate_files() {
	// Here we create a temporary directory and populate it with our sample
	// template definition files; usually the template files would already
	// exist in some location known to the program.
	dir := makeTemporaryFile(temporaryFile{
		"allinone",
		`{{define "zerotime"}}COALESCE({{ . }}::timestamp, '0001-01-01 00:00:00'::timestamp){{end -}}
		{{define "nulltime"}}NULLIF({{ . }}, '0001-01-01 00:00:00'::timestamp){{end -}}
		{{define "zerouuid"}}COALESCE({{ . }}::uuid, '00000000-0000-0000-0000-000000000000'::uuid){{end -}}
		{{define "nulluuid"}}NULLIF({{ . }}, '00000000-0000-0000-0000-000000000000'::uuid){{end -}}
		INSERT INTO author(id, created) VALUES({{template "nulluuid" (hold .ID)}}, {{template "nulltime" (hold .Created)}}) RETURNING {{template "zerouuid" "id"}}, {{template "zerotime" "created"}}`,
	})
	defer os.RemoveAll(dir)

	pattern := filepath.Join(dir, "allinone")

	tmpl := pgtmpl.Must(pgtmpl.ParseFiles(pattern))
	query := pgtmpl.Query{}
	data := map[string]interface{}{
		"ID":      "553d6085-0b28-4f2c-a018-d0f34b03b9e7",
		"Created": time.Time{},
	}
	err := tmpl.Execute(&query, data)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	printQuery(query)
	// Output:
	// TEXT: INSERT INTO author(id, created) VALUES(NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), NULLIF($2, '0001-01-01 00:00:00'::timestamp)) RETURNING COALESCE(id::uuid, '00000000-0000-0000-0000-000000000000'::uuid), COALESCE(created::timestamp, '0001-01-01 00:00:00'::timestamp)
	// ARGUMENT: 553d6085-0b28-4f2c-a018-d0f34b03b9e7 0001-01-01 00:00:00 +0000 UTC
}

func ExampleTemplate_glob() {
	// Here we create a temporary directory and populate it with our sample
	// template definition files; usually the template files would already
	// exist in some location known to the program.
	dir := makeTemporaryDirectory([]temporaryFile{
		{"zerotime", `COALESCE({{ . }}::timestamp, '0001-01-01 00:00:00'::timestamp)`},
		{"nulltime", `NULLIF({{ . }}, '0001-01-01 00:00:00'::timestamp)`},
		{"zerouuid", `COALESCE({{ . }}::uuid, '00000000-0000-0000-0000-000000000000'::uuid)`},
		{"nulluuid", `NULLIF({{ . }}, '00000000-0000-0000-0000-000000000000'::uuid)`},
		{"0_upsert", `INSERT INTO author(id, created) VALUES({{template "nulluuid" (hold .ID)}}, {{template "nulltime" (hold .Created)}}) RETURNING {{template "zerouuid" "id"}}, {{template "zerotime" "created"}}`},
	})
	defer os.RemoveAll(dir)

	pattern := filepath.Join(dir, "*")
	tpl := pgtmpl.Must(pgtmpl.ParseGlob(pattern))
	query := pgtmpl.Query{}
	data := map[string]interface{}{
		"ID":      "553d6085-0b28-4f2c-a018-d0f34b03b9e7",
		"Created": time.Time{},
	}
	err := tpl.Execute(&query, data)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	printQuery(query)
	// Output:
	// TEXT: INSERT INTO author(id, created) VALUES(NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), NULLIF($2, '0001-01-01 00:00:00'::timestamp)) RETURNING COALESCE(id::uuid, '00000000-0000-0000-0000-000000000000'::uuid), COALESCE(created::timestamp, '0001-01-01 00:00:00'::timestamp)
	// ARGUMENT: 553d6085-0b28-4f2c-a018-d0f34b03b9e7 0001-01-01 00:00:00 +0000 UTC
}

func ExampleTemplate_fs() {
	dir := makeTemporaryDirectory([]temporaryFile{
		{"zerotime", `COALESCE({{ . }}::timestamp, '0001-01-01 00:00:00'::timestamp)`},
		{"nulltime", `NULLIF({{ . }}, '0001-01-01 00:00:00'::timestamp)`},
		{"zerouuid", `COALESCE({{ . }}::uuid, '00000000-0000-0000-0000-000000000000'::uuid)`},
		{"nulluuid", `NULLIF({{ . }}, '00000000-0000-0000-0000-000000000000'::uuid)`},
		{"0_upsert", `INSERT INTO author(id, created) VALUES({{template "nulluuid" (hold .ID)}}, {{template "nulltime" (hold .Created)}}) RETURNING {{template "zerouuid" "id"}}, {{template "zerotime" "created"}}`},
	})
	defer os.RemoveAll(dir)

	tpl := pgtmpl.Must(pgtmpl.ParseFS(os.DirFS(dir), "*"))
	query := pgtmpl.Query{}
	data := map[string]interface{}{
		"ID":      "553d6085-0b28-4f2c-a018-d0f34b03b9e7",
		"Created": time.Time{},
	}
	err := tpl.Execute(&query, data)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	printQuery(query)
	// Output:
	// TEXT: INSERT INTO author(id, created) VALUES(NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), NULLIF($2, '0001-01-01 00:00:00'::timestamp)) RETURNING COALESCE(id::uuid, '00000000-0000-0000-0000-000000000000'::uuid), COALESCE(created::timestamp, '0001-01-01 00:00:00'::timestamp)
	// ARGUMENT: 553d6085-0b28-4f2c-a018-d0f34b03b9e7 0001-01-01 00:00:00 +0000 UTC
}
