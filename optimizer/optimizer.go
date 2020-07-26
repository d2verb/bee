package optimizer

import (
	"github.com/d2verb/bee/ir"
)

// ConstantFolding calculates all binary operation of constant value
func ConstantFolding(program *ir.Program) bool {
	var changed = false
	for _, function := range program.Functions {
		for _, basicBlock := range function.BasicBlocks {
			for i := 0; i+2 < len(basicBlock.Irs); i++ {
				ir0, ok := basicBlock.Irs[i].(*ir.ImmIr)
				if !ok {
					continue
				}
				ir1, ok := basicBlock.Irs[i+1].(*ir.ImmIr)
				if !ok {
					continue
				}
				ir2, ok := basicBlock.Irs[i+2].(*ir.BinaryOpIr)
				if !ok {
					continue
				}

				var result int64
				switch ir2.Operator {
				case "+":
					result = ir0.Value + ir1.Value
				case "-":
					result = ir0.Value - ir1.Value
				case "*":
					result = ir0.Value * ir1.Value
				case "/":
					result = ir0.Value / ir1.Value
				}

				basicBlock.Irs[i] = &ir.NopIr{}
				basicBlock.Irs[i+1] = &ir.NopIr{}
				basicBlock.Irs[i+2] = &ir.ImmIr{R: ir2.R0, Value: result}

				changed = true
			}
		}
	}

	return changed
}

// EliminateNop eliminates all NOPs
func EliminateNop(program *ir.Program) {
	for _, function := range program.Functions {
		for _, basicBlock := range function.BasicBlocks {
			newIrs := []ir.Ir{}
			for _, _ir := range basicBlock.Irs {
				if _, ok := _ir.(*ir.NopIr); ok {
					continue
				}
				newIrs = append(newIrs, _ir)
			}
			basicBlock.Irs = newIrs
		}
	}
}

// Rewrite
//
// BPREL r1, a@(rbp - 0)
// STORE [r1] r0
// BPREL r2, a@(rbp - 0)
// LOAD r3, [r2]
//
// to
//
// BPREL r1, a@(rbp - 0)
// STORE [r1] r0
// NOP
// MOV r3, r0
func eliminateRedundantCode(basicBlock *ir.BasicBlock) {
	for i := 0; i+3 < len(basicBlock.Irs); i++ {
		ir0, ok := basicBlock.Irs[i].(*ir.BprelIr)
		if !ok {
			continue
		}
		ir1, ok := basicBlock.Irs[i+1].(*ir.StoreIr)
		if !ok {
			continue
		}
		ir2, ok := basicBlock.Irs[i+2].(*ir.BprelIr)
		if !ok {
			continue
		}
		ir3, ok := basicBlock.Irs[i+3].(*ir.LoadIr)
		if !ok {
			continue
		}
		if ir0.Var != ir2.Var {
			continue
		}
		basicBlock.Irs[i+2] = &ir.NopIr{}
		basicBlock.Irs[i+3] = &ir.MovIr{R0: ir3.R0, R1: ir1.R1}
	}
}

// Peephole does peephole optimization
func Peephole(program *ir.Program) {
	for _, function := range program.Functions {
		for _, basicBlock := range function.BasicBlocks {
			eliminateRedundantCode(basicBlock)
		}
	}
}

// LocalOptimize optimizes program
func LocalOptimize(program *ir.Program) {
	Peephole(program)
	EliminateNop(program)

	for ConstantFolding(program) {
		EliminateNop(program)
	}
}
