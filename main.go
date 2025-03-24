package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"graph/project"
	"graph/types"
	"log"
	"net/http"
)

func GetAdjacentMatrix(devices []string, graph map[string][]types.Edge) [][][]int {
	n := len(devices)
	edges := make([][][]int, n)
	for index, dev := range devices {
		edges[index] = make([][]int, n)
		for i, de := range devices {
			edges[index][i] = []int{}
			if index == i {
				edges[index][i] = []int{0}
			} else {
				edge, ok := graph[dev]
				if ok {
					found := false
					for _, e := range edge {
						if e.Node == de {
							found = true
							//edges[index][i] = []int{e.Weight}
							edges[index][i] = append(edges[index][i], e.Weight)
						}
					}
					if !found {
						//edges[index][i] = []int{1000000}

						edges[index][i] = []int{1e3}
						//edges[index][i] = []int{project.INF}
					}
				} else {
					//edges[index][i] = []int{1000000}

					edges[index][i] = []int{1e3}
					//edges[index][i] = []int{project.INF}
				}

			}
		}
	}
	return edges
}

func main() {
	ws := new(restful.WebService)

	ws.Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	//添加一个路由

	ws.Route(ws.POST("/getResource/{node}").To(getResource))
	//将路由添加到go-restful
	restful.Add(ws)
	fmt.Println("开始监听8080")
	//启动HTTP服务器
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("开始11111111")

}

func getResource(req *restful.Request, resp *restful.Response) {

	//获取节点名称
	nodeName := req.PathParameter("node")
	fmt.Printf("接收到的节点名称: %s\n", nodeName)

	//初始化接收数据结构体， 用于获取输入信息,&containerEntries这一步已经将输入数据输给了结构体
	var containerEntries []types.ContainerQuantityEntry
	if err := req.ReadEntity(&containerEntries); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("无法解析请求体: %v", err))
		return
	}
	fmt.Println("接收到的结构体是:", &containerEntries, "\n")

	resourceMap := project.InitializeResourceMap(nodeName)
	graph := project.InitializeGraph(nodeName)

	//结果列表，存储每个容器的结果
	var results []types.ContainerEntry
	//存储所有设备的数组
	all := []string{}
	for _, devs := range resourceMap {
		all = append(all, devs...)
	}
	//进入循环前存储哪些设备是被使用的
	isValid := make([]bool, len(all))

	for _, container := range containerEntries {
		//标记container
		fmt.Println("现在处理的是:", container, "\n")

		//记录索引的map，设备和脚标的关系
		resourceCountMap := []map[string]int{}

		// 构建 resourceCountMap
		containerResourceMap := make(map[string]int)
		for resourceName, quantity := range container.ResourceQuantity {
			containerResourceMap[resourceName] = int(quantity)
		}

		resourceCountMap = append(resourceCountMap, containerResourceMap)

		fmt.Println("resourceCountMap是:", resourceCountMap, "\n")
		fmt.Println("resourceMap是:", resourceMap, "\n")

		//逻辑处理，将链接表生成为链接矩阵并且调用multi.go文件的算法生成返回的结果
		result := processResourceAllocations(isValid, resourceCountMap, resourceMap, graph, container.Name)

		// 将结果写入响应
		if result != nil {

			results = append(results, *result)

		} else {
			resp.WriteErrorString(http.StatusInternalServerError, "无法生成资源分配结果")
		}

	}
	resp.WriteHeaderAndEntity(http.StatusOK, results)
}

func processResourceAllocations(isValid []bool, resourceCountMap []map[string]int, resourceMap map[string][]string, graph types.Graph, name string) *types.ContainerEntry {
	all, resourceIndexes, resources := prepareResourceData(resourceMap)
	//deviceCount := len(all)
	//isValid := make([]bool, deviceCount)

	for i, resourceCount := range resourceCountMap {
		fmt.Printf("处理资源配置 #%d\n", i+1)
		fmt.Println("resourceCount是:", resourceCount, "\n")
		//if !hasValidResourceConnections(resourceCount, graph) {
		//	fmt.Printf("警告：资源配置 #%d 无效，存在未关联的资源。\n", i+1)
		//	continue
		//}

		updateValidResources(isValid, all, resourceIndexes)

		edges := GetAdjacentMatrix(all, graph.AdjacencyList)

		resourceCounts := getResourceCounts(resources, resourceCount)
		fmt.Println("resourceCounts是:", resourceCounts, "\n")

		fmt.Println("资源索引:", resourceIndexes)
		fmt.Println("资源类型:", resources)
		fmt.Println("资源数量:", resourceCounts)

		minWeight, bestCombination := project.CountNode(edges, resourceCounts, resourceIndexes)

		fmt.Println("最小权重:", minWeight)
		fmt.Println("最佳组合:", bestCombination)

		updateValidCombination(isValid, bestCombination)

		fmt.Println("有效资源:", isValid)
		fmt.Println("全部数组", all, "\n")
		fmt.Println("resources是", resources, "\n")
		result := getResults(name, all, isValid, resources, minWeight, resourceMap)
		fmt.Println()

		return result

	}

	return nil
}

func prepareResourceData(resourceMap map[string][]string) ([]string, [][]int, []string) {
	all := []string{}
	resourceIndexes := [][]int{}
	resources := []string{}
	index := 0

	for resource, devs := range resourceMap {
		resources = append(resources, resource)
		fmt.Printf("resources: %v,\n", resource)
		all = append(all, devs...)

		resourceIndex := make([]int, len(devs))
		for i := range devs {
			resourceIndex[i] = index
			index++
		}
		resourceIndexes = append(resourceIndexes, resourceIndex)
	}

	return all, resourceIndexes, resources
}

func updateValidResources(isValid []bool, all []string, resourceIndexes [][]int) {
	for index, val := range isValid {
		if val {
			all[index] = "nil"
			for i := range resourceIndexes {
				for k, value := range resourceIndexes[i] {
					if index == value {
						resourceIndexes[i][k] = -1
					}
				}
			}
		}
	}
}

func getResourceCounts(resources []string, resourceCount map[string]int) []int {
	counts := make([]int, len(resources))
	for i, resource := range resources {
		counts[i] = resourceCount[resource]
	}
	return counts
}

func updateValidCombination(isValid []bool, bestCombination []int) {
	for _, v := range bestCombination {
		isValid[v] = true
	}
}

func hasValidResourceConnections(resourceCount map[string]int, graph types.Graph) bool {
	requestedResources := []string{}
	for resource, count := range resourceCount {
		if count > 0 {
			requestedResources = append(requestedResources, resource)
		}
	}

	if len(requestedResources) <= 1 {
		return true
	}

	for i := 0; i < len(requestedResources); i++ {
		for j := i + 1; j < len(requestedResources); j++ {
			if !hasPath(graph, requestedResources[i], requestedResources[j]) {
				fmt.Printf("调试: %s 和 %s 之间没有找到有效路径\n", requestedResources[i], requestedResources[j])
				return false
			}
		}
	}

	return true
}

func hasPath(graph types.Graph, resource1, resource2 string) bool {
	visited := make(map[string]bool)
	return dfs(graph, resource1+"0", resource2+"0", visited)
}

func dfs(graph types.Graph, start, target string, visited map[string]bool) bool {
	if start == target {
		return true
	}

	visited[start] = true

	for _, edge := range graph.AdjacencyList[start] {
		if !visited[edge.Node] {
			if dfs(graph, edge.Node, target, visited) {
				return true
			}
		}
	}

	return false
}

func getResults(name string, all []string, isValid []bool, resources []string, minWeight int, resourceMap map[string][]string) *types.ContainerEntry {
	res := []string{}
	var weight float64 = float64(minWeight)

	for i, v := range isValid {
		if v && all[i] != "nil" {
			res = append(res, all[i])
		}
	}
	fmt.Println("this is res", res, "\n")
	fmt.Println("this is resources", resources, "\n")

	r := new(types.ContainerEntry)
	r.Name = name
	r.Resources = make(map[string][]string)
	r.Weight = weight

	//for deviceType, deviceNums := range resourceMap {
	//	for _, resource := range resources {
	//		if resource == deviceType {
	//			for _, totalDevice := range deviceNums {
	//				for _, selectDevice := range res {
	//					if totalDevice == selectDevice {
	//						r.Resources[deviceType] = append(r.Resources[deviceType], selectDevice)
	//					}
	//				}
	//			}
	//		}
	//	}
	//}
	// 假设 res 和 resources 都是 []string 类型
	resourcesMap := make(map[string]struct{}, len(resources))
	resMap := make(map[string]struct{}, len(res))

	// 将 resources 和 res 转换为 map，方便 O(1) 查找
	for _, resource := range resources {
		resourcesMap[resource] = struct{}{}
	}
	for _, selectDevice := range res {
		resMap[selectDevice] = struct{}{}
	}

	// 遍历 resourceMap
	for deviceType, deviceNums := range resourceMap {
		// 仅处理 deviceType 在 resources 中的情况
		if _, found := resourcesMap[deviceType]; found {
			for _, totalDevice := range deviceNums {
				// 仅在 totalDevice 存在于 res 中时进行匹配
				if _, selected := resMap[totalDevice]; selected {
					r.Resources[deviceType] = append(r.Resources[deviceType], totalDevice)
				}
			}
		}
	}
	fmt.Printf("containerEntry: %v \n", r)
	fmt.Printf("===================================================================" + "\n")
	return r
}

