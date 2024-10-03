package exception

import (
	"fmt"
	"os"
)

func PrintExitMsgCtx(msg, ctx string, exitCode int) {
	fmt.Printf("error: %s\n", msg)
	fmt.Printf("   --> %s\n", ctx)
	os.Exit(exitCode)
}
