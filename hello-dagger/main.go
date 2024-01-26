package main

import "context"

type HelloDagger struct{}

func (m *HelloDagger) HelloDagger(ctx context.Context) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "figlet"}).
		WithExec([]string{"figlet", "dagger"}).
		Stdout(ctx)
}
