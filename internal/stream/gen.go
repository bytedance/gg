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

// This program generates *_gen.go. see gen.sh for usage.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/bytedance/gg/gcond"
	"github.com/bytedance/gg/gresult"
	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/iter"
	"github.com/bytedance/gg/internal/rtassert"
	"github.com/bytedance/gg/internal/stream"
)

var (
	parent      string
	parentTypes string
	parentElem  string
	child       string
	childTypes  string
	childElem   string
	importPaths string
	ignorePaths string
	ignoreFuncs string
)

func init() {
	// Init command line flags.

	flag.StringVar(&parent, "parent", "", "name of type to be inherited, if you want to inherit from Stream, use empty string")
	flag.StringVar(&parentTypes, "parent-types", "", "types constraint of parent")
	flag.StringVar(&parentElem, "parent-elem", "", "element type of parent")
	flag.StringVar(&child, "child", "", "name of new type")
	flag.StringVar(&childTypes, "child-types", "", "type constraints of child")
	flag.StringVar(&childElem, "child-elem", "", "element type of child")
	flag.StringVar(&importPaths, "import-paths", "", "space-separated import paths")
	flag.StringVar(&ignorePaths, "ignore-paths", "", "space-separated ignored import paths")
	flag.StringVar(&ignoreFuncs, "ignore-funcs", "", "space-separated ignored functions/methods")
}

var (
	add = gvalue.Add[string]
)

func main() {
	flag.Parse()

	if parent == "" || parentTypes == "" || child == "" || childTypes == "" {
		fmt.Fprintln(os.Stderr, "value of -parent or -parent-types or -child or -child-types can not be empty")
		os.Exit(1)
	}

	if parentElem == "" {
		parentElem = "T"
	}
	if childElem == "" {
		childElem = "T"
	}

	buf := new(bytes.Buffer)
	g := newGenerator()
	g.Write(buf)
	// fmt.Println(string(buf.Bytes()))
	data := gresult.Of(format.Source(buf.Bytes())).Value()
	dst := strings.ToLower(child) + "_gen.go"
	rtassert.ErrMustNil(ioutil.WriteFile(dst, data, 0660))

	fmt.Printf("generated %q from %q\n", dst, iter.ToSlice(iter.FromMapKeys(g.src)))
}

type generator struct {
	parentName  string
	parentTypes []string
	parentElem  string
	childName   string
	childTypes  []string
	childElem   string
	importPaths []string
	ignorePaths []string
	ignoreFuncs []string
	src         map[string]string
	fset        *token.FileSet
	gast        *ast.File // AST of auto-generated code
	hast        *ast.File // AST of handwritten code
}

func newGenerator() *generator {
	fset := token.NewFileSet()

	// Reaed generated  source code.
	gfile := strings.ToLower(parent) + "_gen.go"
	gsrc := string(gresult.Of(ioutil.ReadFile(gfile)).Value())
	// Create the AST by parsing src.
	gnode := gresult.Of(parser.ParseFile(fset, gfile, nil, parser.ParseComments)).Value()

	// Reaed handwritten source code.
	hfile := strings.ToLower(parent) + ".go"
	hsrc := string(gresult.Of(ioutil.ReadFile(hfile)).Value())
	// Create the AST by parsing src.
	hnode := gresult.Of(parser.ParseFile(fset, hfile, nil, parser.ParseComments)).Value()

	g := &generator{
		parentName:  parent,
		parentTypes: strings.Split(parentTypes, " "),
		parentElem:  parentElem,
		childName:   child,
		childTypes:  strings.Split(childTypes, " "),
		childElem:   childElem,
		src: map[string]string{
			hfile: hsrc,
			gfile: gsrc,
		},
		fset: fset,
		gast: gnode,
		hast: hnode,
	}
	if importPaths != "" {
		g.importPaths = strings.Split(importPaths, " ")
	}
	if ignorePaths != "" {
		g.ignorePaths = strings.Split(ignorePaths, " ")
	}
	if ignoreFuncs != "" {
		g.ignoreFuncs = strings.Split(ignoreFuncs, " ")
	}

	return g
}

func (c *generator) nodeToString(n ast.Node) string {
	file := c.fset.File(n.Pos()).Name()
	start := c.fset.Position(n.Pos()).Offset
	stop := c.fset.Position(n.End()).Offset
	return c.src[file][start:stop]
}

func (c *generator) parentType() string {
	return c.parentName + "[" + c.toTypeParams(c.parentTypes) + "]"
}

func (c *generator) childType() string {
	return c.childName + "[" + c.toTypeParams(c.childTypes) + "]"
}

func (c *generator) writeImport(w io.Writer) {
	paths := make(map[string]struct{})

	collectImport := func(n ast.Node) bool {
		if n == nil {
			return true
		}
		if _, ok := n.(*ast.ImportSpec); ok {
			p := gresult.Of(strconv.Unquote(c.nodeToString(n))).Value()
			paths[p] = struct{}{}
		}
		return true
	}
	if c.gast != nil {
		ast.Inspect(c.gast, collectImport)
	}
	ast.Inspect(c.hast, collectImport)

	for _, p := range c.importPaths {
		paths[p] = struct{}{}
	}
	for _, p := range c.ignorePaths {
		if _, ok := paths[p]; ok {
			delete(paths, p)
		}
	}

	fmt.Fprintln(w, "import (")
	for _, v := range stream.FromStringMapKeys(paths).Sort().ToSlice() {
		fmt.Fprintln(w, strconv.Quote(v))
	}
	fmt.Fprintln(w, ")")
}

func (c *generator) writeStruct(w io.Writer) {
	ctx := map[string]string{
		"Parent":           c.parentName,
		"Child":            c.childName,
		"Constraints":      c.toTypeConstraints(c.childTypes),
		"ParentTypeParams": c.toTypeParamsOfParent(),
		"Variant": gcond.If(c.childElem == "T",
			c.childTypes[1],
			c.childElem),
	}
	tmpl := gresult.Of(template.New("struct").Parse(`
	// {{.Child}} is a {{.Variant}} variant of {{.Parent}}.
	type {{.Child}}[{{.Constraints}}] struct {
		{{.Parent}}[{{.ParentTypeParams}}]
	}
	`)).Value()
	rtassert.ErrMustNil(tmpl.Execute(w, ctx))
}

func (c *generator) writeMethods(w io.Writer) {
	rewriteMethod := func(n ast.Node) bool {
		if n == nil {
			return true
		}
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		for _, f := range c.ignoreFuncs {
			if f == funcDecl.Name.Name {
				return true
			}
		}
		//
		// Write function sig.
		//

		funcName := funcDecl.Name.Name
		funcRecv := ""
		if funcDecl.Recv == nil {
			// Replace function name because it is already used.
			rewriteFuncName := func(funcName, prefix string) string {
				if strings.HasPrefix(funcName, prefix+c.parentName) {
					funcName = strings.Replace(funcName, prefix+c.parentName, prefix+c.childName, 1)
				} else if strings.HasPrefix(funcName, prefix) {
					funcName = strings.Replace(funcName, prefix, prefix+c.childName, 1)
				}
				return funcName
			}
			funcName = rewriteFuncName(funcName, "From")
			funcName = rewriteFuncName(funcName, "Steal")
			funcName = rewriteFuncName(funcName, "Repeat")
		} else {
			// Rewrite Receiver
			funcRecv = strings.Replace(c.nodeToString(funcDecl.Recv), c.parentType(), c.childType(), 1)
		}

		// Write function/method signature: replace type name & struct name.
		funcTypeParams := ""
		if funcDecl.Type.TypeParams != nil {
			funcTypeParams = strings.Replace(c.nodeToString(funcDecl.Type.TypeParams),
				c.toTypeConstraints(c.parentTypes),
				c.toTypeConstraints(c.childTypes), 1)
		}

		var needGen bool
		var funcParams []string
		for _, p := range funcDecl.Type.Params.List {
			ps := c.nodeToString(p)
			if strings.Contains(ps, c.parentType()) {
				needGen = true
				ps = strings.ReplaceAll(ps, c.parentType(), c.childType())
			} else {
				ps = strings.ReplaceAll(ps, c.parentElem, c.childElem)
			}
			funcParams = append(funcParams, ps)
		}

		var funcResults []string
		if funcDecl.Type.Results != nil {
			for _, r := range funcDecl.Type.Results.List {
				rs := c.nodeToString(r)
				if strings.Contains(rs, c.parentType()) {
					needGen = true
					rs = strings.ReplaceAll(rs, c.parentType(), c.childType())
				} else {
					rs = strings.ReplaceAll(rs, c.parentElem, c.childElem)
				}
				funcResults = append(funcResults, rs)
			}
		}

		// If the receiver type does not appear in params and results of functions,
		// no need to rewrite it.
		//
		// For example, generating *Comparable[T] from *Stream[T]:
		//
		// - Map(f func(T) T) *Stream[T]: Should be rewritten to Map(f func(T) T) *Comparable[T]
		// - Reduce(f func(T, T) T) goption.O[T]: Does not need to be rewritten
		if !needGen {
			return true
		}

		// Write function document.
		if funcDecl.Doc != nil {
			fmt.Fprintln(w, c.nodeToString(funcDecl.Doc))
		}

		sig := "func " + funcRecv + funcName + funcTypeParams +
			"(" + strings.Join(funcParams, ",") + ")" +
			"(" + strings.Join(funcResults, ",") + ")"
		fmt.Fprint(w, sig)

		//
		// Write function body.
		//

		var body []string
		body = append(body, "{")

		var params []string
		for _, f := range funcDecl.Type.Params.List {
			for _, name := range f.Names {
				var param string
				if typ := c.nodeToString(f.Type); typ == c.parentType() {
					param = name.Name + "." + c.parentName
				} else if typ == "..."+c.parentType() {
					body = append(body, `conv := func(c {{.Child}}[{{.TypeParams}}]) {{.Parent}}[{{.ParentTypeParams}}] {`)
					body = append(body, `	return c.{{.Parent}}`)
					body = append(body, `}`)
					body = append(body, `tmp := iter.ToSlice(iter.Map(conv, iter.FromSlice(`+name.Name+`)))`)
					param = "tmp..."
				} else {
					param = name.Name
				}
				params = append(params, param)
			}
		}

		// Prepare template context.
		ctx := map[string]string{
			"Parent":           c.parentName,
			"ParentTypeParams": c.toTypeParamsOfParent(),
			"Child":            c.childName,
			"Name":             c.nodeToString(funcDecl.Name),
			"Params":           strings.Join(params, ","),
			"TypeParams":       c.toTypeParams(c.childTypes),
		}
		if funcDecl.Recv != nil { // For method: Call parent's corresponding method
			ctx["Receiver"] = funcDecl.Recv.List[0].Names[0].Name

			if noResult := funcDecl.Type.Results == nil; noResult {
				// Only execute parent's method
				body = append(body, `{{.Receiver}}.{{.Parent}}.{{.Name}}({{.Params}})`)
			} else if noNeedWrap := !strings.Contains(c.nodeToString(funcDecl.Type.Results), c.parentType()); noNeedWrap {
				// Execute parent's method and return the result
				body = append(body, `return {{.Receiver}}.{{.Parent}}.{{.Name}}({{.Params}})`)
			} else {
				// Execute parent's method and return the wrapped result
				body = append(body, `return {{.Child}}[{{.TypeParams}}] { {{.Receiver}}.{{.Parent}}.{{.Name}}({{.Params}}) }`)
			}
		} else { // For function: Call parent's corresponding function
			// Execute parent's function and return the wrapped result
			body = append(body, `return {{.Child}}[{{.TypeParams}}] { {{.Name}}({{.Params}}) }`)
		}

		body = append(body, "}")
		tmpl := gresult.Of(template.New("struct").Parse(strings.Join(body, "\n"))).Value()
		rtassert.ErrMustNil(tmpl.Execute(w, ctx))

		// Break lines.
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		return true
	}

	if c.gast != nil {
		ast.Inspect(c.gast, rewriteMethod)
	}
	ast.Inspect(c.hast, rewriteMethod)
}

func (c *generator) Write(w io.Writer) {
	fmt.Fprintln(w, "// code generated by go run gen.go; DO NOT EDIT.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "package stream")
	fmt.Fprintln(w)

	c.writeImport(w)
	c.writeStruct(w)
	c.writeMethods(w)
}

// ["K", "comparable", "V", "any"] => "K comparable, V any"
func (c *generator) toTypeConstraints(s []string) string {
	return stream.FromStringSlice(s).
		Filter(isTypeParam()).
		Zip(add, stream.RepeatString(" ")).
		Zip(add,
			stream.FromStringSlice(s).
				Filter(isTypeConstraint())).
		Join(", ")
}

// ["K", "comparable", "V", "any"] => "K, V"
func (c *generator) toTypeParams(s []string) string {
	return stream.FromStringSlice(s).
		Filter(isTypeParam()).
		Join(", ")
}

func (c *generator) toTypeParamsOfParent() string {
	return gcond.If(len(c.parentTypes) == len(c.childTypes),
		c.toTypeParams(c.childTypes),
		c.childElem)
}

func isTypeParam() func(string) bool {
	var i = 0
	return func(string) bool {
		is := i%2 == 0
		i++
		return is
	}
}

func isTypeConstraint() func(string) bool {
	var i = 0
	return func(string) bool {
		is := i%2 == 1
		i++
		return is
	}
}
