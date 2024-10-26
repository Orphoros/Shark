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

	sharkEmitter := emitter.New(path, os.Stdout, &argConfig.OrpVM)

	_ = sharkEmitter.Compile(&content)

	st := sharkEmitter.GetSymbolTable()

	l := lexer.New(&content)

	var identName string

	for t := l.NextToken(); t.Type != token.EOF; {
		if t.Type == token.IDENT && t.Pos.Line-1 == int(params.Position.Line) && t.Pos.ColFrom-1 <= int(params.Position.Character) && t.Pos.ColTo-1 >= int(params.Position.Character) {
			identName = t.Literal
			break
		}
		t = l.NextToken()
	}

	sym, ok := st.FindIdent(identName)

	if !ok {
		return &protocol.Hover{
			Contents: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "Error: " + identName + " not found",
			},
		}, nil
	}

	var value string

	switch sym.Scope {
	case compiler.BuiltinScope:
		value = "**Shark Function (_builtin_)**\n```shark\n" + sym.Name + "(...args)\n```"
	case compiler.GlobalScope:
		value = "**Shark Identifier (_global expression_)**\n```shark\nlet "
		if sym.Mutable {
			value += "mut "
		}
		value += sym.Name + "\n```"
	case compiler.LocalScope:
		value = "**Shark Identifier (_local expression_)**\n```shark\nlet "
		if sym.Mutable {
			value += "mut "
		}
		value += sym.Name + "\n```"

	}

	return &protocol.Hover{
		Contents: &protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: value,
		},
	}, nil
}
