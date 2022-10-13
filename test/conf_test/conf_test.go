package conf_test

import (
	"chat/conf"
	"testing"
)

func Test(t *testing.T) {
	conf.Init()
	println(conf.MongoDBAddr)

}
