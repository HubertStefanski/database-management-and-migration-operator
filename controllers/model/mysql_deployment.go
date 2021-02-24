package model

import (
	v1alpha1 "github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/constants"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetMysqlDeployment returns a mysql Deployment object
func GetMysqlDeployment(m *v1alpha1.DBMMOMySQL) *appsv1.Deployment {
	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.MysqlDeploymentName,
			Namespace: m.Namespace,
		},
	}
	return dep
}

//GetMysqlInitCommand Initialises the container with the command specified in the deployment spec
func GetMysqlInitCommand(m *v1alpha1.DBMMOMySQL) string {
	if m.Spec.Deployment.TableInitCMD != nil && *m.Spec.Deployment.TableInitCMD != "" {
		return *m.Spec.Deployment.TableInitCMD
	}
	return ""
}

// MysqlDeploymentGetEnvFrom returns the environment variables contained within a secret
func MysqlDeploymentGetEnvFrom(m *v1alpha1.DBMMOMySQL) []v1.EnvFromSource {
	var envFrom []v1.EnvFromSource
	if m.Spec.Deployment != nil && m.Spec.Deployment.EnvFrom != nil {
		for _, v := range m.Spec.Deployment.EnvFrom {
			envFrom = append(envFrom, *v.DeepCopy())
		}
	}
	return envFrom
}
