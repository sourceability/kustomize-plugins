FROM golang:1.12-alpine AS kustomize

RUN apk add git gcc g++

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=1 \
    KUSTOMIZE_VERSION=sigs.k8s.io/kustomize@v0.0.0-20190607032952-51d82bece3dd
ENV KUSTOMIZE_DIR=/go/pkg/mod/$KUSTOMIZE_VERSION

RUN go get $KUSTOMIZE_VERSION \
    && cd $KUSTOMIZE_DIR \
    && go build

ENV PLUGIN_DIR=/root/.config/kustomize/plugin/sourceability.com

COPY plugin/sourceability.com $PLUGIN_DIR

WORKDIR $PLUGIN_DIR/v1/kep897patch

RUN GOROOT=/usr/local/go go test -v Kep897Patch_test.go
RUN go build -buildmode plugin -o Kep897Patch.so Kep897Patch.go \
    && chmod +x Kep897Patch.so

# --------------------

FROM alpine:3.9

COPY --from=kustomize /go/pkg/mod/sigs.k8s.io/kustomize@v0.0.0-20190607032952-51d82bece3dd/kustomize /bin/kustomize
COPY --from=kustomize /root/.config/kustomize /root/.config/kustomize

ENTRYPOINT ["/bin/kustomize"]
