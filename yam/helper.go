package yam

import (
	"github.com/kakaisaname/props/kvs"
	"github.com/prometheus/common/log"
	"strings"
)

func ByYaml(content string) *kvs.MapProperties {
	y := NewYamlProperties()
	err := y.Load(strings.NewReader(content))
	if err != nil {
		log.Error(err)
		return nil
	}
	return &y.MapProperties
}
