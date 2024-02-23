# dmserver

#### 介绍
dm项目后端代码

#### 软件架构
软件架构说明


#### 安装教程

1. 安装go1.20.4
2. 执行 go mod tidy，继承依赖
3. 执行 go build 编译可执行文件
4. 执行 go install 可编译输出release文件，输出目录 /root/go/bin/

#### 使用说明

1. 需要安装 rabbitmq，mysql，redis，相关配置在根路径下 config.yml 中
2. 修改proto文件后，需要执行 protoc --go_out=./ --go-grpc_out=./ ./consumer/proto/streamMessage.proto 输出生成代码
3. rabbitmq

#### 分支介绍

1.  main 稳定分支
2.  develop 开发分支

#### 服务器相关路径
源码路径：/home/project/douyinApi/
服务路径：/home/project/douyinServers/

nginx 配置路径：/www/server/nginx/conf/vhost/

#### 待完成功能

- ✅ web基本框架集成
- ✅ grpc，rabbitmq集成
- ✅ c# grpc demo编写
- ✅ rabbitmq调试
- ✅ 配置apifox自动化测试接口
- ✅ 改写grpc成websocket
- ⭕️ ~~白名单检查~~
- ✅ 登录接口
- ✅ 世界排行榜功能
- ✅ nginx 负载均衡搭建及配置
- ✅ 礼物记录入库
- ✅ 联调测试
- ✅ 完整链路压测
- ⭕️ 完整链路长跑测试


