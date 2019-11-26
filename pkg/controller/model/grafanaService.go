package model

import (
	"github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

func getServiceLabels(cr *v1alpha1.Grafana) map[string]string {
	if cr.Spec.Service == nil {
		return nil
	}
	return cr.Spec.Service.Labels
}

func getServiceAnnotations(cr *v1alpha1.Grafana) map[string]string {
	if cr.Spec.Service == nil {
		return nil
	}
	return cr.Spec.Service.Annotations
}

func getServiceType(cr *v1alpha1.Grafana) v1.ServiceType {
	if cr.Spec.Service == nil {
		return v1.ServiceTypeClusterIP
	}
	if cr.Spec.Service.Type == "" {
		return v1.ServiceTypeClusterIP
	}
	return cr.Spec.Service.Type
}

func GetGrafanaPort(cr *v1alpha1.Grafana) int {
	if cr.Spec.Config.Server.HttpPort == "" {
		return GrafanaHttpPort
	}

	port, err := strconv.Atoi(cr.Spec.Config.Server.HttpPort)
	if err != nil {
		return GrafanaHttpPort
	}

	return port
}

func getServicePorts(cr *v1alpha1.Grafana) []v1.ServicePort {
	return []v1.ServicePort{
		{
			Name:       "grafana",
			Protocol:   "TCP",
			Port:       int32(GetGrafanaPort(cr)),
			TargetPort: intstr.FromString("grafana-http"),
		},
	}
}

func GrafanaService(cr *v1alpha1.Grafana) *v1.Service {
	return &v1.Service{
		ObjectMeta: v12.ObjectMeta{
			Name:        GrafanaServiceName,
			Namespace:   cr.Namespace,
			Labels:      getServiceLabels(cr),
			Annotations: getServiceAnnotations(cr),
		},
		Spec: v1.ServiceSpec{
			Ports: getServicePorts(cr),
			Selector: map[string]string{
				"app": GrafanaPodLabel,
			},
			ClusterIP: "",
			Type:      getServiceType(cr),
		},
	}
}

func GrafanaServiceReconciled(cr *v1alpha1.Grafana, currentState *v1.Service) *v1.Service {
	reconciled := currentState.DeepCopy()
	reconciled.Labels = getServiceLabels(cr)
	reconciled.Annotations = getServiceAnnotations(cr)
	reconciled.Spec.Ports = getServicePorts(cr)
	reconciled.Spec.Type = getServiceType(cr)
	return reconciled
}

func GrafanaServiceSelector(cr *v1alpha1.Grafana) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      GrafanaServiceName,
	}
}
