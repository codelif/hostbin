APP := hostbin
CLI_APP := hbcli
BIN_DIR := bin
BUILD_BINARY := $(BIN_DIR)/$(APP)
BUILD_CLI_BINARY := $(BIN_DIR)/$(CLI_APP)
GO ?= go

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
SYSTEMD_DIR ?= /etc/systemd/system
ETC_DIR ?= /etc/$(APP)
STATE_DIR ?= /var/lib/$(APP)

UNIT_SRC := deploy/systemd/$(APP).service
ENV_SRC := deploy/systemd/$(APP).env.example

INSTALLED_BINARY := $(DESTDIR)$(BINDIR)/$(APP)
INSTALLED_CLI_BINARY := $(DESTDIR)$(BINDIR)/$(CLI_APP)
INSTALLED_UNIT := $(DESTDIR)$(SYSTEMD_DIR)/$(APP).service
INSTALLED_ENV := $(DESTDIR)$(ETC_DIR)/$(APP).env
INSTALLED_STATE := $(DESTDIR)$(STATE_DIR)

.PHONY: help fmt tidy test test-race build build-server build-cli run install install-server install-cli install-config install-state install-systemd install-all uninstall clean

help:
	@printf '%s\n' \
		"Targets:" \
		"  make fmt             Format Go code" \
		"  make tidy            Tidy Go module dependencies" \
		"  make test            Run the Go test suite" \
		"  make test-race       Run the Go test suite with the race detector" \
		"  make build           Build ./bin/$(APP) and ./bin/$(CLI_APP)" \
		"  make build-server    Build ./bin/$(APP)" \
		"  make build-cli       Build ./bin/$(CLI_APP)" \
		"  make run             Run the server locally" \
		"  make install         Install the server binary to $(BINDIR)" \
		"  make install-cli     Install the hbcli binary to $(BINDIR)" \
		"  make install-config  Install $(APP).env if it does not already exist" \
		"  make install-state   Create $(STATE_DIR)" \
		"  make install-systemd Install the systemd unit to $(SYSTEMD_DIR)" \
		"  make install-all     Install server, hbcli, env example, state dir, and unit" \
		"  make uninstall       Remove installed binaries and unit from DESTDIR/PREFIX" \
		"  make clean           Remove built binaries" \
		"" \
		"Install overrides:" \
		"  PREFIX=/usr/local    Binary prefix (default)" \
		"  DESTDIR=/tmp/pkg     Package staging root"

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

build: build-server build-cli

build-server:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BUILD_BINARY) ./cmd/server

build-cli:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BUILD_CLI_BINARY) ./cmd/hbcli

run:
	$(GO) run ./cmd/server

install: install-server

install-server: build-server
	install -d "$(DESTDIR)$(BINDIR)"
	install -m 0755 "$(BUILD_BINARY)" "$(INSTALLED_BINARY)"

install-cli: build-cli
	install -d "$(DESTDIR)$(BINDIR)"
	install -m 0755 "$(BUILD_CLI_BINARY)" "$(INSTALLED_CLI_BINARY)"

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

install-all: install-server install-cli install-config install-state install-systemd

uninstall:
	rm -f "$(INSTALLED_BINARY)" "$(INSTALLED_CLI_BINARY)" "$(INSTALLED_UNIT)"

clean:
	rm -rf $(BIN_DIR)
