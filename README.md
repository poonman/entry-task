# entry-task 说明
- 简介
- 网络通信框架dora v1版说明
- 网络通信框架dora v2版说明
- entry-task服务端实现
- entry-task客户端实现
- 环境准备  
- 测试和Benchmark

### 1 简介
该项目layout有4个部分。
- 网络通信基础框架(命名dora. 取意于哆啦A梦)
- 服务端实现(server)
- 客户端实现(client)
- 协议存储目录(tiger)

### 2 网络通信框架dora v1版说明（可跳过，直接看v2版）  
dora v1 版 主要分为 5 大部分：  
- protocol 传输协议
- server 服务端
- client 客户端
- codec 编解码器
- misc 日志库(log)，脚手架工具(protoc-gen-dora)，配置库(lion)和辅助函数(helper)

### 2.1 消息协议
总览  
- 1字节的magic  
- 4字节的包头长度  
- 4字节的payload长度  
- 包头  
- payload  
  
#### 2.1.1 包头  
采用protobuf协议来序列化
```protobuf
syntax = "proto3";
package protocol;
// 通用头部
message Head{
    // 消息类型定义
    enum MessageType {
        Request = 0;
        Response = 1;
        Heartbeat = 2; // 暂未实现
    }

    // 版本号
    int32 version = 1;
    // 消息类型
    MessageType message_type = 2;
    // payload 的序列化方式 such as, proto，json and xml etc. 可扩展
    string serialize_type = 3;
    // 请求序号
    uint64 seq = 4;
}

message PkgHead {
    // 通用头部
    Head head = 1;
    // 方法名
    string Method = 2;
    // 可扩展元数据
    map<string, string> meta =3;
}
```
### 2.2 服务端  
- 采用同步阻塞的方式处理消息的收发。  
- 支持请求处理前的拦截器  
- 长连接  
  
### 2.3 客户端
- 每条连接采用同步阻塞的方式处理消息收发
- 长连接  
- 实现连接池
- 连接可复用

### 2.4 编解码器
- 定义了编解码器的接口，目前只支持protobuf的编解码，可扩展
```go
package codec
type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	Name() string
}
```

### 2.5 misc
- 基础的日志库，目前只支持打印在标准输出，可实现接口扩展
- 实现了桩代码生成工具，protoc-gen-dora，可根据protobuf协议文件快速生成stub code
- 配置库 lion（使用了第三方库，魔改了一下）
- 辅助函数，依赖注入的shorthand实现


## 3 网络通信框架dora v2版说明

dora v2版在v1版的基础上，新增了transport，变更了dora协议。主要分为 6 大部分：
- transport: 传输层（新增）。更通用的传输层协议(命名nap: new awesome protocol)
- protocol: dora通信的框架层协议（变更）。在transport层协议之上，使用protobuf序列化
- server: 服务端
- client: 客户端
- codec: 编解码器
- misc: 日志库(log)，脚手架工具(protoc-gen-dora)，配置库(lion)和辅助函数(helper)

### 3.1 transport
协议
- 4字节的payload长度
- 4字节的streamID
- 1字节的消息类型
- 1字节的标记位
- payload

```go
package frame

type FrameType uint8

const (
    FrameDate FrameType = 0x0 // 数据
    FramePing FrameType = 0x1 // 心跳
)

type Flags uint8

const (
    FlagPingAck Flags = 0x1
)

type FrameHeader struct {
    Type FrameType
    Flags Flags
    Length uint32
    ID uint32
}
```

### 3.2 dora框架协议
```protobuf
syntax = "proto3";
package protocol;
option go_package = "protocol";

message Head{
    // 版本号
    int32 version = 1;
    // 序列化方式：proto or json
    string serialize_type = 2;
    // 请求序号
    uint64 seq = 3;
    // 请求方法
    string method = 4;
    // 扩展字段
    map<string, string> meta =5;
}

message Pkg {
    Head head = 1;
    bytes payload = 2;
}
```

### 3.3 服务端
- 支持基于方法名的路由
- 采用同步阻塞的方式处理消息的收发。
- 支持拦截器
- 长连接
- 连接支持心跳（新增）

### 3.4 客户端
- 每条连接采用同步阻塞的方式处理消息收发
- 长连接
- 支持心跳（新增）
- 实现连接池，连接可复用

### 3.5 编解码器
- 定义了编解码器的接口，目前只支持protobuf的编解码，可扩展
```go
package codec
type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	Name() string
}
```

### 3.6 misc
- 基础的日志库，目前只支持输出到标准输出，可实现接口扩展
- 实现了protobuf的桩代码生成工具，protoc-gen-dora，可根据protobuf协议文件快速生成stub code
- 配置库 lion（使用了第三方库，魔改了一下）
- 辅助函数，依赖注入的shorthand实现

## 4 entry-task服务端实现
基于dora框架实现的服务端，包含以下几个部分
- DDD建模（api,app,domain,infra）
- 依赖注入(di)
- 业务模型(domain/aggr)
- 拦截器实现（api/interceptor)


### 4.1 DDD（领域驱动设计）
DDD层级划分：接口层（api），应用层（app），领域层（domain），基础设施层（infra）
每一层的职责如下：
- 接口层：顾名思义，接口暴露的入口
- 应用层：usecase的实现，应用层的基本流程主要有3步：加载对象，执行对象的业务逻辑（可能有多个对象交互），对象落地（可选）
- 领域层：业务对象的定义和逻辑实现，以及基础设施层的接口定义
- 基础设施层：主要是仓储的实现，外部网关的实现，配置。

### 4.2 依赖注入
因为大部分的组件都是单例模式，所以采用依赖注入的方式，使得初始化行为不再繁琐。

### 4.3 业务模型
聚合对象主要有以下几个部分：
- account：账号
- kvstore：kv存储
- quota：配额
- session：会话管理(会话暂时放在内存中，未做超时处理，可以考虑写入redis中)
- ratelimiter：限流器。实现依赖于官方包
    
### 4.4 拦截器
实现了两个拦截器(interceptor)：
- 限流 limit
- 会话认证 auth

## 5 entry-task客户端实现
基于dora框架实现的客户端，包含以下几个部分
- DDD建模 
- 依赖注入的方式
- benchmark实现
- 较为丰富的CLI实现


## 6 环境准备
### 6.1 组件部署

#### redis
```shell
brew install redis ## 安装版本为： redis 6.2.3
```
修改密码。redis 6以上的版本增加了ACL
```shell
vim /usr/local/etc/redis.conf
```
``` shell
requirepass: 123456 ## 设置密码 123456
```
以standalone模式运行redis服务端
```shell
redis-server /usr/local/etc/redis.conf
```

#### mysql 
mysql 采用了docker部署
```shell
docker pull mysql:latest
docker run --name mysql-server -v /Users/wen.pan/data/docker/mysql:/var/lib/mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root  -d mysql:latest --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```

### 6.2 数据准备
#### 账户数据（account） 和 配额数据（quota）写入mysql
见链接：[SQL存储过程](https://github.com/poonman/entry-task/-/tree/master/server/sql)

### 证书文件准备
- server.pem
- server.key
- client.pem
- client.key  
参考链接：[证书生成](https://colobu.com/2016/06/07/simple-golang-tls-examples/)

### 服务端配置文件准备
- config.yaml 必须和可执行文件 etserver 在同一级目录下
配置文件内容见链接：[服务端配置文件](https://github.com/poonman/entry-task/-/blob/master/server/bin/config.yaml)

### 客户端配置文件准备
- config.yaml 必须和可执行文件 etclient 在同一级目录下
配置文件内容见链接：[客户端配置文件](https://github.com/poonman/entry-task/-/blob/master/client/bin/config.yaml)
  
## 7 测试和Benchmark

### 7.1 客户端运行参数说明
子命令有 benchmark, login, read, write。子命令可以组合使用（以','分割），比如 etclient login,read -u=1 -k=a

查看帮助说明
```shell
./etclient --help
Usage: etclient COMMAND [OPTIONS]

Support composed command separated by ','. eg: './etclient login,read -u 1 -p 1 -k a' 

Commands:
        benchmark
        login
        read
        write

Run 'etclient COMMAND --help' for more information on a Commands.

```

查看子命令的帮助说明  
- login
```shell
$ ./etclient login --help
Usage of login:
  -p, --password string   password (default "1")
  -u, --user string       username (default "1")
```

- read
```shell
$ ./etclient read --help
Usage of read:
  -k, --key string        key must be in the range [a:z] (default "a")
  -p, --password string   password (default "1")
  -u, --user string       username (default "1")

```

- write
```shell
$ ./etclient write --help
Usage of write:
  -k, --key string        key must be in the range [a:z] (default "a")
  -p, --password string   password (default "1")
  -u, --user string       username (default "1")
  -v, --value string      value must be in the range [a:z] (default "a")

```

- benchmark
```shell
$ ./etclient benchmark --help
Usage: etclient benchmark METHOD [OPTIONS]

Methods:
        read
        write

Run 'kvClient benchmark METHOD --help' for more information on a Commands.

```

- benchmark read
```shell
$ ./etclient benchmark read --help
Usage of etclient benchmark read:
  -c, --concurrency int   concurrency (default 1)
  -k, --key string        key must be in the range [a:z] (default "a")
  -p, --password string   password (default "1")
  -r, --requests int      requests (default 1)
  -u, --user string       username (default "1")

```

- benchmark write
```shell
$ ./etclient benchmark write --help
Usage of etclient benchmark write:
  -c, --concurrency int   concurrency (default 1)
  -k, --key string        key must be in the range [a:z] (default "a")
  -p, --password string   password (default "1")
  -r, --requests int      requests (default 1)
  -u, --user string       username (default "1")
  -v, --value string      value must be in the range [a:z] (default "a")
```

### 7.2 客户端测试示例
- 登陆
```shell
./etclient login -u 100 -p 100
```

- 未登陆读
```shell
./etclient read -u 100 -p 100 -k a
```

- 登陆读
```shell
./etclient login,read -u 100 -p 100 -k a
```

- 未登陆写
```shell
./etclient write -u 100 -p 100 -k a -v m 
```

- 登陆写
```shell
./etclient login,write -u 100 -p 100 -k a -v m 
```

- 登陆读和写
```shell
./etclient login,write,read -u 100 -p 100 -k b -v n 
```

### 7.3 Benchmark
Benchmark 只能测试 read 和 write 接口
```shell
./etclient benchmark read --user 100 --password 100 --key b --concurrency 10 --requests 100 
```

#### 报告
##### 本机环境：
- Macbook Pro-i7 6核 12线程 16G内存
- 运行redis，mysql，entry-task-server, entry-task-kvClient  


##### dora v2 版测试报告
```shell
$./etclient benchmark read -u 100000 -p 100000 -k c -c 12 -r 10000
 Benchmark Report:[{
    "Concurrency": 12,
    "requests": 10000,
    "Success": 120000,
    "Failure": 0,
    "QPS": 27671.387,
    "MaxRT": "45.295218ms",
    "MinRT": "152.847µs",
    "AvgRT": "433.66µs",
    "Duration": "4.345747431s",
    "LatencySummary": {
        "Latencies": [
            " 50% :    327.426µs (cumulative count 60000)",
            " 60% :    336.667µs (cumulative count 72000)",
            " 70% :    346.643µs (cumulative count 84000)",
            " 80% :    358.962µs (cumulative count 96000)",
            " 90% :    378.362µs (cumulative count 108000)",
            " 95% :    393.953µs (cumulative count 114000)",
            " 99% :    414.565µs (cumulative count 118800)",
            "100% :     433.66µs (cumulative count 120000)"
        ]
    },
    "username": "100000",
    "method": "read",
    "key": 99,
    "value": 97
}]

```

