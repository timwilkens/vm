# Iteration var
SET R1 $1
# Accumulator
SET R2 $0
# Stop condition
SET R3 $1000

# See if divisible by 3
# Copy into R4 to overwrite
MOV R4 R1
DIV R4 $3
# Skip add if not divisible
JNZ Q 23

# Divisible by three. Add to accumulator.
ADD R2 R1
# Skip checking by 5
JMP 35

# Check if divisible by 5
MOV R4 R1
DIV R4 $5
JNZ Q 35

ADD R2 R1

# Increment and stop looping when >= 1000
INC R1
CMP R1 R3
JLT 9
# Dump out answer
SHOW R2
STOP
