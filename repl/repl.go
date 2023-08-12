package repl

import (
	"io"
	"shark/emitter"

	"github.com/c-bata/go-prompt"
)

const PROMPT = ">>> "

func Start(textOut io.Writer) {
	in := "stdin"
	sharkEmitter := emitter.New(&in, textOut)

	for {
		p := prompt.New(sharkEmitter.Interpret, completer,
			prompt.OptionPrefix(PROMPT),
			prompt.OptionTitle("Shark"),
			prompt.OptionCompletionOnDown(),
			prompt.OptionPrefixTextColor(prompt.Blue),
			prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSuggestionBGColor(prompt.DarkGray))

		p.Run()
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), false)
}
