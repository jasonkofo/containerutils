package containerutils

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

// ComposeTemplateContent is a struct that provides necessary content that is
// used for creating the docker-compose templates. This struct is fed into Go's
// templating engine and then used to create the returned compose file
type ComposeTemplateContent struct {
	VolumeDirectory string
}

type composeMode int

const (
	developerCompose    composeMode = 0
	orchestratorCompose composeMode = 1 << (iota - 1)
	productionCompose
)

// ComposeServiceConfig represents JSON data that is posted from clients of this
// service. Clients are able to specify the services that they would like
// in their compose files, should they post data correctly.
type ComposeServiceConfig struct {
	WhiteList *[]string `json:"whiteList"`
}

// ComposeFile is a local abstraction of the docker-compose template file
type ComposeFile struct {
	Version  string                            `yaml:"version"`
	Services map[string]*Service               `yaml:"services"`
	Secrets  map[string]map[string]interface{} `yaml:"secrets,omitempty"`
}

// Write marshals the ComposeFile object into the given filename
func (cf *ComposeFile) Write(filename string, options map[string]interface{}) error {
	rp, okRP := options["routerport"]
	ep, okEP := options["suppressports"]
	dbp, okDBP := options["dbport"]
	https, okHTTP := options["https"]
	suppressPorts := false
	var (
		ok       bool
		isHTTPS  = false
		dbPort   int
		dbPort32 int32
		dbPort64 int64
	)
	if okEP {
		if suppressPorts, ok = ep.(bool); !ok {
			suppressPorts = false
		}
	}
	if okHTTP {
		isHTTPS = https.(bool)
	}
	if okDBP {
		dbPort, ok = dbp.(int)
		if !ok {
			dbPort32, ok = dbp.(int32)
		}
		if ok {
			dbPort = int(dbPort32)
		} else {
			dbPort64, ok = dbp.(int64)
		}
		if ok {
			dbPort = int(dbPort64)
		} else {
			dbPort = 5432
		}
	}
	if okRP {
		if svc, okr := cf.Services["router"]; okr {
			// [JKG 2021-04-13] Inside the docker container, we map the router
			// to port 80, so any custom configuration will have to map the
			// external binding on the host to the internal docker port 80
			attr := make(Attributes, 0, len(svc.DockerComposePort))
			for _, port := range svc.DockerComposePort {
				if strings.HasSuffix(string(port), ":80") {
					attr.AndString(
						fmt.Sprintf("127.0.0.1:%v:80", rp),
					)
				} else if strings.HasSuffix(string(port), ":443") {
					if isHTTPS {
						attr.AndString(
							fmt.Sprintf("443:443"),
						)
					}
				} else {
					attr.AndString(string(port))
				}
			}
			svc.DockerComposePort = attr
		}
	}

	if suppressPorts {
		for svcName, svc := range cf.Services {
			if svcName != "router" {
				svc.DockerComposePort = nil
			}
		}
	}

	if dbPort != 5432 {
		dbSvc, ok := cf.Services["db"]
		if ok {
			dbSvc.DockerComposePort = Attributes{
				Attribute(fmt.Sprintf("127.0.0.1:%v:5432", dbPort)),
			}
		}
	}

	_cf, err := yaml.Marshal(cf)
	if err != nil {
		return fmt.Errorf("Could not write '%v': %v", filename, err.Error())
	}
	return ioutil.WriteFile(filename, _cf, 0644)
}

// ReadFromFile attempts to read the contents of a given file into memory.
// This method does NOT support the case where a user wants to merge to 2 sets
// and instead assumes that they are wanting a frresh list every time
func (cf *ComposeFile) ReadFromFile(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Could not remove services from Docker Compose")
	}

	if err = yaml.Unmarshal(file, cf); err != nil {
		return fmt.Errorf("Could not read '%v' into a native yml object: %v", filename, err)
	}
	return nil
}

// DeleteBlacklisted deletes services from the ComposeFile object if they
// are in the given blacklist
func (cf *ComposeFile) DeleteBlacklisted(blacklist []string) {
	_bl := make(map[string]int)
	for _, svc := range blacklist {
		_bl[svc] = 1
	}
	for _svcKey := range cf.Services {
		if _, ok := _bl[_svcKey]; ok {
			delete(cf.Services, _svcKey)
		}
	}

	cf.EnsureDependencies()
}

// EnsureDependencies removes dependencies that are not present in the
// ComposeFile
func (cf *ComposeFile) EnsureDependencies() {
	for _, content := range cf.Services {
		for _, dependency := range content.DependsOn {
			if !cf.HasService(dependency.String()) && content.DependsOn.Contains(dependency.String()) {
				content.DependsOn.Remove(dependency.String())
			}
		}
	}
}

// HasService confirms whether or not the service object is found in the
// ComposeFile. This is more of a debug helper
func (cf *ComposeFile) HasService(service string) bool {
	_, ok := cf.Services[service]
	return ok
}

// Print puts the object out to stdout
func (cf *ComposeFile) Print() {
	_cf, _ := yaml.Marshal(cf)
	fmt.Println(string(_cf))
}
