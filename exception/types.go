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

	SharkErrorDuplicateIdentifier

	// SharkErrorUnknownOperator is the error code when an unknown operator is found.
	SharkErrorUnknownOperator
	// SharkErrorIdentifierExpected is the error code when an identifier is expected, but another token is found.
	SharkErrorIdentifierExpected

	SharkErrorNonNumberIncrement

	SharkErrorNonNumberDecrement

	SharkErrorMismatchedTypes

	SharkErrorUnknownStringOperator

	SharkErrorUnknownBoolOperator

	SharkErrorVMStackOverflow

	SharkErrorVMFrameStackOverflow

	SharkErrorNonFunction

	SharkErrorNonFunctionCall

	SharkErrorArgumentNumberMismatch

	SharkErrorNonHashable

	SharkErrorNonIndexable

	SharkErrorNoDefaultValue

	SharkErrorImmutableValue

	SharkErrorIndexOutOfBounds

	SharkErrorOptionalParameter

	SharkErrorDivisionByZero

	SharkErrorUnknownType

	SharkErrorTypeMismatch

	SharkErrorTupleDeconstructMismatch

	SharkErrorInvalidNumber

	SharkErrorArgumentCount

	SharkErrorNotCallable

	SharkErrorTypeNotFound

	SharkErrorTypeSyntax
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

var errMessages = []struct {
	code   SharkErrorCode
	errMsg string
}{
	{SharkErrorCodeUnknown, "unknown error"},
	{SharkErrorExpectedExpression, "expected to receive an expression that evaluates to a value, but got '%v' instead, which has no value"},
	{SharkErrorExpectedIdentifier, "expected to receive an identifier, but got '%v' instead"},
	{SharkErrorUnexpectedToken, "received an unexpected token '%v'"},
	{SharkErrorEOF, "unexpected end of file"},
	{SharkErrorInteger, "expected an integer, but got '%v' instead"},
	{SharkErrorUnterminatedString, "string is not terminated"},
	{SharkErrorTopLeverReturn, "return statement is in the main frame"},
	{SharkErrorIdentifierNotFound, "identifier '%v' not found"},
	{SharkErrorDuplicateIdentifier, "identifier '%v' is already declared"},
	{SharkErrorUnknownOperator, "unknown operator '%v'"},
	{SharkErrorIdentifierExpected, "expected an identifier, but got '%v' instead"},
	{SharkErrorNonNumberIncrement, "cannot increment non-number value '%v'"},
	{SharkErrorNonNumberDecrement, "cannot decrement non-number value '%v'"},
	{SharkErrorMismatchedTypes, "mismatched types '%v' and '%v'"},
	{SharkErrorUnknownStringOperator, "unknown operator '%v' for strings"},
	{SharkErrorUnknownBoolOperator, "unknown operator '%v' for booleans"},
	{SharkErrorVMStackOverflow, "stack overflow"},
	{SharkErrorVMFrameStackOverflow, "frame stack overflow"},
	{SharkErrorNonFunction, "not a function '%v'"},
	{SharkErrorNonFunctionCall, "cannot call non-function '%v'"},
	{SharkErrorArgumentNumberMismatch, "expected %v arguments, but got %v"},
	{SharkErrorNonHashable, "non-hashable type '%v'"},
	{SharkErrorNonIndexable, "non-indexable type '%v'"},
	{SharkErrorNoDefaultValue, "no value to set for function default parameter"},
	{SharkErrorImmutableValue, "cannot modify immutable value '%v'"},
	{SharkErrorIndexOutOfBounds, "index out of bounds '%v'"},
	{SharkErrorOptionalParameter, "parameter '%v' cannot be after an optional parameter"},
	{SharkErrorDivisionByZero, "division by zero"},
	{SharkErrorUnknownType, "unknown type '%v'"},
	{SharkErrorTypeMismatch, "type mismatch '%v'"},
	{SharkErrorTupleDeconstructMismatch, "cannot deconstruct a tuple with %v elements into %v variables"},
	{SharkErrorInvalidNumber, "invalid number"},
	{SharkErrorArgumentCount, "wrong number of arguments"},
	{SharkErrorNotCallable, "not callable"},
	{SharkErrorTypeNotFound, "type '%v' not found"},
	{SharkErrorTypeSyntax, "syntax error: %v"},
}
