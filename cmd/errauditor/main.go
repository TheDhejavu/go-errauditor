package main

import (
	"github.com/thedhejavu/errauditor/errauditor"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(errauditor.Analyzer)
}
