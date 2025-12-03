# k8s架构设计

1）API server:k8s集群的入口，kubect/client-go
2）etcd数据库：cr（crd表结构，cr数据库中表记录）
3）controller-manager：node的controller、pod的副本数管理对应的controller、service的controller
  读取etcd中cr的信息，进行调谐，将实际信息写入etcd,努力将期望状态变成实际状态
4）scheduler：读取etcd中信息，

# k8s调度原理
如何将pod调度到指定node上
