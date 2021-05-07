package containerutils

var _file = RegistryFile{}

// ConstructServiceRegistry returns the a pretty-printed byte representation
// of the service-registry.json
func ConstructServiceRegistry() ([]byte, error) {
	return _file.filterExcludedServices().jsonBytes()
}

// ConstructDeveloperCompose returns all the services that
// are containerized
func ConstructDeveloperCompose() ComposeFile {
	return _file.Services.ToDockerCompose("", developerCompose, 0)
}

// ConstructOrchestratorCompose returns a compose file for running the test orchestrator
// This means that the only exposed port is Router
func ConstructOrchestratorCompose(routerPort int) ComposeFile {
	return _file.Services.ToDockerCompose("", orchestratorCompose, 0)
}

// ConstructProductionCompose returns a compose file for running the test orchestrator
// This means that the only exposed port is Router
func ConstructProductionCompose(tagName string, routerPort int) ComposeFile {
	return _file.Services.ToDockerCompose(tagName, productionCompose, routerPort)
}
