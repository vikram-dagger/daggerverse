package main

type Trivy struct{}

func (m *Trivy) Check(opts VulnOpts) (*Container, error) {
	c := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "curl", "git"}).
		WithExec([]string{"sh", "-c", "curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin v0.45.1"}).
		WithWorkdir("/tmp").
		WithExec([]string{"git", "clone", opts.Repository, opts.Path}).
		WithExec([]string{"trivy", "fs", "-f", "json", opts.Path})
	return c, nil
}

/*
func (m *Trivy) Trivy(ctr *Container) (*Container, error) {
	return ctr, nil
}

func (ctr *Container) Trivy(path string) (*Container, error) {
	c := ctr.
		WithWorkdir("/tmp").
		WithExec([]string{"sh", "-c", "curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin v0.45.1"}).
		WithExec([]string{"trivy", "fs", "-f", "json", path})
	return c, nil
}
*/
