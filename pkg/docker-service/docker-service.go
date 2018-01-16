package ds

const (
	// BackendZookeeper zookeeper
	BackendZookeeper = "zookeeper"
	// BackendEtcd etcd
	BackendEtcd = "etcd"
)

// Service broker service define
type Service struct {
	Name string `json:"Name"`
	ID   string `json:"ID"`
	IP   string `json:"IP"`
	Port uint32 `json:"Port"`
}
