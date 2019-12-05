package debug_log

import (
	"fmt"
	"os"
)

var Debug bool

func init() {
	debug := os.Getenv("DEBUG")
	if debug != "" && debug != "false" && debug != "0" {
		Debug = true
	}
}

func Info(error string, prefix string) {
	if Debug {
		fmt.Printf("[INFO][%s] %s\n", prefix, error)
	}
}

func Warning(error string, prefix string) {
	if Debug {
		fmt.Printf("[WARNING][%s] %s\n", prefix, error)
	}
}

func Error(error string, prefix string) {
	if Debug {
		fmt.Printf("[ERROR][%s] %s\n", prefix, error)
	}
}
