package polaris

import (
	"net"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/polarismesh/polaris-go/pkg/model"
)

type polarisKitexInstance struct {
	kitexInstance   discovery.Instance
	polarisInstance model.Instance
	polarisOptions  Options
}

func (i *polarisKitexInstance) Address() net.Addr {
	return i.kitexInstance.Address()
}

func (i *polarisKitexInstance) Weight() int {
	return i.kitexInstance.Weight()
}

func (i *polarisKitexInstance) Tag(key string) (value string, exist bool) {
	value, exist = i.kitexInstance.Tag(key)
	return
}
