package flowctl

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
)

const dnsTmpl = "{{.ServiceName}}-{{.Number}}.{{.ServiceName}}.{{.Namespace}}.{{.DNSSuffix}}:{{.Port}}"

func createStatefulSet(jobName, image, clusterSpec string, replicas int) error {
	client, err := getClient(isInCluster)
	if err != nil {
		return err
	}
	appsClient := client.Apps().StatefulSets(defaultNamespace)
	job := statefulSetSpec(jobName, image, clusterSpec, replicas)
	_, err = appsClient.Create(job)
	if err != nil {
		return err
	}
	glog.Infof("Created stateful set for job %s\n", jobName)
	return nil
}

func statefulSetSpec(jobName, image, clusterSpec string, replicas int) *v1beta1.StatefulSet {
	replica := int32(replicas)
	terminationPeriodSeconds := int64(0)
	return &v1beta1.StatefulSet{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1beta1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				"app": "tensorflow",
			},
		},
		Spec: v1beta1.StatefulSetSpec{
			Replicas:    &replica,
			ServiceName: fmt.Sprintf("%s-%s", prefix, jobName),
			Template: v1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						prefix: jobName,
						"app":  "tensorflow",
					},
					Annotations: map[string]string{
						"pod.alpha.kubernetes.io/initialized": "true",
					},
				},
				Spec: v1.PodSpec{
					TerminationGracePeriodSeconds: &terminationPeriodSeconds,
					Containers: []v1.Container{
						{
							Name:            fmt.Sprintf("%s-%s", prefix, jobName),
							Image:           image,
							ImagePullPolicy: v1.PullIfNotPresent,
							Command: []string{
								"/bin/bash",
								"-c",
							},
							Args: []string{
								"/tf_grpc.py --task_id=${HOSTNAME##*-}" + fmt.Sprintf(" --cluster_spec='%s'", clusterSpec) + fmt.Sprintf(" --job_name=%s", jobName),
							},
							Ports: []v1.ContainerPort{
								{
									ContainerPort: 2222,
								},
							},
						},
					},
				},
			},
		},
	}
}
