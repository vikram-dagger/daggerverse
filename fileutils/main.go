// A module with various filesystem utility functions

package main

import (
	"context"
)

type Fileutils struct{}

// Returns a tree representation of the directory provided as a string
func (m *Fileutils) Tree(ctx context.Context, dir *Directory) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", dir).
		WithWorkdir("/mnt").
		WithExec([]string{"tree"}).
		Stdout(ctx)
}
