package clusters

import (
	"fmt"
	"github.com/GabeCordo/keitt/processor/components/channel"
	"github.com/GabeCordo/keitt/processor/components/cluster"
	"github.com/GabeCordo/keitt/processor/components/helper"
)

type MetaDataCluster struct {
	helper helper.Helper
}

func (mdc *MetaDataCluster) SetHelper(helper helper.Helper) {
	mdc.helper = helper
}

func (mdc *MetaDataCluster) ExtractFunc(m cluster.M, c channel.OneWay) {

	key := m.GetKey("test")

	if key == "" {
		fmt.Println("key was not passed successfully")
	} else {
		c.Push(key)
	}
}

func (mdc *MetaDataCluster) TransformFunc(m cluster.M, in any) (out any, success bool) {
	return in, true
}

func (mdc *MetaDataCluster) LoadFunc(m cluster.M, in any) {
	key := (in).(string)
	mdc.helper.Log(key)
}
