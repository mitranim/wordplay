MAKEFLAGS := --silent --always-make
TEST      := test $(if $(filter $(verb), true), -v,) -count=1 -short -run=$(run)
BENCH     := test -count=1 -short -bench=$(or $(run),.) -benchmem

test_w:
	gow -c -v $(TEST)

test:
	go $(TEST)

bench_w:
	gow -c -v $(BENCH)

bench:
	go $(BENCH)

lint:
	golangci-lint run

run_w:
	gow -c -v run . $(run)

run:
	go run . $(run)

norm:
	$(MAKE) run run=norm

# Usage:
#   make log search="<word>"
log:
	git log --format=reference --date=iso-strict -S "$(search)"
