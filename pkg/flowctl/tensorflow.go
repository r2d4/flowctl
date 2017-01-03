package flowctl

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"strings"
)

func encodeClusterSpec(jobs map[string]int) (string, error) {
	var buffer bytes.Buffer
	tmpl := template.Must(template.New("dnsTmpl").Parse(dnsTmpl))

	// Sort the jobs alphabetically for
	// reproducible cluster_specs
	jobNames := []string{}
	for j, _ := range jobs {
		jobNames = append(jobNames, j)
	}
	sort.Strings(jobNames)

	for i, jobName := range jobNames {
		replicas := jobs[jobName]
		buffer.WriteString(jobName + "|")
		serviceName := fmt.Sprintf("%s-%s", prefix, jobName)
		for j := 0; j < replicas; j++ {
			data := struct {
				ServiceName string
				Number      int
				Namespace   string
				DNSSuffix   string
				Port        int
			}{
				serviceName,
				j,
				defaultNamespace,
				dnsSuffix,
				2222,
			}
			err := tmpl.Execute(&buffer, data)
			if err != nil {
				return "", err
			}
			if j != replicas-1 {
				buffer.WriteString(";")
			}
		}
		if i != len(jobs)-1 {
			buffer.WriteString(",")
		}
	}
	return buffer.String(), nil
}

func decodeClusterSpec(spec string) (map[string]int, error) {
	var res map[string]int
	jobs := strings.Split(spec, ",")
	for _, job := range jobs {
		if job == "" {
			return nil, fmt.Errorf("Incorrect cluster spec, empty job: %s", spec)
		}
		servers := strings.Split(job, ";")
		if len(servers) == 0 {
			return nil, fmt.Errorf("Incorrect cluster spec, 0 servers for job %s", job)
		}
		res[job] = len(servers)
	}
	return res, nil
}
