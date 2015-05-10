MOV R1 $10
MOV R2 $30
# The assembler should catch this.
# Jumps must go to the start of an instruction.
CMP R1 R2
JMP 1
STOP
