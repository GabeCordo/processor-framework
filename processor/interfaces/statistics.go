package interfaces

import "time"

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
