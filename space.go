package main

import (
	_ "github.com/cool-team-official/cool-admin-go/modules/space/controller"
	_ "github.com/cool-team-official/cool-admin-go/modules/space/middleware"
	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"google.golang.org/grpc"
)

func init() {
	var (
		ctx = gctx.GetInitCtx()
	)
	g.Log().Debug(ctx, "module space init start ...")
	g.Log().Debug(ctx, "module space init finished ...")
}

func main() {
	grpcx.Resolver.Register(etcd.New("127.0.0.1:2379"))
	c := grpcx.Server.NewConfig()
	c.Options = append(c.Options, []grpc.ServerOption{
		grpcx.Server.ChainUnary(
			grpcx.Server.UnaryValidate,
		)}...,
	)
	s := grpcx.Server.New(c)
	// space.Register(s)
	s.Run()
}
