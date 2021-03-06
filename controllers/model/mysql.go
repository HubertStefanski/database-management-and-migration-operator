package model

import (
	cachev1alpha1 "github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/constants"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetMysqlPvc returns the mysql Persistent volume claim for mysql
func GetMysqlPvc(m *cachev1alpha1.DBMMOMySQL) *corev1.PersistentVolumeClaim {
	//var className = new(string)
	//*className = constants.MysqlStorageClassName

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.MysqlClaimName,
			Namespace: m.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				constants.MysqlPVAccessModes,
			},
			Resources: corev1.ResourceRequirements{
				Limits: nil,
				Requests: corev1.ResourceList{
					"storage": resource.MustParse(getMysqlStorageSize(m)),
				},
			},
			//TODO remove this after local development
			//StorageClassName: className,
		},
	}

	// Set Mysql instance as the owner and controller
	return pvc
}

func getMysqlStorageSize(m *cachev1alpha1.DBMMOMySQL) string {
	if m.Spec.Deployment.StorageCapacity != nil {
		return *m.Spec.Deployment.StorageCapacity
	}
	return constants.MysqlCapacityStorageTest // TODO change me after testing

}

// GetMysqlService returns the mysql service
func GetMysqlService(m *cachev1alpha1.DBMMOMySQL) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.MysqlServiceName,
			Namespace: m.Namespace,
			Labels:    GetLabels(m.Name),
		},
	}
	// Set Mysql instance as the owner and controller
	return service

}

// GetLabels returns the labels for selecting the resources
// belonging to the given mysql CR name.
func GetLabels(name string) map[string]string {
	return map[string]string{"app": constants.MysqlAppSelector, "mysql_cr": name}
}

// GetPvcNames returns the pv names of mysql
func GetPvcNames(pvcs []corev1.PersistentVolumeClaim) []string {
	var persistentVolumesClaimNames []string
	for _, pvc := range pvcs {
		persistentVolumesClaimNames = append(persistentVolumesClaimNames, pvc.Name)
	}
	return persistentVolumesClaimNames
}

// GetPvNames returns the pv names of mysql
func GetPvNames(pvs []corev1.PersistentVolume) []string {
	var persistentVolumesNames []string
	for _, pv := range pvs {
		persistentVolumesNames = append(persistentVolumesNames, pv.Name)
	}
	return persistentVolumesNames
}

// GetPodNames returns the pod names of the array of pods passed in
func GetPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// GetServiceNames returns the pod names of the array of pods passed in
func GetServiceNames(services []corev1.Service) []string {
	var serviceNames []string
	for _, service := range services {
		serviceNames = append(serviceNames, service.Name)
	}
	return serviceNames
}
