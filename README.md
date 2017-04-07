This is a incomplete experimentation of creating a C++ binding of Kubernetes go client library.

##### Environment
Go  
C++ 14   
protobuf

##### clone the code
``` bash
mkdir -p $GOPATH/src/github.com/jianhuiz/k8s-client-cpp
git clone https://github.com/jianhuiz/k8s-client-cpp $GOPATH/src/github.com/jianhuiz/k8s-client-cpp
```

##### get go client and the dependencies
``` bash
go get k8s.io/client-go
go get k8s.io/apiserver/pkg/apis/example/v1
```

##### generate c++ protobuf code
``` bash
protoc -I ./src --cpp_out ./src/github.com/jianhuiz/k8s-client-cpp/cpp/ \
    ./src/k8s.io/client-go/pkg/api/v1/generated.proto \
    ./src/k8s.io/apimachinery/pkg/api/resource/generated.proto \
    ./src/k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto \
    ./src/k8s.io/apimachinery/pkg/runtime/generated.proto \
    ./src/k8s.io/apimachinery/pkg/runtime/schema/generated.proto \
    ./src/k8s.io/apimachinery/pkg/util/intstr/generated.proto \
    ./src/k8s.io/apiserver/pkg/apis/example/v1/generated.proto
```

##### build go library
``` bash
cd $GOPATH/src/github.com/jianhuiz/k8s-client-cpp/go
go build -buildmode=c-archive
```
##### build c++ sample
``` bash
cd $GOPATH/src/github.com/jianhuiz/k8s-client-cpp/cpp
make
```
