TARGET:=gray

.PHONY: all
all: $(TARGET)

.PHONY: test
test: all
	go test ./...

.PHONY: clean
clean:
	$(RM) $(TARGET) parser/parser.go

$(TARGET): parser/parser.go
	go build

parser/parser.go: parser/parser.go.y
	go generate ./...
