MAKEFLAGS := --silent --always-make
TEST      := test -v -count=1 -short
BENCH     := test -count=1 -short -bench=. -benchmem

test-w:
	gow -c -v $(TEST)

test:
	go $(TEST)

bench-w:
	gow -c -v $(BENCH)

bench:
	go $(BENCH)

log:
	git log --format=reference --date=iso-strict -S "$(search)"
