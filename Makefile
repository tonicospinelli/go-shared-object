# file names
GO_SRC = cmd/cshared/main.go
SO_TARGET = libcsv.so
C_SRC = libcsv.c
C_TARGET = libcsv

# compile commands
GO_BUILD = go build
GCC = gcc
RM = rm -f

# flags
GO_FLAGS = -buildmode=c-shared
C_FLAGS = -L.

# default target
all: $(SO_TARGET) $(C_TARGET)

# compile the shared library from Golang code
$(SO_TARGET): $(GO_SRC)
	$(GO_BUILD) $(GO_FLAGS) -o $(SO_TARGET) $(GO_SRC)

# compile the C program linking with the shared library
$(C_TARGET): $(SO_TARGET) $(C_SRC)
	$(GCC) -o $(C_TARGET) $(C_SRC) $(C_FLAGS) -lcsv

# clean built files
clean:
	$(RM) $(SO_TARGET) $(C_TARGET)

# target to execute the program
run: all
	@LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:." ./$(C_TARGET)

.PHONY: all clean run
