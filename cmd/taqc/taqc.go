package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	g "github.com/moznion/gowrtr/generator"
	"github.com/moznion/taqc/cmd/taqc/internal"
)

var (
	version  string
	revision string
)

func main() {
	var typeName string
	var output string
	var showVersion bool

	flag.StringVar(&typeName, "type", "", "[mandatory] a type name")
	flag.StringVar(&output, "output", "", `[optional] output file name (default "srcdir/<type>_gen.go")`)
	flag.BoolVar(&showVersion, "version", false, "show the version information")

	flag.Parse()

	if showVersion {
		versionJSON, _ := json.Marshal(map[string]string{
			"version":  version,
			"revision": revision,
		})
		fmt.Printf("%s\n", versionJSON)
		return
	}

	args := flag.Args()
	if len(args) <= 0 {
		args = []string{"."}
	}

	pkg, err := internal.ParsePackage(args)
	if err != nil {
		log.Fatal(fmt.Errorf("[error] failed to parse a package: %w", err))
	}

	astFiles, err := internal.ParseFiles(pkg.GoFiles)
	if err != nil {
		log.Fatal(fmt.Errorf("[error] failed to parse a file: %w", err))
	}

	fields, err := internal.CollectQueryParameterFieldsFromAST(typeName, astFiles)
	if err != nil {
		log.Fatal(fmt.Errorf("[error] failed to collect fields from files: %w", err))
	}

	rootStmt := g.NewRoot(
		g.NewComment(fmt.Sprintf(" Code generated by gonstructor %s; DO NOT EDIT.", strings.Join(os.Args[1:], " "))),
		g.NewNewline(),
		g.NewPackage(pkg.Name),
		g.NewNewline(),
		g.NewImport("fmt", "net/url"),
	)

	f := g.NewFunc(g.NewFuncReceiver("v", typeName), g.NewFuncSignature("ToQueryParameters").ReturnTypes("url.Values")).AddStatements(
		g.NewRawStatement("qp := url.Values{}"),
	)
	for _, field := range fields {
		fieldName := field.FieldName
		paramName := field.ParamName
		fieldType := field.FieldType
		switch fieldType {
		case "string":
			f = f.AddStatements(g.NewRawStatementf(`qp.Set("%s", v.%s)`, paramName, fieldName))
		case "int64":
			f = f.AddStatements(g.NewRawStatementf(`qp.Set("%s", fmt.Sprintf("%%d", v.%s))`, paramName, fieldName))
		case "float64":
			f = f.AddStatements(g.NewRawStatementf(`qp.Set("%s", fmt.Sprintf("%%f", v.%s))`, paramName, fieldName))
		case "bool":
			f = f.AddStatements(
				g.NewIf(
					fmt.Sprintf("v.%s", fieldName),
					g.NewRawStatementf(`qp.Set("%s", "1")`, paramName),
				),
			)
		case "*string":
			f = f.AddStatements(
				g.NewIf(
					fmt.Sprintf("v.%s != nil", fieldName),
					g.NewRawStatementf(`qp.Set("%s", *v.%s)`, paramName, fieldName),
				),
			)
		case "*int64":
			f = f.AddStatements(
				g.NewIf(
					fmt.Sprintf("v.%s != nil", fieldName),
					g.NewRawStatementf(`qp.Set("%s", fmt.Sprintf("%%d", *v.%s))`, paramName, fieldName),
				),
			)
		case "*float64":
			f = f.AddStatements(
				g.NewIf(
					fmt.Sprintf("v.%s != nil", fieldName),
					g.NewRawStatementf(`qp.Set("%s", fmt.Sprintf("%%f", *v.%s))`, paramName, fieldName),
				),
			)
		case "*bool":
			f = f.AddStatements(
				g.NewIf(
					fmt.Sprintf("v.%s != nil && *v.%s", fieldName, fieldName),
					g.NewRawStatementf(`qp.Set("%s", "1")`, paramName),
				),
			)
		case "[]string":
			sliceFieldName := strcase.ToLowerCamel(fmt.Sprintf("%s_slice", fieldName))
			f = f.AddStatements(
				g.NewRawStatementf("%s := v.%s", sliceFieldName, fieldName),
				g.NewFor(
					fmt.Sprintf("i := 0; i < len(%s); i++", sliceFieldName),
					g.NewRawStatementf(`qp.Add("%s", %s[i])`, paramName, sliceFieldName),
				),
			)
		case "[]int64":
			sliceFieldName := strcase.ToLowerCamel(fmt.Sprintf("%s_slice", fieldName))
			f = f.AddStatements(
				g.NewRawStatementf("%s := v.%s", sliceFieldName, fieldName),
				g.NewFor(
					fmt.Sprintf("i := 0; i < len(%s); i++", sliceFieldName),
					g.NewRawStatementf(`qp.Add("%s", fmt.Sprintf("%%d", %s[i]))`, paramName, sliceFieldName),
				),
			)
		case "[]float64":
			sliceFieldName := strcase.ToLowerCamel(fmt.Sprintf("%s_slice", fieldName))
			f = f.AddStatements(
				g.NewRawStatementf("%s := v.%s", sliceFieldName, fieldName),
				g.NewFor(
					fmt.Sprintf("i := 0; i < len(%s); i++", sliceFieldName),
					g.NewRawStatementf(`qp.Add("%s", fmt.Sprintf("%%f", %s[i]))`, paramName, sliceFieldName),
				),
			)
		default:
			log.Fatalf("[error] unsupported field type: %s", fieldType)
		}
	}

	f = f.AddStatements(g.NewReturnStatement("qp"))

	code, err := rootStmt.AddStatements(f).Gofmt("-s").Generate(0)
	if err != nil {
		log.Fatalf("[error] failed to generate code: %s", err)
	}

	err = ioutil.WriteFile(getFilenameToGenerate(args, typeName, output), []byte(code), 0644)
	if err != nil {
		log.Fatal(fmt.Errorf("[error] failed output generated code to a file: %w", err))
	}
}

func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func getFilenameToGenerate(args []string, typeName string, output string) string {
	if output != "" {
		return output
	}

	var dir string
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		dir = filepath.Dir(args[0])
	}
	return fmt.Sprintf("%s/%s_gen.go", dir, strcase.ToSnake(typeName))
}
