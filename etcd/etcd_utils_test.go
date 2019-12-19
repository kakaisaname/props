package etcd

import (
	"github.com/kakaisaname/props/kvs"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var etcd_mock_started = false
var etcdAddress string

func init() {
	//etcdAddress = "http://172.16.1.248:2379"
	etcdAddress = "http://127.0.0.1:2379"
	GetOrNewMockTestEtcd(etcdAddress)
	go kvs.ExecCommand("pwd", "-LP")
	if !etcd_mock_started {
		go testEtcd.StartMockEtcd()
	}
	testEtcd.WaitingForEtcdStarted()
}

var testEtcd *MockTestEtcd

type MockTestEtcd struct {
	Address string
}

func GetOrNewMockTestEtcd(address string) *MockTestEtcd {
	if testEtcd == nil {
		testEtcd = &MockTestEtcd{
			Address: address,
		}
	}

	return testEtcd
}

func (m *MockTestEtcd) StartMockEtcd() <-chan int {
	ec := make(chan int, 1)
	isStarted := m.CheckEtcdIsStarted()

	if isStarted {
		ec <- 1
		return ec
	}
	command := "etcd"

	started := kvs.ExecCommand(command, "--data-dir=./temp/etcd")

	if started {
		ec <- 1
	} else {
		ec <- 0
	}
	return ec
}

func (m *MockTestEtcd) WaitingForEtcdStarted() {
	for {
		isStarted := m.CheckEtcdIsStarted()
		if isStarted {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func (m *MockTestEtcd) CheckEtcdIsStarted() bool {
	res, err := http.Get(etcdAddress + "/version")
	log.Info(res, err)

	if err != nil {
		return false
	}
	if res != nil && res.StatusCode == 200 {
		res.Body.Close()
		return true
	}
	return false
}
