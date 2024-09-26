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

func (err *SharkError) String() string {
	var str bytes.Buffer

	var errType string

	switch err.ErrType {
	case SharkErrorTypeParser:
		errType = "parser"
	case SharkErrorTypeLexer:
		errType = "lexer"
	case SharkErrorTypeCompiler:
		errType = "compiler"
	case SharkErrorTypeRuntime:
		errType = "runtime"
	}

	str.WriteString(fmt.Sprintf("%s_error[%04d]: %s\n", errType, err.ErrCode, err.ErrMsg))

	str.WriteString("  --> ")

	if err.InputName != nil {
		str.WriteString(*err.InputName)
	} else {
		str.WriteString("std")
	}

	if len(err.ErrCause) == 0 {
		return str.String()
	}

	emptySpace := strings.Repeat(" ", len(strconv.Itoa(err.ErrCause[len(err.ErrCause)-1].Line))+1)

	str.WriteString(fmt.Sprintf(":%d:%d\n%s|\n", err.ErrCause[0].Line, err.ErrCause[0].Col, emptySpace))

	if err.InputContent == nil {
		return str.String()
	}

	if len(err.ErrCause) > 0 {
		lines := strings.Split(*err.InputContent, "\n")
		for _, cause := range err.ErrCause {
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

	if err.ErrHelpMsg != nil {
		str.WriteString(fmt.Sprintf("%shelp: %s\n", emptySpace, *err.ErrHelpMsg))
	}

	return str.String()
}
