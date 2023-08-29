package generator

type WorkloadType string

const (
	WT_Deployment  WorkloadType = "deployment"
	WT_DaemonSet   WorkloadType = "daemonset"
	WT_StatefulSet WorkloadType = "statefulset"
	WT_ReplicaSet  WorkloadType = "replicaset"
	WT_Job         WorkloadType = "job"
	WT_CronJob     WorkloadType = "cronjob"
)

type Base struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels"`
}

type App struct {
	Base             `json:"-"`
	Common           `yaml:",inline"`
	Version          string                `yaml:"version"`
	DefaultNamespace string                `yaml:"namespace,omitempty"`
	Components       map[string]*Component `yaml:"components,omitempty"`
}

type Component struct {
	Base             `json:"-"`
	Common           `yaml:",inline"`
	DefaultNamespace string       `yaml:"namespace,omitempty"`
	Type             WorkloadType `yaml:"type"`

	Deployment  *CompDeployment  `yaml:"deployment,omitempty"`
	DaemonSet   *CompDaemonSet   `yaml:"daemonSet,omitempty"`
	StatefulSet *CompStatefulSet `yaml:"statefulSet,omitempty"`
	ReplicaSet  *CompReplicaSet  `yaml:"replicaSet,omitempty"`
	Job         *CompJob         `yaml:"job,omitempty"`
	CronJob     *CompCronJob     `yaml:"cronJob,omitempty"`
}

type Common struct {
	ConfigMaps []*ConfigMap
	Secrets    []*Secret
	Services   []*Service

	ServiceAccount     *ServiceAccount
	ClusterRole        *ClusterRole
	ClusterRoleBinding *ClusterRoleBinding
	Role               *Role
	RoleBinding        *RoleBinding
	Ingress            *Ingress
}

type CompDeployment struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type CompDaemonSet struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type CompStatefulSet struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type CompReplicaSet struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type CompJob struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type CompCronJob struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type ConfigMap struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type Secret struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type Service struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type ServiceAccount struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type ClusterRole struct {
	Base `json:"-"`
}

type ClusterRoleBinding struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type Role struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type RoleBinding struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

type Ingress struct {
	Base             `json:"-"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}
