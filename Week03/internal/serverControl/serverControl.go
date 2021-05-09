package serverControl

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/google/uuid"
	merrors "github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// App is an application components lifecycle manager
type App struct {
	opts     options			//函数选项模式
	ctx      context.Context
	cancel   func()
	instance *registry.ServiceInstance	//注册
	log      *log.Helper
}

//提供当前APP的ctx的读取方法，用于为子APP提供context树
func (thisapp *App)GetctxFromApp() context.Context {
	//进行深拷贝并传出
	return thisapp.ctx
}

// New create an application lifecycle manager.
func New(opts ...Option) *App {
	options := options{
		ctx:    context.Background(),		//基础context
		logger: log.DefaultLogger,
		sigs:   []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	if id, err := uuid.NewUUID(); err == nil {		//生成唯一辨识资讯
		options.id = id.String()
	}
	for _, o := range opts {		//对New出来的APP中的options进行设置
		o(&options)
	}
	ctx, cancel := context.WithCancel(options.ctx)		//从当前APP中衍生一个用于取消子context的上下文ctx
	return &App{
		opts:     options,
		ctx:      ctx,
		cancel:   cancel,
		instance: buildInstance(options),
		log:      log.NewHelper("app", options.logger),
	}
}

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (a *App) Run() error {
	a.log.Infow(
		"service_id", a.opts.id,
		"service_name", a.opts.name,
		"version", a.opts.version,
	)
	g, ctx := errgroup.WithContext(a.ctx)
	for _, srv := range a.opts.servers {
		srv := srv
		g.Go(func() error {
			<-ctx.Done() // wait for stop signal
			return srv.Stop()
		})
		g.Go(func() error {
			return srv.Start()
		})
	}
	if a.opts.registrar != nil {		//注册服务
		if err := a.opts.registrar.Register(a.opts.ctx, a.instance); err != nil {
			return err
		}
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("APP quit by parent ctx.Done \n")
				return ctx.Err()
			case <-c:
				fmt.Printf("APP stop by os.signal \n")
				a.Stop()
			}
		}
	})
	//除正常停止情况 以及 关闭父服务触发context树关闭的情况外 则返回对应err
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return merrors.Wrapf(err, "Server Control Unexpected Failed")
	}
	return nil
}

// Stop gracefully stops the application.
func (a *App) Stop() error {
	if a.opts.registrar != nil {
		if err := a.opts.registrar.Deregister(a.opts.ctx, a.instance); err != nil {
			return err
		}
	}
	if a.cancel != nil {
		a.cancel()		//取消当前服务以及所有子服务
	}
	return nil
}

func buildInstance(o options) *registry.ServiceInstance {
	if len(o.endpoints) == 0 {
		for _, srv := range o.servers {
			if e, err := srv.Endpoint(); err == nil {
				o.endpoints = append(o.endpoints, e)
			}
		}
	}
	return &registry.ServiceInstance{
		ID:        o.id,
		Name:      o.name,
		Version:   o.version,
		Metadata:  o.metadata,
		Endpoints: o.endpoints,
	}
}
