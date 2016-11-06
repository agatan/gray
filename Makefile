TARGET := gray
SRCS := $(shell find . -name "*.go") parser/parser.go.y
GOFLAGS := -ldflags "-s"

.PHONY: all
all: $(TARGET)

.PHONY: run
run: $(TARGET)
	./$(TARGET)

.PHONY: test
test: all
	go test $(GOFLAGS) ./...

.PHONY: clean
clean:
	$(RM) $(TARGET) parser/parser.go

$(TARGET): $(SRCS)
	go generate ./...
	go build $(GOFLAGS)

