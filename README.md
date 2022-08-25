# registry-polaris (*This is a community driven project*)

## Notice 

This repository has been migrated to [polaris](https://github.com/kitex-contrib/polaris) and welcome to use.

~~Some application runtime use [polaris](https://github.com/polarismesh/polaris) for service discovery. Polaris is a cloud-native service discovery and governance center. 
It can be used to solve the problem of service connection, fault tolerance, traffic control and secure in distributed and microservice architecture.~~

## ~~How to install registry-polaris?~~
```
go get -u github.com/kitex-contrib/registry-polaris
```

## ~~How to use with Kitex server?~~

```go
import (
        ...
   	"context"
   	"log"
   	"net"
   
   	"github.com/cloudwego/kitex/pkg/registry"
   	"github.com/polarismesh/polaris-go/pkg/config"
   	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api"
   	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api/hello"
   	"github.com/cloudwego/kitex/server"
   	polaris "github.com/kitex-contrib/registry-polaris"
        ...
)

const (
	confPath       = "polaris.yaml"
	Namespace      = "Polaris"
	// At present,polaris server tag is v1.4.0，can't support auto create namespace,
	// If you want to use a namespace other than default,Polaris ,before you register an instance,
	// you should create the namespace at polaris console first.
)

func main() {
    ...
	r, err := polaris.NewPolarisRegistry(confPath)
	if err != nil {
		log.Fatal(err)
	}
	Info := &registry.Info{
		ServiceName: "echo",
		Tags: map[string]string{
			"namespace": Namespace,
		},
	}
        // https://www.cloudwego.io/docs/kitex/tutorials/framework-exten/service_discovery/#usage-example
	newServer := hello.NewServer(new(HelloImpl), server.WithRegistry(r), server.WithRegistryInfo(Info),
		server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8888}))

	err = newServer.Run()
	if err != nil {
		log.Fatal(err)
	}
	...
}
```


## ~~How to use with Kitex client?~~

```go
import (
        ...
	"context"
	"log"
	"time"

	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api"
	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api/hello"
	"github.com/cloudwego/kitex/client"
	polaris "github.com/kitex-contrib/registry-polaris"
	"github.com/polarismesh/polaris-go/pkg/config"
        ...
)

const (
	confPath       = "polaris.yaml"
	Namespace      = "Polaris"
	// At present,polaris server tag is v1.4.0，can't support auto create namespace,
	// if you want to use a namespace other than default,Polaris ,before you register an instance,
	// you should create the namespace at polaris console first.
)

func main() {
    ...
	r, err := polaris.NewPolarisResolver(confPath)
	if err != nil {
		log.Fatal(err)
	}
        // https://www.cloudwego.io/docs/kitex/tutorials/framework-exten/service_discovery/#usage-example
	// client.WithTag sets the namespace tag for service discovery
	newClient := hello.MustNewClient("echo", client.WithTag("namespace", Namespace),
		client.WithResolver(r), client.WithRPCTimeout(time.Second*60))
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	resp, err := newClient.Echo(ctx, &api.Request{Message: "Hi,polaris!"})
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
	...
	}
}

```
## ~~How to install polaris?~~
Polaris support stand-alone and cluster. More information can be found in [install polaris](https://polarismesh.cn/zh/doc/%E5%BF%AB%E9%80%9F%E5%85%A5%E9%97%A8/%E5%AE%89%E8%A3%85%E6%9C%8D%E5%8A%A1%E7%AB%AF/%E5%AE%89%E8%A3%85%E5%8D%95%E6%9C%BA%E7%89%88.html#%E5%8D%95%E6%9C%BA%E7%89%88%E5%AE%89%E8%A3%85)

## ~~Todolist~~
~~Welcome to contribute your ideas~~

## ~~Use polaris with Kitex~~

~~See example and test~~
  
## ~~Compatibility~~

~~Compatible with polaris.~~

maintained by: [liu-song](https://github.com/liu-song)
