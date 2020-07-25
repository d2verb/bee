package ir

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/d2verb/bee/ast"
)

// Ir is the interface for all intermediate representation
type Ir interface {
	ir()
	String() string
}

// Program represents program in IR
type Program struct {
	Functions []*Function
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, function := range p.Functions {
		out.WriteString(function.String())
	}

	return out.String()
}

// Function represents function in IR
type Function struct {
	Node        *ast.Function
	BasicBlocks []*BasicBlock
}

func (f *Function) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("[%s]\n", f.Node.Name))
	for _, basicBlock := range f.BasicBlocks {
		out.WriteString(basicBlock.String())
	}
	out.WriteString("\n")
	return out.String()
}

// Register represents the virtual register
type Register struct {
	VirtualNo int
}

func (reg *Register) String() string {
	return fmt.Sprintf("r%d", reg.VirtualNo)
}

// BasicBlock represents a basic block
type BasicBlock struct {
	Label int
	Irs   []Ir
}

func (bb *BasicBlock) String() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf(".L%d:\n", bb.Label))
	for _, ir := range bb.Irs {
		out.WriteString(fmt.Sprintf("  %s\n", ir.String()))
	}

	return out.String()
}

// BinaryOpIr represents `r0 = r1 OP r2`
type BinaryOpIr struct {
	Operator string
	R0       *Register
	R1       *Register
	R2       *Register
}

func (ir *BinaryOpIr) ir() {}
func (ir *BinaryOpIr) String() string {
	return fmt.Sprintf("%s = %s %s %s",
		ir.R0.String(), ir.R1.String(), ir.Operator, ir.R2.String())
}

// UnaryOpIr represents `r0 = OP r1`
type UnaryOpIr struct {
	Operator string
	R0       *Register
	R1       *Register
}

func (ir *UnaryOpIr) ir() {}
func (ir *UnaryOpIr) String() string {
	return fmt.Sprintf("%s = %s %s",
		ir.R0.String(), ir.Operator, ir.R1.String())
}

// BrIr represents `BR r .L0 .L1`
// If condition is true, then jump to L0, else jump tp L1
type BrIr struct {
	R           *Register
	Consequence *BasicBlock
	Alternative *BasicBlock
}

func (ir *BrIr) ir() {}
func (ir *BrIr) String() string {
	return fmt.Sprintf("BR r%d, .L%d, .L%d",
		ir.R.VirtualNo, ir.Consequence.Label, ir.Alternative.Label)
}

// ImmIr represents `IMM r VALUE`
type ImmIr struct {
	R     *Register
	Value int64
}

func (ir *ImmIr) ir() {}
func (ir *ImmIr) String() string {
	return fmt.Sprintf("IMM r%d, %d", ir.R.VirtualNo, ir.Value)
}

// JmpIr represents `JMP .L0`
type JmpIr struct {
	Target *BasicBlock
}

func (ir *JmpIr) ir() {}
func (ir *JmpIr) String() string {
	return fmt.Sprintf("JMP .L%d", ir.Target.Label)
}

// PutsIr represents `PUTS r`
type PutsIr struct {
	R *Register
}

func (ir *PutsIr) ir() {}
func (ir *PutsIr) String() string {
	return fmt.Sprintf("PUTS r%d", ir.R.VirtualNo)
}

// RetIr represents `RET r`
type RetIr struct {
	R *Register
}

func (ir *RetIr) ir() {}
func (ir *RetIr) String() string {
	return fmt.Sprintf("RET r%d", ir.R.VirtualNo)
}

// CallIr represents `CALL r0 r1 r2 ...`
type CallIr struct {
	Function  string
	Return    *Register
	Arguments []*Register
}

func (ir *CallIr) ir() {}
func (ir *CallIr) String() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("CALL %s ", ir.Function))

	rs := []string{}
	for _, r := range ir.Arguments {
		rs = append(rs, fmt.Sprintf("r%d", r.VirtualNo))
	}
	out.WriteString(strings.Join(rs, ", "))

	return out.String()
}

// BprelIr represents `BPREL r rbp - var.offset` to calculate the address of local variable
type BprelIr struct {
	R   *Register
	Var *ast.Variable
}

func (ir *BprelIr) ir() {}
func (ir *BprelIr) String() string {
	return fmt.Sprintf("BPREL r%d, %s@(rbp - %d)", ir.R.VirtualNo, ir.Var.Name, ir.Var.Offset)
}

// LoadIr represents `LOAD r0 [r1]`
type LoadIr struct {
	R0 *Register
	R1 *Register
}

func (ir *LoadIr) ir() {}
func (ir *LoadIr) String() string {
	return fmt.Sprintf("LOAD r%d, [r%d]", ir.R0.VirtualNo, ir.R1.VirtualNo)
}

// StoreIr represents `STORE [r0] r1`
type StoreIr struct {
	R0 *Register
	R1 *Register
}

func (ir *StoreIr) ir() {}
func (ir *StoreIr) String() string {
	return fmt.Sprintf("STORE [r%d] r%d", ir.R0.VirtualNo, ir.R1.VirtualNo)
}

// StoreArgIr represents `STORE_ARG <index> <var>`
type StoreArgIr struct {
	Index int
	Var   *ast.Variable
}

func (ir *StoreArgIr) ir() {}
func (ir *StoreArgIr) String() string {
	return fmt.Sprintf("STORE_ARG %d %s", ir.Index, ir.Var.Name)
}
