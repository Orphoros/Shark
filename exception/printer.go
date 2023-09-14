package exception

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func PrintSharkLineError(err *SharkError, inputContent, inputName *string) {
	// TODO: Accept io.Writer as output instead of defaulting error print to console
	emptySpace := strings.Repeat(" ", len(strconv.Itoa(err.ErrCause[len(err.ErrCause)-1].Line))+1)

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

	fmt.Printf("%s_error[%04d]: %s\n", errType, err.ErrCode, err.ErrMsg)

	fmt.Printf("  --> %s:%d:%d\n%s|\n", *inputName, err.ErrCause[0].Line, err.ErrCause[0].Col, emptySpace)

	if len(err.ErrCause) > 0 {
		lines := strings.Split(*inputContent, "\n")
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
					errorLineMarker = strings.Repeat("^", utf8.RuneCountInString((curLineContent)))
				}

				fmt.Printf("%d |\t%s\n", i, curLineContent)

				fmt.Printf("%s|\t%s %s\n", emptySpace, errorLineMarker, msg)
				fmt.Printf("%s|\n", emptySpace)
			}
		}
	}

	if err.ErrHelpMsg != "" {
		fmt.Printf("%shelp: %s\n", emptySpace, err.ErrHelpMsg)
	}
}

func PrintSharkRuntimeError(err *SharkError) {

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

	fmt.Printf("<red1>%s_error[%04d]: %s\n", errType, err.ErrCode, err.ErrMsg)

	if err.ErrHelpMsg != "" {
		fmt.Printf("\t-->help: %s\n", err.ErrHelpMsg)
	}
}

func PrintExitMsg(msg string, exitCode int) {
	fmt.Printf("error: %s\n", msg)
	os.Exit(exitCode)
}

func PrintExitMsgCtx(msg, ctx string, exitCode int) {
	fmt.Printf("error: %s\n", msg)
	fmt.Printf("   --> %s\n", ctx)
	os.Exit(exitCode)
}
