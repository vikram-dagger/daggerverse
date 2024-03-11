// A module to deploy the Dagger documentation as a Google Cloud Run service
//
// This module contains functions to build a static site containing the Dagger
// documentation, and then deploy the static site to Google Cloud Run
//

package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type DaggerDocs struct{}

// Deploys a source Directory containing a static website to an
// existing Google Cloud Run service
//
// example:
// dagger call deploy --source ./docs --project user-project --location us-central1 --repository user-repo --credential env:GOOGLE_CREDENTIAL
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

// Builds a static docs website from a source Directory
// and returns a Directory containing the build result
//
// example:
// dagger call build --source ./docs
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
