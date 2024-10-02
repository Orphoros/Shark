package emitter

import (
	"io"
	"shark/bytecode"
	"shark/compiler"
	"shark/config"
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
	vmConf      *config.VmConf
}

func New(sourceName *string, out io.Writer, vmConf *config.VmConf) *Emitter {
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

func (i *Emitter) Compile(sharkCode *string) *bytecode.Bytecode {
	l := lexer.New(sharkCode)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		i.printParserErrors(p.Errors(), i.sourceName, sharkCode)
		return nil
	}

	comp := compiler.NewWithState(i.symbolTable, i.constants)
	if err := comp.Compile(program); err != nil {
		i.printCompilerError(err, i.sourceName, sharkCode)
		return nil
	}

	return comp.Bytecode()
}

func (i *Emitter) Exec(bytecode *bytecode.Bytecode) {
	i.constants = bytecode.Constants
	machine := vm.NewWithGlobalsStore(bytecode, i.globals, i.vmConf)

	if err := machine.Run(); err != nil {
		i.printCompilerError(err, i.sourceName, nil)
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
		i.printCompilerError(err, i.sourceName, &in)
		return
	}
	code := comp.Bytecode()
	i.constants = code.Constants
	machine := vm.NewWithGlobalsStore(code, i.globals, i.vmConf)

	if err := machine.Run(); err != nil {
		i.printCompilerError(err, i.sourceName, nil)
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
		err.SetInputName(*filename)
		err.SetInputContent(content)
		io.WriteString(i.output, err.String())
		io.WriteString(i.output, "\n")
	}
}

func (i *Emitter) printCompilerError(err *exception.SharkError, filename *string, content *string) {
	err.SetInputName(*filename)
	err.SetInputContent(content)
	io.WriteString(i.output, err.String())
	io.WriteString(i.output, "\n")
}
