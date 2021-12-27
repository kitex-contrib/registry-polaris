# registry-polaris (*This is a community driven project*)

Some application runtime use [polaris](https://github.com/polarismesh/polaris) for service discovery.Polaris is a cloud-native service discovery and governance center. 
It can be used to solve the problem of service connection, fault tolerance, traffic control and secure in distributed and microservice architecture.

## How to use with registry-polaris?
```
go get -u -v github.com/kitex-contrib/registry-polaris
```

## How to use with Kitex server?

```go
import (
    ...
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
    polaris "github.com/kitex-contrib/registry-polaris"
    ...
)

func main() {
    ...
    r, err := polaris.NewPolarisRegistry([]string{"127.0.0.1:8091"}) // r should not be reused.
    if err != nil {
        log.Fatal(err)
    }
    // https://www.cloudwego.io/docs/kitex/tutorials/framework-exten/registry/#integrate-into-kitex
    server, err := echo.NewServer(new(EchoImpl), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "echo"}, server.WithRegistry(r)))
    if err != nil {
        log.Fatal(err)
    }
    err = server.Run()
    if err != nil {
        log.Fatal(err)
    }
    ...
}
```


## How to use with Kitex client?

```go
import (
    ...
    "github.com/cloudwego/kitex/client"
    polaris "github.com/kitex-contrib/registry-polaris"
    ...
)

func main() {
    ...
    r, err := polaris.NewPolarisResolver([]string{"127.0.0.1:8091"})
    if err != nil {
        log.Fatal(err)
    }
    client, err := echo.NewClient("echo", client.WithResolver(r))
    if err != nil {
        log.Fatal(err)
    }
    ...
}
```
## How to install Polaris?
Polaris support stand-alone and cluster. More information can be found in [install Polaris](https://polarismesh.cn/zh/doc/%E5%BF%AB%E9%80%9F%E5%85%A5%E9%97%A8/%E5%AE%89%E8%A3%85%E6%9C%8D%E5%8A%A1%E7%AB%AF/%E5%AE%89%E8%A3%85%E5%8D%95%E6%9C%BA%E7%89%88.html#%E5%8D%95%E6%9C%BA%E7%89%88%E5%AE%89%E8%A3%85)

## Todolist
   Use polaris's watch mechanism to monitor registered service changes

## Use polaris with Kitex

See example
  
## Compatibility

Compatible with polaris.

maintained by: [liu-song](https://github.com/liu-song)
