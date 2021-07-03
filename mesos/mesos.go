package mesos

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/extras/scheduler/callrules"
	"github.com/mesos/mesos-go/api/v1/lib/extras/scheduler/controller"
	"github.com/mesos/mesos-go/api/v1/lib/extras/scheduler/eventrules"
	"github.com/mesos/mesos-go/api/v1/lib/extras/store"
	"github.com/mesos/mesos-go/api/v1/lib/httpcli"
	"github.com/mesos/mesos-go/api/v1/lib/httpcli/httpsched"
	"github.com/mesos/mesos-go/api/v1/lib/scheduler"
	"github.com/mesos/mesos-go/api/v1/lib/scheduler/calls"
	"github.com/mesos/mesos-go/api/v1/lib/scheduler/events"
	"github.com/seoyhaein/caleb-middle/config"
	log "github.com/sirupsen/logrus"
)

var (
	registrationMinBackoff = 1 * time.Second
	registrationMaxBackoff = 15 * time.Second
)

// FrameworkId 를 어디에다가 저장할까???  일단은 zk 에다가 집어 넣는 방향으로
// 환경설정 부터 시작해야함 일단 여기서는 해당 기능을 생략하고 진행한다.
// 임시로 메모리에다가 저장함.

// TODO 일단 context 살펴봐야함 5/1
// TODO 여기를 함수를 일단 3부분으로 나누는 것에 생각하기, 에러의 처리가 약하다. 7/3
func MesosRun(ctx context.Context, conf *config.Conf) error {

	frameworkIdStore, err := newFrameworkIDStore()
	if err != nil {
		return err
	}

	cli, err := newClient(&conf.Mesos, frameworkIdStore)

	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(ctx)

	logger := controller.LogEvents(func(e *scheduler.Event) {
		log.Printf("Event: %s", e)
	}).Unless(conf.Mesos.LogAllEvents)

	handler := eventrules.New(
		logAllEvents().If(conf.Mesos.LogAllEvents),
		controller.LiftErrors(),
	).Handle(events.Handlers{
		scheduler.Event_HEARTBEAT: events.HandlerFunc(func(ctx context.Context, e *scheduler.Event) error {
			log.Info("Heartbeat")
			return nil
		}),
		scheduler.Event_ERROR: events.HandlerFunc(func(ctx context.Context, e *scheduler.Event) error {
			log.Info(e.GetError().Message)
			return nil
		}),
		scheduler.Event_SUBSCRIBED: buildSubscribedEventHandler(frameworkIdStore, conf.Mesos.FailoverTimeout, func(e *scheduler.Event) {
			log.Info("subscribed...")
		}),
		scheduler.Event_OFFERS: buildOffersEventHandler(cli),
		scheduler.Event_UPDATE: buildUpdateEventHandler(cli),
	}.Otherwise(logger.HandleEvent))

	err = controller.Run(
		ctx,
		newFrameworkInfo(&conf.Mesos, frameworkIdStore),
		cli,
		/*controller.WithRegistrationTokens(
			backoff.Notifier(registrationMinBackoff, registrationMaxBackoff, ctx.Done()),
		),*/
		controller.WithEventHandler(handler),
		controller.WithSubscriptionTerminated(func(err error) {
			log.Info("7/3 이부분 살펴보자.")
			cancel()
			if err == io.EOF { //res.Decode error 발생시 io.EOF, 테스트 필
				log.Info("disconnected")
			}
		}),
	)
	return err
}

func newFrameworkIDStore() (store.Singleton, error) {
	return store.DecorateSingleton(
		store.NewInMemorySingleton(),
		store.DoSet().AndThen(func(_ store.Setter, v string, _ error) error {
			log.Info("Framework ID: %s", v)
			return nil
		})), nil
}

func newFrameworkInfo(conf *config.Mesos, idStore store.Singleton) *mesos.FrameworkInfo {
	labels := make([]mesos.Label, len(conf.Labels))
	for k, v := range conf.Labels {
		func(v string) {
			labels = append(labels, mesos.Label{Key: k, Value: &v})
		}(v)
	}
	frameworkInfo := &mesos.FrameworkInfo{
		User:       conf.User,
		Name:       conf.FrameworkName,
		Checkpoint: &conf.Checkpoint,
		Capabilities: []mesos.FrameworkInfo_Capability{
			{Type: mesos.FrameworkInfo_Capability_MULTI_ROLE},
		},
		Labels:          &mesos.Labels{labels},
		FailoverTimeout: func() *float64 { ft := conf.FailoverTimeout.Seconds(); return &ft }(),
		//WebUiURL:        &config.WebUIURL,
		//Hostname:        &config.Hostname,
		//Principal:       &config.Principal,
		Roles: conf.Roles,
	}
	//TODO 이코드는 조기에 머지 않을 텐데..
	id, _ := idStore.Get()
	frameworkInfo.ID = &mesos.FrameworkID{Value: *proto.String(id)}
	return frameworkInfo
}

// TODO 7/3 여기 코드 부분에서 버그 있음
func newClient(c *config.Mesos, frameworkId store.Singleton) (calls.Caller, error) {
	if len(c.Addrs) == 0 {
		return nil, errors.New("List of Mesos addresses is empty")
	}
	//TODO 7/3
	/*
	  일단 아래 코드의 경우는 c.Addrs[0] 이런식으로 강제로 첫번째 마스터의 ip 를 세팅하도록 하였는데, 중복될 경우는 바꿔줘야함.
	*/

	cli := httpcli.New(
		httpcli.Endpoint(fmt.Sprintf("%s/api/v1/scheduler", c.Addrs[0])),
		httpcli.Do(httpcli.With(
			httpcli.Timeout(time.Second*10),
		)))

	return callrules.New(
		logCalls(map[scheduler.Call_Type]string{scheduler.Call_SUBSCRIBE: "Connecting..."}),
		callrules.WithFrameworkID(store.GetIgnoreErrors(frameworkId)),
	).Caller(httpsched.NewCaller(cli)), nil
}

func logCalls(messages map[scheduler.Call_Type]string) callrules.Rule {
	return func(ctx context.Context, c *scheduler.Call, r mesos.Response, err error, ch callrules.Chain) (context.Context, *scheduler.Call, mesos.Response, error) {
		if message, ok := messages[c.GetType()]; ok {
			log.Info(message)
		}
		return ch(ctx, c, r, err)
	}
}

func logAllEvents() eventrules.Rule {
	return func(ctx context.Context, e *scheduler.Event, err error, ch eventrules.Chain) (context.Context, *scheduler.Event, error) {
		log.Info("%+v", *e)
		return ch(ctx, e, err)
	}
}
