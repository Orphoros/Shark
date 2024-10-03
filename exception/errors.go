package exception

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type SharkErrorCause struct {
	CauseMsg string
	Line     int
	LineTo   int
	Col      int
	ColTo    int
}

func NewSharkErrorCause(causeMsg string, line, lineTo, col, colTo int) SharkErrorCause {
	return SharkErrorCause{
		CauseMsg: causeMsg,
		Line:     line,
		LineTo:   lineTo,
		Col:      col,
		ColTo:    colTo,
	}
}

type SharkError struct {
	ErrMsg       string
	ErrCause     []SharkErrorCause
	ErrHelpMsg   *string
	ErrCode      SharkErrorCode
	ErrType      SharkErrorType
	InputName    *string
	InputContent *string
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

	emptySpace := strings.Repeat(" ", len(strconv.Itoa(e.ErrCause[len(e.ErrCause)-1].Line))+1)

	str.WriteString(fmt.Sprintf(":%d:%d\n%s|\n", e.ErrCause[0].Line, e.ErrCause[0].Col, emptySpace))

	if e.InputContent == nil {
		return str.String()
	}

	if len(e.ErrCause) > 0 {
		lines := strings.Split(*e.InputContent, "\n")
		for _, cause := range e.ErrCause {
			for i := cause.Line; i <= cause.LineTo; i++ {
				curLineContent := strings.ReplaceAll(lines[i-1], "\t", " ")
				msg := cause.CauseMsg
				if i != cause.LineTo {
					msg = ""
				}
				var errorLineMarker string
				if cause.Line == cause.LineTo {
					errorLineMarker = strings.Repeat(" ", cause.Col-1) + strings.Repeat("^", cause.ColTo-cause.Col)
				} else if cause.Line == i {
					errorLineMarker = strings.Repeat(" ", cause.Col-1) + strings.Repeat("^", utf8.RuneCountInString(curLineContent)-cause.Col+1)
				} else if cause.LineTo == i {
					errorLineMarker = strings.Repeat("^", cause.ColTo-1)
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
