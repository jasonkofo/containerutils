package containerutils

import (
	"errors"
)

// ResolveServiceDependencies recursively traces image dependencies
// and ensures that the required image has been added
func ResolveServiceDependencies(out *ComposeFile, template *ComposeFile) (bool, error) {
	serviceAdded := false
	if template == nil {
		return false, errors.New("No services found in the template input")
	}
	if out == nil {
		return false, errors.New("No services to resolve")
	}

	getServiceFromCompose := func(file ComposeFile, service string) (*Service, error) {
		s, ok := file.Services[service]
		if !ok {
			return nil, errors.New(service + " not found in compose template")
		}
		return s, nil
	}

	for k := range out.Services {
		dependsOn := template.Services[k].DependsOn
		if dependsOn == nil {
			continue
		}

		for _, dependentService := range dependsOn {
			if _, ok := out.Services[dependentService.String()]; ok {
				continue
			}
			service, err := getServiceFromCompose(*template, dependentService.String())
			if err != nil || service == nil {
				return false, err
			}
			out.Services[dependentService.String()] = service
			serviceAdded = true
			for serviceAdded == true {
				serviceAdded, err = ResolveServiceDependencies(out, template)
				if err != nil {
					return false, err
				}
			}
		}
	}
	return serviceAdded, nil
}
