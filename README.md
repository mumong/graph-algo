# graph-algo
This is a k8s resource scheduler development project based on graph search algorithm

这是一个基于图搜索算法来实现k8s集群内部资源发现的项目，使用的是迪节斯搜索算法和弗洛伊德算法。
将集群内的资源视为图的节点，将资源之间的关系视为图的边，将这个整体看做一个有向无环图，集群内的资源可以是fpga,gpu,dpu等。

在main函数中通过 /getResource 指定路由方式进入主逻辑。

下面是请求的格式，通过传入容器名称+所需资源组合如：container1, gpu:1, fpga:1, cola:4。这样的组合后，项目会通过获取集群中硬件资源的数量与关系构建为关系列表，进而转化为图结构。
`curl -X POST  http://localhost:8080/getResource/master  -H "Content-Type: application/json" -d "[{\"name\": \"container1\", \"resourceQuantity\": {\"nvidia.com/gpu\": 1, \"fpga\": 1, \"myway5.com/cola\": 4}}, {\"name\": \"container2\", \"resourceQuantity\": {\"nvidia.com/gpu\": 1, \"fpga\": 1, \"myway5.com/cola\": 2}}]"`


通过迪杰斯算法得到一个最优组合，主要是计算不同组合的得分计算得分最小的组合，代表了这一组资源获取的路径最小。下图为实际的运行结果示意。
![ad7122589e06d03b865c84a28ff9e6c](https://github.com/user-attachments/assets/551ed513-10ff-4acc-9a93-ab38d8c3d1b7)
