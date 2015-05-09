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
	STOP // End
} InstructionSet;

typedef enum {
	R1, R2,  R3,  R4,  R5,  R6,  R7,  R8,
	R9, R10, R11, R12, R13, R14, R15, R16,
	NUM_OF_REGISTERS
} Registers;

int regs[NUM_OF_REGISTERS];

// Simple decreasing while loop
/*const int program[] = {
	SET, R1, 50,
	SET, R2, 5,
	SUB, R1, R2,
	SHOW, R1,
	JNZ, R1, 6,
	LOAD, R1, 
	STOP
}; */

// Compute 10 factorial
const int program[] = {
	SET, R1, 1, // Accumulate here
	SET, R2, 10, 
	SET, R3, 1, // For SUB
	MULT, R1, R2,
	SUB, R2, R3,
	JNZ, R2, 6,
	SHOW, R1,
	STOP,
};

int NUM_INSTR = (int)(sizeof(program) / sizeof(program[0]));

// Current instruction pointer.
// Index into the program array.
int ip = 0;

// Stack pointer. -1 indicates not set.
int sp = -1;

// Fixed stack size.
const int STACK_SIZE = 256;
int stack[STACK_SIZE];

bool running = true;

int fetch() {
	return program[ip];
}

// Ensure ++ip is valid.
// Die if it isn't.
void validate_ip(char *func) {
	if (++ip >= NUM_INSTR) {
		printf("*** Not enough args to %s\n", func);
		exit(1);
	}
}

void validate_r(int n) {
	if (n >= NUM_OF_REGISTERS) {
		printf("*** Illegal register: %d\n", n);
		exit(1);
	}
}

void eval(int instr) {
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
				printf("*** Stack overflow");
				exit(1);
			}
			validate_ip("PUSH");
			stack[sp] = program[ip];
			break;
		}
		case POP: {
			if (sp < 0) {
				printf("*** Pop from empty stack\n");
				exit(1);
			}
			int popped = stack[sp--];
			printf("POPPED: %d\n", popped);
			break;
		}
		case ADD: {
			validate_ip("ADD");
			int r_one = program[ip];
			validate_r(r_one);

			validate_ip("ADD");
			int r_two = program[ip];
			validate_r(r_two);

			regs[r_one] = regs[r_one] + regs[r_two];
			break;
		}

		case SUB: {
			validate_ip("SUB");
			int r_one = program[ip];
			validate_r(r_one);

			validate_ip("SUB");
			int r_two = program[ip];
			validate_r(r_two);

			regs[r_one] = regs[r_one] - regs[r_two];
			break;
		}

		case DIV: {
			validate_ip("DIV");
			int r_one = program[ip];
			validate_r(r_one);

			validate_ip("DIV");
			int r_two = program[ip];
			validate_r(r_two);

			regs[r_one] = regs[r_one] / regs[r_two];
			break;
		}
		case MULT: {
			validate_ip("MULT");
			int r_one = program[ip];
			validate_r(r_one);

			validate_ip("MULT");
			int r_two = program[ip];
			validate_r(r_two);

			regs[r_one] = regs[r_one] * regs[r_two];
			break;
		}
		case SET: {
			validate_ip("SET");
			int dest = program[ip];
			validate_r(dest);

			validate_ip("SET");
			int val = program[ip];

			regs[dest] = val;
			break;
		}
		case SHOW: {
			validate_ip("SHOW");
			validate_r(program[ip]);

			int val = regs[program[ip]];
			printf("REG: %d VAL: %d\n", program[ip], val);
			break;
		}
		case MOV: {
			validate_ip("MOV");
			int r_one = program[ip];
			validate_r(r_one);

			validate_ip("MOV");
			int r_two = program[ip];
			validate_r(r_two);

			regs[r_one] = regs[r_two];
			break;
		}
		case LOAD: {
			validate_ip("LOAD");
			int r = program[ip];
			validate_r(r);

			if (++sp >= STACK_SIZE) {
				printf("*** Stack overflow");
				exit(1);
			}
			stack[sp] = regs[r];
			break;
		}
		case STORE: {
			int val = stack[sp--];
			validate_ip("STORE");
			int r = program[ip];
			validate_r(r);

			regs[r] = val;
			break;
		}
		case JMP: {
			validate_ip("JMP");
			int addr = program[ip];
			ip = (addr-1); // Eval loop increments for us
			break;
		}
		case JZ: {
			validate_ip("JZ");
			validate_r(program[ip]);

			int val = regs[program[ip]];
			validate_ip("JZ");
			if (val == 0) {
				int addr = program[ip];
				ip = (addr-1);
			}
			break;
		}
		case JNZ: {
			validate_ip("JNZ");
			validate_r(program[ip]);

			int val = regs[program[ip]];
			validate_ip("JNZ");
			if (val != 0) {
				int addr = program[ip];
				ip = (addr-1);
			}
			break;
		}
	}
}

int main() {
	while (running) {
		eval(fetch());
		ip++;
	}

	return 0;
}
