# xmicro

### WISH
* 基础的微服务
* 有多种配置中心方式
* 有多种服务发现方式
* 多种网关选择
* grpc

### DID
* 封装微服务服务注册，支持Nacos、ETCD，可扩展✅
  * 用的单例模式+工厂模式，用工厂函数实现而不是工厂类，用的`GetInstance`创建和返回 `Registry` 实例
* 自定义返回体结构 ✅
* zap日志 & lumberjack 记录日志 ✅
* nacos 配置中心 ✅
* nacos 服务注册与发现 ✅
* nacos 停止程序后取消注册服务 ✅

### TODO
* 封装微服务发现&启动程序时发现多个服务
* etcd和nacos的注册发现目前是分的两个方法，看能不能用设计模式优化下
* 链路追踪
* etcd 配置中心
* 多房间聊天室服务


### 项目配置
1 执行 docker-compose.yml 安装基础服务
2 安装nacos服务端，创建命名空间ID，将config下面的nacos配置改为新的配置
3 将"nacos上配置文件的备份"里的服务配置文件放在nacos上，可根据需求进行调整