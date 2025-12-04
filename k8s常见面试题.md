# k8s架构设计

1）API server: 负责整个k8s集群的认证和授权，访问k8s集群资源的入口，访问方式：kubect/client-go
2）etcd数据库：cr（crd表结构，cr数据库中表记录）
3）controller-manager：node的controller、pod的副本数管理对应的controller、service的controller
  读取etcd中cr的信息，进行调谐，将实际信息写入etcd,努力将期望状态变成实际状态
4）scheduler：读取etcd中信息，按照一定的调度算法将pod调度到指定节点


# 容器网络
网络插件 network.knitter.io，负责容器与容器之间的网络通信。
主要功能：负责动态分配唯一的IP地址，负责pod之间跨节点通信等

网络平面   网卡      哪层网络协议
lan       eth1      layer3   
net_api   eth0      layer3
动态分配ipv4和ipv6

七层网络协议
L7  应用层      如：HTTP/FTP
L6  表示层      数据转换
L5  会话层      建立、管理和终止会话
L4  传输层      数据传输控制，如：TCP/UDP
L3  网络层      ip地址寻址和路由选择，如：IPV4和IpV6
L2 数据链路层   Mac地址寻址，如以太网
L1 物理层       电缆、光纤

# k8s调度原理
如何将pod调度到指定node上
## 调度原理
1、informer：通过API server监听k8s集群资源的变化
  1）list：全量、周期查询
  2）watch：http长连接，监听k8s资源（pod、nodedeng）的变化
2、本地会维护scheduling pod queue（待调度的pod的队列）、scheduling node queue（待调度的node的队列）
3、从待调度的pod队列中取出一个pod，先对node进行过滤（predict算法，比如端口占用、cpu、memory、污点、节点亲和性等），然后对node进行打分（priority算法，比如:cpu、memory剩余量等），评选出最优的node，将node与pod进行绑定
  LOOP
## 调度方式
### nodeNode方式调度
很少使用，一般通过标签进行匹配

### nodeSelector方式调度
```yaml
nodeSelector:
  labels:
    app: test
```
匹配节点的标签
```bash
kubectl label node node1 app=test
kubectl label node node1 app-
```
### taints和tolerations方式调度
taints作用在node上，tolerations作用在pod上
类比：node比作男生，身上有很多缺点（taints），pod比作女生，如果俩人需要在一起，需要女生（pod）去容忍（tolerations）男生（node）的缺点（taints）

##### taints
taints类型
1、NoScheduler: 待运行的pod不会调度到该node上
2、PreferNoScheduler：尽量不要调度该node,如果找不到合适的node，还是会调度到该node上
3、NoExecute: 已经运行在该node上pod会被驱逐（可以在pod上配置容忍时间）
taints数据格式   key:value:type（value可以为空）
```bash
kubectl taint node node1 app=test:NoSchedule
kubectl taint node node1 app=test:NoSchedule-
```

##### tolerations
tolerations
tolerations数据格式   key:value:operator:type：tolerations time（value可以为空）
operator枚举值：Exist|Equal
```yaml
tolerations:
  - key: app
    value: test
    operator: Equal
    effective: Noschedule
```
备注：
1、NoExecute污点类型使用场景：想要驱逐pod，需要在node上添加NoExecute污点

### Affinity调度
#### node亲和调度
##### node硬亲和
```yaml
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
        - matchExpressions:
            - key: app
              operator: In
              values:
              - test
```
operator运算符枚举值：In、NotIn、Exists、NotExists

需要强制满足以下规则，pod才会调度到该node上；已经运行在该node上的pod则忽略
语法更具表达力，可以进行组合，满足各种使用场景，比如：可以将pod调度到有某个标签，又没有某个标签的节点

##### node软亲和
```yaml
affinity:
  nodeAffinity:
    preferedDuringSchedulingIgnoredDuringExecution:
    - weight:80
      preference:
        - matchExpressions:
            - key: app
              operator: In
              values:
              - test
```
优先选择策略：用户可以配置**权重优先级**

#### pod亲和调度
##### pod硬亲和
```yaml
affinity:
  podAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
    - topologykey: kubernetes.io/hostname
      labelSelector:
        - matchExpressions:
            - key: app
              operator: In
              values:
              - test
```
使用场景，比如：podA调度到node1，此时想要将podB调度到与podA同一个节点，此时可以通过pod亲和性，达到的效果是podA运行在哪个节点，podB也会运行在该节点上
podA -> node1
podA -> PodB

#### pod反亲和调度
和不要满足pod的调度到同一个node
```yaml
affinity:
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
    - topologykey: kubernetes.io/hostname
      labelSelector:
        - matchExpressions:
            - key: app
              operator: In
              values:
              - test
```
podA -> node1
podB -> 除node1外其他节点

# 服务发现
传统方式：端口映射方式
问题痛点：暴露太多端口

## Service 
为一组pod提供统一的对外访问的入口，实现负载均衡
通过在service的spec.Selector下配置标签，将service和一组pod进行绑定

### NodePort方式
type: NodePort, 通常比较大30080

### LoaderBalance方式
动态分配ExternalIp（局域网的虚拟Ip）

## Ingress
技术选型：Ingress Controller Traefik
为一组service提供一个统一的对外访问的endpoint

# 工作负载
管理一组pod
## Deployment
控制pod的副本数replicas

## StatefulSet
replicas 2
1、保证有序性 
启动顺序：pod-0、pod-1
启动顺序：pod-1、pod-0
2、保证有状态性
id不变  
3、稳定的服务发现
pod-0  对指定pod进行访问

## DaemonSet
在每个节点上运行一个pod,作为守护进程

## Job
一次性的任务

## CronJob
周期性的任务






