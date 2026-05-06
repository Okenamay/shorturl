package linter

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestExitCheckAnalyzer(t *testing.T) {
	// analysistest.TestData() находит директорию testdata
	testdata := analysistest.TestData()

	// analysistest.Run запускает анализатор на пакете "main"
	// (который находится в testdata/src/main)
	// Директивы // want в файле main.go будут проверены
	analysistest.Run(t, testdata, ExitCheckAnalyzer, "mainpkg")
}
