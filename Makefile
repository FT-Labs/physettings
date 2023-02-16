PROG = physettings

ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

rwildcard=$(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))

SRC = $(call rwildcard,.,*.go)

PREFIX ?= /usr

GOCMD = go

VERSION ?= 1.0.0

build: $(SRC)
	GOOS=linux GOARCH=amd64 ${GOCMD} build -o ${PROG} main.go

deps: $(SRC)
	${GOCMD} get github.com/FT-Labs/tview@9d459cd
	${GOCMD} get github.com/gdamore/tcell/v2@latest

clean:
	rm -f ${PROG}

all: deps build

install:
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f ${PROG} $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/${PROG}

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/${PROG}

.PHONY: all build deps clean install uninstall
