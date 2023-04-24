MAKEFLAGS := --silent --always-make
TEST := test $(if $(filter $(verb),true),-v,) -count=1 -short -run=$(run)
BENCH := test -count=1 -short -run=- -bench=$(or $(run),.) -benchmem -benchtime=128ms -timeout=1s
GOW := gow -c -v -i=testdata

test_w:
	$(GOW) $(TEST)

test:
	go $(TEST)

bench_w:
	$(GOW) $(BENCH)

bench:
	go $(BENCH)

lint:
	golangci-lint run

run_w:
	$(GOW) run . $(run)

run:
	go run . $(run)

norm:
	$(MAKE) run run=norm

discord_download:
	$(MAKE) run run=discord_download

# Usage:
#   make log search="<word>"
log:
	git log --format=reference --date=iso-strict -S "$(search)"
