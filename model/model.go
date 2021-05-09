package model

import (
	"fmt"
	"time"

	"github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/resources"
)

// Container 지원.
// ContainerType defines containerizer genre.
type ContainerType string

// JobDocker defines fields for job running inside Docker container.
type JobDocker struct {
	Image          string
	ForcePullImage bool
}

// JobMesos defines fields for job running inside Mesos container.
// https://mesos.apache.org/documentation/latest/mesos-containerizer/
type JobMesos struct {
	Image string
}

// JobContainer defines containerizer-related fields for job.
type JobContainer struct {
	Type   ContainerType
	Docker *JobDocker `json:",omitempty"`
	Mesos  *JobMesos  `json:",omitempty"`
}

const (
	// Docker denotes Docker containerizer.
	Docker ContainerType = "Docker"
	// Mesos denotes Mesos containerizer.
	// https://mesos.apache.org/documentation/latest/container-image/
	Mesos = "Mesos"
)

// State defines job position.
type State string

// 아래 구문 파악해보자
const (
	// IDLE denotes not running job which either hasn't been scheduled yet or its last run was successful.
	IDLE State = "Idle"
	// STAGING denotes job which has been scheduled to run in response to offer.
	STAGING = "Staging"
	// STARTING denotes job which has been lanunched by executor.
	STARTING = "Starting"
	// RUNNING denotes job which has been started.
	RUNNING = "Running"
	// FAILED denotes job whose last run failed.
	FAILED = "Failed"
)

// JobID defines job identifier.
type JobID struct {
	Group   string
	Project string
	ID      string
}

// JobConf defines job's configuration fields.
// 추가될 필드들이 있음. mesos.TaskInfo 에 일부 전달됨.
type JobConf struct {
	JobID
	//Schedule   JobSchedule
	Env        map[string]string // environment
	Secrets    map[string]string
	Container  JobContainer
	CPUs       float64
	Mem        float64
	Disk       float64
	Cmd        string
	User       string
	Shell      bool
	Role       string // 추가됨.
	Arguments  []string
	Labels     map[string]string
	MaxRetries int

	// 특정노드로 Job 을 보내기 위해 agent ip
	IP *string
}

// JobRuntime defines job's runtime fields.
type JobRuntime struct {
	State          State
	LastStart      time.Time
	CurrentTaskID  string
	CurrentAgentID string
	Retries        int
}

// Job encompasses fields for job's configuration and runtime.
type Job struct {
	JobConf
	JobRuntime
}

// Fully qualified identifier unique across jobs from all groups and projects.
func (jid *JobID) String() string {
	return fmt.Sprintf("%s:%s:%s", jid.Group, jid.Project, jid.ID)
}

// Resources returns resources required by job.
func (j *JobConf) Resources() mesos.Resources {
	res := mesos.Resources{}
	res.Add(
		resources.NewCPUs(j.CPUs).Resource,
		resources.NewMemory(j.Mem).Resource,
		resources.NewDisk(j.Disk).Resource,
	)
	return res
}

func ResourcesT(c float64, m float64) mesos.Resources {
	res := mesos.Resources{}
	res.Add(
		resources.NewCPUs(c).Resource,
		resources.NewMemory(m).Resource,
	)
	return res
}
