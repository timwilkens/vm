#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

// Instructions set values into the first reg passed in.

typedef enum{
	NOP, // Do nothing
	PUSH, // Push value to stack
	ADD, // Add values in regs and store in first
	ADDV,
	SUB, // Sub reg 1 from reg 2 and store in first
	SUBV,
	MULT, // Mult values in regs and store in first
	MULTV,
	DIV, // Div reg 1 by reg 2, store in first
	DIVV,
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
	DEC, // decrement value in reg by 1
	STOP // End
} InstructionSet;

typedef enum {
	R1, R2,  R3,  R4,  R5,  R6,  R7,  R8,
	R9, R10, R11, R12, R13, R14, R15, R16,
	Q, // Store remainder from DIV calls
	Z, // Store result of last CMP
	NUM_OF_REGISTERS
} Registers;

int64_t regs[NUM_OF_REGISTERS];

int64_t NUM_INSTR;

// Current instruction pointer.
// Index into the program array.
int64_t IP = 0;

// Stack pointer. -1 indicates not set.
int64_t SP = -1;

// Fixed stack size.
const int64_t STACK_SIZE = 256;
int64_t stack[STACK_SIZE];


void run(int64_t program[]) {

	bool running = true;

	while (running) {
		int64_t instr = program[IP];

		switch (instr) {
			case NOP: {
				break;
			}
			case STOP: {
				running = false;
				break;
			}
			case PUSH: {
				if (++SP >= STACK_SIZE) {
					printf("*** Stack overflow\n");
					exit(1);
				}
				IP++;
				stack[SP] = program[IP];
				break;
			}
			case POP: {
				if (SP < 0) {
					printf("*** Pop from empty stack\n");
					exit(1);
				}
				int64_t popped = stack[SP--];
				printf("%lld\n", popped);
				break;
			}
			case ADD: {
				int64_t r_one = program[++IP];
				int64_t r_two = program[++IP];
				regs[r_one] = regs[r_one] + regs[r_two];
				break;
			}
			case ADDV: {
				int64_t r_one = program[++IP];
				int64_t val = program[++IP];
				regs[r_one] = regs[r_one] + val;
				break;
			}
			case SUB: {
				int64_t r_one = program[++IP];
				int64_t r_two = program[++IP];
				regs[r_one] = regs[r_one] - regs[r_two];
				break;
			}
			case SUBV: {
				int64_t r_one = program[++IP];
				int64_t val = program[++IP];
				regs[r_one] = regs[r_one] - val;
				break;
			}
			case DIV: {
				int64_t r_one = program[++IP];
				int64_t r_two = program[++IP];
				regs[Q] = regs[r_one] % regs[r_two];
				regs[r_one] = regs[r_one] / regs[r_two];
				break;
			}
			case DIVV: {
				int64_t r_one = program[++IP];
				int64_t val = program[++IP];
				regs[Q] = regs[r_one] % val;
				regs[r_one] = regs[r_one] / val;
				break;
			}
			case MULT: {
				int64_t r_one = program[++IP];
				int64_t r_two = program[++IP];
				regs[r_one] = regs[r_one] * regs[r_two];
				break;
			}
			case MULTV: {
				int64_t r_one = program[++IP];
				int64_t val = program[++IP];
				regs[r_one] = regs[r_one] * val;
				break;
			}
			case SET: {
				int64_t dest = program[++IP];
				int64_t val = program[++IP];
				regs[dest] = val;
				break;
			}
			case SHOW: {
				int64_t val = regs[program[++IP]];
				printf("%lld\n", val);
				break;
			}
			case MOV: {
				int64_t r_one = program[++IP];
				int64_t r_two = program[++IP];
				regs[r_one] = regs[r_two];
				break;
			}
			case LOAD: {
				int64_t r = program[++IP];
	
				if (++SP >= STACK_SIZE) {
					printf("*** Stack overflow");
					exit(1);
				}
				stack[SP] = regs[r];
				break;
			}
			case STORE: {
				int64_t val = stack[SP--];
				int64_t r = program[++IP];
				regs[r] = val;
				break;
			}
			case JMP: {
				int64_t addr = program[++IP];
				IP = (addr-1); // Eval loop increments for us
				break;
			}
			case JZ: {
				int64_t val = regs[program[++IP]];
				IP++;
				if (val == 0) {
					int64_t addr = program[IP];
					IP = (addr-1);
				}
				break;
			}
			case JNZ: {
				int64_t val = regs[program[++IP]];
				IP++;
				if (val != 0) {
					int64_t addr = program[IP];
					IP = (addr-1);
				}
				break;
			}
			case CMP: {
				int64_t val_one = regs[program[++IP]];
				int64_t val_two = regs[program[++IP]];
	
				regs[Z] = val_one < val_two ? -1 :
					val_one == val_two ? 0 : 1;
				break;
			}
			case JE: {
				IP++;
				if (regs[Z] == 0) {
					IP = (program[IP]-1);
				}
				break;
			}
			case JNE: {
				IP++;
				if (regs[Z] != 0) {
					IP = (program[IP]-1);
				}
				break;
			}
			case JLT: {
				IP++;
				if (regs[Z] == -1) {
					IP = (program[IP]-1);
				}
				break;
			}
			case JGT: {
				IP++;
				if (regs[Z] == 1) {
					IP = (program[IP]-1);
				}
				break;
			}
			case INC: {
				IP++;
				regs[program[IP]] += 1;
				break;
			}
			case DEC: {
				IP++;
				regs[program[IP]] -= 1;
				break;
			}
		}
		IP++;
	}
}

int64_t *load_program(char *file) {
    	FILE *fp = fopen(file, "rb");
    	if (fp == NULL) {
        	printf("Failed to open: %s\n", file);
        	exit(1);
    	}

	fseek(fp, 0, SEEK_END);
        NUM_INSTR = ftell(fp);
        fseek(fp, 0, SEEK_SET);

	if (NUM_INSTR % 8 != 0) {
		printf("Invalid binary size: %lld\n", NUM_INSTR);
		exit(1);
	}

	int64_t *program = malloc(NUM_INSTR*sizeof(int64_t));

	int64_t instruction;
	int i = 0;
    	while (fread(&instruction, sizeof(instruction), 1, fp) == 1) {
		program[i++] = instruction;
	}

	return program;
}

int main(int argc, char *argv[]) {

	if (argc < 2) {
		printf("File required\n");
        	return 1;
	}

	int64_t *program = load_program(argv[1]);
	run(program);
	free(program);

	return 0;
}
