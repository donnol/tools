package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"runtime"
)

func main() {
	// Parse a single source file.
	const input = `
package temperature
type S struct {
}

func (s S) A(i int) {

}

type myint int

type Myint myint

func (s S) B(i Myint) {

}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "celsius.go", input, 0)
	if err != nil {
		log.Fatal(err)
	}

	// imp := importer.Default()
	imp := importer.ForCompiler(fset, "source", nil)
	conf := types.Config{Importer: imp}
	pkg, err := conf.Check("temperature", fset, []*ast.File{f}, nil)
	if err != nil {
		log.Fatal(err)
	}

	celsius := pkg.Scope().Lookup("S").Type()
	for _, t := range []types.Type{celsius, types.NewPointer(celsius)} {
		fmt.Printf("Method set of %s:\n", t)
		mset := types.NewMethodSet(t)
		for i := 0; i < mset.Len(); i++ {
			fmt.Println(mset.At(i))
		}
		fmt.Println()
	}
	fmt.Println("Hello, playground", runtime.Version())
}
