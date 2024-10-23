### THANK U KUBE-VIP

include .env

SHELL := /bin/sh
OUT_DIR = out

TARGET := scrapper
CSV := malstat.csv
DB := ${DATABASE}

# These will be provided to the target
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Build=$(BUILD)"

.PHONY: all build clean install uninstall check run deploy

all: check install

$(TARGET):
	@go build $(LDFLAGS) -o $(OUT_DIR)/$(TARGET)

build: $(TARGET)
	@true

clean:
	rm -f $(OUT_DIR)/$(TARGET)
	rm -f $(OUT_DIR)/$(CSV)

install:
	@go install $(LDFLAGS)

uninstall: clean
	rm -f $$(which ${TARGET})

check:
	go mod tidy

run: install
	@$(TARGET) scrap --top 100 --csv $(OUT_DIR)/$(CSV) --db $(DB)

deploy: build
	ansible-playbook deployments/ansible/deploy.yml -vv 
	