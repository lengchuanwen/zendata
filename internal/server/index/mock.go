package index

import (
	"github.com/easysoft/zendata/internal/server/controller"
	"github.com/easysoft/zendata/internal/server/core/module"
	"github.com/kataras/iris/v12"
)

type MockModule struct {
	MockCtrl *controller.MockCtrl `inject:""`
}

func NewMockModule() *DataModule {
	return &DataModule{}
}

// Party 执行
func (m *MockModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Get("/{paths:path}", m.MockCtrl.Mock).Name = "API服务模拟"
	}

	return module.NewModule("/", handler)
}
