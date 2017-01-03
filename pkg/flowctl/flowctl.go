package flowctl

import (
	"github.com/golang/glog"
	"github.com/pkg/errors"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const defaultImage = "gcr.io/k8s-minikube/tensorflow_grpc:0.1"
const prefix = "tf"
const dnsSuffix = "svc.cluster.local"
const defaultNamespace = "default"
const isInCluster = false

func CreateServers(jobs map[string]int) error {
	clusterSpec, err := encodeClusterSpec(jobs)
	if err != nil {
		return err
	}
	for jobName, replicas := range jobs {
		err := createService(jobName)
		if err != nil {
			return err
		}
		err = createStatefulSet(jobName, defaultImage, clusterSpec, replicas)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteServers() error {
	client, err := getClient(isInCluster)
	if err != nil {
		return err
	}
	listOptions := v1.ListOptions{
		LabelSelector: "app=tensorflow",
	}
	deleteOptions := &v1.DeleteOptions{}

	err = client.
		Apps().
		StatefulSets(defaultNamespace).
		DeleteCollection(deleteOptions, listOptions)

	if err != nil {
		return errors.Wrap(err, "deleting statefulsets")
	}
	glog.Infoln("Deleted statefulsets")

	err = client.
		CoreV1().
		Pods(defaultNamespace).
		DeleteCollection(deleteOptions, listOptions)

	if err != nil {
		return errors.Wrap(err, "deleting pods")
	}
	glog.Infoln("Deleted pods")

	svcs, err := client.
		CoreV1().
		Services(defaultNamespace).
		List(listOptions)

	if err != nil {
		return errors.Wrap(err, "getting svcs")
	}

	for _, svc := range svcs.Items {
		err = client.
			CoreV1().
			Services(defaultNamespace).
			Delete(svc.ObjectMeta.Name, deleteOptions)
		if err != nil {
			return errors.Wrapf(err, "deleting svc %s", svc.ObjectMeta.Name)
		}
		glog.Infof("Deleted svc %s", svc.ObjectMeta.Name)
	}

	glog.Infoln("Deleted services")

	return nil
}

func getClient(inCluster bool) (*kubernetes.Clientset, error) {
	var cfg *rest.Config
	var err error
	if inCluster {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		kubeconfigPath := clientcmd.RecommendedHomeFile
		cfg, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
			&clientcmd.ConfigOverrides{}).ClientConfig()
		if err != nil {
			return nil, err
		}
	}
	return kubernetes.NewForConfig(cfg)
}
