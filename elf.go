package libbpfgo

import (
	"debug/elf"
	"fmt"
)

func symbolOffset(f *elf.File, funcName string) uint32 {
	if f == nil {
		panic("need init elf file by call ")
	}

	regularSymbols, regularSymbolsErr := f.Symbols()
	dynamicSymbols, dynamicSymbolsErr := f.DynamicSymbols()

	// Only if we failed getting both regular and dynamic symbols - then we abort.
	if regularSymbolsErr != nil && dynamicSymbolsErr != nil {
		panic(fmt.Errorf("could not open regular or dynamic symbol sections to resolve symbol offset: %w %s", regularSymbolsErr, dynamicSymbolsErr))
	}

	// Concatenating into a single list.
	// The list can have duplications, but we will find the first occurrence which is sufficient.
	syms := append(regularSymbols, dynamicSymbols...)

	sectionsToSearchForSymbol := []*elf.Section{}

	for i := range f.Sections {
		if f.Sections[i].Flags == elf.SHF_ALLOC+elf.SHF_EXECINSTR {
			fmt.Printf("%s\n", f.Sections[i].Name)
			sectionsToSearchForSymbol = append(sectionsToSearchForSymbol, f.Sections[i])
		}
	}
	var executableSection *elf.Section
	for j := range syms {
		if syms[j].Name == funcName {
			// Find what section the symbol is in by checking the executable section's
			// addr space.
			for m := range sectionsToSearchForSymbol {
				if syms[j].Value > sectionsToSearchForSymbol[m].Addr &&
					syms[j].Value < sectionsToSearchForSymbol[m].Addr+sectionsToSearchForSymbol[m].Size {
					executableSection = sectionsToSearchForSymbol[m]
				}
			}

			if executableSection == nil {
				panic("could not find symbol in executable sections of binary")
			}

			return uint32(syms[j].Value - executableSection.Addr + executableSection.Offset)
		}
	}
	panic("no such function in this elf file!")
}
