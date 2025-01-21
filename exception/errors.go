package exception

import (
	"bytes"
	"fmt"
	"shark/token"
	"strconv"
	"strings"
	"unicode/utf8"
)

type SharkErrorCause struct {
	CauseMsg string
	Pos      token.Position
}

func NewSharkErrorCause(causeMsg string, pos token.Position) SharkErrorCause {
	return SharkErrorCause{
		CauseMsg: causeMsg,
		Pos:      pos,
	}
}

type SharkError struct {
	ErrHelpMsg   *string
	InputName    *string
	InputContent *string
	ErrMsg       string
	ErrCause     []SharkErrorCause
	ErrCode      SharkErrorCode
	ErrType      SharkErrorType
}

func NewSharkError(errType SharkErrorType, errCode SharkErrorCode, parameters ...interface{}) *SharkError {
	if errMessages[errCode].errMsg == "" {
		errCode = SharkErrorCodeUnknown
	}

	return &SharkError{
		ErrMsg:   fmt.Sprintf(errMessages[errCode].errMsg, parameters...),
		ErrCode:  errCode,
		ErrType:  errType,
		ErrCause: make([]SharkErrorCause, 0),
	}
}

func (e *SharkError) Error() string {
	return e.ErrMsg
}

func (e *SharkError) AddCause(cause SharkErrorCause) {
	e.ErrCause = append(e.ErrCause, cause)
}

func (e *SharkError) SetHelpMsg(helpMsg string) {
	e.ErrHelpMsg = &helpMsg
}

func (e *SharkError) SetInputName(inputName string) {
	e.InputName = &inputName
}

func (e *SharkError) SetInputContent(inputContent *string) {
	e.InputContent = inputContent
}

func (e *SharkError) String() string {
	var str bytes.Buffer

	var errType string

	switch e.ErrType {
	case SharkErrorTypeParser:
		errType = "parser"
	case SharkErrorTypeLexer:
		errType = "lexer"
	case SharkErrorTypeCompiler:
		errType = "compiler"
	case SharkErrorTypeRuntime:
		errType = "runtime"
	}

	str.WriteString(fmt.Sprintf("%s_error[%04d]: %s\n", errType, e.ErrCode, e.ErrMsg))

	str.WriteString("  --> ")

	if e.InputName != nil {
		str.WriteString(*e.InputName)
	} else {
		str.WriteString("std")
	}

	if len(e.ErrCause) == 0 {
		return str.String()
	}

	emptySpace := strings.Repeat(" ", len(strconv.Itoa(e.ErrCause[len(e.ErrCause)-1].Pos.Line))+1)

	str.WriteString(fmt.Sprintf(":%d:%d\n%s|\n", e.ErrCause[0].Pos.Line, e.ErrCause[0].Pos.ColFrom, emptySpace))

	if e.InputContent == nil {
		return str.String()
	}

	if len(e.ErrCause) > 0 {
		lines := strings.Split(*e.InputContent, "\n")
		for _, cause := range e.ErrCause {
			for i := cause.Pos.Line; i <= cause.Pos.LineTo; i++ {
				curLineContent := strings.ReplaceAll(lines[i-1], "\t", " ")
				msg := cause.CauseMsg
				if i != cause.Pos.LineTo {
					msg = ""
				}
				var errorLineMarker string
				if cause.Pos.Line == cause.Pos.LineTo {
					errorLineMarker = strings.Repeat(" ", cause.Pos.ColFrom-1) + strings.Repeat("^", cause.Pos.ColTo-cause.Pos.ColFrom)
				} else if cause.Pos.Line == i {
					errorLineMarker = strings.Repeat(" ", cause.Pos.ColFrom-1) + strings.Repeat("^", utf8.RuneCountInString(curLineContent)-cause.Pos.ColFrom+1)
				} else if cause.Pos.LineTo == i {
					errorLineMarker = strings.Repeat("^", cause.Pos.ColTo-1)
				} else {
					errorLineMarker = strings.Repeat("^", utf8.RuneCountInString(curLineContent))
				}

				str.WriteString(fmt.Sprintf("%d |\t%s\n", i, curLineContent))

				str.WriteString(fmt.Sprintf("%s|\t%s %s\n", emptySpace, errorLineMarker, msg))
				str.WriteString(fmt.Sprintf("%s|\n", emptySpace))
			}
		}
	}

	if e.ErrHelpMsg != nil {
		str.WriteString(fmt.Sprintf("%shelp: %s\n", emptySpace, *e.ErrHelpMsg))
	}

	return str.String()
}
