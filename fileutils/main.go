package main

import (
	"context"
)

type Fileutils struct{}

// example usage: "dagger call grep-dir --directory . --pattern xyz"
func (m *Fileutils) GrepDir(ctx context.Context, directory *Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}

// example usage: "dagger call tree --directory ."
func (m *Fileutils) Tree(ctx context.Context, directory *Directory) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"tree"}).
		Stdout(ctx)
}
