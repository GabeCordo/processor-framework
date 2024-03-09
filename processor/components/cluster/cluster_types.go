package cluster

import (
	"github.com/GabeCordo/processor-framework/processor/interfaces"
	"time"
)

type Segment int8

const (
	Extract   Segment = 0
	Transform         = 1
	Load              = 2
)

type OnCrash string

const (
	Restart   OnCrash = "Restart"
	DoNothing         = "DoNothing"
)

type OnLoad string

const (
	CompleteAndPush OnLoad = "CompleteAndPush"
	WaitAndPush            = "WaitAndPush"
)

type EtlMode string

const (
	Batch  EtlMode = "mode/batch"  // The cluster is provisioned when invoked by an operator or application.
	Stream         = "mode/stream" // The cluster is provisioned automatically when the system is started.
)

type Transformer func(in []any) (out any, success bool)

// M contains metadata about the running supervisor including any state
// information that a developer might need to interact with the Supervisor.
type M interface {
	GetKey(key string) string
}

type H interface {
	IsDebugEnabled() bool
	SaveToCache(data string) (string, error)
	LoadFromCache(identifier string) (string, error)
	Log(message string) error
	Logf(format string, data ...any) error
	Warning(message string) error
	Warningf(format string, data ...any) error
	Fatal(message string) error
	Fatalf(format string, data ...any) error
}

type Out interface {
	Push(any) bool
}

type Lambda func(helper H, metadata M) (out any, success bool)

// Cluster is a set of functions that define an ETL process
// cluster functions are provisioned on goroutines to run in parallel and
// process data.
type Cluster interface {
	ExtractFunc(helper H, metadata M, out Out)
	TransformFunc(helper H, metadata M, in any) (out any, success bool)
}

type LoadAll interface {
	LoadFunc(helper H, metadata M, in []any)
}

type LoadOne interface {
	LoadFunc(helper H, metadata M, in any)
}

type SystemFunctions interface {
	Setup(curr time.Time, h H)
	Teardown(curr time.Time, h H)
}

type VerifiableET interface {
	VerifyETFunction(in any) (valid bool)
}

type VerifiableTL interface {
	VerifyTLFunction(in any) (valid bool)
}

// Test
// TODO : needs to be implemented
type Test interface {
	MockExtractFunc(metadata M, out Out)
	VerifyTransformOutput(metadata M, in any) (success bool)
	MockLoadFunc(metadata M, in any)
}

type Config struct {
	Identifier                  string  `json:"identifier" yaml:"identifier"`
	OnLoad                      OnLoad  `json:"on-load" yaml:"on-load"`
	OnCrash                     OnCrash `json:"on-crash" yaml:"on-crash"`
	StartWithNTransformClusters int     `json:"start-with-n-t-channels" yaml:"start-with-n-t-clusters"`
	StartWithNLoadClusters      int     `json:"start-with-n-l-channels" yaml:"start-with-n-l-clusters"`
	ETChannelThreshold          int     `json:"et-channel-threshold" yaml:"et-channel-threshold"`
	ETChannelGrowthFactor       int     `json:"et-channel-growth-factor" yaml:"et-channel-growth-factor"`
	TLChannelThreshold          int     `json:"tl-channel-threshold" yaml:"tl-channel-threshold"`
	TLChannelGrowthFactor       int     `json:"tl-channel-growth-factor" yaml:"tl-channel-growth-factor"`
}

func (config Config) ToStandard() *interfaces.Config {

	dst := new(interfaces.Config)

	dst.Identifier = config.Identifier
	dst.OnLoad = interfaces.OnLoad(config.OnLoad)
	dst.OnCrash = interfaces.OnCrash(config.OnCrash)
	dst.StartWithNTransformClusters = config.StartWithNTransformClusters
	dst.StartWithNLoadClusters = config.StartWithNLoadClusters
	dst.ETChannelThreshold = config.ETChannelThreshold
	dst.ETChannelGrowthFactor = config.ETChannelGrowthFactor
	dst.TLChannelThreshold = config.TLChannelThreshold
	dst.TLChannelGrowthFactor = config.TLChannelGrowthFactor

	return dst
}

type DataTiming struct {
	ETIn  time.Time
	ETOut time.Time
	TLIn  time.Time
	TLOut time.Time
}

type TimingStatistics struct {
	MinTimeBeforePop time.Duration `json:"min-time-before-pop-ns"`
	MaxTimeBeforePop time.Duration `json:"max-time-before-pop-ns"`
	AverageTime      time.Duration `json:"average-time-ns"`
	MedianTime       time.Duration `json:"median-time-ns"`
}

func (ts TimingStatistics) ToStandard() *interfaces.TimingStatistics {
	standard := new(interfaces.TimingStatistics)

	standard.MinTimeBeforePop = ts.MinTimeBeforePop
	standard.MaxTimeBeforePop = ts.MaxTimeBeforePop
	standard.MedianTime = ts.MedianTime
	standard.AverageTime = ts.AverageTime

	return standard
}

type Statistics struct {
	Threads struct {
		NumProvisionedExtractRoutines int `json:"num-provisioned-extract-routines"`
		NumProvisionedTransformRoutes int `json:"num-provisioned-transform-routes"`
		NumProvisionedLoadRoutines    int `json:"num-provisioned-load-routines"`
	} `json:"threads"`
	Channels struct {
		NumEtThresholdBreaches int `json:"num-et-threshold-breaches"`
		NumTlThresholdBreaches int `json:"num-tl-threshold-breaches"`
	} `json:"channels"`
	Data struct {
		TotalProcessed     int `json:"total-processed"`
		TotalOverETChannel int `json:"total-over-et"`
		TotalOverTLChannel int `json:"total-over-tl"`
		TotalDropped       int `json:"total-dropped"`
	} `json:"data"`
	Timing struct {
		ET               TimingStatistics `json:"et-channel"`
		etSet            bool
		TL               TimingStatistics `json:"tl-channel"`
		tlSet            bool
		MaxTotalTime     time.Duration `json:"max-total-time-ns"`
		MinTotalTime     time.Duration `json:"min-total-time-ns"`
		AverageTotalTime time.Duration `json:"avg-total-time-ns"`
		MedianTotalTime  time.Duration `json:"med-total-time-ns"`
		totalSet         bool
	} `json:"timing"`
}

func (statistics *Statistics) ToStandard() *interfaces.Statistics {

	standard := new(interfaces.Statistics)

	standard.Threads = statistics.Threads
	standard.Channels = statistics.Channels
	standard.Data = statistics.Data
	standard.Timing.TL = *statistics.Timing.TL.ToStandard()
	standard.Timing.ET = *statistics.Timing.ET.ToStandard()
	standard.Timing.MinTotalTime = statistics.Timing.MinTotalTime
	standard.Timing.MaxTotalTime = statistics.Timing.MaxTotalTime
	standard.Timing.MedianTotalTime = statistics.Timing.MedianTotalTime
	standard.Timing.AverageTotalTime = statistics.Timing.AverageTotalTime

	return standard
}

type Status uint8

const (
	Registered = iota
	UnMounted
	Mounted
	InUse
	MarkedForDeletion
)

type Event uint8

const (
	Register = iota
	Mount
	UnMount
	Use
	Delete
)

type Response struct {
	Config     Config        `json:"core"`
	Stats      *Statistics   `json:"stats"`
	LapsedTime time.Duration `json:"lapsed-time"`
	DidItCrash bool          `json:"crashed"`
}
