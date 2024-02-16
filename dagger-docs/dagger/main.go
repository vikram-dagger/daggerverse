package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type DaggerDocs struct{}

// example usage
// dagger -m github.com/vikram-dagger/daggerverse/dagger-docs call deploy --source ./dagger --project user-experiments --location us-central1 --repository user-test --credential env:GOOGLE_CREDENTIAL
func (m *DaggerDocs) Deploy(source *Directory, project string, location string, repository string, credential *Secret) (string, error) {
	ctx := context.Background()

	build := m.Build(source)

	registry := fmt.Sprintf("%s-docker.pkg.dev/%s/%s/dagger-docs", location, project, repository)
	split := strings.Split(registry, "/")
	address, err := dag.Container().From("nginx:1.25-alpine").
		WithDirectory("/usr/share/nginx/html", build).
		WithExposedPort(80).
		WithRegistryAuth(split[0], "_json_key", credential).
		Publish(ctx, registry)
	if err != nil {
		panic(err)
	}

	return dag.GoogleCloudRun().CreateService(ctx, project, location, address, 80, credential)
}

// example usage
// dagger -m github.com/vikram-dagger/daggerverse/dagger-docs call build --source ./dagger
func (m *DaggerDocs) Build(source *Directory) *Directory {
	return dag.Container().
		From("node:21").
		WithDirectory("/home/node", source).
		WithWorkdir("/home/node/docs").
		WithMountedCache("/home/node/docs/node_modules", dag.CacheVolume("node-21-modules")).
		WithMountedCache("/home/node/.npm", dag.CacheVolume("npm-cache")).
		WithExec([]string{"npm", "install"}).
		WithEnvVariable("CACHE", time.Now().String()).
		WithExec([]string{"npm", "run", "build"}).
		Directory("./build")
}
