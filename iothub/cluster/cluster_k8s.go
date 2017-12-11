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
	"flag"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
)

type k8sClusterMgr struct {
	config    core.Config
	mutex     sync.Mutex
	tenants   map[string]*tenant
	products  map[string]*product
	brokers   map[string]*broker
	clientset *kubernetes.Clientset
}

// newClusterManager retrieve clustermanager instance connected with clustermgr
func newK8sClusterManager(c core.Config) (*k8sClusterMgr, error) {
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
	return &k8sClusterMgr{
		config:    c,
		mutex:     sync.Mutex{},
		tenants:   make(map[string]*tenant),
		products:  make(map[string]*product),
		brokers:   make(map[string]*broker),
		clientset: clientset,
	}, nil
}

// CreateBrokers create a number of brokers for tenant and product
func (p *k8sClusterMgr) CreateBrokers(tid string, pid string, count int32) ([]string, error) {
	podname := fmt.Sprintf("%s-%s", tid, pid)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.products[pid]; found {
		return nil, fmt.Errorf("product '%s' already existed in iothub", pid)
	}
	deploymentsClient := p.clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sentel-broker",
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32ptr(count),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "sentel-broker",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  podname,
							Image: "sentel-broker:1.00",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "broker",
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
		return nil, err
	}
	glog.Infof("broker deployment created:%q.\n", result.GetObjectMeta().GetName())

	// maybe we shoud wait pod to be started
	time.Sleep(5 * time.Second) // TODO

	// get pod list
	pods, err := p.clientset.CoreV1().Pods(podname).List(metav1.ListOptions{})
	if err != nil {
		glog.Errorf("Failed to get pod list for tenant(%s)", tid)
		return nil, err
	}
	pp := &product{
		pid:       pid,
		tid:       tid,
		createdAt: time.Now(),
		brokers:   make(map[string]*broker),
	}

	// get all created pods, create broker for each pod
	names := []string{}
	for _, pod := range pods.Items {
		b := &broker{
			bid:         uuid.NewV4().String(),
			status:      brokerStatusStarted,
			createdAt:   time.Now(),
			lastUpdated: time.Now(),
			context:     &pod,
		}
		pp.brokers[b.bid] = b
		names = append(names, b.bid)
		p.brokers[b.bid] = b
	}
	p.products[pid] = pp
	if _, found := p.tenants[tid]; !found {
		p.tenants[tid] = &tenant{
			tid:       tid,
			createdAt: time.Now(),
			products:  make(map[string]*product),
		}
	} else {
		p.tenants[tid].products[pid] = pp
	}
	return names, nil
}

// startBroker start specified broker
func (p *k8sClusterMgr) StartBroker(bid string) error {
	return nil
}

// stopBroker stop specified node
func (p *k8sClusterMgr) StopBroker(bid string) error {
	return nil
}

// startBroker start specified broker
func (p *k8sClusterMgr) StartBrokers(tid, bid string) error {
	return nil
}

// stopBroker stop specified node
func (p *k8sClusterMgr) StopBrokers(tid, bid string) error {
	return nil
}

// deleteBrokers stop and delete brokers for tenant
func (p *k8sClusterMgr) DeleteBrokers(tid, pid string) error {
	podname := fmt.Sprintf("%s-%s", tid, pid)
	deletePolicy := metav1.DeletePropagationForeground
	deploymentsClient := p.clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	return deploymentsClient.Delete(podname, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

// deleteBroker stop and delete specified broker
func (p *k8sClusterMgr) DeleteBroker(bid string) error {
	return nil
}

// rollbackBrokers rollback tenant's brokers
func (p *k8sClusterMgr) RollbackBrokers(tid, bid string, replicas int32) error {
	podname := fmt.Sprintf("%s-%s", tid, bid)
	deploymentsClient := p.clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(podname, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}

		result.Spec.Replicas = int32ptr(replicas)
		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})
	return retryErr
}

func int32ptr(i int32) *int32 { return &i }
