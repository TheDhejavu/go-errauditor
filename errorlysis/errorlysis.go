package errorlysis

import (
	"fmt"
	"go/ast"
	"log"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/analysis"
)

type ErrorType string

const (
	Error   ErrorType = "Error"
	Default ErrorType = "Default"
)

var (
	logger   *logrus.Logger
	Analyzer = &analysis.Analyzer{
		Name: "errorlysis",
		Doc:  "Reports returned errors",
		Run:  run,
	}
)

// ExtractFuncType extracts and returns the func returned type
func ExtractFuncType(funcType *ast.FuncType) (ErrorType, int) {
	for idx, r := range funcType.Results.List {
		if etype, ok := r.Type.(*ast.Ident); ok {
			if etype.Name == "error" {
				return Error, idx
			}
		}
	}
	return Default, -1
}

func ReportSelFromExpr(pass *analysis.Pass, expr ast.Expr, arg string) {

	if selExpr, ok := expr.(*ast.SelectorExpr); ok {
		pass.Report(analysis.Diagnostic{
			Pos:     selExpr.Pos(),
			Message: fmt.Sprintf("-- %s(%s) --", selExpr.Sel.Name, arg),
		})
	}
}

func ExtarctArgFromExpr(pass *analysis.Pass, expr []ast.Expr) string {
	if arg, ok := expr[0].(*ast.BasicLit); ok {
		return arg.Value
	}
	return ""
}

// ExtractReturnedErrorFromStmt extracts all instance of returned errors and string.
func ExtractReturnedErrorFromStmt(pass *analysis.Pass, etypePosIdx int, body *ast.BlockStmt) {
	ast.Inspect(body, func(node ast.Node) bool {
		if rtrnStmt, ok := node.(*ast.ReturnStmt); ok {
			expr := rtrnStmt.Results[etypePosIdx]
			// handle call expression
			if callExpr, ok := expr.(*ast.CallExpr); ok {
				ReportSelFromExpr(
					pass,
					callExpr.Fun,
					ExtarctArgFromExpr(pass, callExpr.Args),
				)
			}
			// handle selExpr
			ReportSelFromExpr(pass, expr, "")
		}
		return true
	})
}

// WalkThroughExpr work through the file nodes
func WalkThroughExpr(pass *analysis.Pass, file *ast.File) {
	for _, d := range file.Decls {
		if funcCall, ok := d.(*ast.FuncDecl); ok {
			name := funcCall.Name.Obj.Name
			returnedType, posIdx := ExtractFuncType(funcCall.Type)
			// check the returned type and position index
			if returnedType == Error && posIdx != -1 {
				// Report Extract
				pass.Report(analysis.Diagnostic{
					Pos:     funcCall.Pos(),
					Message: fmt.Sprintf(" -- %s() ---", name),
				})
				ExtractReturnedErrorFromStmt(pass, posIdx, funcCall.Body)
			}
			// ignore if func return type is not an error.
		}
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	logger = logrus.New()

	lvl, err := logrus.ParseLevel("info")
	if err != nil {
		log.Panic(err)
	}

	logger.SetLevel(lvl)

	for _, f := range pass.Files {
		ast.Inspect(f, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.File:
				// walk though expression
				WalkThroughExpr(pass, n)
			}
			// ast.Print(pass.Fset, f)
			return true
		})
	}
	return nil, nil
}
