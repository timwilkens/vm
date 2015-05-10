package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var opCodes map[string]int64 = map[string]int64{
	"NOP":   0,
	"PUSH":  1,
	"ADD":   2,
	"ADDV":  3, // PRIVATE
	"SUB":   4,
	"SUBV":  5, // PRIVATE
	"MULT":  6,
	"MULTV": 7, // PRIVATE
	"DIV":   8,
	"DIVV":  9, // PRIVATE
	"POP":   10,
	"SET":   11,
	"MOV":   12,
	"SHOW":  13,
	"LOAD":  14,
	"STORE": 15,
	"JMP":   16,
	"JZ":    17,
	"JNZ":   18,
	"JE":    19,
	"JNE":   20,
	"JLT":   21,
	"JGT":   22,
	"CMP":   23,
	"INC":   24,
	"DEC":   25,
	"PRINT": 26,
	"STOP":  27,
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
	"Z":   17,
}

var numRegisters = len(regs)

var jmpCodes []int64
var jmpLabels map[string]int64

func regError(r string) error {
	return errors.New(fmt.Sprintf("Invalid register: %s", r))
}

func stripValue(s string) (int64, error) {
	s = strings.TrimPrefix(s, "$")
	val, err := strconv.Atoi(s)
	return int64(val), err
}

var isLabel = regexp.MustCompile(`[A-Z]$`)

func toIntCodes(parts []string) ([]int64, error) {
	op := parts[0]

	var codes []int64
	var err error

	switch op {
	case "NOP", "POP", "STOP":
		codes, err = parseNoArg(parts)
	case "JMP", "JE", "JNE", "JLT", "JGT":
		codes, err = parseAddr(parts)
		if err == nil {
			jmpCodes = append(jmpCodes, codes[1])
		}
	case "SHOW", "LOAD", "STORE", "INC", "DEC", "PRINT":
		codes, err = parseReg(parts)
	case "PUSH", "SET":
		codes, err = parseRegAndVal(parts)
	case "JZ", "JNZ":
		codes, err = parseRegAndAddr(parts)
		if err == nil {
			jmpCodes = append(jmpCodes, codes[2])
		}
	case "ADD", "SUB", "MULT", "DIV":
		codes, err = parseArithmetic(parts)
	case "MOV", "CMP":
		codes, err = parseTwoReg(parts)
	default:
		codes, err = nil, errors.New(fmt.Sprintf("Unknown op: %s", op))
	}

	if err != nil {
		return nil, err
	}

	return codes, err
}

func parseNoArg(parts []string) ([]int64, error) {
	if len(parts) != 1 {
		return nil, errors.New("Invalid arguments")
	}
	return []int64{opCodes[parts[0]]}, nil
}

func parseAddr(parts []string) ([]int64, error) {
	if len(parts) != 2 {
		return nil, errors.New("Invalid arguments")
	}

	// Jmp label
	if isLabel.MatchString(parts[1]) {
		if addr, ok := jmpLabels[parts[1]]; ok {
			return []int64{opCodes[parts[0]], addr}, nil
		} else {
			return nil, errors.New(fmt.Sprintf("Invalid jmp label: %s", parts[1]))
		}
	} else {
		fmt.Println("NON LABEL")
		val, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		return []int64{opCodes[parts[0]], int64(val)}, nil
	}

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

func parseRegAndAddr(parts []string) ([]int64, error) {
	if len(parts) != 3 {
		return nil, errors.New("Invalid arguments")
	}
	r := parts[1]

	if reg, ok := regs[r]; ok {
		if isLabel.MatchString(parts[2]) {
			if addr, ok := jmpLabels[parts[2]]; ok {
				instrs := []int64{
					opCodes[parts[0]],
					reg,
					addr}
				return instrs, nil
			} else {
				return nil, errors.New(fmt.Sprintf("Invalid jmp label: %s", parts[1]))
			}
		} else {
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
		}
	} else {
		return nil, regError(r)
	}
}

func parseRegAndVal(parts []string) ([]int64, error) {
	if len(parts) != 3 {
		return nil, errors.New("Invalid arguments")
	}
	r := parts[1]

	if reg, ok := regs[r]; ok {
		if strings.HasPrefix(parts[2], "$") {
			val, err := stripValue(parts[2])
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
			return nil, errors.New(fmt.Sprintf("Malformed value: %s", parts[2]))
		}
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

func parseArithmetic(parts []string) ([]int64, error) {
	if len(parts) != 3 {
		return nil, errors.New("Invalid arguments")
	}
	r1 := parts[1]

	if reg1, ok := regs[r1]; ok {
		// Use value rather than register
		if strings.HasPrefix(parts[2], "$") {
			val, err := stripValue(parts[2])
			if err != nil {
				return nil, err
			}
			op := (parts[0] + "V")
			instrs := []int64{
				opCodes[op],
				reg1,
				int64(val),
			}
			return instrs, nil

		} else {
			r2 := parts[2]
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

var labelPrefix = "!"
var commentPrefix = "#"

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

	// Record the offsets for all instructions.
	// Any jump instruction must be to one of these.
	jmpPoints := make(map[int64]bool)

	lineNumber := 0

	var lines []string

	// Init global labels
	jmpLabels = make(map[string]int64)

	instructionNum := int64(0)

	// Loop through once and record label locations
	// for expansion later
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")
		lineNumber++

		if err != nil {
			break
		}

		// Support blank lines and comments
		if line == "" || strings.HasPrefix(line, commentPrefix) {
			// Needed for correct line numbers
			lines = append(lines, line)
			continue
		}

		jmpPoints[instructionNum] = true

		// Record ip location for labels
		// Strip label from instruction
		if strings.HasPrefix(line, labelPrefix) {
			parts := strings.Split(line, " ")
			label := parts[0]
			label = strings.TrimPrefix(label, labelPrefix)
			jmpLabels[label] = instructionNum
			line = strings.Join(parts[1:], " ")
		}

		lines = append(lines, line)

		instructionNum += int64(len(strings.Split(line, " ")))
	}

	var instructions []int64

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, " ")
		codes, err := toIntCodes(parts)
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
