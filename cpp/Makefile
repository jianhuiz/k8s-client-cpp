# Declaration of variables
CC = g++
CC_FLAGS = -I . -std=c++14
LD_FLAGS = ../go/go.a -lprotobuf -lpthread
 
# File names
EXEC = k8s-cpp-client
SOURCES = $(shell find . -iname "*.cc")
OBJECTS = $(SOURCES:.cc=.o)
 
# Main target
$(EXEC): $(OBJECTS)
	$(CC) $(OBJECTS) $(LD_FLAGS) -o $(EXEC)
 
# To obtain object files
%.o: %.cc
	$(CC) -c $(CC_FLAGS) $< -o $@
 
# To remove generated files
clean:
	rm -f $(EXEC) $(OBJECTS)
