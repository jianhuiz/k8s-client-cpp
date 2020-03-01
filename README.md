This is a incomplete experimentation of creating a C++ binding of Kubernetes go client library.

##### Environment
Go
C++ 14
protobuf

##### clone the code
``` bash
git clone https://github.com/jianhuiz/k8s-client-cpp
```

*Note: All blocks below must be executed from the root of the repository.*

##### get the dependencies
``` bash
cd ./go
go mod vendor
```

##### generate c++ protobuf code
``` bash
protoc -I $PWD/go/vendor --cpp_out $PWD/cpp/ \
    $PWD/go/vendor/k8s.io/api/core/v1/generated.proto \
    $PWD/go/vendor/k8s.io/apimachinery/pkg/api/resource/generated.proto \
    $PWD/go/vendor/k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto \
    $PWD/go/vendor/k8s.io/apimachinery/pkg/runtime/generated.proto \
    $PWD/go/vendor/k8s.io/apimachinery/pkg/runtime/schema/generated.proto \
    $PWD/go/vendor/k8s.io/apimachinery/pkg/util/intstr/generated.proto
```

##### build go library
``` bash
cd ./go
go build -buildmode=c-archive -o go.a
```
##### build c++ sample
``` bash
cd ./cpp
make
```

#### macOS users
Change the LDFLAGS to:
```
LD_FLAGS = ../go/go.a -lprotobuf -lpthread -framework CoreFoundation -framework Security
```
