# Iteration var
SET R1 1
# Accumulator
SET R2 0
# Stop condition
SET R3 1000
SET R5 0

# Store in Regs for now until we support OPs on values
SET R6 3
SET R7 5

# See if divisible by 3
# Copy into R4 to overwrite
MOV R4 R1
DIV R4 R6
# Skip add if not divisible
JNZ Q 32

# Divisible by three. Add to accumulator.
ADD R2 R1
# Skip checking by 5
JMP 44

# Check if divisible by 5
MOV R4 R1
DIV R4 R7
JNZ Q 44

ADD R2 R1

# Increment and stop looping when >= 1000
INC R1
CMP R1 R3
JLT 18
# Dump out answer
SHOW R2
STOP
