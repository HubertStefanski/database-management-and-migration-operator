package model

import (
	"github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/constants"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetMysqlIngress returns the default ingress configuration for a mysql isntance
func GetMysqlIngress(m *v1alpha1.DBMMOMySQL) *netv1.Ingress {
	ingr := &netv1.Ingress{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.MysqlIngressName,
			Namespace: m.Namespace,
		},
	}
	return ingr
}
