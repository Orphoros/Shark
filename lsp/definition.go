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

func definition(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
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

	// FIXME: For function params, it gives the global scope definition, and not the local scope definition
	sym, ok := st.FindIdent(identName)

	if ok && sym.Scope != compiler.BuiltinScope {
		return &protocol.Location{
			URI: params.TextDocument.URI,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(sym.Pos.Line - 1),
					Character: uint32(sym.Pos.ColFrom - 1),
				},
				End: protocol.Position{
					Line:      uint32(sym.Pos.Line - 1),
					Character: uint32(sym.Pos.ColTo - 1),
				},
			},
		}, nil
	} else {
		return nil, nil
	}
}
