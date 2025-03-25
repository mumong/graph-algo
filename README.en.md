# README.md
- [English](README.en.md)
- [Chinese](README.md)
- [Japanese](README.jp.md)

# 🚀 K8s Resource Scheduler Based on Graph Search Algorithm  

## 📌 Project Overview  
This project is a **Kubernetes resource scheduling system** based on **graph search algorithms**, utilizing **Dijkstra’s algorithm** and **Floyd’s algorithm** for **cluster resource discovery**.  
The project treats **resources within the cluster as nodes in a graph** and **relationships between resources as edges**, forming a **Directed Acyclic Graph (DAG)** to achieve optimal resource allocation.  

### **Supported Resource Types**
- **FPGA**
- **GPU**
- **DPU**
- **Custom Resources (e.g., Cola)**

---

## 🎯 Project Core Logic
1. **Resource Request** → Submit container and resource requirements via `curl`  
2. **Resource Graph Construction** → Parse Kubernetes resources to build a **graph structure**  
3. **Path Calculation** → Use **Dijkstra’s algorithm** to compute the **optimal resource path**  
4. **Resource Allocation** → Allocate the calculated **optimal combination**  

---

## 🛠 API Interface
### **Get Optimal Resource Allocation**
**Request Method**: `POST`  
**Request Path**: `/getResource/master`  

**Request Example**
```bash
curl -X POST http://localhost:8080/getResource/master \
     -H "Content-Type: application/json" \
     -d '[{"name": "container1", "resourceQuantity": {"nvidia.com/gpu": 1, "fpga": 1, "myway5.com/cola": 4}}, 
          {"name": "container2", "resourceQuantity": {"nvidia.com/gpu": 1, "fpga": 1, "myway5.com/cola": 2}}]'
```

The optimal combination is obtained through the Dijkstra algorithm, which mainly calculates the scores of different combinations and the combination with the smallest score represents the smallest path for obtaining this group of resources. The following figure shows the actual operation results.

![ad7122589e06d03b865c84a28ff9e6c](https://github.com/user-attachments/assets/551ed513-10ff-4acc-9a93-ab38d8c3d1b7)

---
## 🛠️ Build & Deploy
1. Run `make docker` to compile and build the Docker image.
2. Run `bash ./tools/upload.sh` to upload the Docker image to the local containerd repository.
3. Run `kubectl apply -f deployment.yaml` to deploy the service using the application file.


