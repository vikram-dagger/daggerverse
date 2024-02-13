package main

import (
	"context"
	"fmt"
	"strings"
)

type DaggerDocs struct {}

// example usage
// dagger -m github.com/vikram-dagger/daggerverse/dagger-docs call deploy --project vikram-experiments --location us-central1 --repository vikram-test --credential env:GOOGLE_CREDENTIAL
func (m *DaggerDocs) Deploy(source *Directory, project string, location string, repository string, credential *Secret) (string, error) {

	ctx := context.Background()

	build := dag.Container().
		From("node:21").
		WithDirectory("/home/node", source).
		WithWorkdir("/home/node").
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
