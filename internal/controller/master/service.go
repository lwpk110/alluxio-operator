package master

import (
	stackv1alpha1 "github.com/zncdata-labs/alluxio-operator/api/v1alpha1"
	"github.com/zncdata-labs/alluxio-operator/internal/common"
	"github.com/zncdata-labs/alluxio-operator/internal/controller/role"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceReconciler struct {
	common.GeneralResourceStyleReconciler[*stackv1alpha1.Alluxio, *stackv1alpha1.MasterRoleGroupSpec]
}

// NewService New a ServiceReconciler
func NewService(
	scheme *runtime.Scheme,
	instance *stackv1alpha1.Alluxio,
	client client.Client,
	mergedLabels map[string]string,
	mergedCfg *stackv1alpha1.MasterRoleGroupSpec,
) *ServiceReconciler {
	return &ServiceReconciler{
		GeneralResourceStyleReconciler: *common.NewGeneraResourceStyleReconciler[*stackv1alpha1.Alluxio,
			*stackv1alpha1.MasterRoleGroupSpec](
			scheme,
			instance,
			client,
			mergedLabels,
			mergedCfg,
		),
	}
}

func (s *ServiceReconciler) Build(data common.ResourceBuilderData) (client.Object, error) {
	instance := s.Instance
	mergedGroupCfg := s.MergedCfg
	roleGroupName := data.GroupName
	roleName := role.Master
	masterPorts := getMasterPorts(mergedGroupCfg)
	jobMasterPorts := getJobMasterPorts(mergedGroupCfg)
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createMasterSvcName(instance.Name, string(roleName), roleGroupName, "0"),
			Namespace: instance.Namespace,
			Labels:    s.MergedLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "rpc",
					Port: masterPorts.Rpc,
				},
				{
					Name: "web",
					Port: masterPorts.Web,
				},
				{
					Name: "embedded",
					Port: masterPorts.Embedded,
				},
				{
					Name: "job-rpc",
					Port: jobMasterPorts.Rpc,
				},
				{
					Name: "job-web",
					Port: jobMasterPorts.Web,
				},
				{
					Name: "job-embedded",
					Port: jobMasterPorts.Embedded,
				},
			},
			Selector:  s.MergedLabels,
			ClusterIP: "None",
		},
	}
	return svc, nil
}
