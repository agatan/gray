TARGET := gray
SRCS := $(shell find . -name "*.go") parser/parser.go.y
GOFLAGS := -ldflags "-s"

RUNTIME_TARGET := libruntime.a
RUNTIME_SRCS := $(shell find ./runtime/src -name "*.rs")

.PHONY: all
all: $(TARGET) $(RUNTIME_TARGET)

.PHONY: run
run: $(TARGET)
	./$(TARGET)

.PHONY: test
test: all
	go test $(GOFLAGS) ./...

.PHONY: clean
clean: runtime_clean
	$(RM) $(TARGET) parser/parser.go

$(TARGET): $(SRCS)
	go generate ./...
	go build $(GOFLAGS)

$(RUNTIME_TARGET): $(RUNTIME_SRCS)
	cd runtime && cargo build --release && cp ./target/release/$(RUNTIME_TARGET) ./$(RUNTIME_TARGET)

.PHONY: runtime_clean
runtime_clean:
	cd runtime && cargo clean && $(RM) ./$(RUNTIME_TARGET)
