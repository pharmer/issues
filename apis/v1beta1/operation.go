package v1beta1

type OperationState int
const (
	OperationDone OperationState = iota // 0
	OperationPending
	OperationProgress
)


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Operation struct {
	ID int64
	UserID string
	ClusterID string
	Code string
	State OperationState
}