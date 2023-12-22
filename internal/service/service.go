package space_service

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cool-team-official/cool-admin-go/cool"
	v1 "github.com/cool-team-official/cool-admin-go/modules/space/pb/v1"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type Controller struct {
	v1.UnimplementedCrudServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterCrudServer(s.Server, &Controller{})
}

func (*Controller) Page(ctx context.Context, req *v1.PageReq) (res *v1.PageRes, err error) {
	return
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) Add(ctx context.Context, req *v1.AddReq) (res *v1.AddRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) Delet(ctx context.Context, req *v1.DeletReq) (res *v1.DeletRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) Info(ctx context.Context, req *v1.InfoReq) (res *v1.InfoRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) SignUrl(ctx context.Context, req *v1.OssSignReq) (res *v1.OssSignRes, err error) {
	var (
		CachePrefix         = "space:sign:"
		OssHost             = cool.GetCfgWithDefault(ctx, "OSS_HOST", gvar.New("")).String()
		OssAccessKeySecret  = cool.GetCfgWithDefault(ctx, "OSS_ACCESS_KEY_SECRET", gvar.New("")).String()
		OssExpires          = cool.GetCfgWithDefault(ctx, "OSS_EXPIRES", gvar.New(120)).Int64()
		SuccessActionStatus = cool.GetCfgWithDefault(ctx, "SUCCESS_ACTION_STATUS", gvar.New(200)).Int32()
		OssAccessKeyId      = cool.GetCfgWithDefault(ctx, "OSS_ACCESS_KEY_ID", gvar.New("")).String()
		OssContentLength    = cool.GetCfgWithDefault(ctx, "OSS_CONTENT_LENGTH", gvar.New(1048576000)).Int64()
		Key                 = fmt.Sprintf("%s/%d/", req.Dir, req.Uid)
	)
	// 查看缓存是否存在
	cacheResp := &v1.OssSignRes{}
	cache, err := cool.CacheManager.Get(ctx, fmt.Sprintf("%s%s", CachePrefix, Key))
	cache.Scan(cacheResp)

	if err != nil {
		g.Log().Error(ctx, "读取缓存出错", err)
	}
	if cacheResp.GetPolicy() != "" {
		g.Log().Info(ctx, "当前签名信息数据来源RedisCache")
		return cacheResp, nil
	}
	if OssAccessKeyId == "" {
		// vera.ToolKit.HandleErrorCtx(ctx, "", true, fmt.Errorf("OSS_ACCESS_KEY_ID Is Empty %s", "请设置环境变量OSS_ACCESS_KEY_ID"))
		g.Log().Error(ctx, fmt.Errorf("OSS_ACCESS_KEY_ID Is Empty %s", "请设置环境变量OSS_ACCESS_KEY_ID"))

	}
	if OssAccessKeySecret == "" {
		g.Log().Error(ctx, fmt.Errorf("OSS_ACCESS_KEY_SECRET Is Empty %s", "请设置环境变量OSS_ACCESS_KEY_SECRET"))
	}
	// 获取当前时间并加上环境变量的过期秒作为过期时间
	OssExpiresTimestamp := time.Now().UTC().Add(time.Duration(OssExpires) * time.Second)
	// 构建策略
	expiration := OssExpiresTimestamp.Format("2006-01-02T15:04:05.000Z")
	policy := map[string]interface{}{
		"expiration": expiration,
		"conditions": []interface{}{
			[]interface{}{"content-length-range", 0, OssContentLength},
		},
	}

	// 将策略转换为 JSON 并编码为 Base64
	policyJSON, err := json.Marshal(policy)
	if err != nil {
		g.Log().Error(ctx, "将策略转换为 JSON 并编码为 Base64 失败", err)
	}
	policyBase64 := base64.StdEncoding.EncodeToString(policyJSON)

	// 计算签名
	h := hmac.New(sha1.New, []byte(OssAccessKeySecret))
	h.Write([]byte(policyBase64))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// 返回RPC结果
	resp := &v1.OssSignRes{
		Policy:                policyBase64,
		OSSAccessKeyId:        OssAccessKeyId,
		Signature:             signature,
		Host:                  OssHost,
		Key:                   Key,
		SUCCESS_ACTION_STATUS: SuccessActionStatus,
	}
	// 缓存当前签名
	err = cool.CacheEPS.Set(ctx, fmt.Sprintf("%s%s", CachePrefix, Key), resp, time.Duration(OssExpires)*time.Second)
	if err != nil {
		g.Log().Error(ctx, "签名信息存储失败")
	}
	return resp, nil

}
