// only for test
package jobtester

import (
	"encoding/json"
	"io/ioutil"

	"github.com/seoyhaein/caleb-middle/model"
)

func CreateJob(path string) (*model.Job, error) {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	/*jobConf := &model.JobConf{
		JobID: model.JobID{},
		Schedule: model.JobSchedule{},
		}
	jobRuntime := &model.JobRuntime{}
	job := &model.Job{JobConf: *jobConf, JobRuntime: *jobRuntime}*/

	// error 확인하자!!
	var job model.Job
	//job := &model.Job{}
	err = json.Unmarshal(file, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
