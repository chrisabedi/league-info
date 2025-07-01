

BINARY_NAME := league-info

GO_FILES := main.go 

all: run

build: 
	go build -o $(BINARY_NAME) $(GO_FILES)


run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)