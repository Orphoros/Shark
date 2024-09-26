package exception

import (
	"fmt"
	"os"
)

func PrintExitMsg(msg string, exitCode int) {
	fmt.Printf("error: %s\n", msg)
	os.Exit(exitCode)
}

func PrintExitMsgCtx(msg, ctx string, exitCode int) {
	fmt.Printf("error: %s\n", msg)
	fmt.Printf("   --> %s\n", ctx)
	os.Exit(exitCode)
}
