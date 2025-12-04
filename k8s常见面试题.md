# k8s架构设计

1）API server: 负责整个k8s集群的认证和授权，访问k8s集群资源的入口，访问方式：kubect/client-go
2）etcd数据库：cr（crd表结构，cr数据库中表记录）
3）controller-manager：node的controller、pod的副本数管理对应的controller、service的controller
  读取etcd中cr的信息，进行调谐，将实际信息写入etcd,努力将期望状态变成实际状态
4）scheduler：读取etcd中信息，按照一定的调度算法将pod调度到指定节点

# k8s调度原理
如何将pod调度到指定node上
## 原理
1、informer：通过API server监听k8s集群资源的变化
  1）list：全量、周期查询
  2）watch：http长连接，监听k8s资源（pod、nodedeng）的变化
2、本地会维护scheduling pod queue（待调度的pod的队列）、scheduling node queue（待调度的node的队列）
3、从待调度的pod队列中取出一个pod，先对node进行过滤（predict算法，比如端口占用、cpu、memory、污点、节点亲和性等），然后对node进行打分（priority算法，比如:cpu、memory剩余量等），评选出最优的node，将node与pod进行绑定
  LOOP
## 实现方式
### nodeSelector方式调度
匹配节点的标签
```bash
kubectl label node1 app=test
```
### taints和tolerations方式调度
节点上有污点，需要pod上
