package jobsscheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/resources"
	"github.com/seoyhaein/caleb-middle/model"
)

// offer 에서 GetURL.Address.Ip 을 제공해주고 있는데, 이 노드 생성 시점에서 IP를 아는게 중요함. 이것으로 직접 agent 에 연결하는 방식으로 하는 것 구상.
// 개발 순서를 일단 문서로 정리하자.

const RFC3339a = "20060102T150405Z0700"

// Scheduler decides which jobs to run in response to received offers.
type Scheduler struct {
	offers []mesos.Offer
}

// FindTasksForOffer returns tasks to run for passed offer.
// IP 로 찾는 함수로 구성해야함.
func (sched *Scheduler) FindTasksForOffer(ctx context.Context, offer *mesos.Offer) []mesos.TaskInfo {
	/*rs := mesos.Resources(offer.Resources)
	log.Debugf("Finding tasks for offer: %s", rs)
	// 리소스에서 job 을 찾고
	jobs, jobsRs := sched.findJobsForResources(rs)
	log.Debugf("Found %d tasks for offer", len(jobs))
	// 찾은 job을 task 로 바꿔줌
	tasks := sched.buildTasksForOffer(jobs, jobsRs, offer)
	return tasks*/
	return nil
}

// offer 가 들어오면 Task로 만든다.
// job 의 리소스와 비교하는 구문이 들어가야한다.
// decline 관련 자료 조사 필요.
// 리소스 비교하는 것 넣자.!!!!
func TransformJobToTaskInfo(ctx context.Context, path string, offer *mesos.Offer) (*mesos.TaskInfo, error) {

	return nil, nil
}

// job.json 형태의 json 에서 가져와서 job 을 만든다.
func createJobs(path string) ([]*model.Job, error) {
	// 아래 코드 검증하자. 에러가 좀 있늗듯.
	// 여기서는 일단 job 이 하나라는 것을 가정하고 시작한다.
	// cap 을 일단 10개로만 두자.
	jobs := make([]*model.Job, 0, 10)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// error 확인하자!!
	var job model.Job
	//job := &model.Job{}
	err = json.Unmarshal(file, &job)
	if err != nil {
		return nil, err
	}
	jobs = append(jobs, &job)
	return jobs, nil
}
func createJob(path string) (*model.Job, error) {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// error 확인하자!!
	var job model.Job
	err = json.Unmarshal(file, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

// 일단, job 이 여러개 들어오고, offer 도 여러개 들어오는 경우를 상정하자.
// index 는 decline 할 offer 들임

// 향후, offer 를 가지고 있는 구조가 있어야함. 그리고 여기서 파라미터로 받아야함.을 수 있
// 이러한 리소스 방식은 버그가 있음.
// 일단 인스턴스를 만들때, job 의 리소스를 통해서 적절한 인스턴스를 생성한다.
// 그리고 해당 인스턴스를 인식하는 방식을 고려해서(일단은..) 해당 인스턴스의 offer 를 확인하고(offer 에서 ip 가 있음.)
// 최종적으로 해당 offer 와 해당 task 의 리소스를 최종확인한다.
// 이와 같은 작업은 클라우드에서 가능하다.
// 따라서, 일단 job의 리소스를 통해서 taskinfo 를 만들자.
// decline 의 경우, 생략해도 될듯하다. 그 이유는 일단 생성되는 offer 자체(인스턴스)가 job 의 갯수와 동일하게 생성된다.
// instance 를 생성하는 알고리즘의 경우도 생각해봐야 할 문제이다. 즉 한번에 여러개를 만드는 것이 아니라 offer 를 소모하는 것에 따라서 instance 를 생성해줘야 한다.

// 작업 단위를 잘 쪼개서 함수로 구현해야한다.

/*func findResourceforTask(jobs []*model.Job, offers []mesos.Offer) (offers.Index, []*mesos.TaskInfo, error) {

	var (
		matched mesos.Resources
		matchedRs []mesos.Resources
	)

	for _, offer := range offers {
		offeredResources := mesos.Resources{}.Plus(offer.Resources...)


		for _, o := range offeredResources {
			for _, job := range jobs {
				jr := job.Resources()
				for _, j1 := range jr {
				if o.Contains(j1) {
					// job resource 가 match 되는게 있다면,
					offeredLessWants := mesos.Resources{o}.Minus(j1)
					matched = append(matched, mesos.Resources{o}.Minus(offeredLessWants...)...)
					jr.Subtract1(j1)										// 루프를 적게돌게하기 위함.
					offeredResources.Subtract1(matched[len(matched)-1]) 	// 루프를 적게돌게하기 위함.
					}
				}
			}

			matchedRs = append(matchedRs, matched) // 하나의 오퍼기준으로 집어넣는다.

		}
	}


	// jobs 의 하나의 job 은 offer 와 대응된다.
	// offer.Resources 는 []Resource 이다. 이것은 Cpus, Mem, Disks, Port... 등의 리소스들의 set 이다.


	return nil, nil, nil

}*/

// 테스트를 위해 만든 함수
// 한개의 offer, 하나의 job
// context 는 생략 생각해보자

// FindTasksForOffer returns tasks to run for passed offer.
func FindTasksForOfferT(ctx context.Context, offer *mesos.Offer) []mesos.TaskInfo {

	rs := mesos.Resources(offer.Resources)
	fmt.Println("Finding tasks for offer: %s", rs)
	// 리소스에서 job 을 찾고
	//jobs, jobsRs := sched.findJobsForResources(rs)
	jobsRs := findJobsForResourcesT(rs)
	//log.Debugf("Found %d tasks for offer", len(jobs))
	// 찾은 job을 task 로 바꿔줌
	//tasks := sched.buildTasksForOffer(jobs, jobsRs, offer)
	tasks := buildTasksForOfferT(offer, jobsRs)

	return tasks
}

func findJobsForResourcesT(res mesos.Resources) []mesos.Resources {
	var tasksRes []mesos.Resources

	// 리소스 초기화 - 상세히 분석필요
	res = res.Unallocate() // resource 의 role 을 nil 로 만들어 버림.....
	//resUnreserved := res.ToUnreserved()
	// 리소스에서 뭔가 예약이 걸려있는거 같다.
	jobRes := model.ResourcesT(4, 6864)

	taskRes := withAllocationRole("*", jobRes)
	tasksRes = append(tasksRes, taskRes)

	return tasksRes
}

func ResourcesT(c float64, m float64) mesos.Resources {
	res := mesos.Resources{}
	res.Add(
		resources.NewCPUs(c).Resource,
		resources.NewMemory(m).Resource,
	)
	return res
}

func buildTasksForOfferT(offer *mesos.Offer, ress []mesos.Resources) []mesos.TaskInfo {
	var tasks []mesos.TaskInfo
	task, _ := newTaskInfoT()
	task.AgentID = offer.AgentID
	validateAll(ress[0])
	task.Resources = ress[0]
	tasks = append(tasks, *task)

	return tasks
}

func newTaskInfoT() (*mesos.TaskInfo, error) {
	ts := time.Now().Format(RFC3339a)
	task := mesos.TaskInfo{

		TaskID: mesos.TaskID{Value: ts},
		Name:   "Task ",
		Command: &mesos.CommandInfo{
			Value: proto.String("echo"),
			//Environment: &env,					// 위에서 설정한 값을 세팅해줌.
			//User:      proto.String("someone"),
			Shell:     proto.Bool(false), // Shell true 면 Arguments 무시, false 면 Arguments 사용
			Arguments: []string{"echo", "helloworld"},
		},
		//Container: &containerInfo,
		//Labels:    &mesos.Labels{labels},
	}

	return &task, nil
}

func validateAll(r mesos.Resources) {
	for i := range r {
		rr := &r[i]
		if err := rr.Validate(); err != nil {
			panic(err)
		}
	}
}

func FindTasksForOffer( /*ctx context.Context,*/ offer *mesos.Offer, path string) (mesos.TaskInfo, error) {
	job, err := createJob(path)
	ts := time.Now().Format(RFC3339a)
	task := mesos.TaskInfo{
		TaskID: mesos.TaskID{Value: ts},
		Name:   "Task ",
		Command: &mesos.CommandInfo{
			Value:     proto.String("echo"),
			User:      proto.String("root"),
			Shell:     proto.Bool(false), // Shell true 면 Arguments 무시, false 면 Arguments 사용
			Arguments: []string{"echo", "hello"},
		},
		AgentID:   offer.AgentID,
		Resources: withAllocationRole("*", job.Resources()),
		//Resources: withAllocationRole(job.Role, job.Resources()),
	}
	return task, err
}

func withAllocationRole(role string, r mesos.Resources) mesos.Resources {
	result := make(mesos.Resources, 0, len(r))
	for i := range r {
		rr := &r[i]
		if rr.GetAllocationInfo().GetRole() != role {
			rr = proto.Clone(rr).(*mesos.Resource) // 머리가 돌대가리인가.. 계속 잊어버림.
			rr.AllocationInfo = &mesos.Resource_AllocationInfo{
				Role: proto.String(role),
			}
		}
		result = append(result, *rr)
	}
	return result
}

func newTaskInfo(job *model.Job) (*mesos.TaskInfo, error) {
	tid, err := newTaskID(&job.JobID) // taskID 만들어줌.
	if err != nil {
		/* 로그는 일치시키자 향후에...*/
		return nil, fmt.Errorf("Getting task ID failed: %s", err)
	}
	// mesos.Environment 는 CommandInfo 사용시 환경변수로 쓰임.
	env := mesos.Environment{
		Variables: []mesos.Environment_Variable{
			//	{Name: "RHYTHM_TASK_ID", Value: &tid},
			//	{Name: "RHYTHM_MEM", Value: proto.String(fmt.Sprintf("%g", job.Mem))},
			//	{Name: "RHYTHM_DISK", Value: proto.String(fmt.Sprintf("%g", job.Disk))},
			//	{Name: "RHYTHM_CPU", Value: proto.String(fmt.Sprintf("%g", job.CPUs))},
		},
	}
	// 추가적으로 job.Env 에 있는 녀석들을 넣어줌. job 생성 기준 살펴보자.
	for k, v := range job.Env {
		envvar := mesos.Environment_Variable{Name: k, Value: proto.String(v)}
		env.Variables = append(env.Variables, envvar)
	}
	// 여기까지 환경변수 설정.

	//secretes 는 일단 주석.
	/*	for k, v := range job.Secrets {
		path := fmt.Sprintf("%s/%s/%s", job.Group, job.Project, v)
		secret, err := sched.secrets.Read(path)
		if err != nil {
			return nil, fmt.Errorf("Reading secret failed: %s", err)
		}
		envvar := mesos.Environment_Variable{Name: k, Value: &secret}
		env.Variables = append(env.Variables, envvar)
	}*/

	// 컨테이너 설정
	// job 에서 컨테이너를 docker 타입인지, 메소스 타입인지 결정.
	/*var containerInfo mesos.ContainerInfo
	switch job.Container.Type {
	case model.Docker:
		containerInfo = mesos.ContainerInfo{
			Type: mesos.ContainerInfo_DOCKER.Enum(),
			Docker: &mesos.ContainerInfo_DockerInfo{
				Image:          job.Container.Docker.Image,
				ForcePullImage: &job.Container.Docker.ForcePullImage,
			},
		}
	case model.Mesos:
		containerInfo = mesos.ContainerInfo{
			Type: mesos.ContainerInfo_MESOS.Enum(),
			Docker: &mesos.ContainerInfo_DockerInfo{
				Image: job.Container.Mesos.Image,
			},
		}
	default:
		containerInfo = mesos.ContainerInfo{}
	}*/
	// label 설정.
	// label 설정은 좀더 파악하자.
	labels := make([]mesos.Label, len(job.Labels))
	for k, v := range job.Labels {
		func(v string) {
			labels = append(labels, mesos.Label{Key: k, Value: &v})
		}(v)
	}
	arg := []string{"echo,", "hello world!"}
	// task 생성.
	task := mesos.TaskInfo{
		TaskID: mesos.TaskID{Value: tid},
		Name:   "Task " + tid,
		Command: &mesos.CommandInfo{
			Value: proto.String(job.Cmd),
			//Environment: &env,					// 위에서 설정한 값을 세팅해줌.
			User:      proto.String(job.User),
			Shell:     proto.Bool(false), // Shell true 면 Arguments 무시, false 면 Arguments 사용
			Arguments: arg,
		},
		//Container: &containerInfo,
		//Labels:    &mesos.Labels{labels},
	}
	return &task, nil
}

/* taskID 만들어주는 녀석들 */

type taskID struct {
	jid  *model.JobID
	uuid string
}

func (tid *taskID) String() string {
	return fmt.Sprintf("%s:%s", tid.jid.String(), tid.uuid)
}

func newTaskID(jid *model.JobID) (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	tid := taskID{jid: jid, uuid: u4.String()}
	return tid.String(), nil
}
