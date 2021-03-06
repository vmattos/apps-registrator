package etcd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/client"
	"github.com/vtex/apps-registrator/models"
)

type Etcd struct {
	client client.KeysAPI
	Prefix string
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

	return Etcd{
		client: client.NewKeysAPI(c),
		Prefix: "/vulcand",
	}
}

func (self *Etcd) SetPrefix(prefix string) {
	self.Prefix = prefix
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

func (self *Etcd) SetRoute(route *models.Route) {
	self.setBackend(route.Backend)
	self.setServer(route.Backend)
	self.setFrontend(route)
	fmt.Println("setting route on etcd")
}

func (self *Etcd) setBackend(bckID string) {
	backend := models.Backend{
		Type: "http",
	}
	value, err := json.Marshal(backend)
	if err != nil {
		panic(err)
	}
	id := strings.Split(bckID, "http://")[1]
	key := self.Prefix + "/backends/" + id + "/backend"
	stringValue := string(value[:])
	self.Set(key, stringValue)
}

func (self *Etcd) setServer(bckID string) {
	server := models.Server{
		URL: bckID,
	}
	value, err := json.Marshal(server)
	if err != nil {
		panic(err)
	}
	id := strings.Split(bckID, "http://")[1]
	key := self.Prefix + "/backends/" + id + "/servers/srv"
	stringValue := string(value[:])
	self.Set(key, stringValue)
}

func (self *Etcd) setFrontend(route *models.Route) {
	path := "PathRegexp(`" + route.Path + "`)"
	bckID := strings.Split(route.Backend, "http://")[1]
	frontend := models.Frontend{
		Type:      "http",
		BackendId: bckID,
		Route:     path,
	}
	value, err := json.Marshal(frontend)
	if err != nil {
		panic(err)
	}
	key := self.Prefix + "/frontends/" + route.Name + "/frontend"
	stringValue := string(value[:])
	self.Set(key, stringValue)
}
