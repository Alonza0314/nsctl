package main

import (
	"runtime"

	"github.com/Alonza0314/nsctl/cmd"
)

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cmd.Execute()
}
