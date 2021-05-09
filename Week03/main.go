package main

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/http"
	merrors "github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"main.go/internal/serverControl"
	"time"
)

func main(){
	//创建一个errGroup用于处理上下文关系
	mainctx := context.Background()
	basectx , mainCancelFunc := context.WithCancel(mainctx)
	g, ctx := errgroup.WithContext(basectx)

	//构建父app
	hs := http.NewServer(
		http.Address(":8080"),
		http.Timeout(10 * time.Second),
	)

	app := serverControl.New(
		serverControl.Server(hs),
		serverControl.Context(ctx),
		serverControl.Name("TestServer"),
	)

	time.AfterFunc(5 * time.Second, func() {
		mainCancelFunc()
	})

	g.Go(func() error {
		<-ctx.Done() // wait for stop signal
		return nil
	})
	g.Go(func() error {
		return app.Run()
	})

	if err := g.Wait(); err != nil && !merrors.Is(err, context.Canceled) {
		fmt.Printf("APP error.Cause Is : %+v \n", merrors.Cause(err))
	}
}

