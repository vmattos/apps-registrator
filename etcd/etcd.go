package etcd

import (
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/client"
)

type Etcd struct {
	client client.KeysAPI
}

func New() Etcd {

	etcdIP := "http://" + os.Getenv("ETCD_IP") + ":4001"
	etcdCfg := client.Config{
		Endpoints:               []string{etcdIP},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(etcdCfg)
	if err != nil {
		panic(err)
	}

	return Etcd{client: client.NewKeysAPI(c)}
}

func (self *Etcd) Get(key string) (string, error) {
	resp, err := self.client.Get(context.Background(), key, nil)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func (self *Etcd) Set(key, value string) (string, error) {
	resp, err := self.client.Set(context.Background(), key, value, nil)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}
