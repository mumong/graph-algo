package types

// edge represents an edge in the graph
type Edge struct {
	Topology string
	Node     string
	Weight   int
}

// Graph represents the graph using an adjacency list
type Graph struct {
	AdjacencyList map[string][]Edge
}

// 请求参数 []ContainerQuantityEntry
type ContainerQuantityEntry struct {
	Name             string           `json:"name"`             // container name
	ResourceQuantity map[string]int64 `json:"resourceQuantity"` // key: resource name, value: resource quantity
}

// 返回值
// container allocated resource informations
type ContainerEntry struct {
	Name      string              `json:"name"`      // container name
	Resources map[string][]string `json:"resources"` // key: resource name  value: []string{deviceId,}
	Weight    float64             `json:"weight"`    // container resources weight
}

// The list of resources allocated by the node to the Pod
type NodeResourceData struct {
	Containers []ContainerEntry `json:"containers"` // All container allocated resource information of pod
	WeightSum  float64          `json:"weightSum"`  // The weight sum of all container resources
}
