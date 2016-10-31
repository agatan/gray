TARGET:=gray

.PHONY: all
all:
	go generate ./...
	go build

.PHONY: run
run: all
	./$(TARGET)

.PHONY: test
test: all
	go test ./...

.PHONY: clean
clean:
	$(RM) $(TARGET) parser/parser.go
