# Sourceability Kustomize Plugins

## KEP 897

This plugin addresses [kubernetes-sigs/kustomize#720](https://github.com/kubernetes-sigs/kustomize/issues/720) by implementing the
`StrategicMergePatch` from [KEP 897](https://github.com/kubernetes/enhancements/pull/897).

Note that I've encountered issues applying patches to resources with prefixes/suffixes.
I'm assuming that other issues like exist with this plugin.

This repository contains one very simple example kustomization, that can be built with:
```
docker run \
    -it --rm \
    -w /app -v $PWD:/app \
    sourceability/kustomize-plugins \
    build --enable_alpha_plugins examples/deployment-patch
```

Which should output:
```
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: elasticsearch
  name: elasticsearch
spec:
  selector:
    matchLabels:
      app: elasticsearch
  template:
    metadata:
      labels:
        app: elasticsearch
    spec:
      containers:
      - image: elasticsearch
        name: elasticsearch
      nodeSelector:
        env: prod
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
      nodeSelector:
        env: prod
```
