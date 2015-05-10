# Iteration var
MOV R1 $1
# Accumulator
MOV R2 $0
# Stop condition
MOV R3 $1000

# See if divisible by 3
# Copy into R4 to overwrite
!LOOP MOV R4 R1
DIV R4 $3
# Skip add if not divisible
JNZ Q FIVE

# Divisible by three. Add to accumulator.
ADD R2 R1
# Skip checking by 5
JMP ENDLOOP

# Check if divisible by 5
!FIVE MOV R4 R1
DIV R4 $5
JNZ Q ENDLOOP

ADD R2 R1

# Increment and stop looping when >= 1000
!ENDLOOP INC R1
CMP R1 R3
JLT LOOP
# Dump out answer
SHOW R2
STOP
