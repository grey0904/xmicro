package app

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"reflect"
	"runtime"
)

type Application struct {
	Redis       *redis.Client
	DB          *gorm.DB
	NacosClient config_client.IConfigClient
}

var App = new(Application)

// Startup ..
func (app *Application) Startup(fns ...func() error) error {

	//app.initialize()

	return SerialUntilError(fns...)()
}

// 创建一个迭代器
func SerialUntilError(fns ...func() error) func() error {
	return func() error {
		for _, fn := range fns {
			if err := try(fn, nil); err != nil {
				return err
				// return errors.Wrap(err, xstring.FunctionName(fn))
			}
		}
		return nil
	}
}

func try(fn func() error, cleaner func()) (ret error) {
	if cleaner != nil {
		defer cleaner()
	}
	defer func() {
		if err := recover(); err != nil {
			_, _, line, _ := runtime.Caller(2)

			if _, ok := err.(error); ok {
				ret = err.(error)
			} else {
				ret = fmt.Errorf("%+v", err)
			}
			ret = errors.Wrap(ret, fmt.Sprintf("%s:%d", FunctionName(fn), line))
		}
	}()
	return fn()
}

func FunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
