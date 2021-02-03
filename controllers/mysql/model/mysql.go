package model

import (
	cachev1alpha1 "github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/constants"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetMysqlDeployment returns a mysql Deployment object
func GetMysqlDeployment(m *cachev1alpha1.DBMMOMySQL) *appsv1.Deployment {
	ls := GetLabels(m.Name)
	replicas := m.Spec.Size

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.MysqlDeploymentName,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: constants.MysqlStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: constants.MysqlClaimName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: constants.MysqlClaimName,
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Name:  constants.MysqlContainerName,
						Image: constants.MysqlContainerImage,
						Ports: []corev1.ContainerPort{{
							ContainerPort: constants.MysqlContainerPort,
							Name:          constants.MysqlContainerPortName,
						}},
						EnvFrom: nil,
						Env: []corev1.EnvVar{
							{
								Name:  constants.MysqlSecretEnvName,
								Value: constants.MysqlSecretEnvVal,
								// TODO add Secret here
								//ValueFrom: nil,
							},
						},
						ImagePullPolicy: "IfNotPresent",
					},
					}},
			},
		},
	}
	return dep
}

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
					"storage": resource.MustParse(constants.MysqlCapacityStorageTest),
				},
			},
			//TODO remove this after local development
			//StorageClassName: className,
		},
		Status: corev1.PersistentVolumeClaimStatus{},
	}
	// Set Mysql instance as the owner and controller
	return pvc
}

// GetMysqlService returns the mysql service
func GetMysqlService(m *cachev1alpha1.DBMMOMySQL) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.MysqlServiceName,
			Namespace: m.Namespace,
			Labels:    GetLabels(m.Name),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: constants.MysqlContainerPort,
				},
			},
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
