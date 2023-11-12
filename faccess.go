package fderef

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

const doc = "fderef is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "fderef",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs

	for _, fn := range funcs {
		for _, b := range fn.Blocks {
			for i, instr := range b.Instrs {
				checkInstr(pass, i, instr)
			}
		}
	}

	return nil, nil
}

func checkInstr(pass *analysis.Pass, id int, instr ssa.Instruction) {
	defer func() {
		recover()
	}()
	call := instr.(*ssa.Call)
	ptrFound := false
	errFound := false
	for j := id; j < len(instr.Block().Instrs); j++ {
		func() {
			defer func() {
				recover()
			}()
			extract := instr.Block().Instrs[j].(*ssa.Extract)
			if extract.Tuple == call {
				if types.Implements(extract.Type(), errType) {
					errFound = true
					return
				}
				if extract.Type().Underlying().(*types.Pointer) != nil {
					ptrFound = true
					return
				}
			}
		}()
	}
	if ptrFound && errFound {
		fmt.Println("😗", call)
	}
}

var errType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

// func returnsPtrsAndErr(tuple *types.Tuple) bool {
// 	if tuple.Len() < 2 {
// 		return false
// 	}
// 	ptrFound := false
// 	errFound := false
// }

// type returnedPtrsAndErr struct {
// 	ptrs []any
// 	err  ssa.Value
// }
