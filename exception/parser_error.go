package exception

type SharkParserError struct {
	ErrMsg     string
	ErrCause   []SharkErrorCause
	ErrHelpMsg string
	ErrCode    SharkErrorCode
}

func (pe *SharkParserError) Cause() []SharkErrorCause {
	return pe.ErrCause
}

func (pe *SharkParserError) Msg() string {
	return pe.ErrMsg
}

func (pe *SharkParserError) SuggestionMsg() string {
	return pe.ErrHelpMsg
}

func (pe *SharkParserError) Code() SharkErrorCode {
	return pe.ErrCode
}
