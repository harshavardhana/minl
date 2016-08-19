CC=gcc
CFLAGS=-Wall

all: sandbox

include syscall-reporter.mk

sandbox: sandbox.o syscall-reporter.o

.PHONY: clean
clean:
	rm -f sandbox sandbox.o syscall-reporter.o

