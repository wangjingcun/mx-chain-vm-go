package elrondapigenerate

import (
	"bufio"
	"fmt"
	"os"
)

// WriteWASMOpcodeCostConfigHelpers generates code for config.txt
// (to be copied manually in config/config.toml)
func WriteWASMOpcodeCostConfigHelpers(out *eiGenWriter) {
	out.WriteString(`// Code generated by elrondapi generator. DO NOT EDIT.

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!!!!!!!!!!!! AUTO-GENERATED FILE !!!!!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

// !!!!!!!!!!!!!!!!!! USE IN config/config.toml !!!!!!!!!!!!!!!!!!!!

`)
	out.WriteString("[WASMOpcodeCost]\n")
	readFile, err := os.Open("generate/cmd/input/wasm_opcodes.txt")
	if err != nil {
		panic(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		opcode := fileScanner.Text()
		out.WriteString(fmt.Sprintf("    %s = 1\n", opcode))
	}
}
