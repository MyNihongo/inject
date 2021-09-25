package main

import (
	"context"
	"testing"
)

func TestAll(t *testing.T) {
	const wd = ``
	ctx := context.Background()

	execute(ctx, wd)
}
