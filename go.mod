module github.com/deitch/eci-distribution

go 1.12

require (
	github.com/containerd/containerd v1.3.0
	github.com/deislabs/oras v0.7.0
	github.com/docker/docker v0.7.3-0.20190826074503-38ab9da00309
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
)

replace github.com/docker/docker => github.com/docker/docker v0.7.3-0.20190826074503-38ab9da00309
