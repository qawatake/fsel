package fderef

import (
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
				ptrerrs := returnsPtrsAndErr(pass, i, instr)
			}
		}
	}

	return nil, nil
}

func returnsPtrsAndErr(pass *analysis.Pass, id int, instr ssa.Instruction) []ptrerr {
	defer func() { recover() }()
	call := instr.(*ssa.Call)
	numResults := call.Common().Signature().Results().Len()
	// -1 to omit error
	ptrerrs := make([]ptrerr, 0, numResults-1)
	var errValue ssa.Value
	for j := id; j < len(instr.Block().Instrs); j++ {
		func() {
			defer func() { recover() }()
			extract := instr.Block().Instrs[j].(*ssa.Extract)
			if extract.Tuple == call {
				if types.Identical(extract.Type(), typesErr) {
					errValue = extract
					return
				}
				if extract.Type().Underlying().(*types.Pointer) != nil {
					ptrerrs = append(ptrerrs, ptrerr{
						ptr: extract,
					})
					return
				}
			}
		}()
	}
	if errValue != nil {
		for _, p := range ptrerrs {
			p.err = errValue
		}
	}
	return ptrerrs
}

type ptrerr struct {
	ptr ssa.Value
	err ssa.Value
}

var typesErr = types.Universe.Lookup("error").Type()
