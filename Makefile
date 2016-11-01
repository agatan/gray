TARGET := gray
SRCS := $(shell find . -name "*.go")

.PHONY: all
all: $(TARGET) parser/parser.go

.PHONY: run
run: all
	./$(TARGET)

.PHONY: test
test: all
	go test ./...

.PHONY: clean
clean:
	$(RM) $(TARGET) parser/parser.go

$(TARGET): $(SRCS) parser/parser.go
	go build

parser/parser.go: parser/parser.go.y
	go generate ./...
