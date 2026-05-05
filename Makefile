BINARY_NAME=ops-pilot
PREFIX=/usr/local/bin

.PHONY: build install uninstall clean test

build:
	go build -o $(BINARY_NAME) ./cmd/ops-pilot

install: build
	@echo "Installing $(BINARY_NAME) to $(PREFIX)..."
	sudo cp $(BINARY_NAME) $(PREFIX)/$(BINARY_NAME)
	sudo chmod +x $(PREFIX)/$(BINARY_NAME)
	@echo "Done! You can now run '$(BINARY_NAME)' from anywhere."

uninstall:
	@echo "Removing $(BINARY_NAME) from $(PREFIX)..."
	sudo rm $(PREFIX)/$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

test:
	go test ./...
