PRINT 'H'
PRINT 'E'
PRINT 'L'
PRINT 'L'
PRINT 'O'

# Parsing can't handle spaces
MOV R1 $32
PRINT R1

PRINT 'W'
PRINT 'O'
PRINT 'R'
PRINT 'L'
PRINT 'D'
PRINT '!'

# Don't support metacharacters yet
MOV R1 $10
PRINT R1
STOP
