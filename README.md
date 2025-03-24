# 🚀 K8s Resource Scheduler Based on Graph Search Algorithm  

## 📌 项目简介  
本项目是一个 **基于图搜索算法** 的 **Kubernetes 资源调度系统**，使用 **Dijkstra 算法** 和 **Floyd 算法** 进行 **集群资源发现**。  
项目将 **集群内的资源视为图的节点**，**资源之间的关系视为图的边**，形成一个 **有向无环图（DAG）** 进行最优资源分配。  

### **支持的资源类型**
- **FPGA**
- **GPU**
- **DPU**
- **自定义资源（如 Cola）**

---

## 🎯 项目核心逻辑
1. **资源请求** → 通过 `curl` 提交容器及资源需求  
2. **资源图构建** → 解析 Kubernetes 资源，构建 **图结构**  
3. **路径计算** → 使用 **Dijkstra** 算法计算 **最优资源路径**  
4. **资源分配** → 分配计算出的 **最优组合**  

---

## 🛠 API 接口
### **获取最优资源分配**
**请求方式**：`POST`  
**请求路径**：`/getResource/master`  

**请求示例**
```bash
curl -X POST http://localhost:8080/getResource/master \
     -H "Content-Type: application/json" \
     -d '[{"name": "container1", "resourceQuantity": {"nvidia.com/gpu": 1, "fpga": 1, "myway5.com/cola": 4}}, 
          {"name": "container2", "resourceQuantity": {"nvidia.com/gpu": 1, "fpga": 1, "myway5.com/cola": 2}}]'
```

通过迪杰斯算法得到一个最优组合，主要是计算不同组合的得分计算得分最小的组合，代表了这一组资源获取的路径最小。下图为实际的运行结果示意。
![ad7122589e06d03b865c84a28ff9e6c](https://github.com/user-attachments/assets/551ed513-10ff-4acc-9a93-ab38d8c3d1b7)

---
## 构建方式
1. 执行`make docker` 编译构建镜像。
2. 执行`bash ./tools/upload.sh` 将docker镜像上传到本地containerd仓库中
3. 执行`kubectl apply -f deployment.yaml` 应用文件部署服务，需要配置对应的rbac确保硬件资源的发现。
