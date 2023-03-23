MAKEFLAGS := --silent --always-make
TEST := test $(if $(filter $(verb),true),-v,) -count=1 -short -run=$(run)
BENCH := test -count=1 -short -run=- -bench=$(or $(run),.) -benchmem -benchtime=128ms -timeout=1s
GOW := gow -c -v -i=testdata

.PHONY: test_w
test_w:
	$(GOW) $(TEST)

.PHONY: test
test:
	go $(TEST)

.PHONY: bench_w
bench_w:
	$(GOW) $(BENCH)

.PHONY: bench
bench:
	go $(BENCH)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: run_w
run_w:
	$(GOW) run . $(run)

.PHONY: run
run:
	go run . $(run)

.PHONY: norm
norm:
	$(MAKE) run run=norm

.PHONY: discord_download
discord_download:
	$(MAKE) run run=discord_download

# Usage:
#   make log search="<word>"
.PHONY: log
log:
	git log --format=reference --date=iso-strict -S "$(search)"
