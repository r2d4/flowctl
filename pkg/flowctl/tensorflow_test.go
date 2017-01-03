package flowctl

import "testing"

type testCase struct {
	Jobs     map[string]int
	Expected string
	Err      bool
}

var testCases = []testCase{
	{
		Jobs: map[string]int{
			"worker": 2,
			"ps":     1,
		},
		Expected: "ps|tf-ps-0.tf-ps.default.svc.cluster.local:2222,worker|tf-worker-0.tf-worker.default.svc.cluster.local:2222;tf-worker-1.tf-worker.default.svc.cluster.local:2222",
		Err:      false,
	},
}

func TestDecodeClusterSpec(t *testing.T) {
	for _, tt := range testCases {
		actual, err := encodeClusterSpec(tt.Jobs)
		if err != nil && !tt.Err {
			t.Errorf("Cluster spec valid but returned error: %s", err)
		}
		if actual != tt.Expected {
			t.Errorf(`Cluster spec
				actual: %s
				expected: %s
				`, actual, tt.Expected)
		}
	}
}
