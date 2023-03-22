MAKEFLAGS := --silent --always-make
TEST := test $(if $(filter $(verb),true),-v,) -count=1 -short -run=$(run)
BENCH := test -count=1 -short -bench=$(or $(run),.) -benchmem -timeout=1s

.PHONY: test_w
test_w:
	gow -c -v $(TEST)

.PHONY: test
test:
	go $(TEST)

.PHONY: bench_w
bench_w:
	gow -c -v $(BENCH)

.PHONY: bench
bench:
	go $(BENCH)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: run_w
run_w:
	gow -c -v run . $(run)

.PHONY: run
run:
	go run . $(run)

.PHONY: norm
norm:
	$(MAKE) run run=norm

# Usage:
#   make log search="<word>"
.PHONY: log
log:
	git log --format=reference --date=iso-strict -S "$(search)"
