package emitter

import (
	"bytes"
	"fmt"
	"io"
	"shark/compiler"
	"shark/exception"
	"shark/lexer"
	"shark/object"
	"shark/parser"
	"shark/vm"
)

type Emitter struct {
	constants   []object.Object
	globals     []object.Object
	symbolTable *compiler.SymbolTable
	output      io.Writer
	sourceName  *string
	vmConf      *vm.VmConf
}

func New(sourceName *string, out io.Writer, vmConf *vm.VmConf) *Emitter {
	emitter := &Emitter{
		constants:   []object.Object{},
		globals:     make([]object.Object, vmConf.GlobalsSize),
		symbolTable: compiler.NewSymbolTable(),
		output:      out,
		sourceName:  sourceName,
		vmConf:      vmConf,
	}
	for i, v := range object.Builtins {
		emitter.symbolTable.DefineBuiltin(i, v.Name)
	}

	return emitter
}

func (i *Emitter) Compile(sharkCode *string) *compiler.Bytecode {
	l := lexer.New(sharkCode)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		i.printParserErrors(p.Errors(), i.sourceName, sharkCode)
		return nil
	}

	comp := compiler.NewWithState(i.symbolTable, i.constants)
	if err := comp.Compile(program); err != nil {
		i.printCompilerError(err)
		return nil
	}

	return comp.Bytecode()
}

func (i *Emitter) Exec(bytecode *compiler.Bytecode) {
	i.constants = bytecode.Constants
	machine := vm.NewWithGlobalsStore(bytecode, i.globals, i.vmConf)

	if err := machine.Run(); err != nil {
		i.printCompilerError(err)
		return
	}

	lastPopped := machine.LastPoppedStackElem()

	if lastPopped != nil {
		if lastPopped.Type() == object.ERROR_OBJ {
			io.WriteString(i.output, "\tERROR: "+lastPopped.Inspect()+"\n")

		}
	}
}

func (i *Emitter) Interpret(in string) {
	l := lexer.New(&in)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		i.printParserErrors(p.Errors(), i.sourceName, &in)
		return
	}

	comp := compiler.NewWithState(i.symbolTable, i.constants)
	err := comp.Compile(program)
	if err != nil {
		i.printCompilerError(err)
		return
	}
	code := comp.Bytecode()
	i.constants = code.Constants
	machine := vm.NewWithGlobalsStore(code, i.globals, i.vmConf)

	if err := machine.Run(); err != nil {
		i.printCompilerError(err)
		return
	}

	lastPopped := machine.LastPoppedStackElem()

	if lastPopped != nil {
		if lastPopped.Type() == object.ERROR_OBJ {
			io.WriteString(i.output, "\tERROR: "+lastPopped.Inspect()+"\n")

		}
	}
}

func (i *Emitter) printParserErrors(errors []exception.SharkError, filename, content *string) {
	for _, err := range errors {
		exception.PrintSharkLineError(&err, content, filename)
	}
}

func (i *Emitter) printCompilerError(err *exception.SharkError) {
	exception.PrintSharkRuntimeError(err)
}

func EmitInstructionsTable(bytecode *compiler.Bytecode, out io.Writer) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "@instructions:\n")

	ins := bytecode.Instructions

	txt := ins.String()

	lines := bytes.Split([]byte(txt), []byte("\n"))

	for i, line := range lines {
		if i == len(lines)-1 {
			break
		}
		fmt.Fprintf(&buf, "| %s\n", line)
	}

	fmt.Fprintf(&buf, "\n@constants:\n")

	for i, constant := range bytecode.Constants {
		switch constant.Type() {
		case object.COMPILED_FUNCTION_OBJ:
			fmt.Fprintf(&buf, "| %04d FUNC {\n", i)
			cf := constant.(*object.CompiledFunction)
			txt := cf.Instructions.String()
			lines := bytes.Split([]byte(txt), []byte("\n"))
			for i, line := range lines {
				if i == len(lines)-1 {
					break
				}
				fmt.Fprintf(&buf, "| \t%s\n", line)
			}
			fmt.Fprintf(&buf, "| }\n")
		case object.STRING_OBJ:
			fmt.Fprintf(&buf, "| %04d %s: \"%s\"\n", i, constant.Type(), constant.Inspect())
		default:
			fmt.Fprintf(&buf, "| %04d %s: %s\n", i, constant.Type(), constant.Inspect())
		}
	}

	out.Write(buf.Bytes())
}
