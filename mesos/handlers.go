package mesos

import (
	"context"
	"time"

	"github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/extras/scheduler/controller"
	"github.com/mesos/mesos-go/api/v1/lib/extras/scheduler/eventrules"
	"github.com/mesos/mesos-go/api/v1/lib/extras/store"
	"github.com/mesos/mesos-go/api/v1/lib/scheduler"
	"github.com/mesos/mesos-go/api/v1/lib/scheduler/calls"
	"github.com/mesos/mesos-go/api/v1/lib/scheduler/events"
	"github.com/seoyhaein/caleb-middle/mesos/jobsscheduler"
	log "github.com/sirupsen/logrus"
)

// eventHandler 를 만드는 방식은 몇가지가 있다. 이것을 기억해야함.
// eventrules.Rule 타입
func buildSubscribedEventHandler(fidStore store.Singleton, failoverTimeout time.Duration, onSuccess func(*scheduler.Event)) eventrules.Rule {
	return eventrules.New(
		controller.TrackSubscription(fidStore, failoverTimeout),
		func(ctx context.Context, e *scheduler.Event, err error, ch eventrules.Chain) (context.Context, *scheduler.Event, error) {
			if err == nil {
				onSuccess(e)
			}
			return ch(ctx, e, err)
		},
	)
}

// offer 맞는 것은 적용하고 맞지 않는 것은 decline 해야함.
// events.HandlerFunc 타입
// 파라미터 생략 알아보기 적용하자.일단 주석처리

// offer 를 따로 저장하고 잡을 수행하는 시점을 해당 이벤트가 왔을때로 하면 안됨.
// 그리고 필요없는 offer 들은 decline 시켜야함.
// 여기서는 offer 를 저장하는역활만 하면됨.
func buildOffersEventHandler(cli calls.Caller) events.HandlerFunc {
	return func(ctx context.Context, e *scheduler.Event) error {
		offers := e.GetOffers().GetOffers()

		log.Info("Number of received offers: %d\n", len(offers))

		for i := range offers {
			if ctx.Err() != nil {
				break
			}
			offer := offers[i]
			s := "/home/dev-comd/go/src/github.com/seoyhaein/caleb-middle/bin/echo.json"
			task, e := jobsscheduler.FindTasksForOffer(&offer, s)
			if e != nil {
				log.Info("Failed to find tasks for offer: %s\n", e)
				return e
			}
			//tasks := jobsscheduler.FindTasksForOfferT(ctx, &offer)
			// tasks 가 배열로 들어가면서 offer.ID 는 한개라면 좀 이해가 안됨. Offer가 여러개 들어가야함. 그것도 맞게 들어가야함.
			accept := calls.Accept(calls.OfferOperations{calls.OpLaunch(task)}.WithOffers(offer.ID))
			// task 정보가 잘못되었다는 소리인가???? mesos-go 다시 샆려보기
			err := calls.CallNoData(ctx, cli, accept.With(calls.RefuseSeconds(time.Hour)))

			if err != nil {
				log.Info("Failed to accept offer: %s", err)
				return err
			}
			//for _, task := range tasks {
			log.Info("Task staged: %s", task.TaskID.Value)
			//}

		}
		return nil
	}
}

func buildUpdateEventHandler(cli calls.Caller) eventrules.Rule {
	var h eventrules.Rule
	h = func(ctx context.Context, e *scheduler.Event, err error, ch eventrules.Chain) (context.Context, *scheduler.Event, error) {

		status := e.GetUpdate().GetStatus()
		switch st := status.GetState(); st {
		case mesos.TASK_FINISHED, mesos.TASK_RUNNING, mesos.TASK_STAGING, mesos.TASK_STARTING:
			log.Info("status update from agent %q: %v", status.GetAgentID().GetValue(), st)
			if st == mesos.TASK_RUNNING && status.AgentID != nil {

				log.Info("AgentID: %s, TaskID: %s, %v\n", status.AgentID.Value, status.TaskID.Value, st)
				return ch(ctx, e, err)
				/*cid := status.GetContainerStatus().GetContainerID()
				if cid != nil {
					fmt.Println("attaching for interactive session to agent %q container %q", status.AgentID.Value, cid.Value)
				}*/

			}
			if st != mesos.TASK_FINISHED {
				log.Info("AgentID: %s, TaskID: %s, %v\n", status.AgentID.Value, status.TaskID.Value, st)
				return ch(ctx, e, err)
			}

		case mesos.TASK_LOST:
			log.Info("AgentID: %s, TaskID: %s, %v\n", status.AgentID.Value, status.TaskID.Value, st)
			return ch(ctx, e, err)
		case mesos.TASK_FAILED:
			log.Info("AgentID: %s, TaskID: %s, %v\n", status.AgentID.Value, status.TaskID.Value, st)
			return ch(ctx, e, err)
		case mesos.TASK_KILLED:
			log.Info("AgentID: %s, TaskID: %s, %v\n", status.AgentID.Value, status.TaskID.Value, st)
			return ch(ctx, e, err)
			//fallthrough
		case mesos.TASK_ERROR:
			log.Info("AgentID: %s, TaskID: %s, %v\n", status.AgentID.Value, status.TaskID.Value, st)
			return ch(ctx, e, err)
		default:
			log.Info("unexpected task state, aborting %v\n", st)
			return ch(ctx, e, err)
		}
		return ch(ctx, e, err)
	}
	return h.AndThen(controller.AckStatusUpdates(cli))
}
