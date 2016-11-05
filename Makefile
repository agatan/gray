TARGET := gray
SRCS := $(shell find . -name "*.go") parser/parser.go.y

.PHONY: all
all: $(TARGET)

.PHONY: run
run: $(TARGET)
	./$(TARGET)

.PHONY: test
test: all
	go test ./...

.PHONY: clean
clean:
	$(RM) $(TARGET) parser/parser.go

$(TARGET): $(SRCS)
	go generate ./...
	go build

