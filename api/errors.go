package api

import "fmt"

type missingAppYamlError struct {
	path string
}

func (e missingAppYamlError) Error() string {
	return fmt.Sprintf("could not find %s at %s", appYamlFilename, e.path)
}

type projectNetworkExists struct {
	id string
}

func (e projectNetworkExists) Error() string {
	return fmt.Sprintf("project %s already has a network", e.id)
}
