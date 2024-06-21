package exception

type SharkErrorCode uint16

type SharkErrorType uint8

const (
	// SharkErrorCodeUnknown is the default error code.
	SharkErrorCodeUnknown SharkErrorCode = iota
	// SharkErrorUnexpectedToken is the error code when an expression is expected
	// but a token or a statement is found.
	SharkErrorExpectedExpression
	// SharkErrorExpectedIdentifier is the error code when an identifier is expected
	// but a token or a statement is found.
	SharkErrorExpectedIdentifier
	// SharkErrorUnexpectedToken is the error code when a token is expected, but
	// another token is found.
	SharkErrorUnexpectedToken
	// SharkErrorUnexpectedEOF is the error code when the file ends unexpectedly.
	SharkErrorEOF
	// SharkErrorInteger is the error code when an integer is invalid.
	SharkErrorInteger
	// SharkErrorUnterminatedString is the error code when a string is not terminated.
	SharkErrorUnterminatedString
	// SharkErrorTopLeverReturn is the error code when return statement is in the main frame.
	SharkErrorTopLeverReturn
	// SharkErrorIdentifierNotFound is the error code when an identifier is not found.
	SharkErrorIdentifierNotFound
	// SharkErrorUnknownOperator is the error code when an unknown operator is found.
	SharkErrorUnknownOperator
	// SharkErrorIdentifierExpected is the error code when an identifier is expected, but another token is found.
	SharkErrorIdentifierExpected

	SharkErrorNonNumberIncrement

	SharkErrorMismatchedTypes

	SharkErrorUnknownStringOperator

	SharkErrorUnknownBoolOperator

	SharkErrorVMStackOverflow

	SharkErrorNonFunction

	SharkErrorNonFunctionCall

	SharkErrorArgumentNumberMismatch

	SharkErrorNonHashable

	SharkErrorNonIndexable

	SharkErrorNoDefaultValue
)

const (
	// SharkErrorTypeParser is the error type for parser errors.
	SharkErrorTypeParser SharkErrorType = iota
	// SharkErrorTypeLexer is the error type for lexer errors.
	SharkErrorTypeLexer
	// SharkErrorTypeCompiler is the error type for compiler errors.
	SharkErrorTypeCompiler
	// SharkErrorTypeRuntime is the error type for runtime errors.
	SharkErrorTypeRuntime
)

type SharkErrorCause struct {
	CauseMsg string
	Line     int
	LineTo   int
	Col      int
	ColTo    int
}

type SharkError struct {
	ErrMsg     string
	ErrCause   []SharkErrorCause
	ErrHelpMsg string
	ErrCode    SharkErrorCode
	ErrType    SharkErrorType
}
