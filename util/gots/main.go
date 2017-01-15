// package gots (or go-ts) is responsible for converting structs into
// typescript models. Has limited support for protobuf generated structs,
// specifically `oneof` and `enum` types

package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	router_api "github.com/octavore/nagax/proto/nagax/router/api"

	"github.com/octavore/press/proto/press/api"
	"github.com/octavore/press/proto/press/models"
	"github.com/octavore/press/proto/press/packages"
)

func main() {
	out := os.Args[1]
	f, err := os.OpenFile(out, os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("error opening file:", err)
	}

	g := newGenerator(f)
	g.RegisterMany(
		models.Page{},
		models.Content{},
		models.Route{},
		models.Timestamp{},
		models.Theme{},
		models.ThemeTemplate{},
		models.ThemePlaceholder{},
		models.ThemeAsset{},
		models.ContentMultiple{},
		models.ContentText{},
		models.ContentString{},
		packages.Package{},
		packages.PackageRelease{},
		packages.Registry{},

		api.TLSSettingsReponse{},
		api.EnableTLSRequest{},

		router_api.Error{},
		router_api.ErrorResponse{},
	)
	g.Write()
	g.Close()
}

func newGenerator(w io.WriteCloser) *Generator {
	return &Generator{
		out:   w,
		enums: map[string]string{},
	}
}

type Generator struct {
	models []reflect.Type
	enums  map[string]string
	out    io.WriteCloser
}

func (g *Generator) RegisterMany(l ...interface{}) {
	for _, i := range l {
		g.Register(i)
	}
}
func (g *Generator) Register(i interface{}) {
	v := reflect.ValueOf(i).Type()
	if v.Kind() != reflect.Struct {
		panic("can only register struct types")
	}
	g.models = append(g.models, v)
}

func (g *Generator) Write() {
	g.p(0, "// DO NOT EDIT! This file is generated automatically by util/gots/main.go\n")
	for _, i := range g.models {
		g.convert(i)
	}
	for t, e := range g.enums {
		g.convertEnum(t, e)
	}
}

func (g *Generator) Close() {
	_ = g.out.Close()
}

func (g *Generator) p(indent int, s string) {
	spaces := strings.Repeat(" ", indent)
	fmt.Fprint(g.out, spaces, s, "\n")
}

func (g *Generator) convertEnum(typeName, enumName string) {
	enumMap := proto.EnumValueMap(enumName)
	enums := []string{}
	for enum := range enumMap {
		enums = append(enums, fmt.Sprintf("'%s'", enum))
	}
	if len(enums) > 0 {
		g.p(0, fmt.Sprintf("export type %s = %s;", typeName, strings.Join(enums, " | ")))
	}
}

const typeFmt = "%s?: %s;"

func (g *Generator) subconvertFields(v reflect.Type) []string {
	fields := []string{}
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)

		// name
		name := tsFieldname(f)
		if name == "-" {
			continue
		}

		// type
		typ := g.goTypeToTSType(f.Type, &f.Tag)
		if typ != "" {
			fields = append(fields, name)
			g.p(2, fmt.Sprintf(typeFmt, name, typ))
		} else {
			g.p(2, "// skipped field: "+name)
		}
	}
	return fields
}

func (g *Generator) convert(v reflect.Type) {
	g.p(0, "export abstract class "+v.Name()+" {")
	fields := g.subconvertFields(v)

	// handle oneof fields
	sp := proto.GetProperties(v)
	if len(sp.OneofTypes) > 0 {
		g.p(2, "")
		g.p(2, "// oneof types:")
		for _, prop := range sp.OneofTypes {
			// merge oneof fields into parent
			f2 := g.subconvertFields(prop.Type.Elem())
			fields = append(fields, f2...)
		}
	}
	g.generateCopyFunction(v.Name(), fields)
	g.p(0, "}\n")
}

// helper function to extract subtags, e.g. `protobuf:"json=foo"`
func lookupSubTag(tag reflect.StructTag, tagName, subTag string) (string, bool) {
	t, ok := tag.Lookup(tagName)
	if !ok {
		return "", false
	}
	tParts := strings.Split(t, ",")
	prefix := subTag + "="
	for _, part := range tParts {
		if strings.HasPrefix(part, prefix) {
			return strings.TrimPrefix(part, prefix), true
		}
	}
	return "", false
}

// extract the field name from the field. prefers protobuf
// declared json name if it exists.
func tsFieldname(f reflect.StructField) string {
	proto, ok := lookupSubTag(f.Tag, "protobuf", "json")
	if ok {
		return proto
	}
	json, ok := f.Tag.Lookup("json")
	if ok {
		return strings.Split(json, ",")[0]
	}
	return strings.ToLower(f.Name)
}

// converts native go types to native ts types
var typeMap = map[string]string{
	"int32": "number",
	"int64": "string",
	"bool":  "boolean",
}

type protoEnum interface {
	EnumDescriptor() ([]byte, []int)
}

var protoEnumType = reflect.TypeOf((*protoEnum)(nil)).Elem()

// convert a go type to a TS type.
// note: protobuf "oneof" is not supported
func (g *Generator) goTypeToTSType(t reflect.Type, tag *reflect.StructTag) string {
	if tag != nil {
		// keep track of enums for later generation
		// AssignableTo is not strictly speaking necessary, rather it is a
		// helper to avoid unnecessary tag checks.
		if t.Name() != "" && t.AssignableTo(protoEnumType) {
			enum, ok := lookupSubTag(*tag, "protobuf", "enum")
			if ok {
				g.enums[t.Name()] = enum
			}
		}

		// do not generate oneof types
		if _, ok := tag.Lookup("protobuf_oneof"); ok {
			return ""
		}
	}

	switch t.Kind() {
	case reflect.Ptr:
		return g.goTypeToTSType(t.Elem(), tag)
	case reflect.Slice:
		typ := g.goTypeToTSType(t.Elem(), tag)
		typ += "[]"
		return typ
	case reflect.Struct:
		return t.Name()
	case reflect.Interface:
		return "any"
	case reflect.Map:
		return fmt.Sprintf("{ [key: %s]: %s; }",
			g.goTypeToTSType(t.Key(), tag),
			g.goTypeToTSType(t.Elem(), tag))
	default:
		typ := t.Name()
		if alt, ok := typeMap[typ]; ok {
			return alt
		}
		return typ
	}
}

func (g *Generator) generateCopyFunction(class string, fields []string) {
	g.p(2, fmt.Sprintf("static copy(from: %s, to?: %s): %s {", class, class, class))
	g.p(4, "to = to || {};")
	for _, field := range fields {
		g.p(4, fmt.Sprintf("to.%s = from.%s;", field, field))
	}
	g.p(4, "return to;")
	g.p(2, "}")
}
