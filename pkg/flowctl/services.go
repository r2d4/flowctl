package flowctl

func serviceSpec(jobName string) *v1.Service {
	return &v1.Service{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				prefix: jobName,
				"app":  "tensorflow",
			},
		},
		Spec: v1.ServiceSpec{
			ClusterIP: "None",
			Selector: map[string]string{
				prefix: jobName,
			},
			Ports: []v1.ServicePort{
				{Port: 2222},
			},
		},
	}
}

func createService(jobName string) error {
	client, err := getClient(isInCluster)
	if err != nil {
		return err
	}
	srv := serviceSpec(jobName)
	_, err = client.
		CoreV1().
		Services(defaultNamespace).
		Create(srv)

	if err != nil {
		return err
	}

	glog.Infof("Created service for job %s\n", jobName)

	return nil
}
