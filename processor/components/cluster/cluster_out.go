package cluster

import (
	"errors"
	"github.com/GabeCordo/processor-framework/processor/components/channel"
	"math/rand"
)

type OneWayManagedChannel struct {
	channel *channel.ManagedChannel
}

func NewOneWayManagedChannel(c *channel.ManagedChannel) (Out, error) {

	if c == nil {
		return nil, errors.New("ManagedChannel passed to NewONeWayManagedChannel was nil")
	}

	oneWayManagedChannel := new(OneWayManagedChannel)
	oneWayManagedChannel.channel = c

	return oneWayManagedChannel, nil
}

func (c OneWayManagedChannel) Push(data any) bool {

	didPush := c.channel.Push(channel.DataWrapper{Id: rand.Uint64(), Data: data})
	return didPush
}
