// Copyright 2025 Bytedance Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build ignore
// +build ignore

package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"log"
	"os"
	"text/template"
)

// Inspired by sort/gen_sort_variants.go
type Variant struct {
	// Package is the package name.
	Package string

	// Name is the variant name: should be unique among variants.
	Name string

	// Path is the file path into which the generator will emit the code for this
	// variant.
	Path string

	// Imports is the imports needed for this package.
	Imports string

	StructPrefix    string
	StructPrefixLow string
	StructSuffix    string
	ExtraFields     string

	// Basic type. T or "".
	Type string

	// Basic type argument. [T] or "".
	TypeArgument string

	// TypeParam is the optional type parameter for the function.
	TypeParam string // e.g. [T any]

	// Funcs is a map of functions used from within the template. The following
	// functions are expected to exist:
	Funcs template.FuncMap
}

type TypeReplacement struct {
	Type string
	Desc string
}

func main() {
	// For New.
	base := &Variant{
		Package:         "skipset",
		Name:            "ordered",
		Path:            "gen_ordered.go",
		Imports:         "\"sync\"\n\"sync/atomic\"\n\"unsafe\"\n\n\"github.com/bytedance/gg/internal/constraints\"\n",
		Type:            "T",
		TypeArgument:    "[T]",
		TypeParam:       "[T constraints.Ordered]",
		StructPrefix:    "Ordered",
		StructPrefixLow: "ordered",
		StructSuffix:    "",
		Funcs: template.FuncMap{
			"Less": func(i, j string) string {
				return fmt.Sprintf("(%s < %s)", i, j)
			},
			"Equal": func(i, j string) string {
				return fmt.Sprintf("%s == %s", i, j)
			},
		},
	}
	generate(base)
	base.Name += "Desc"
	base.StructSuffix += "Desc"
	base.Path = "gen_ordereddesc.go"
	base.Funcs = template.FuncMap{
		"Less": func(i, j string) string {
			return fmt.Sprintf("(%s > %s)", i, j)
		},
		"Equal": func(i, j string) string {
			return fmt.Sprintf("%s == %s", i, j)
		},
	}
	generate(base)

	// For NewFunc.
	basefunc := &Variant{
		Package:         "skipset",
		Name:            "func",
		Path:            "gen_func.go",
		Imports:         "\"sync\"\n\"sync/atomic\"\n\"unsafe\"\n",
		Type:            "T",
		TypeArgument:    "[T]",
		TypeParam:       "[T any]",
		ExtraFields:     "\nless func(a,b T)bool\n",
		StructPrefix:    "Func",
		StructPrefixLow: "func",
		StructSuffix:    "",
		Funcs: template.FuncMap{
			"Less": func(i, j string) string {
				return fmt.Sprintf("s.less(%s,%s)", i, j)
			},
			"Equal": func(i, j string) string {
				return fmt.Sprintf("!s.less(%s,%s)", j, i)
			},
		},
	}
	generate(basefunc)
}

// generate generates the code for variant `v` into a file named by `v.Path`.
func generate(v *Variant) {
	// Parse templateCode anew for each variant because Parse requires Funcs to be
	// registered, and it helps type-check the funcs.
	tmpl, err := template.New("gen").Funcs(v.Funcs).Parse(templateCode)
	if err != nil {
		log.Fatal("template Parse:", err)
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, v)
	if err != nil {
		log.Fatal("template Execute:", err)
	}

	os.WriteFile(v.Path, out.Bytes(), 0644)

	formatted, err := format.Source(out.Bytes())
	if err != nil {
		log.Fatal("format:", err)
	}

	if err := os.WriteFile(v.Path, formatted, 0644); err != nil {
		log.Fatal("WriteFile:", err)
	}
}

//go:embed skipset.tpl
var templateCode string
