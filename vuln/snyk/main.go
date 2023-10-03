package main

type Snyk struct{}

func (m *Snyk) Checks(opts VulnOpts) (*Container, error) {
	c := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "curl", "git"}).
		WithExec([]string{"sh", "-c", "curl --compressed https://static.snyk.io/cli/latest/snyk-alpine -o snyk && chmod +x ./snyk && mv ./snyk /usr/local/bin/"}).
		WithWorkdir("/tmp").
		WithExec([]string{"git", "clone", opts.Repository, opts.Path}).
		WithEnvVariable("SNYK_TOKEN", opts.Token).
		WithExec([]string{"snyk", "test", "--json", opts.Path})
	return c, nil
}
