package main

import (
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/http"
	merrors "github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"time"
)

var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server, rr registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.Name("beer.cart.service"),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
		kratos.Registrar(rr),
	)
}

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

