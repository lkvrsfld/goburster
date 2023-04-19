package main

import (
	"os"
	"os/exec"
    "runtime"
)

var runtimeos string

var clear map[string]func()

func InitClear() {
    runtimeos = runtime.GOOS
}

func clear() {
    var clearCmd
    clearCmd := []string{["clear"]}

    if runtimeos == "windows" {
        clearCmd = ["cmd", "/c", "cls"]
    }
    cmd := exec.Command("clear") 
    cmd.Stdout = os.Stdout
    cmd.Run()
}