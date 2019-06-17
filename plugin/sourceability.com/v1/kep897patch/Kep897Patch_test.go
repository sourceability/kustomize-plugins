package main_test

import (
	"strings"
	"testing"

	"sigs.k8s.io/kustomize/pkg/kusttest"
	"sigs.k8s.io/kustomize/plugin"
)

func TestKep897SimplePatch(t *testing.T) {
	tc := plugin.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"sourceability.com", "v1", "Kep897Patch")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")
	th.WriteF("/app/deployment-patch.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: notImportantHere
spec:
  template:
    spec:
      nodeSelector:
        type: prod
`)
	rm := th.LoadAndRunTransformer(`
apiVersion: sourceability.com/v1
kind: Kep897Patch
metadata:
  name: notImportantHere
patches:
- path: deployment-patch.yaml
  target:
    kind: Deployment
  type: StrategicMergePatch
`, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - image: elasticsearch
        name: elasticsearch
---
apiVersion: v1
kind: Pod
metadata:
  name: elasticsearch
spec:
  containers:
  - image: elasticsearch
    name: elasticsearch
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - image: elasticsearch
        name: elasticsearch
      nodeSelector:
        type: prod
---
apiVersion: v1
kind: Pod
metadata:
  name: elasticsearch
spec:
  containers:
  - image: elasticsearch
    name: elasticsearch
`)
}

func TestKep897PatchMultiple(t *testing.T) {
	tc := plugin.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"sourceability.com", "v1", "Kep897Patch")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")
	th.WriteF("/app/deployment-patch.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: notImportantHere
spec:
  template:
    spec:
      nodeSelector:
        type: prod
`)
	th.WriteF("/app/pod-patch.yaml", `
apiVersion: v1
kind: Pod
metadata:
  name: notImportantHere
spec:
  nodeSelector:
    type: prod
`)
	rm := th.LoadAndRunTransformer(`
apiVersion: sourceability.com/v1
kind: Kep897Patch
metadata:
  name: notImportantHere
patches:
- path: deployment-patch.yaml
  target:
    kind: Deployment
  type: StrategicMergePatch
- path: pod-patch.yaml
  target:
    kind: Pod
  type: StrategicMergePatch
`, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - image: elasticsearch
        name: elasticsearch
---
apiVersion: v1
kind: Pod
metadata:
  name: elasticsearch
spec:
  containers:
  - image: elasticsearch
    name: elasticsearch
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - image: elasticsearch
        name: elasticsearch
      nodeSelector:
        type: prod
---
apiVersion: v1
kind: Pod
metadata:
  name: elasticsearch
spec:
  containers:
  - image: elasticsearch
    name: elasticsearch
  nodeSelector:
    type: prod
`)
}

func TestKep897PatchLabels(t *testing.T) {
	tc := plugin.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"sourceability.com", "v1", "Kep897Patch")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")
	th.WriteF("/app/deployment-patch.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: notImportantHere
spec:
  template:
    spec:
      nodeSelector:
        type: prod
`)
	rm := th.LoadAndRunTransformer(`
apiVersion: sourceability.com/v1
kind: Kep897Patch
metadata:
  name: notImportantHere
patches:
- path: deployment-patch.yaml
  target:
    kind: Deployment
    LabelSelector:
      env: prod
  type: StrategicMergePatch
`, `
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    env: prod
  name: nginx
spec:
  template:
    spec:
      containers:
      - image: elasticsearch
        name: elasticsearch
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    env: prod
  name: elasticsearch
spec:
  containers:
  - image: elasticsearch
    name: elasticsearch
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    env: prod
  name: nginx
spec:
  template:
    spec:
      containers:
      - image: elasticsearch
        name: elasticsearch
      nodeSelector:
        type: prod
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    env: prod
  name: elasticsearch
spec:
  containers:
  - image: elasticsearch
    name: elasticsearch
`)
}

func TestKep897PatchAnnotations(t *testing.T) {
	tc := plugin.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"sourceability.com", "v1", "Kep897Patch")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")
	th.WriteF("/app/deployment-patch.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: notImportantHere
spec:
  template:
    spec:
      nodeSelector:
        type: prod
`)
	rm := th.LoadAndRunTransformer(`
apiVersion: sourceability.com/v1
kind: Kep897Patch
metadata:
  name: notImportantHere
patches:
- path: deployment-patch.yaml
  target:
    kind: Deployment
    matchAnnotations:
      env: prod
  type: StrategicMergePatch
`, `
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    env: prod
  name: nginx
spec:
  template:
    spec:
      containers:
      - image: elasticsearch
        name: elasticsearch
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    env: prod
  name: elasticsearch
spec:
  containers:
  - image: elasticsearch
    name: elasticsearch
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    env: prod
  name: nginx
spec:
  template:
    spec:
      containers:
      - image: elasticsearch
        name: elasticsearch
      nodeSelector:
        type: prod
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    env: prod
  name: elasticsearch
spec:
  containers:
  - image: elasticsearch
    name: elasticsearch
`)
}

func TestKep897PatchComplex(t *testing.T) {
	tc := plugin.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"sourceability.com", "v1", "Kep897Patch")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")
	th.WriteF("/app/deployment-patch.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: notImportantHere
spec:
  template:
    spec:
      nodeSelector:
        type: prod
      containers:
      - image: istio
        name: istio
`)
	rm := th.LoadAndRunTransformer(`
apiVersion: sourceability.com/v1
kind: Kep897Patch
metadata:
  name: notImportantHere
patches:
- path: deployment-patch.yaml
  target:
    kind: Deployment
    group: apps
    version: v1
    labelSelector:
      env: prod
    matchAnnotations:
      istio: enabled
  type: StrategicMergePatch
`, `
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    env: prod
  name: nginx
spec:
  template:
    spec:
      containers:
      - image: nginx
        name: nginx
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    istio: enabled
  labels:
    env: prod
  name: nginx-2
spec:
  template:
    spec:
      containers:
      - image: nginx
        name: nginx
---
apiVersion: apps/v2
kind: Deployment
metadata:
  annotations:
    istio: enabled
  labels:
    env: prod
  name: nginx-3
spec:
  template:
    spec:
      containers:
      - image: nginx
        name: nginx
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    istio: enabled
  labels:
    env: prod
  name: nginx-4
spec:
  template:
    spec:
      containers:
      - image: nginx
        name: nginx
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    env: prod
  name: elasticsearch
spec:
  containers:
  - image: elasticsearch
    name: elasticsearch
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    istio: enabled
  labels:
    env: prod
  name: nginx-2
spec:
  template:
    spec:
      containers:
      - image: istio
        name: istio
      - image: nginx
        name: nginx
      nodeSelector:
        type: prod
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    istio: enabled
  labels:
    env: prod
  name: nginx-4
spec:
  template:
    spec:
      containers:
      - image: istio
        name: istio
      - image: nginx
        name: nginx
      nodeSelector:
        type: prod
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    env: prod
  name: nginx
spec:
  template:
    spec:
      containers:
      - image: nginx
        name: nginx
---
apiVersion: apps/v2
kind: Deployment
metadata:
  annotations:
    istio: enabled
  labels:
    env: prod
  name: nginx-3
spec:
  template:
    spec:
      containers:
      - image: nginx
        name: nginx
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    env: prod
  name: elasticsearch
spec:
  containers:
  - image: elasticsearch
    name: elasticsearch
`)
}

func TestKep897PatchOnlyAcceptsStrategicMerge(t *testing.T) {
	tc := plugin.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"sourceability.com", "v1", "Kep897Patch")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")
	th.WriteF("/app/deployment-patch.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: notImportantHere
spec:
  template:
    spec:
      nodeSelector:
        type: prod
`)
	err := th.ErrorFromLoadAndRunTransformer(`
apiVersion: sourceability.com/v1
kind: Kep897Patch
metadata:
  name: notImportantHere
patches:
- path: deployment-patch.yaml
  target:
    kind: Deployment
  type: JsonPatch
`, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - image: elasticsearch
        name: elasticsearch
`)
	if err == nil {
		t.Fatalf("expected an error")
	}
	if !strings.Contains(err.Error(),
		"StrategicMergePatch is the only supported patch type, \"JsonPatch\" given") {
		t.Fatalf("incorrect error %v", err)
	}
}
