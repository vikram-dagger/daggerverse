package main

type Vuln struct{}

type VulnOpts struct {
	Token      string
	Repository string
	Path       string
}

func (m *Vuln) Check(scanner string, opts VulnOpts) (*Container, error) {
	if scanner == "trivy" {
		return dag.Trivy().Check(opts), nil
	}
	if scanner == "snyk" {
		return dag.Snyk().Checks(opts), nil
	}
	return dag.Container(), nil
}
