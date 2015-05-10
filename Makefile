SRC_FILES = vm.c
CC_FLAGS = -Wall -Wextra -g -std=c11
CC = clang

all:
	${CC} ${SRC_FILES} ${CC_FLAGS} -o vm
	go build assembler.go

clean:
	rm -f vm
	rm -f assembler
