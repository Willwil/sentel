package ds

const (
	// BackendZookeeper zookeeper
	BackendZookeeper = "zookeeper"
	// BackendEtcd etcd
	BackendEtcd = "etcd"
)

// Service broker service define
type Service struct {
	ServiceName string `json:"ServiceName"`
	ServiceID   string `json:"ServiceId"`
	IP          string `json:"IP"`
	Port        uint32 `json:"Port"`
}
