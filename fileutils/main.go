package main

import (
	"context"
)

type Fileutils struct{}

// example usage: "dagger call grep-dir --dir . --pattern xyz"
func (m *Fileutils) GrepDir(ctx context.Context, dir *Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", dir).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}

// example usage: "dagger call tree --dir ."
func (m *Fileutils) Tree(ctx context.Context, dir *Directory) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", dir).
		WithWorkdir("/mnt").
		WithExec([]string{"tree"}).
		Stdout(ctx)
}
