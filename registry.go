package containerutils

import (
	"encoding/json"
	"fmt"

	mapper "github.com/PeteProgrammer/go-automapper"
)

// RegistryFile is an abstraction of the service-registry.json file
type RegistryFile struct {
	Services          Services `json:"services,omitempty"`
	DeposedServices   []string `json:"deposedServices,omitempty"`
	UnmanagedServices []string `json:"unmanagedServices,omitempty"`
}

func (rf *RegistryFile) clone() RegistryFile {
	_rf := RegistryFile{}
	copy(_rf.Services, rf.Services)
	copy(_rf.DeposedServices, rf.DeposedServices)
	copy(_rf.UnmanagedServices, rf.UnmanagedServices)
	return *rf
}

func (rf *RegistryFile) filterExcludedServices() *RegistryFile {
	_rf := rf.clone()
	_services := Services{}
	for _, service := range rf.Services {
		if !service.ExcludeFromServiceRegistry {
			_services.Append(service)
		}
	}
	_rf.Services = _services

	return &_rf
}

// Services is a slice of Services - has methods on it
type Services []Service

// ToSliceKind returns the underlying slice type
func (s *Services) ToSliceKind() []Service {
	if s == nil {
		return []Service{}
	}
	return []Service(*s)
}

// Append appends Service s1 to the services object
func (s *Services) Append(s1 Service) {
	_s := Services(append(s.ToSliceKind(), s1))
	s = &_s
}

// Service is an abstraction a service in the service-registry.json file
// The "Container" native go field is used as "container" in service-registry.json, whereas
// it is used as "name" in the docker-compose.yml
// This struct is partially based on the compose spec found here:
// https://github.com/compose-spec/compose-spec/blob/master/spec.md#services-top-level-element
type Service struct {
	Name                       string                      `json:"name,omitempty" yaml:"-"`
	Image                      string                      `json:"-" yaml:"image"`
	Container                  string                      `json:"container,omitempty" yaml:"-"`
	RepoName                   string                      `json:"-" yaml:"-"`
	Port                       []int                       `json:"port,omitempty" yaml:"-"`
	DockerComposePort          Attributes                  `json:"-" yaml:"ports,omitempty"`
	Port80InDocker             bool                        `json:"port80InDocker,omitempty" yaml:"-"`
	PingCustom                 bool                        `json:"pingCustom,omitempty" yaml:"-"`
	CustomCreate               string                      `json:"customCreate,omitempty" yaml:"-"`
	CustomDelete               string                      `json:"customDelete,omitempty" yaml:"-"`
	URL                        string                      `json:"url,omitempty" yaml:"-"`
	HasPing                    bool                        `json:"hasPing,omitempty" yaml:"-"`
	InstallType                string                      `json:"installType,omitempty" yaml:"-"`
	BinPath                    string                      `json:"binPath,omitempty" yaml:"-"`
	Command                    string                      `json:"-" yaml:"command,omitempty"`
	CommandKeyPhrase           string                      `json:"commandKeyPhrase,omitempty" yaml:"-"`
	PGConnectionManager        *ServicePGConnectionManager `json:"pgConnectionManager,omitempty" yaml:"-"`
	Dependencies               []string                    `json:"dependencies,omitempty" yaml:"-"`
	DependsOn                  Attributes                  `json:"-" yaml:"depends_on,omitempty"`
	Logs                       []ServiceLogs               `json:"logs,omitempty" yaml:"-"`
	IsExclusivelyLinux         bool                        `json:"isExclusivelyLinux,omitempty" yaml:"-"`
	DefaultTag                 string                      `json:"defaultTag,omitempty" yaml:"-"`
	Volumes                    []Attribute                 `json:"-" yaml:"volumes,omitempty"`
	Environment                map[string]interface{}      `json:"-" yaml:"environment,omitempty"`
	Restart                    string                      `json:"-" yaml:"restart,omitempty"`
	ExcludeFromServiceRegistry bool                        `json:"-" yaml:"-"`
	IsExternalImage            bool                        `json:"-" yaml:"-"`
	Deploy                     map[string]interface{}      `json:"-" yaml:"deploy,omitempty"` // This might need to be fleshed out
	Build                      string                      `json:"-" yaml:"build,omitempty"`  // Build can be fleshed out more but it might be a good idea to just force it to call a script
	CapAdd                     Attributes                  `json:"-" yaml:"cap_add,omitempty"`
	CapDrop                    Attributes                  `json:"-" yaml:"cap_drop,omitempty"`
	CGroupParent               string                      `json:"-" yaml:"cgroup_parent,omitempty"`
	Configs                    Attributes                  `json:"-" yaml:"configs,omitempty"`
	ContainerName              string                      `json:"-" yaml:"container_name,omitempty"`
	DNS                        Attributes                  `json:"-" yaml:"dns,omitempty"`
	DNSOpt                     Attributes                  `json:"-" yaml:"dns_opt,omitempty"`
	DNSSearch                  Attributes                  `json:"-" yaml:"dns_search,omitempty"`
	Domainname                 string                      `json:"-" yaml:"domainname,omitempty"`
	Entrypoint                 string                      `json:"-" yaml:"entrypoint,omitempty"`
	EnvFile                    Attributes                  `json:"-" yaml:"env_file,omitempty"`
	Expose                     Attributes                  `json:"-" yaml:"expose,omitempty"`
	Extends                    map[string]interface{}      `json:"-" yaml:"extends,omitempty"`
	ExternalLinks              Attributes                  `json:"-" yaml:"external_links,omitempty"`
	ExtraHosts                 Attributes                  `json:"-" yaml:"extra_hosts,omitempty"`
	GroupAdd                   Attributes                  `json:"-" yaml:"group_add,omitempty"`
	Healthcheck                map[string]interface{}      `json:"-" yaml:"healthcheck,omitempty"`
	Hostname                   string                      `json:"-" yaml:"hostname,omitempty"`
	Init                       bool                        `json:"-" yaml:"init,omitempty"`
	Links                      Attributes                  `json:"-" yaml:"links,omitempty"`
	Logging                    map[string]interface{}      `json:"-" yaml:"logging,omitempty"`
	NetworkMode                string                      `json:"-" yaml:"network_mode,omitempty"`
	Networks                   Attributes                  `json:"-" yaml:"networks,omitempty"`
	Profiles                   Attributes                  `json:"-" yaml:"profiles,omitempty"`
	PullPolicy                 string                      `json:"-" yaml:"pull_policy,omitempty"`
	Secrets                    Attributes                  `json:"-" yaml:"secrets,omitempty"`
	ShmSize                    string                      `json:"-" yaml:"shm_size,omitempty"`
	StopGracePeriod            string                      `json:"-" yaml:"stop_grace_period,omitempty"`
	Sysctls                    Attributes                  `json:"-" yaml:"sysctls,omitempty"`
	ULimits                    map[string]interface{}      `json:"-" yaml:"ulimits,omitempty"`
}

func (s *Service) clone() {

}

func (s *Service) transformPort(routerPort int) {
	if s.DockerComposePort != nil {
		return
	}
	s.DockerComposePort = s.GetDockerComposePort(routerPort)
}

// toDockerCompose transforms the Service object into something that
// would be suited for a docker-compose.yml
func (s *Service) toDockerCompose(mode composeMode, routerPort int) *Service {
	_s := &Service{}
	mapper.Map(s, _s)
	// We expose the router ports for all of our test orchestrator images, as
	// well as potentially doing so in the case of production machines
	if mode == developerCompose || mode == orchestratorCompose && s.Container == "router" {
		s.transformPort(routerPort)
	}
	return _s
}

// GetDockerComposePort returns the "port" entry for the docker-compose
// service entry
func (s *Service) GetDockerComposePort(routerPort int) Attributes {
	_p := s.Port
	if s.Container == "router" && routerPort != 0 {
		_p = []int{routerPort}
	}
	attr := Attributes{}
	for _, port := range _p {
		_str := fmt.Sprintf("127.0.0.1:%v:%v", port, s.pickInnerPort())
		attr = attr.AndString(_str)
	}
	return attr
}

// containerized returns the services that have the "container"
// field populated in the service-registry.json
func (s Services) containerized() Services {
	ss := Services{}
	for _, svc := range s {
		if svc.Container == "" {
			continue
		}
		ss = append(ss, svc)
	}
	return ss
}

// GetServicesAsYML returns a YML map of the services and their subsisting
// information
func (s *Services) GetServicesAsYML(tag string, mode composeMode, routerPort int) map[string]*Service {
	_s := make(map[string]*Service)
	for _, svc := range s.containerized() {
		_svc := svc.SetDockerComposeImage(tag).toDockerCompose(mode, routerPort)
		if svc.DependsOnDB() {
			_svc = _svc.waitForPostgres()
		}
		_s[svc.Container] = _svc
	}
	return _s
}

// ToDockerCompose performs a set of transformations on the services
// to turn them into a form that can be used for docker-compose.yml
func (s Services) ToDockerCompose(tag string, mode composeMode, routerPort int) ComposeFile {
	_s := s.GetServicesAsYML(tag, mode, routerPort)
	_c := ComposeFile{}
	_c.Version = "3.2"
	_c.Services = _s
	return _c
}

// DependsOnDB establishes whether a service depends on the db or dbpool docker image.
func (s Service) DependsOnDB() bool {
	return s.DependsOn.Contains("db") || s.DependsOn.Contains("dbpool")
}

func (s *Service) pickInnerPort() string {
	if s.Port80InDocker {
		return "80"
	}
	return fmt.Sprintf("%v", s.Port)
}

func (s *Service) waitForPostgres() *Service {
	s.Command = "wait-for-nc.sh config:80 -- wait-for-postgres.sh db /opt/" + pickBinaryName(s.Container)
	return s
}

func pickBinaryName(serviceName string) string {
	switch serviceName {
	case "job":
		return "imqs-jobservice"
	default:
		return serviceName
	}
}

// SetDockerComposeImage constructs the DockerCompose image of the given service
// If the image property is not specified (which is a convention we follow for
// services that are produced inhouse), it tries to construct the correct image
// tag given the "desiredTag" property. If an empty string is passed into the
// method, and the "DefaultTag" is unspecified, the tag "latest" will instead
// be passed on to the tag
func (s Service) SetDockerComposeImage(desiredTag string) *Service {
	if s.Image != "" || s.IsExternalImage {
		return &s
	}

	_t := desiredTag

	if _t == "" && s.DefaultTag == "" {
		_t = "latest"
	} else if _t == "" && s.DefaultTag != "" {
		_t = s.DefaultTag
	}

	s.Image = fmt.Sprintf("imqs/%v:%v", s.Container, _t)
	return &s
}

func (rf *RegistryFile) jsonBytes() ([]byte, error) {
	b, err := json.MarshalIndent(rf, " ", "\t")
	if err != nil {
		return nil, fmt.Errorf("Could not marshal JSON: %v", err)
	}
	return b, nil
}

// ServicePGConnectionManager is an abstracton of "pgConnectionManager" object that belongs to the service
// in the service-registry.json
type ServicePGConnectionManager struct {
	MaxPGConnectionSource     string `json:"maxPGConnectionSource"`
	MaxPGConnection           int    `json:"maxPGConnection,omitEmpty"`
	MaxPGConnectionTextFile   string `json:"maxPGConnectionTextFile,omitempty"`
	MaxPGConnectionKeyPhrase  string `json:"maxPGConnectionKeyPhrase,omitempty"`
	MaxPGConnectionMultiplier int    `json:"maxPGConnectionMultiplier,omitempty"`
}

// ServiceLogs is an abstraction of "logs" object that belongs to the service
// in the service-registry.json
type ServiceLogs struct {
	Name     string `json:"name,omitempty"`
	Filename string `json:"filename,omitempty"`
	Parser   string `json:"parser,omitempty"`
}
