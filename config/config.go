package config

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"time"
)

// Conf defines server options.
type Conf struct {
	API         API
	Storage     Storage     // 저장소
	Coordinator Coordinator // coordinator
	Secrets     Secrets     // vault 향후 이것의 쓰임을 살펴보자
	Mesos       Mesos
	Logging     Logging
}

// API defines API server options.
type API struct {
	Addr     string
	CertFile string
	KeyFile  string
	Auth     APIAuth
}

// API server authz backends.
const (
	APIAuthBackendGitLab = "gitlab"
	APIAuthBackendNone   = "none"
	APIAuthBackendLDAP   = "ldap"
)

// APIAuth defines API server authz options.
type APIAuth struct {
	Backend string
	GitLab  APIAuthGitLab
	LDAP    APIAuthLDAP
}

// APIAuthGitLab defines options of GitLab authz backend.
type APIAuthGitLab struct {
	Addr   string
	CACert string
}

// APIAuthLDAP defines options of LDAP authz backend.
type APIAuthLDAP struct {
	Addrs              []string
	UserDN             string
	UserAttr           string
	CACert             string
	UserACL            map[string]map[string]string
	GroupACL           map[string]map[string]string
	BindDN             string
	BindPassword       string
	Timeout            time.Duration
	GroupFilter        string
	GroupDN            string
	GroupAttr          string
	CaseSensitiveNames bool
}

// Storage defines server storage options.
type Storage struct {
	Backend   string
	ZooKeeper StorageZK
}

// StorageBackendZK denotes ZooKeeper storage backend.
const StorageBackendZK = "zookeeper"

// StorageZK defines ZooKeeper storage backend options.
type StorageZK struct {
	Dir     string
	Addrs   []string
	Timeout time.Duration // 살펴보자
	Auth    ZKAuth
	TaskTTL time.Duration // 살펴보자
}

// CoordinatorBackendZK denotes ZooKeeper coordinator backend.
const CoordinatorBackendZK = "zookeeper"

// Coordinator defines server coordinator options.
type Coordinator struct {
	Backend   string
	ZooKeeper CoordinatorZK
}

// CoordinatorZK defines ZooKeeper coordinator backend options.
type CoordinatorZK struct {
	Dir         string
	Addrs       []string
	Timeout     time.Duration
	Auth        ZKAuth
	ElectionDir string
}

// ZooKeeper authn methods.
const (
	ZKAuthSchemeDigest = "digest"
	ZKAuthSchemeWorld  = "world"
)

// ZKAuth defines ZooKeeper authn options.
type ZKAuth struct {
	Scheme string
	Digest ZKAuthDigest
}

// ZKAuthDigest defines options of ZooKeeper's digest authn scheme.
type ZKAuthDigest struct {
	User     string
	Password string
}

// Secrets backends.
const (
	SecretsBackendVault = "vault"
	SecretsBackendNone  = "none"
)

// Secrets defines server secrets options.
type Secrets struct {
	Backend string
	Vault   SecretsVault
}

// SecretsVault defines Vault secrets backend options.
type SecretsVault struct {
	Token   string
	Addr    string
	Timeout time.Duration
	Root    string
	CACert  string
}

// Mesos defines Mesos-related options.
// Mesos 에서 Roles 과 Principal
// http://mesos.apache.org/documentation/latest/authorization/
type Mesos struct {
	FrameworkName   string
	Auth            MesosAuth
	Addrs           []string
	CACert          string
	Checkpoint      bool
	FailoverTimeout time.Duration
	Hostname        string
	User            string
	WebUIURL        string
	Principal       string
	Labels          map[string]string // k,v http://mesos.apache.org/documentation/latest/app-framework-development-guide/
	Roles           []string
	LogAllEvents    bool
}

// Mesos authn schemes.
const (
	MesosAuthTypeBasic = "basic"
	MesosAuthTypeNone  = "none"
)

// MesosAuth defines Mesos authn options.
type MesosAuth struct {
	Type  string
	Basic MesosAuthBasic
}

// MesosAuthBasic defines options of Mesos' basic authn method.
type MesosAuthBasic struct {
	Username string
	Password string
}

// Logging backends.
const (
	LoggingBackendNone   = "none"
	LoggingBackendSentry = "sentry"
)

// Logging levels.
const (
	LoggingLevelDebug = "debug"
	LoggingLevelInfo  = "info"
	LoggingLevelWarn  = "warn"
	LoggingLevelError = "error"
)

// Logging defines server logging options.
type Logging struct {
	Level   string
	Backend string
	Sentry  LoggingSentry
}

// LoggingSentry defines Sentry logging backend options.
type LoggingSentry struct {
	DSN    string
	CACert string
	Tags   map[string]string
}

// 7/3 이 함수의 의미를 파악하자.
func millisecondFieldsToDuration(v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			millisecondFieldsToDuration(v.Field(i))
		} else {
			if v.Field(i).Type() == reflect.TypeOf(time.Second) {
				if v.Field(i).CanSet() {
					d := v.Field(i).Interface().(time.Duration)
					d *= time.Millisecond
					v.Field(i).Set(reflect.ValueOf(d))
				}
			}
		}
	}
}

// New creates new configuration struct.
// 7/3 일단 여기를 차후 정리를 해야함.
func New(args []string) *Conf {

	if args[0] != "config" {
		return nil
	}
	file, err := ioutil.ReadFile(args[1])
	if err != nil {
		return &Conf{
			Mesos: Mesos{
				FrameworkName:   "test",
				Addrs:           []string{"http://localhost:5050"},
				FailoverTimeout: 1000 * 3600 * 24 * 7, // 7d
				Roles:           []string{"*"},
				Auth: MesosAuth{
					Type: MesosAuthTypeNone,
				},
			},
		}

		//return nil, err
	}
	var conf = &Conf{
		API: API{
			Addr: "localhost:8000",
			Auth: APIAuth{
				Backend: APIAuthBackendNone,
				LDAP: APIAuthLDAP{
					Timeout:     5000,
					GroupFilter: "(|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))",
					GroupAttr:   "cn",
				},
			},
		},
		Storage: Storage{
			Backend: StorageBackendZK,
			ZooKeeper: StorageZK{
				Addrs:   []string{"127.0.0.1"},
				Timeout: 10000, // 10s
				Dir:     "rhythm",
				Auth: ZKAuth{
					Scheme: ZKAuthSchemeWorld,
				},
				TaskTTL: 1000 * 3600 * 24, // 24h
			},
		},
		Coordinator: Coordinator{
			Backend: CoordinatorBackendZK,
			ZooKeeper: CoordinatorZK{
				Addrs:   []string{"127.0.0.1"},
				Timeout: 10000, // 10s
				Dir:     "rhythm",
				Auth: ZKAuth{
					Scheme: ZKAuthSchemeWorld,
				},
			},
		},
		Secrets: Secrets{
			Backend: SecretsBackendNone,
			Vault: SecretsVault{
				Timeout: 0, // no timeout
				Root:    "secret/rhythm/",
			},
		},
		Mesos: Mesos{
			FailoverTimeout: 1000 * 3600 * 24 * 7, // 7d
			Roles:           []string{"*"},
			Auth: MesosAuth{
				Type: MesosAuthTypeNone,
			},
		},
		Logging: Logging{
			Backend: LoggingBackendNone,
			Level:   LoggingLevelDebug,
			//Level:   LoggingLevelInfo,
		},
	}
	err = json.Unmarshal(file, conf)

	conf.Coordinator.ZooKeeper.ElectionDir = "election/mesos_scheduler"
	// All time.Duration fields from Conf should be in milliseconds so
	// conversion to time elapsed in nanoseconds (represented by time.Duration)
	// is needed.
	millisecondFieldsToDuration(reflect.ValueOf(conf).Elem())
	return conf
}
