MOV R1 $2
# Accumulate
MOV R11 $1

CALL FACTORIAL
SHOW R11
STOP

!FACTORIAL CMP R1 $10
JGT END
MULT R11 R1
INC R1
CALL FACTORIAL
!END RET
