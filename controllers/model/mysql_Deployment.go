package model

import (
	"github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

// MysqlDeploymentGetEnvFrom returns the environment variables contained within a secret
func MysqlDeploymentGetEnvFrom(m *v1alpha1.DBMMOMySQL) []v1.EnvFromSource {
	var envFrom []v1.EnvFromSource
	if m.Spec.DBMMOMYSQLDeployment != nil && m.Spec.DBMMOMYSQLDeployment.EnvFrom != nil {
		for _, v := range m.Spec.DBMMOMYSQLDeployment.EnvFrom {
			envFrom = append(envFrom, *v.DeepCopy())
		}
	}
	return envFrom
}
