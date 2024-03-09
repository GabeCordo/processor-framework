package clusters

import (
	"github.com/GabeCordo/processor-framework/processor/components/cluster"
)

type HelloCluster struct {
}

func (cluster HelloCluster) ExtractFunc(h cluster.H, m cluster.M, out cluster.Out) {

	out.Push("hello")
}

func (cluster HelloCluster) TransformFunc(h cluster.H, m cluster.M, in any) (out any, success bool) {

	message := (in).(string)
	message += " world"

	return message, true
}

func (cluster HelloCluster) LoadFunc(h cluster.H, m cluster.M, in any) {

	message := (in).(string)
	h.Log(message)
}
