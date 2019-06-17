FROM golang:1.12-alpine AS build

RUN apk add git gcc g++

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=1

COPY . /app
WORKDIR /app

RUN go install sigs.k8s.io/kustomize

RUN cd plugin/sourceability.com/v1/kep897patch \
    && go build -buildmode plugin -o Kep897Patch.so Kep897Patch.go \
    && chmod +x Kep897Patch.so

# --------------------

FROM alpine:3.9

COPY --from=build /go/bin/kustomize /bin/kustomize
COPY --from=build /app/plugin /root/.config/kustomize/plugin/

ENTRYPOINT ["/bin/kustomize"]
