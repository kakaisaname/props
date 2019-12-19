package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/kakaisaname/props/kvs"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	CONSUL_WAIT_TIME = time.Second * 10
)

//通过key/value来组织，过滤root prefix后，替换/为.作为properties key
//Deprecated
//看看ConsulConfigSource
type ConsulKeyValueConfigSource struct {
	kvs.MapProperties
	name   string
	root   string
	client *api.Client
	kv     *api.KV
	config *api.Config
}

//Deprecated
func NewConsulKeyValueConfigSource(address, root string) *ConsulKeyValueConfigSource {
	return NewConsulKeyValueConfigSourceByName("", address, root, CONSUL_WAIT_TIME)
}

//Deprecated
func NewConsulKeyValueConfigSourceByName(name, address, root string, timeout time.Duration) *ConsulKeyValueConfigSource {
	s := &ConsulKeyValueConfigSource{}
	if name == "" {
		name = strings.Join([]string{"consul", address, root}, ":")
	}

	s.name = name
	s.Values = make(map[string]string)
	s.root = root
	s.config = api.DefaultConfig()
	s.config.Address = address
	s.config.WaitTime = timeout
	client, err := api.NewClient(s.config)
	if err != nil {
		panic(err)
	}
	s.client = client
	s.kv = client.KV()
	s.init()

	return s
}

//Deprecated
func NewConsulKeyValueCompositeConfigSource(contexts []string, address string) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "ConsulKevValue"
	for _, context := range contexts {
		c := NewConsulKeyValueConfigSource(address, context)
		s.Add(c)
	}

	return s
}

func (s *ConsulKeyValueConfigSource) init() {
	s.findProperties(s.root, nil)
}

func (s *ConsulKeyValueConfigSource) watchContext() {

}

func (s *ConsulKeyValueConfigSource) Close() {
}

func (s *ConsulKeyValueConfigSource) findProperties(parentPath string, children []string) {
	prefix := s.root
	q := &api.QueryOptions{}

	keys, _, err := s.kv.Keys(prefix, "", q)
	if err != nil {
		log.Error(err)
		return
	}
	for _, k := range keys {
		kv, _, err := s.kv.Get(k, q)
		if err != nil {
			log.Error(err)
			continue
		}
		value := string(kv.Value)
		s.registerKeyValue(k, value)
	}

}

func (s *ConsulKeyValueConfigSource) sanitizeKey(path string, context string) string {
	key := strings.Replace(path, context+"/", "", -1)
	key = strings.Replace(key, "/", ".", -1)
	return key
}

func (s *ConsulKeyValueConfigSource) registerKeyValue(path, value string) {
	key := s.sanitizeKey(path, s.root)
	s.Set(key, value)

}

func (s *ConsulKeyValueConfigSource) Name() string {
	return s.name
}
