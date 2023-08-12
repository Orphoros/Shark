package exception

type SharkLexerError struct {
	ErrMsg     string
	ErrCause   []SharkErrorCause
	ErrHelpMsg string
	ErrCode    SharkErrorCode
}

func (le *SharkLexerError) Cause() []SharkErrorCause { return le.ErrCause }
func (le *SharkLexerError) Msg() string              { return le.ErrMsg }
func (le *SharkLexerError) SuggestionMsg() string    { return le.ErrHelpMsg }
func (le *SharkLexerError) Code() SharkErrorCode     { return le.ErrCode }
