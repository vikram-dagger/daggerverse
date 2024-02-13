package main

import (
	"context"
	"fmt"
	"strings"
)

type DaggerDocs struct {}

// example usage:
func (m *DaggerDocs) Deploy(project string, location string, repository string, credential *Secret) (string, error) {

	ctx := context.Background()

	tree := dag.Git("https://github.com/dagger/dagger").
		Branch("main").
		Tree()

	build := dag.Container().
		From("node:21").
		WithDirectory("/home/node", tree).
		WithWorkdir("/home/node/docs").
		WithMountedCache("/src/node_modules", dag.CacheVolume("node-21-modules")).
		WithExec([]string{"npm", "install"}).
		WithExec([]string{"npm", "run", "build"}).
		Directory("./build")

	registry := fmt.Sprintf("%s-docker.pkg.dev/%s/%s/dagger-docs", location, project, repository)
	split := strings.Split(registry, "/")
	address, err := dag.Container().From("nginx:1.25-alpine").
		WithDirectory("/usr/share/nginx/html", build).
		WithExposedPort(80).
		WithRegistryAuth(split[0], "_json_key", credential).
		Publish(ctx, registry)
	if (err != nil) {
		panic(err)
	}

	return dag.GoogleCloudRun().Deploy(ctx, project, location, address, 80, credential)
}
