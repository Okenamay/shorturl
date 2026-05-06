package main

import (
	"github.com/Okenamay/shorturl.git/cmd/linter/linter"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	// singlechecker.Main запускает анализатор
	singlechecker.Main(linter.ExitCheckAnalyzer)
}
