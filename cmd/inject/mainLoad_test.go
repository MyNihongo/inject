package main

import (
	"context"
	"testing"
)

func TestLoad(t *testing.T) {
	ctx := context.Background()
	const wd = `C:\tfs_mku\__github\MyNihongo.NuGet\go-inject\examples\`

	loadFileContent(ctx, wd)
}
