// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/theparanoids/ashirt-server/backend/dtos"
)

func main() {
	fmt.Println("// Copyright 2022, Yahoo Inc.")
	fmt.Println("// Licensed under the terms of the MIT. See LICENSE file in project root for terms.")
	fmt.Println()
	fmt.Println("// Types in this file are generated by backend/dtos/gentypes")
	fmt.Println("// Changes made will be overridden on next generate")
	gen(dtos.APIKey{})
	gen(dtos.Evidence{})
	gen(dtos.EvidenceMetadata{})
	gen(dtos.Finding{})
	gen(dtos.Operation{})
	gen(dtos.Query{})
	gen(dtos.Tag{})
	gen(dtos.DefaultTag{})
	gen(dtos.TagWithUsage{})
	gen(dtos.User{})
	gen(dtos.UserOwnView{})
	gen(dtos.AuthenticationInfo{})
	gen(dtos.UserAdminView{})
	gen(dtos.UserOperationRole{})
	gen(dtos.DetailedAuthenticationInfo{})
	gen(dtos.SupportedAuthScheme{})
	gen(dtos.TagPair{})
	gen(dtos.TagDifference{})
	gen(dtos.TagByEvidenceDate{})
	gen(dtos.FindingCategory{})
	gen(dtos.CheckConnection{})
	gen(dtos.NewUserCreatedByAdmin{})
	gen(dtos.CreateUserOutput{})
	gen(dtos.ServiceWorker{})

	// Since this file only contains typescript types, webpack doesn't pick up the
	// changes unless there is some actual executable javascript referenced from
	// the app itself. By exporting an empty function and calling it in the app
	// https://github.com/TypeStrong/ts-loader/issues/808
	fmt.Println("export function cacheBust() {}")
}

func gen(dtoStruct interface{}) {
	fmt.Println()
	t := reflect.TypeOf(dtoStruct)

	definition := []string{}
	fieldDefinitions := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			definition = append(definition, toTypescriptType(field.Type))
		} else {
			fieldDefinitions = append(fieldDefinitions, genFieldDefinition(field))
		}
	}

	definition = append(definition, "{\n"+strings.Join(fieldDefinitions, "")+"}")
	fmt.Printf("export type %s = ", t.Name())
	fmt.Println(strings.Join(definition, " & "))
}

func genFieldDefinition(field reflect.StructField) string {
	jsonKey := strings.Split(field.Tag.Get("json"), ",")[0]
	if jsonKey != "-" {
		return fmt.Sprintf("  %s: %s,\n", jsonKey, toTypescriptType(field.Type))
	}
	return ""
}

func toTypescriptType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return "number"
	case reflect.String:
		return "string"
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			// []byte is serialized as base64 encoded string
			return "string /* base64 encoded */"
		}
		return "Array<" + toTypescriptType(t.Elem()) + ">"
	case reflect.Ptr:
		return toTypescriptType(t.Elem()) + " | undefined"
	case reflect.Struct:
		if t.PkgPath() == "time" && t.Name() == "Time" {
			return "string"
		}
		if strings.HasSuffix(t.PkgPath(), "ashirt-server/backend/dtos") {
			return t.Name()
		}
		panic(fmt.Errorf("Type from unknown package: %s, (%s)", t.PkgPath(), t.Name()))
	}
	panic(fmt.Errorf("Unknown kind: %s", t.Kind()))
}
