
syntax = "proto3";

package <包名>;

option go_package = "<Go的包名>";

// 分页
message PageReq {
    // 当前页
    int64 Page = 1;
    // 当前页条数
    int64 Size = 2;
    // 排序字段
    string Order = 3;
    // 升序/降序
    string Sort = 4;
    // ...其他参数自行添加
    
}

message Pagination {
    // 当前页
    int64 Page = 1;
    // 当前页条数
    int64 Size = 2;
    // 总数
    int64 Total = 3;
}

message ListItem {
    oneof item_type {
        int64 Id = 1;
        // 列表每一项类型
    }
}
message PageRes {
    repeated ListItem List = 1;
    Pagination Pagination = 2;
}
// 列表
message ListReq {
    // 排序字段
    string Order = 1;
    // 升序/降序
    string Sort = 2;
}
message ListRes {
    repeated ListItem List = 1;
}

// 删除
message DeletReq {}
message DeletReq {}
// 新增
message AddReq {}
message AddRes {}
// 详情
message InfoReq {
    int64 Id = 1;
}
message InfoReq {}
// 更新
message UpdateReq {
    int64 Id = 1;
}
message UpdateRes {}
service Crud {
    rpc Page(PageReq) returns(PageRes);
    rpc Page(PageReq) returns(PageRes);

}