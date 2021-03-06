//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package cluster

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/cloustone/sentel/pkg/config"
	sd "github.com/cloustone/sentel/pkg/service-discovery"
	"github.com/golang/glog"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
)

type k8sCluster struct {
	config           config.Config
	mutex            sync.Mutex
	pods             map[string]string
	clientset        *kubernetes.Clientset
	serviceDiscovery sd.ServiceDiscovery
}

// newClusterManager retrieve clustermanager instance connected with clustermgr
func newK8sCluster(c config.Config) (*k8sCluster, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolue path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolue path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &k8sCluster{
		config:    c,
		mutex:     sync.Mutex{},
		pods:      make(map[string]string),
		clientset: clientset,
	}, nil
}

func (p *k8sCluster) Initialize() error {
	return nil
}
func (p *k8sCluster) SetServiceDiscovery(s sd.ServiceDiscovery) {
	p.serviceDiscovery = s
}

func (p *k8sCluster) CreateNetwork(name string) (string, error) {
	return "", nil
}
func (p *k8sCluster) RemoveNetwork(name string) error {
	return nil
}

// CreateBrokers create a number of brokers for tenant and product
func (p *k8sCluster) CreateService(spec ServiceSpec) (string, error) {
	podname := fmt.Sprintf("tenant-%s", spec.TenantId)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	deploymentsClient := p.clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: spec.ServiceName,
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32ptr(spec.Replicas),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": spec.ServiceName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  podname,
							Image: spec.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          spec.ServiceName,
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		return "", err
	}
	glog.Infof("broker deployment created:%q.\n", result.GetObjectMeta().GetName())
	return podname, nil
}

func (p *k8sCluster) RemoveService(serviceName string) error {
	deletePolicy := metav1.DeletePropagationForeground
	deploymentsClient := p.clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	return deploymentsClient.Delete(serviceName, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func (p *k8sCluster) UpdateService(serviceId string, spec ServiceSpec) error {
	deploymentsClient := p.clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(spec.ServiceName, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}

		result.Spec.Replicas = int32ptr(spec.Replicas)
		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})
	return retryErr
}
func (p *k8sCluster) IntrospectService(serviceId string) (ServiceIntrospec, error) {
	return ServiceIntrospec{}, errors.New("no implemented")
}

func int32ptr(i int32) *int32 { return &i }
