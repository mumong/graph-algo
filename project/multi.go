package project

import (
	"fmt"
	"math"
)

const INF = int(1e4)

// 获取所有 n 选 r 的组合
// [-1-1 2] 1
func combinations(elements []int, r int) [][]int {
	var result [][]int
	var helper func([]int, int, int)

	helper = func(comb []int, start int, r int) {
		if r == 0 {
			result = append(result, append([]int{}, comb...))
			return
		}

		for i := start; i <= len(elements)-r; i++ {
			if elements[i] != -1 {
				helper(append(comb, elements[i]), i+1, r-1)
			}

		}
	}

	helper([]int{}, 0, r)
	return result
}

// 使用 Floyd-Warshall 算法计算每两个节点之间的最短路径
func floydWarshall(edges [][][]int, n int) [][]int {
	// 初始化最短路径矩阵
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, n)
		for j := range dist[i] {
			if i == j {
				dist[i][j] = 0
			} else {
				// 初始化为某个极大的值
				dist[i][j] = INF
			}
		}
	}

	// 初始化为直接的边权
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if len(edges[i][j]) > 0 {
				dist[i][j] = getMinEdgeWeight(edges, i, j)
			}
		}
	}

	fmt.Printf("haha:%v", dist[0])
	// 三重循环更新最短路径
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				//if dist[i][k] < INF && dist[k][j] < INF {
				//	newDist := dist[i][j] + dist[k][j]
				//	if newDist < dist[i][j] && newDist < INF {
				//		dist[i][j] = newDist
				//	}
				//}
				if dist[i][k] < INF && dist[k][j] < INF && dist[i][j] > dist[i][k]+dist[k][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}

	return dist
}

// 获取两节点之间最小权重
func getMinEdgeWeight(edges [][][]int, u int, v int) int {
	if u == v {
		return 0
	}
	minWeight := math.MaxInt32
	for _, weight := range edges[u][v] {
		if weight < minWeight {
			minWeight = weight
		}
	}
	return minWeight
}

// 计算一个组合的最优路径权重
func calculateCombinationWeight(dist [][]int, combination []int) (int, bool) {
	totalWeight := 0
	for i := 0; i < len(combination); i++ {
		for j := i + 1; j < len(combination); j++ {
			weight := dist[combination[i]][combination[j]]
			//if weight == INF {
			if weight >= INF {
				return totalWeight, false // 不可达路径
			}
			totalWeight += weight
		}
	}
	return totalWeight, true
}

//1 资源满足，类型满足，权重可达。
//2 资源满足，类型满足，权重不可达，（选出来）
//3 资源，类型都不满足。不选了

// 动态寻找最优资源组合
func findOptimalCombination(dist [][]int, resources [][]int, resourceCounts []int) (int, []int) {
	minWeight := math.MaxInt32
	bestCombination := []int{}

	var dfs func(int, []int)
	dfs = func(index int, currentCombination []int) {
		if index == len(resources) {
			totalWeight, isValid := calculateCombinationWeight(dist, currentCombination)
			if isValid && totalWeight < minWeight {
				minWeight = totalWeight
				bestCombination = make([]int, len(currentCombination))
				copy(bestCombination, currentCombination)
			}
			return
		}

		for _, comb := range combinations(resources[index], resourceCounts[index]) {
			newCombination := append(currentCombination, comb...)
			dfs(index+1, newCombination)
		}
	}

	dfs(0, []int{})

	return minWeight, bestCombination
}

func CountNode(edges [][][]int, resourceCounts []int, resourcesIndexes [][]int) (int, []int) {
	// 资源分类
	//resources := [][]int{gpus, fpgas, rdma}
	//resourceCounts := []int{2, 1, 1} // 选择 2 个 GPU，1 个 FPGA 和 1 个 RDMA
	// 使用 Floyd-Warshall 计算最短路径
	fmt.Printf("运行算法前")

	dist := floydWarshall(edges, len(edges))
	fmt.Println("Floyd-Warshall 计算的最短路径矩阵:")

	//dist := dijkstra(edges, len(edges))
	//fmt.Println("dijkstra 计算的最短路径矩阵:")
	// 打印 Floyd-Warshall 结果
	for i := range dist {
		fmt.Println(dist[i])
	}
	fmt.Printf("运行算法后")
	// 动态查找最优组合
	minWeight, bestCombination := findOptimalCombination(dist, resourcesIndexes, resourceCounts)

	return minWeight, bestCombination

}

// 使用dijkstra计算两个节点之间的最短路径
func dijkstra(edges [][][]int, n int) [][]int {
	// 初始化最短路径矩阵和访问标记
	dist := make([][]int, n)
	graph := make([][]int, n)
	visited := make([]bool, n) // 一维数组用于标记节点是否被访问

	for i := range dist {
		dist[i] = make([]int, n)
		graph[i] = make([]int, n)
		for j := range dist[i] {
			if i == j {
				dist[i][j] = 0
				graph[i][j] = 0
			} else {
				// 初始化为无穷大，表示初始时不知道节点之间的距离
				dist[i][j] = INF
				graph[i][j] = INF
			}
		}
	}

	// 初始化图的边权
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if len(edges[i][j]) > 0 {
				graph[i][j] = getMinEdgeWeight(edges, i, j)
			}
		}
	}

	// 针对每一个节点，运行Dijkstra算法
	for start := 0; start < n; start++ {
		// 重置访问标记
		for i := 0; i < n; i++ {
			visited[i] = false
		}
		dist[start][start] = 0

		// 计算从节点 start 到其他节点的最短距离
		for k := 0; k < n-1; k++ {
			currentNode := -1
			minDistance := INF

			// 找到距离起始节点最近的未访问节点
			for j := 0; j < n; j++ {
				if !visited[j] && dist[start][j] < minDistance {
					minDistance = dist[start][j]
					currentNode = j
				}
			}

			// 如果没有找到合适的节点，则跳出循环
			if currentNode == -1 {
				break
			}

			// 标记当前节点已访问
			visited[currentNode] = true

			// 更新从 start 节点到其他所有节点的最短距离
			for j := 0; j < n; j++ {
				if graph[currentNode][j] != INF && !visited[j] {
					newDistance := dist[start][currentNode] + graph[currentNode][j]
					if newDistance < dist[start][j] {
						dist[start][j] = newDistance
					}
				}
			}
		}
	}

	return dist
}
