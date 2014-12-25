# kube-register

Register Kubernetes Kubelet machines with the Kubernetes API server using Fleet data.

## Usage

```
kube-register -metadata="kubelet=true" -fleet-endpoint="http://127.0.0.1:4002" -apiserver-endpoint="http://127.0.0.1:8080" -reverse-lookup=false
```

By default kube-register registers new machines with their public IP addresses, found in fleet. By setting `-reverse-lookup=true` kube-register will do a reverse DNS lookup to get the machine's hostname and registers the machine with its hostname instead with its IP address. Make sure the machine's IP resolves to the same hostname that is printed by `hostname -f` on the machine itself.

## Building

```
mkdir -p "${GOPATH}/src/github.com/kelseyhightower"
cd "${GOPATH}/src/github.com/kelseyhightower"
git clone https://github.com/kelseyhightower/kube-register.git
cd kube-register
godep go build .
```

Slightly smaller binary.

```
CGO_ENABLED=0 GOOS=linux godep go build -a -tags netgo -ldflags '-w' .
```

To build the binary inside a docker container.

```
docker run -v $SRC:/opt/kube-register -i -t google/golang /bin/bash -c "cd /opt/kube-register && go get github.com/tools/godep && godep go build ."
```

Where $SRC is the absolute path to the kube-register directory.
