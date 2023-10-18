package main

import (
	"testing"
)

func TestLoadMapFromPath(t *testing.T) {
	module := NewModuleFromFile("./main.bpf.o")
	module.LoadPinnedMapFromPath("heap")
}
