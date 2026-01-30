// Command timemanip is a standalone linter that detects usage of time.Time manipulation methods.
package main

import (
	"github.com/yoshihiro-shu/timemanip"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(timemanip.Analyzer)
}
