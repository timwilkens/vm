package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var opCodes map[string]int64 = map[string]int64{
	"NOP":   0,
	"PUSH":  1,
	"ADD":   2,
	"SUB":   3,
	"MULT":  4,
	"DIV":   5,
	"POP":   6,
	"SET":   7,
	"MOV":   8,
	"SHOW":  9,
	"LOAD":  10,
	"STORE": 11,
	"JMP":   12,
	"JZ":    13,
	"JNZ":   14,
	"JE":    15,
	"JNE":   16,
	"JLT":   17,
	"JGT":   18,
	"CMP":   19,
	"INC":   20,
	"STOP":  21,
}

var regs map[string]int64 = map[string]int64{
	"R1":  0,
	"R2":  1,
	"R3":  2,
	"R4":  3,
	"R5":  4,
	"R6":  5,
	"R7":  6,
	"R8":  7,
	"R9":  8,
	"R10": 9,
	"R11": 10,
	"R12": 11,
	"R13": 12,
	"R14": 13,
	"R15": 14,
	"R16": 15,
	"Q":   16,
}

var numRegisters = len(regs)

var jmpCodes []int64

func regError(r string) error {
	return errors.New(fmt.Sprintf("Invalid register: %s", r))
}

func toIntCodes(line string) ([]int64, error) {
	parts := strings.Split(line, " ")
	op := parts[0]

	var codes []int64
	var err error

	switch op {
	case "NOP", "POP", "STOP":
		codes, err = parseNoArg(parts)
	case "JMP", "JE", "JNE", "JLT", "JGT":
		codes, err = parseValue(parts)
	case "SHOW", "LOAD", "STORE", "INC":
		codes, err = parseReg(parts)
	case "PUSH", "SET", "JZ", "JNZ":
		codes, err = parseRegAndVal(parts)
	case "ADD", "SUB", "MULT", "DIV", "MOV", "CMP":
		codes, err = parseTwoReg(parts)
	default:
		codes, err = nil, errors.New(fmt.Sprintf("Unknown op: %s", op))
	}

	if err != nil {
		return nil, err
	}

	// Store jump addresses.
	// Check after all codes have been parsed.
	switch op {
	case "JMP", "JE", "JNE", "JTL", "JGT":
		jmpCodes = append(jmpCodes, codes[1])
	case "JZ", "JNZ":
		jmpCodes = append(jmpCodes, codes[2])
	}

	return codes, err
}

func parseNoArg(parts []string) ([]int64, error) {
	if len(parts) != 1 {
		return nil, errors.New("Invalid arguments")
	}
	return []int64{opCodes[parts[0]]}, nil
}

func parseValue(parts []string) ([]int64, error) {
	if len(parts) != 2 {
		return nil, errors.New("Invalid arguments")
	}

	val, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	return []int64{opCodes[parts[0]], int64(val)}, nil
}

func parseReg(parts []string) ([]int64, error) {
	if len(parts) != 2 {
		return nil, errors.New("Invalid arguments")
	}

	if reg, ok := regs[parts[1]]; ok {
		return []int64{opCodes[parts[0]], reg}, nil

	} else {
		return nil, regError(parts[1])
	}
}

func parseRegAndVal(parts []string) ([]int64, error) {
	if len(parts) != 3 {
		return nil, errors.New("Invalid arguments")
	}
	r := parts[1]

	if reg, ok := regs[r]; ok {
		val, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, err
		}

		instrs := []int64{
			opCodes[parts[0]],
			reg,
			int64(val),
		}
		return instrs, nil
	} else {
		return nil, regError(r)
	}
}

func parseTwoReg(parts []string) ([]int64, error) {
	if len(parts) != 3 {
		return nil, errors.New("Invalid arguments")
	}
	r1 := parts[1]
	r2 := parts[2]

	if reg1, ok := regs[r1]; ok {
		if reg2, ok := regs[r2]; ok {

			instrs := []int64{
				opCodes[parts[0]],
				reg1,
				reg2,
			}
			return instrs, nil
		} else {
			return nil, regError(r2)
		}
	} else {
		return nil, regError(r1)
	}
}

func maybeDie(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func die(s string) {
	fmt.Println(s)
	os.Exit(1)
}

func main() {
	in := flag.String("in", "", "Input file")
	out := flag.String("out", "", "Output file")
	flag.Parse()

	if *in == "" {
		die("in is required")
	}

	if *out == "" {
		die("out is required")
	}

	file, err := os.Open(*in)
	maybeDie(err)

	defer file.Close()
	reader := bufio.NewReader(file)

	var instructions []int64

	// Record the offsets for all instructions.
	// Any jump instruction must be to one of these.
	jmpPoints := make(map[int64]bool)

	lineNumber := 0

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")
		lineNumber++

		if err != nil {
			break
		}

		// Support blank lines
		if line == "" {
			continue
		}

		// Skip Comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		jmpPoints[int64(len(instructions))] = true

		codes, err := toIntCodes(line)
		if err != nil {
			die(fmt.Sprintf("Line %d - ERROR: %s", lineNumber, err.Error()))
		}

		for _, c := range codes {
			instructions = append(instructions, c)
		}
	}

	// Check that jumps are valid.
	for _, jc := range jmpCodes {
		if !jmpPoints[jc] {
			die(fmt.Sprintf("Invalid jump addr: %d", jc))
		}
	}

	binFile, err := os.Create(*out)
	maybeDie(err)
	defer binFile.Close()

	buf := bufio.NewWriter(binFile)
	defer buf.Flush()

	for _, instr := range instructions {
		err = binary.Write(buf, binary.LittleEndian, instr)
		maybeDie(err)
	}
}
