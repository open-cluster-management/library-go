package client

import (
	"context"
	"fmt"

	"github.com/open-cluster-management/library-go/pkg/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

//NewDefaultClient returns a client.Client for the current-context in kubeconfig
//kubeconfig: The path of the kubeconfig, see (../config/config.go#LoadConfig) for more information
func NewDefaultClient(kubeconfig string, options client.Options) (client.Client, error) {
	return NewClient("", kubeconfig, "", options)
}

//url: The url of the server
//kubeconfig: The path of the kubeconfig, see (../config/config.go#LoadConfig) for more information
//context: The context to connect to
func NewClient(url, kubeconfig, context string, options client.Options) (client.Client, error) {
	klog.V(5).Infof("Create kubeclient for url %s using kubeconfig path %s\n", url, kubeconfig)
	config, err := config.LoadConfig(url, kubeconfig, context)
	if err != nil {
		return nil, err
	}

	client, err := client.New(config, options)
	if err != nil {
		return nil, err
	}

	return client, nil
}

//NewDefaultKubeClient returns a kubernetes.Interface for the current-context in kubeconfig
//kubeconfig: The path of the kubeconfig, see (../config/config.go#LoadConfig) for more information
func NewDefaultKubeClient(kubeconfig string) (kubernetes.Interface, error) {
	return NewKubeClient("", kubeconfig, "")
}

//NewKubeClient returns a kubernetes.Interface based on the provided url, kubeconfig and context
//url: The url of the server
//kubeconfig: The path of the kubeconfig, see (../config/config.go#LoadConfig) for more information
//context: The context to connect to
func NewKubeClient(url, kubeconfig, context string) (kubernetes.Interface, error) {
	klog.V(5).Infof("Create kubeclient for url %s using kubeconfig path %s\n", url, kubeconfig)
	config, err := config.LoadConfig(url, kubeconfig, context)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

//NewDefaultKubeClientDynamic returns a dynamic.Interface for the current-context in kubeconfig
//kubeconfig: The path of the kubeconfig, see (../config/config.go#LoadConfig) for more information
func NewDefaultKubeClientDynamic(kubeconfig string) (dynamic.Interface, error) {
	return NewKubeClientDynamic("", kubeconfig, "")
}

//NewKubeClientDynamic returns a dynamic.Interface based on the provided url, kubeconfig and context
//url: The url of the server
//kubeconfig: The path of the kubeconfig, see (../config/config.go#LoadConfig) for more information
//context: The context to connect to
func NewKubeClientDynamic(url, kubeconfig, context string) (dynamic.Interface, error) {
	klog.V(5).Infof("Create kubeclient dynamic for url %s using kubeconfig path %s\n", url, kubeconfig)
	config, err := config.LoadConfig(url, kubeconfig, context)
	if err != nil {
		return nil, err
	}

	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

//NewDefaultKubeClientAPIExtension returns a clientset.Interface for the current-context in kubeconfig
//kubeconfig: The path of the kubeconfig, see (../config/config.go#LoadConfig) for more information
func NewDefaultKubeClientAPIExtension(kubeconfig string) (clientset.Interface, error) {
	return NewKubeClientAPIExtension("", kubeconfig, "")
}

//NewKubeClientAPIExtension returns a clientset.Interface based on the provided url, kubeconfig and context
//url: The url of the server
//kubeconfig: The path of the kubeconfig, see (../config/config.go#LoadConfig) for more information
//context: The context to connect to
func NewKubeClientAPIExtension(url, kubeconfig, context string) (clientset.Interface, error) {
	klog.V(5).Infof("Create kubeclient apiextension for url %s using kubeconfig path %s\n", url, kubeconfig)
	config, err := config.LoadConfig(url, kubeconfig, context)
	if err != nil {
		return nil, err
	}

	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

//HaveServerResources returns an error if all provided APIGroups are not installed
//client: the client to use
//expectedAPIGroups: The list of expected APIGroups
func HaveServerResources(client clientset.Interface, expectedAPIGroups []string) error {
	clientDiscovery := client.Discovery()
	for _, apiGroup := range expectedAPIGroups {
		klog.V(1).Infof("Check if %s exists", apiGroup)
		_, err := clientDiscovery.ServerResourcesForGroupVersion(apiGroup)
		if err != nil {
			klog.V(1).Infof("Error while retrieving server resource %s: %s", apiGroup, err.Error())
			return err
		}
	}
	return nil
}

//HaveCRDs returns an error if all provided CRDs are not installed
//client: the client to use
//expectedCRDs: The list of expected CRDS to find
func HaveCRDs(client clientset.Interface, expectedCRDs []string) error {
	clientAPIExtensionV1beta1 := client.ApiextensionsV1beta1()
	for _, crd := range expectedCRDs {
		klog.V(1).Infof("Check if %s exists", crd)
		_, err := clientAPIExtensionV1beta1.CustomResourceDefinitions().Get(context.TODO(), crd, metav1.GetOptions{})
		if err != nil {
			klog.V(1).Infof("Error while retrieving crd %s: %s", crd, err.Error())
			return err
		}
	}
	return nil
}

//HaveDeploymentsInNamespace returns an error if all provided deployment are not installed in the given namespace
//client: the client to use
//namespace: The namespace to search in
//expectedDeploymentNames: The deployment names to search
func HaveDeploymentsInNamespace(client kubernetes.Interface, namespace string, expectedDeploymentNames []string) error {
	versionInfo, err := client.Discovery().ServerVersion()
	if err != nil {
		return err
	}
	klog.V(1).Infof("Server version info: %v", versionInfo)

	deployments := client.AppsV1().Deployments(namespace)

	for _, deploymentName := range expectedDeploymentNames {
		klog.V(1).Infof("Check if deployment %s exists", deploymentName)
		deployment, err := deployments.Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			klog.V(1).Infof("Error while retrieving deployment %s: %s", deploymentName, err.Error())
			return err
		}
		if deployment.Status.Replicas != deployment.Status.ReadyReplicas {
			err = fmt.Errorf("Expect %d but got %d Ready replicas", deployment.Status.Replicas, deployment.Status.ReadyReplicas)
			klog.Errorln(err)
			return err
		}
		for _, condition := range deployment.Status.Conditions {
			if condition.Reason == "MinimumReplicasAvailable" {
				if condition.Status != corev1.ConditionTrue {
					err = fmt.Errorf("Expect %s but got %s", condition.Status, corev1.ConditionTrue)
					klog.Errorln(err)
					return err
				}
			}
		}
	}

	return nil
}
