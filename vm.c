#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

// Instructions set values into the first reg passed in.

typedef enum{
	NOP, // Do nothing
	PUSH, // Push value to stack
	ADD, // Add values in regs and store in first
	SUB, // Sub reg 1 from reg 2 and store in first
	MULT, // Mult values in regs and store in first
	DIV, // Div reg 1 by reg 2, store in first
	POP, // Pop from stack and print
	SET, // Set reg to value
	MOV, // Copy from reg 2 to reg 1
	SHOW, // Print value in reg
	LOAD, // Push from register to stack
	STORE, // Pop from stack and sore in register
	JMP, // Move ip to the value passed in
	JZ, // Compare passed in reg to zero and jump if equal.
	JNZ, // Compare passed in reg to zero and jump if not equal.
	JE, // Jump to the instruction if last compare = 0
	JNE, // Jump to the instruction if last compare != 0
	JLT, // Jump to the instruction if last compare = -1
	JGT, // Jump to the instruction if last compare = 1
	CMP, // Compare reg 1 to reg 2. Store result in z.
	INC, // increment value in reg by 1
	STOP // End
} InstructionSet;

typedef enum {
	R1, R2,  R3,  R4,  R5,  R6,  R7,  R8,
	R9, R10, R11, R12, R13, R14, R15, R16,
	Q,
	NUM_OF_REGISTERS
} Registers;

int64_t regs[NUM_OF_REGISTERS];

int64_t NUM_INSTR;

// Current instruction pointer.
// Index into the program array.
int64_t ip = 0;

// Stack pointer. -1 indicates not set.
int64_t sp = -1;

// Fixed stack size.
const int64_t STACK_SIZE = 256;
int64_t stack[STACK_SIZE];

// Set on CMP instructions.
// -1 for less, 0 for equal, 1 for greater.
int z = 0;

// Set on DIV instructions.
// Contains the quotient.
int q = 1;

bool running = true;

void eval(int64_t program[]) {

	int64_t instr = program[ip];

	switch (instr) {
		case NOP: {
			break;
		}
		case STOP: {
			running = false;
			break;
		}
		case PUSH: {
			if (++sp >= STACK_SIZE) {
				printf("*** Stack overflow\n");
				exit(1);
			}
			ip++;
			stack[sp] = program[ip];
			break;
		}
		case POP: {
			if (sp < 0) {
				printf("*** Pop from empty stack\n");
				exit(1);
			}
			int64_t popped = stack[sp--];
			printf("%lld\n", popped);
			break;
		}
		case ADD: {
			int64_t r_one = program[++ip];
			int64_t r_two = program[++ip];
			regs[r_one] = regs[r_one] + regs[r_two];
			break;
		}

		case SUB: {
			int64_t r_one = program[++ip];
			int64_t r_two = program[++ip];
			regs[r_one] = regs[r_one] - regs[r_two];
			break;
		}

		case DIV: {
			int64_t r_one = program[++ip];
			int64_t r_two = program[++ip];
			regs[Q] = regs[r_one] % regs[r_two];
			regs[r_one] = regs[r_one] / regs[r_two];
			break;
		}
		case MULT: {
			int64_t r_one = program[++ip];
			int64_t r_two = program[++ip];
			regs[r_one] = regs[r_one] * regs[r_two];
			break;
		}
		case SET: {
			int64_t dest = program[++ip];
			int64_t val = program[++ip];
			regs[dest] = val;
			break;
		}
		case SHOW: {
			int64_t val = regs[program[++ip]];
			printf("%lld\n", val);
			break;
		}
		case MOV: {
			int64_t r_one = program[++ip];
			int64_t r_two = program[++ip];
			regs[r_one] = regs[r_two];
			break;
		}
		case LOAD: {
			int64_t r = program[++ip];

			if (++sp >= STACK_SIZE) {
				printf("*** Stack overflow");
				exit(1);
			}
			stack[sp] = regs[r];
			break;
		}
		case STORE: {
			int64_t val = stack[sp--];
			int64_t r = program[++ip];
			regs[r] = val;
			break;
		}
		case JMP: {
			int64_t addr = program[++ip];
			ip = (addr-1); // Eval loop increments for us
			break;
		}
		case JZ: {
			int64_t val = regs[program[++ip]];
			ip++;
			if (val == 0) {
				int64_t addr = program[ip];
				ip = (addr-1);
			}
			break;
		}
		case JNZ: {
			int64_t val = regs[program[++ip]];
			ip++;
			if (val != 0) {
				int64_t addr = program[ip];
				ip = (addr-1);
			}
			break;
		}
		case CMP: {
			int64_t val_one = regs[program[++ip]];
			int64_t val_two = regs[program[++ip]];

			z = val_one < val_two ? -1 :
				val_one == val_two ? 0 : 1;
			break;
		}
		case JE: {
			ip++;
			if (z == 0) {
				ip = (program[ip]-1);
			}
			break;
		}
		case JNE: {
			ip++;
			if (z != 0) {
				ip = (program[ip]-1);
			}
			break;
		}
		case JLT: {
			ip++;
			if (z == -1) {
				ip = (program[ip]-1);
			}
			break;
		}
		case JGT: {
			ip++;
			if (z == 1) {
				ip = (program[ip]-1);
			}
			break;
		}
		case INC: {
			ip++;
			regs[program[ip]] += 1;
			break;
		}
	}
}

int main(int argc, char *argv[]) {

	if (argc < 2) {
		printf("File required\n");
        	return 1;
	}

    	FILE *fp = fopen(argv[1], "rb");
    	if (fp == NULL) {
        	printf("Failed to open: %s\n", argv[1]);
        	return 1;
    	}

	fseek(fp, 0, SEEK_END);
        NUM_INSTR = ftell(fp);
        fseek(fp, 0, SEEK_SET);

	if (NUM_INSTR % 8 != 0) {
		printf("Invalid binary size: %lld\n", NUM_INSTR);
		exit(1);
	}

	int64_t program[NUM_INSTR];

	int64_t instruction;
	int i = 0;
    	while (fread(&instruction, sizeof(instruction), 1, fp) == 1) {
		program[i++] = instruction;
	}

	while (running) {
		eval(program);
		ip++;
	}

	return 0;
}
