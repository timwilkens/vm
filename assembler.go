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
	"STOP":  20,
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
}

var numRegisters = 16

func regError(r string) error {
	return errors.New(fmt.Sprintf("Invalid register: %s", r))
}

func toIntCodes(line string) ([]int64, error) {
	parts := strings.Split(line, " ")
	op := parts[0]

	switch op {
	case "NOP", "POP", "STOP":
		return parseNoArg(parts)
	case "JMP", "JE", "JNE", "JLT", "JGT":
		return parseValue(parts)
	case "SHOW", "LOAD", "STORE":
		return parseReg(parts)
	case "PUSH", "SET", "JZ", "JNZ":
		return parseRegAndVal(parts)
	case "ADD", "SUB", "MULT", "DIV", "MOV", "CMP":
		return parseTwoReg(parts)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown op: %s", op))
	}
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

func main() {
	in := flag.String("in", "", "Input file")
	out := flag.String("out", "", "Output file")
	flag.Parse()

	if *in == "" {
		maybeDie(errors.New("in is required"))
	}

	if *out == "" {
		maybeDie(errors.New("out is required"))
	}

	file, err := os.Open(*in)
	maybeDie(err)

	defer file.Close()
	reader := bufio.NewReader(file)

	var instructions []int64

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")

		if err != nil {
			break
		}
		codes, err := toIntCodes(line)
		maybeDie(err)

		for _, c := range codes {
			instructions = append(instructions, c)
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
