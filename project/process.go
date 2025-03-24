package project

import (
	"context"
	"fmt"
	"graph/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

// const filepath = "D:\\desktop\\mconfig"
// 获取 Kubernetes Client
//func getKubernetesClient() (*kubernetes.Clientset, dynamic.Interface, error) {
//	var config *rest.Config
//	var err error
//
//	// 如果在集群内运行，使用 inClusterConfig
//	if runtime.GOOS == "linux" {
//		config, err = rest.InClusterConfig() // 使用集群内部配置
//		if err != nil {
//			return nil, nil, fmt.Errorf("error creating in-cluster config: %v", err)
//		}
//	} else {
//		// 在本地运行时使用 kubeconfig 文件
//		kubeconfigPath := "D:\\\\desktop\\\\mconfig"
//		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
//		if err != nil {
//			return nil, nil, fmt.Errorf("error building config from kubeconfig: %v", err)
//		}
//	}
//
//	clientset, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		return nil, nil, fmt.Errorf("error creating Kubernetes client: %v", err)
//	}
//
//	dynamicClient, err := dynamic.NewForConfig(config)
//	if err != nil {
//		return nil, nil, fmt.Errorf("error creating dynamic client: %v", err)
//	}
//
//	return clientset, dynamicClient, nil
//}

// 获取 Kubernetes Client
func getKubernetesClient() (*kubernetes.Clientset, dynamic.Interface, error) {
	var config *rest.Config
	var err error

	// 判断是否有 KUBERNETES_SERVICE_HOST 环境变量，以此判断是否在 Kubernetes 集群内部
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		// 集群内使用 inClusterConfig
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, nil, fmt.Errorf("error creating in-cluster config: %v", err)
		}
	} else {
		// 本地开发使用 kubeconfig
		kubeconfigPath := "D:\\\\desktop\\\\mconfig" // 替换为你的 kubeconfig 路径
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, nil, fmt.Errorf("error building config from kubeconfig: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating dynamic client: %v", err)
	}

	return clientset, dynamicClient, nil
}

func InitializeResourceMap(nodeName string) map[string][]string {
	_, dynamicClient, err := getKubernetesClient()
	if err != nil {
		fmt.Printf("Error initializing Kubernetes client: %v\n", err)
		return nil
	}

	// 定义 DeviceTree 的 GVR
	deviceTreeGVR := schema.GroupVersionResource{
		Group:    "ham.xnet.com", // 确保使用了正确的 group
		Version:  "v1",
		Resource: "devicetrees",
	}

	// 从集群中获取 DeviceTree
	deviceTree, err := dynamicClient.Resource(deviceTreeGVR).
		Namespace("xnet").
		Get(context.TODO(), nodeName, metav1.GetOptions{})

	if err != nil {
		fmt.Printf("Error fetching DeviceTree: %v\n", err)
		return nil
	}

	// 定义资源映射表
	resourceMap := map[string][]string{}

	// 从 DeviceTree 中提取 unused 字段
	unused, found, err := unstructured.NestedMap(deviceTree.Object, "spec", "unused")
	if err != nil || !found {
		fmt.Printf("Error getting unused resources: %v\n", err)
		return nil
	}

	// 解析 unused 资源
	for resourceType, devices := range unused {
		deviceList := []string{}
		if deviceSlice, ok := devices.([]interface{}); ok {
			for _, device := range deviceSlice {
				if deviceStr, ok := device.(string); ok {
					deviceList = append(deviceList, deviceStr)
				}
			}
		}
		resourceMap[resourceType] = deviceList
	}

	return resourceMap
}

// 初始化 Graph，获取 DeviceTopology 并解析邻接表
func InitializeGraph(nodeName string) types.Graph {
	_, dynamicClient, err := getKubernetesClient()
	if err != nil {
		fmt.Printf("Error initializing Kubernetes client: %v\n", err)
		return types.Graph{}
	}

	// 定义 DeviceTopology 的 GVR
	deviceTopoGVR := schema.GroupVersionResource{
		Group:    "ham.xnet.com",
		Version:  "v1",
		Resource: "devicetopologies",
	}

	// 从集群中获取 DeviceTopology
	deviceTopo, err := dynamicClient.Resource(deviceTopoGVR).
		Namespace("xnet").
		Get(context.TODO(), nodeName, metav1.GetOptions{})

	if err != nil {
		fmt.Printf("Error fetching DeviceTopology: %v\n", err)
		return types.Graph{}
	}

	// 定义邻接列表
	adjacencyList := map[string][]types.Edge{}

	// 从 DeviceTopology 中提取 adjacencylist 字段
	adjacencyMap, found, err := unstructured.NestedMap(deviceTopo.Object, "spec", "adjacencylist")
	if err != nil || !found {
		fmt.Printf("Error getting adjacencylist: %v\n", err)
		return types.Graph{}
	}

	// 解析 adjacencylist 生成图的邻接表
	for device, edges := range adjacencyMap {
		var edgeList []types.Edge
		if edgeSlice, ok := edges.([]interface{}); ok {
			for _, edge := range edgeSlice {
				if edgeMap, ok := edge.(map[string]interface{}); ok {
					// 将 weight 解析为 int64，然后转换为 int
					weight, _ := edgeMap["weight"].(int64) // 解析为 int64
					edgeList = append(edgeList, types.Edge{
						Topology: edgeMap["topology"].(string),
						Node:     edgeMap["node"].(string),
						Weight:   int(weight), // 将 int64 转换为 int
					})
				}
			}
		}
		adjacencyList[device] = edgeList
	}

	return types.Graph{AdjacencyList: adjacencyList}
}
