package mysql

import (
	"context"
	"errors"
	"github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/model"
	appsv1 "k8s.io/api/apps/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

const (
	//ConditionStatusSuccess is the expected return keyword for a successful resources
	ConditionStatusSuccess = "True"
)

func (r *DBMMOMySQLReconciler) getDeploymentReadiness(ctx context.Context, mysql *v1alpha1.DBMMOMySQL) (bool, error) {
	dep := model.GetMysqlDeployment(mysql)
	if err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: dep.Namespace,
		Name:      dep.Name,
	}, dep); err != nil {
		return false, err
	}
	// A deployment has an array of conditions
	for _, condition := range dep.Status.Conditions {
		// One failure condition exists, if this exists, return the Reason
		if condition.Type == appsv1.DeploymentReplicaFailure {
			return false, errors.New(condition.Reason)
			// A successful deployment will have the progressing condition type as true
		} else if condition.Type == appsv1.DeploymentProgressing && condition.Status != ConditionStatusSuccess {
			return false, nil
		}
	}

	return dep.Status.ReadyReplicas == dep.Status.Replicas, nil
}

// TODO this one might need a bit of rethinking
func (r *DBMMOMySQLReconciler) getPVReadiness(ctx context.Context, mysql *v1alpha1.DBMMOMySQL) (bool, error) {
	pvc := model.GetMysqlPvc(mysql)
	if err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: pvc.Namespace,
		Name:      pvc.Name,
	}, pvc); err != nil {
		return false, err
	}

	return true, nil
}

func (r *DBMMOMySQLReconciler) getIngressReadiness(ctx context.Context, mysql *v1alpha1.DBMMOMySQL) (bool, error) {
	ingress := model.GetMysqlIngress(mysql)
	if err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: ingress.Namespace,
		Name:      ingress.Name,
	}, ingress); err != nil {
		return false, err
	}
	if ingress == nil {
		return false, nil
	}

	return len(ingress.Status.LoadBalancer.Ingress) > 0, nil

}

func (r *DBMMOMySQLReconciler) getCollectiveReadiness(ctx context.Context, mysql *v1alpha1.DBMMOMySQL) (bool, error) {

	ready, err := r.getPVReadiness(ctx, mysql)
	if ready != true || k8serr.IsNotFound(err) {
		r.Log.Info("NOT READY", "resource", "persistentVolume")
		return false, nil
	}
	if err != nil {
		return false, err
	}
	r.Log.Info("READY", "resource", "persistentVolume")

	ready, err = r.getDeploymentReadiness(ctx, mysql)
	if ready != true || k8serr.IsNotFound(err) {
		r.Log.Info("NOT READY", "resource", "deployment")
		return false, nil
	}
	if err != nil {
		return false, err
	}
	r.Log.Info("READY", "resource", "deployment")

	if mysql.Spec.Deployment.Ingress != nil && mysql.Spec.Deployment.Ingress.Enabled != nil && *mysql.Spec.Deployment.Ingress.Enabled != false {
		ready, err = r.getIngressReadiness(ctx, mysql)

		if ready != true || k8serr.IsNotFound(err) {
			r.Log.Info("NOT READY", "resource", "ingress")
			return false, nil
		}
		if err != nil {
			return false, err
		}
		r.Log.Info("READY", "resource", "ingress")

	}

	return true, nil
}
