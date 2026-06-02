package util

import (
	"fmt"
	"os"
	"runtime"
)

func AbnormalExit(err error) {
	fmt.Println("abnormal exit: " + err.Error())
	PrintStack()
	os.Exit(1)
}

func PrintStack() {
	pc := make([]uintptr, 10)
	n := runtime.Callers(0, pc)
	for i := 0; i < n; i++ {
		funcName := runtime.FuncForPC(pc[i])
		file, line := funcName.FileLine(pc[i])
		fmt.Printf("%s:%d %s\n", file, line, funcName.Name())
	}
}
