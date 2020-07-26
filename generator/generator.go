package generator

import (
	"github.com/d2verb/bee/ast"
	"github.com/d2verb/bee/ir"
)

// IrGenerator represents IR generator and contains internal state
type IrGenerator struct {
	program  *ast.Program
	bbLabel  int
	regLabel int
	out      *ir.BasicBlock // current basic block to appending ir
	function *ir.Function
}

// New returns a new IR generator
func New(program *ast.Program) *IrGenerator {
	ig := &IrGenerator{
		program:  program,
		bbLabel:  0,
		regLabel: 0,
	}
	return ig
}

// Generate generates IR from AST
func (ig *IrGenerator) Generate() *ir.Program {
	program := &ir.Program{}

	for _, function := range ig.program.Functions {
		ig.function = &ir.Function{Node: function}

		// empty basic block to making analysis easy
		ig.setCurrentBasicBlock(ig.newBasicBlock())
		bb := ig.newBasicBlock()
		ig.jmp(bb)

		// actually function entry point
		ig.setCurrentBasicBlock(bb)
		ig.generateStoreArguments(function.Parameters)
		ig.generateStatement(function.Body)

		// always return 0 at the end of function
		ig.ret(ig.imm(0))

		program.Functions = append(program.Functions, ig.function)
	}

	return program
}

func (ig *IrGenerator) generateStoreArguments(params []*ast.Variable) {
	for i, param := range params {
		ig.storeArg(i, param)
	}
}

func (ig *IrGenerator) generateStatement(node ast.Statement) {
	switch node := node.(type) {
	case *ast.IfStatement:
		consequence := ig.newBasicBlock()
		alternative := ig.newBasicBlock()
		last := ig.newBasicBlock()

		ig.br(ig.generateExpression(node.Condition), consequence, alternative)

		ig.setCurrentBasicBlock(consequence)
		ig.generateStatement(node.Consequence)
		ig.jmp(last)

		ig.setCurrentBasicBlock(alternative)
		if node.Alternative != nil {
			ig.generateStatement(node.Alternative)
		}
		ig.jmp(last)

		ig.setCurrentBasicBlock(last)
	case *ast.WhileStatement:
		cond := ig.newBasicBlock()
		body := ig.newBasicBlock()
		last := ig.newBasicBlock()

		ig.jmp(cond)

		ig.setCurrentBasicBlock(cond)
		ig.br(ig.generateExpression(node.Condition), body, last)

		ig.setCurrentBasicBlock(body)
		ig.generateStatement(node.Body)
		ig.jmp(cond)

		ig.setCurrentBasicBlock(last)
	case *ast.ExpressionStatement:
		ig.generateExpression(node.Expression)
	case *ast.PutsStatement:
		r := ig.generateExpression(node.Value)
		ig.puts(r)
	case *ast.ReturnStatement:
		r := ig.generateExpression(node.Value)
		ig.ret(r)
	case *ast.BlockStatement:
		for _, statement := range node.Statements {
			ig.generateStatement(statement)
		}
	}
}

func (ig *IrGenerator) generateExpression(node ast.Expression) *ir.Register {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return ig.imm(node.Value)
	case *ast.CallExpression:
		arguments := []*ir.Register{}
		for _, argument := range node.Arguments {
			arguments = append(arguments, ig.generateExpression(argument))
		}
		return ig.call(node.Function, arguments)
	case *ast.InfixExpression:
		if node.Operator == "=" {
			from := ig.generateExpression(node.Right)
			to := ig.bprel(node.Left.(*ast.Identifier).Var)
			return ig.store(to, from)
		}
		return ig.binop(node.Operator,
			ig.generateExpression(node.Left),
			ig.generateExpression(node.Right))
	case *ast.PrefixExpression:
		return ig.unary(node.Operator, ig.generateExpression(node.Right))
	case *ast.Identifier:
		return ig.load(ig.bprel(node.Var))
	}
	return nil
}

func (ig *IrGenerator) storeArg(index int, variable *ast.Variable) {
	ir := &ir.StoreArgIr{
		Index: index,
		Var:   variable,
	}
	ig.out.Irs = append(ig.out.Irs, ir)
}

func (ig *IrGenerator) store(to *ir.Register, from *ir.Register) *ir.Register {
	ir := &ir.StoreIr{
		R0: to,
		R1: from,
	}
	ig.out.Irs = append(ig.out.Irs, ir)
	return ir.R0
}

func (ig *IrGenerator) load(from *ir.Register) *ir.Register {
	ir := &ir.LoadIr{
		R0: ig.newRegister(),
		R1: from,
	}
	ig.out.Irs = append(ig.out.Irs, ir)
	return ir.R0
}

func (ig *IrGenerator) bprel(variable *ast.Variable) *ir.Register {
	ir := &ir.BprelIr{
		R:   ig.newRegister(),
		Var: variable,
	}
	ig.out.Irs = append(ig.out.Irs, ir)
	return ir.R
}

func (ig *IrGenerator) unary(op string, r1 *ir.Register) *ir.Register {
	ir := &ir.UnaryOpIr{
		Operator: op,
		R0:       ig.newRegister(),
		R1:       r1,
	}
	ig.out.Irs = append(ig.out.Irs, ir)
	return ir.R0
}

func (ig *IrGenerator) binop(op string, r1 *ir.Register, r2 *ir.Register) *ir.Register {
	ir := &ir.BinaryOpIr{
		Operator: op,
		R0:       ig.newRegister(),
		R1:       r1,
		R2:       r2,
	}
	ig.out.Irs = append(ig.out.Irs, ir)
	return ir.R0
}

func (ig *IrGenerator) call(function string, rs []*ir.Register) *ir.Register {
	ir := &ir.CallIr{
		Function:  function,
		Arguments: rs,
		Return:    ig.newRegister(),
	}
	ig.out.Irs = append(ig.out.Irs, ir)
	return ir.Return
}

func (ig *IrGenerator) br(r *ir.Register, consequence *ir.BasicBlock, alternative *ir.BasicBlock) ir.Ir {
	ir := &ir.BrIr{
		R:           r,
		Consequence: consequence,
		Alternative: alternative,
	}
	ig.out.Irs = append(ig.out.Irs, ir)
	return ir
}

func (ig *IrGenerator) imm(value int64) *ir.Register {
	ir := &ir.ImmIr{
		R:     ig.newRegister(),
		Value: value,
	}
	ig.out.Irs = append(ig.out.Irs, ir)
	return ir.R
}

func (ig *IrGenerator) jmp(target *ir.BasicBlock) ir.Ir {
	ir := &ir.JmpIr{
		Target: target,
	}
	ig.out.Irs = append(ig.out.Irs, ir)
	return ir
}

func (ig *IrGenerator) puts(r *ir.Register) ir.Ir {
	ir := &ir.PutsIr{
		R: r,
	}
	ig.out.Irs = append(ig.out.Irs, ir)
	return ir
}

func (ig *IrGenerator) ret(r *ir.Register) ir.Ir {
	ir := &ir.RetIr{
		R: r,
	}
	ig.out.Irs = append(ig.out.Irs, ir)
	ig.setCurrentBasicBlock(ig.newBasicBlock())
	return ir
}

func (ig *IrGenerator) newBasicBlock() *ir.BasicBlock {
	bb := &ir.BasicBlock{
		Label: ig.nextBBLabel(),
		Irs:   []ir.Ir{},
	}
	return bb
}

func (ig *IrGenerator) newRegister() *ir.Register {
	reg := &ir.Register{
		VirtualNo: ig.nextRegLabel(),
	}
	return reg
}

func (ig *IrGenerator) nextBBLabel() int {
	label := ig.bbLabel
	ig.bbLabel++
	return label
}

func (ig *IrGenerator) nextRegLabel() int {
	label := ig.regLabel
	ig.regLabel++
	return label
}

func (ig *IrGenerator) setCurrentBasicBlock(basicBlock *ir.BasicBlock) {
	ig.function.BasicBlocks = append(ig.function.BasicBlocks, basicBlock)
	ig.out = basicBlock
}
