package main

import (
	"github.com/thedhejavu/go-error-analyzer/errorlysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(errorlysis.Analyzer)
}
