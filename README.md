# registry-polaris (*This is a community driven project*)

Some application runtime use [polaris](https://github.com/polarismesh/polaris) for service discovery.Polaris is a cloud-native service discovery and governance center. 
It can be used to solve the problem of service connection, fault tolerance, traffic control and secure in distributed and microservice architecture.

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
    // https://www.cloudwego.io/docs/tutorials/framework-exten/registry/#integrate-into-kitex
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

## Todolist
   Use polaris's watch mechanism to monitor registered service changes

## More info

 See example

## Compatibility

Compatible with polaris.

maintained by: [liu-song](https://github.com/liu-song)
