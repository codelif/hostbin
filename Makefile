APP := hostbin
SERVER_DIR := server
BIN_DIR := bin
BUILD_BINARY := $(BIN_DIR)/$(APP)
GO ?= go

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
SYSTEMD_DIR ?= /etc/systemd/system
ETC_DIR ?= /etc/$(APP)
STATE_DIR ?= /var/lib/$(APP)

UNIT_SRC := deploy/systemd/$(APP).service
ENV_SRC := deploy/systemd/$(APP).env.example

INSTALLED_BINARY := $(DESTDIR)$(BINDIR)/$(APP)
INSTALLED_UNIT := $(DESTDIR)$(SYSTEMD_DIR)/$(APP).service
INSTALLED_ENV := $(DESTDIR)$(ETC_DIR)/$(APP).env
INSTALLED_STATE := $(DESTDIR)$(STATE_DIR)

.PHONY: help fmt tidy test test-race build run install install-config install-state install-systemd install-all uninstall clean

help:
	@printf '%s\n' \
		"Targets:" \
		"  make fmt             Format Go code" \
		"  make tidy            Tidy Go module dependencies" \
		"  make test            Run the Go test suite" \
		"  make test-race       Run the Go test suite with the race detector" \
		"  make build           Build ./bin/$(APP)" \
		"  make run             Run the server locally" \
		"  make install         Install the binary to $(BINDIR)" \
		"  make install-config  Install $(APP).env if it does not already exist" \
		"  make install-state   Create $(STATE_DIR)" \
		"  make install-systemd Install the systemd unit to $(SYSTEMD_DIR)" \
		"  make install-all     Install binary, env example, state dir, and unit" \
		"  make uninstall       Remove installed binary and unit from DESTDIR/PREFIX" \
		"  make clean           Remove built binaries" \
		"" \
		"Install overrides:" \
		"  PREFIX=/usr/local    Binary prefix (default)" \
		"  DESTDIR=/tmp/pkg     Package staging root"

fmt:
	$(GO) -C $(SERVER_DIR) fmt ./...

tidy:
	$(GO) -C $(SERVER_DIR) mod tidy

test:
	$(GO) -C $(SERVER_DIR) test ./...

test-race:
	$(GO) -C $(SERVER_DIR) test -race ./...

build:
	mkdir -p $(BIN_DIR)
	$(GO) -C $(SERVER_DIR) build -o ../$(BUILD_BINARY) ./cmd/server

run:
	$(GO) -C $(SERVER_DIR) run ./cmd/server

install: build
	install -d "$(DESTDIR)$(BINDIR)"
	install -m 0755 "$(BUILD_BINARY)" "$(INSTALLED_BINARY)"

install-config:
	install -d "$(DESTDIR)$(ETC_DIR)"
	@if [ -f "$(INSTALLED_ENV)" ]; then \
		printf '%s\n' "Keeping existing $(INSTALLED_ENV)"; \
	else \
		install -m 0640 "$(ENV_SRC)" "$(INSTALLED_ENV)"; \
	fi

install-state:
	install -d -m 0750 "$(INSTALLED_STATE)"

install-systemd:
	install -d "$(DESTDIR)$(SYSTEMD_DIR)"
	install -m 0644 "$(UNIT_SRC)" "$(INSTALLED_UNIT)"
	@printf '%s\n' "Installed $(INSTALLED_UNIT)"
	@printf '%s\n' "Run: sudo systemctl daemon-reload && sudo systemctl enable --now $(APP)"

install-all: install install-config install-state install-systemd

uninstall:
	rm -f "$(INSTALLED_BINARY)" "$(INSTALLED_UNIT)"

clean:
	rm -rf $(BIN_DIR)
