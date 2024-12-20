package lsp

import (
	"fmt"
	"os"
	"shark/compiler"
	"shark/config"
	"shark/emitter"
	"shark/internal"
	"shark/lexer"
	"shark/token"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func hover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	path, err := internal.GetFilePathFromURI(params.TextDocument.URI)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	bytes, err := internal.ReadFile(*path)
	if err != nil {
		return nil, err
	}
	content := string(bytes)

	argConfig, err := config.LocateConfig(nil, path)
	if err != nil {
		return nil, err
	}

	sharkEmitter := emitter.New(path, os.Stdout, &argConfig.NidumVM)

	l := lexer.New(&content)

	var identName string
	var identPos token.Position

	for t := l.NextToken(); t.Type != token.EOF; {
		if t.Type == token.IDENT && t.Pos.Line-1 == int(params.Position.Line) && t.Pos.ColFrom-1 <= int(params.Position.Character) && t.Pos.ColTo-1 >= int(params.Position.Character) {
			identName = t.Literal
			identPos = t.Pos
			break
		}
		t = l.NextToken()
	}

	_ = sharkEmitter.Compile(&content, identPos)

	st := sharkEmitter.GetSymbolTable()

	sym, ok := st.FindIdent(identName)

	if !ok {
		return nil, nil
	}

	var value string

	switch sym.Scope {
	case compiler.BuiltinScope:
		value = "**Shark Function (_builtin_)**\n```shark\n" + sym.Name + "(...args)\n```"
	case compiler.GlobalScope:
		value = "**Shark Identifier (_global expression_)**\n```shark\n"
		if sym.VariadicType {
			value += "var "
		} else {
			value += "let "
		}
		if sym.Mutable {
			value += "mut "
		}
		value += sym.Name + " : " + sym.ObjType.String() + "\n```"
	case compiler.LocalScope:
		value = "**Shark Identifier (_local expression_)**\n```shark\n"
		if sym.VariadicType {
			value += "var "
		} else {
			value += "let "
		}

		if sym.Mutable {
			value += "mut "
		}

		value += sym.Name + " : " + sym.ObjType.String() + "\n```"
	}

	return &protocol.Hover{
		Contents: &protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: value,
		},
	}, nil
}
