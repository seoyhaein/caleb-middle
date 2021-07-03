package command

import (
	"context"

	"github.com/seoyhaein/caleb-middle/config"
	"github.com/seoyhaein/caleb-middle/mesos"
)

// 7/3 일단 메소스 만 넣었음.
// 이 함수에서 일단 메소스 시작하고 grpc 서버 시작해야 함.
// 그리고 이 둘 함수 사이에서 채널을 연결해서 통신을 해야 할 듯하다.
func Run(ctx context.Context, args []string) error {

	// mesos 환경설정 파일 읽어와서 세팅해줌.
	conf := config.New(args) // 수정 검토. readability 관점
	// error 처리 해야함.
	err := mesos.MesosRun(ctx, conf)

	return err
}
