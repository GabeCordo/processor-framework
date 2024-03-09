package supervisor

import (
	"fmt"
	"github.com/GabeCordo/processor-framework/processor/components/channel"
	"github.com/GabeCordo/processor-framework/processor/components/cluster"
	"log"
	"time"
)

const (
	DefaultNumberOfClusters       = 1
	DefaultMonitorRefreshDuration = 1
	DefaultChannelThreshold       = 10
	DefaultChannelGrowthFactor    = 2
)

func (supervisor *Supervisor) Event(event Event) bool {
	supervisor.mutex.Lock()
	defer supervisor.mutex.Unlock()

	if supervisor.State == UnTouched {
		if event == Startup {
			supervisor.State = Running
		} else if (event == Suspend) || (event == TearedDown) {
			supervisor.State = Stopping
		} else {
			return false
		}
	} else if supervisor.State == Running {
		if event == StartProvision {
			supervisor.State = Provisioning
		} else if event == Error {
			supervisor.State = Failed
		} else if event == Suspend {
			supervisor.State = Stopping
		} else if event == TearedDown {
			supervisor.State = Terminated
		} else {
			return false
		}
	} else if supervisor.State == Provisioning {
		if event == EndProvision {
			supervisor.State = Running
		} else if event == Error {
			supervisor.State = Failed
		} else if event == Suspend {
			supervisor.State = Stopping
		} else {
			return false
		}
	} else if supervisor.State == Stopping {
		if event == TearedDown {
			supervisor.State = Terminated
		} else {
			return false
		}
	} else if (supervisor.State == Failed) || (supervisor.State == Terminated) {
		return false
	}

	return true // represents a boolean ~ hasStateChanged?
}

func (supervisor *Supervisor) IsAlive() bool {

	supervisor.mutex.RLock()
	defer supervisor.mutex.RUnlock()

	return (supervisor.State != Failed) && (supervisor.State != Terminated)
}

func (supervisor *Supervisor) Start() (response *cluster.Response) {
	supervisor.Event(Startup)

	defer supervisor.Event(TearedDown)

	defer func() {
		// has the user defined function crashed during runtime?
		if r := recover(); r != nil {
			// yes => return a response that identifies that the cluster crashed
			response = cluster.NewResponse(
				supervisor.Config,
				supervisor.Stats,
				time.Now().Sub(supervisor.StartTime),
				true,
			)
		}
	}()

	supervisor.StartTime = time.Now()

	if sysFunc, ok := (supervisor.group).(cluster.SystemFunctions); ok {
		sysFunc.Setup(supervisor.StartTime, supervisor.helper)
	}

	//// start creating the default frontend goroutines
	supervisor.Provision(cluster.Extract)

	// the common specifies the number of transform functions to start running in parallel
	for i := 0; i < supervisor.Config.StartWithNTransformClusters; i++ {
		supervisor.Provision(cluster.Transform)
	}

	// the common specifies the number of load functions to start running in parallel
	for i := 0; i < supervisor.Config.StartWithNLoadClusters; i++ {
		supervisor.Provision(cluster.Load)
	}

	//// end creating the default frontend goroutines

	// every N seconds we should check if the ETChannel or TLChannel is congested
	// and requires us to provision additional nodes
	go supervisor.Runtime()

	supervisor.waitGroup.Wait() // wait for the Extract-Transform-Load (ETL) Cycle to Complete

	// calculate the timings produced by data being fed across each of the channels
	supervisor.CalculateTiming()

	response = cluster.NewResponse(
		supervisor.Config,
		supervisor.Stats,
		time.Now().Sub(supervisor.StartTime),
		false,
	)

	return response
}

func (supervisor *Supervisor) Teardown() {

	if sysFunc, ok := (supervisor.group).(cluster.SystemFunctions); ok {
		sysFunc.Teardown(time.Now(), supervisor.helper)
	}

	supervisor.Event(Suspend)
}

func (supervisor *Supervisor) Runtime() {
	for {
		if supervisor.State == Terminated {
			break
		}

		supervisor.ETChannel.GetState()

		if (supervisor.State == Stopping) && supervisor.ETChannel.Accepting() {
			supervisor.ETChannel.StopPushes()
		}

		// is ETChannel congested?
		if supervisor.ETChannel.GetState() == channel.Congested {
			supervisor.Stats.Channels.NumEtThresholdBreaches++
			n := supervisor.Stats.Threads.NumProvisionedTransformRoutes
			for n > 0 {
				supervisor.Provision(cluster.Transform)
				n--
			}
			supervisor.Stats.Threads.NumProvisionedTransformRoutes *= supervisor.ETChannel.GetGrowthFactor()
		}

		supervisor.TLChannel.GetState()

		// TODO : I don't think this is needed
		//if (supervisor.State == Stopping) && supervisor.TLChannel.Accepting() {
		//	supervisor.TLChannel.StopPushes()
		//}

		// is TLChannel congested?
		if supervisor.TLChannel.GetState() == channel.Congested {
			supervisor.Stats.Channels.NumTlThresholdBreaches++
			n := supervisor.Stats.Threads.NumProvisionedLoadRoutines
			for n > 0 {
				supervisor.Provision(cluster.Load)
				n--
			}
			supervisor.Stats.Threads.NumProvisionedLoadRoutines *= supervisor.TLChannel.GetGrowthFactor()
		}

		// check if the channel is congested after DefaultMonitorRefreshDuration seconds
		time.Sleep(DefaultMonitorRefreshDuration * time.Second)
	}
}

func (supervisor *Supervisor) ExtractWrapper(h cluster.H, m cluster.M, out cluster.Out) <-chan struct{} {
	done := make(chan struct{})

	// the function always finishes till completion unless a direct shutdown is called on the server
	// which stops data collection from some source
	go func() {
		defer func() {
			done <- struct{}{}
			close(done)
		}()
		supervisor.group.ExtractFunc(h, m, out)
	}()
	return done
}

func (supervisor *Supervisor) ExtractShutdownWrapper() <-chan struct{} {
	done := make(chan struct{})

	// we need to create a separate goroutine otherwise it will block the current
	// thread from re-evaluating the select statement wherever the ExtractShutdownWrapper is called
	go func() {
		defer close(done)
		for {
			// the IsAlive clause ensures that once a supervisor is dead, we will not leak memory
			// with a forever-running goroutine
			if (supervisor.State == Stopping) || (!supervisor.IsAlive()) {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return done
}

func (supervisor *Supervisor) Provision(segment cluster.Segment) {
	supervisor.Event(StartProvision)
	defer supervisor.Event(EndProvision)

	go func() {
		switch segment {
		case cluster.Extract:
			{
				defer func() {
					if r := recover(); r != nil {
						log.Println(r)
						log.Println("cluster.Extract function raised error")
						supervisor.ETChannel.ProducerDone()
						supervisor.waitGroup.Done()
					}
				}()
				supervisor.Stats.Threads.NumProvisionedExtractRoutines++
				oneWayChannel, _ := cluster.NewOneWayManagedChannel(supervisor.ETChannel)

				supervisor.ETChannel.AddProducer()

				select {
				case <-supervisor.ExtractWrapper(supervisor.helper, supervisor.Metadata, oneWayChannel):
					break
				case <-supervisor.ExtractShutdownWrapper():
					fmt.Println("shutdown caused extract to finish early")
					break
				}
				supervisor.ETChannel.ProducerDone()
			}
		case cluster.Transform:
			{
				defer func() {
					if r := recover(); r != nil {
						log.Println("cluster.Transform function raised error")
						log.Println(r)
						supervisor.TLChannel.ProducerDone()
						supervisor.waitGroup.Done()
					}
				}()

				supervisor.Stats.Threads.NumProvisionedTransformRoutes++
				supervisor.TLChannel.AddProducer()

				for request := range supervisor.ETChannel.GetChannel() {

					// associates a TimeOut to the data being removed from the channel and decrements
					// the data counter for the current pipe
					supervisor.ETChannel.DataPopped(request.Id)

					supervisor.Stats.Data.TotalOverETChannel++

					if i, ok := (supervisor.group).(cluster.VerifiableET); ok && !i.VerifyETFunction(request) {
						continue
					}

					data, success := supervisor.group.TransformFunc(supervisor.helper, supervisor.Metadata, request.Data)
					if success {
						supervisor.Stats.Data.TotalOverTLChannel++
						supervisor.TLChannel.Push(channel.DataWrapper{Id: request.Id, Data: data})
					}
				}

				supervisor.TLChannel.ProducerDone()
			}
		case cluster.Load:
			{
				defer func() {
					if r := recover(); r != nil {
						log.Println(r)
						log.Println("cluster.Load function raised error")
						supervisor.waitGroup.Done()
					}
				}()

				supervisor.Stats.Threads.NumProvisionedLoadRoutines++

				aggregatedData := make([]any, 0)

				for request := range supervisor.TLChannel.GetChannel() {
					supervisor.Stats.Data.TotalProcessed++

					// associates a TimeOut to the data being removed from the channel and decrements
					// the data counter for the current pipe
					supervisor.TLChannel.DataPopped(request.Id)

					if i, ok := (supervisor.group).(cluster.VerifiableTL); ok && !i.VerifyTLFunction(request) {
						continue
					}

					if a, success := (supervisor.group).(cluster.LoadOne); success {
						a.LoadFunc(supervisor.helper, supervisor.Metadata, request.Data)
					} else if _, success := (supervisor.group).(cluster.LoadAll); success {
						aggregatedData = append(aggregatedData, request.Data)
					}
				}

				if a, success := (supervisor.group).(cluster.LoadAll); success {
					a.LoadFunc(supervisor.helper, supervisor.Metadata, aggregatedData)
				}
			}
		}

		// notify the wait group a process has completed ~ if all are finished we close the monitor
		supervisor.waitGroup.Done()
	}()

	// a new function (E, T, or L) is provisioned
	// we should inform the wait group that the supervisor isn't finished until the wg is done
	supervisor.waitGroup.Add(1)
}

func (supervisor *Supervisor) Deletable() bool {
	return (supervisor.State == Terminated) || (supervisor.State == Failed)
}

func (supervisor *Supervisor) CalculateTiming() {

	supervisor.Stats.Data.TotalDropped = supervisor.Stats.Data.TotalOverETChannel - supervisor.Stats.Data.TotalOverTLChannel
	supervisor.Stats.CalculateTiming(supervisor.ETChannel.Timestamps, supervisor.TLChannel.Timestamps)
}

func (supervisor *Supervisor) Print() {
	fmt.Printf("Id: %d\n", supervisor.Id)
	fmt.Printf("Cluster: %s\n", supervisor.Config.Identifier)
}

func (status Status) ToString() string {
	switch status {
	case UnTouched:
		return "UnTouched"
	case Running:
		return "Running"
	case Provisioning:
		return "Provisioning"
	case Failed:
		return "Failed"
	case Terminated:
		return "Terminated"
	default:
		return "None"
	}
}
