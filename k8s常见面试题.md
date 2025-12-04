# k8s架构设计

1）API server: 负责整个k8s集群的认证和授权，访问k8s集群资源的入口，访问方式：kubect/client-go
2）etcd数据库：cr（crd表结构，cr数据库中表记录）
3）controller-manager：node的controller、pod的副本数管理对应的controller、service的controller
  读取etcd中cr的信息，进行调谐，将实际信息写入etcd,努力将期望状态变成实际状态
4）scheduler：读取etcd中信息，按照一定的调度算法将pod调度到指定节点

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




