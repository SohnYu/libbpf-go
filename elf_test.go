package main

import (
	"debug/elf"
	"fmt"
	"testing"
)

func TestSymbols(t *testing.T) {
	open, err := elf.Open("/home/ysh/Nexus/dockerFile/Gindemo/main")
	if err != nil {
		panic(err)
	}
	fmt.Println(open)
}
