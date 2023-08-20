package threads

import (
	"github.com/GabeCordo/mango-go/processor/threads/common"
	"github.com/GabeCordo/mango-go/processor/threads/http"
	"github.com/GabeCordo/mango-go/processor/threads/provisioner"
	"github.com/GabeCordo/mango/core"
	"github.com/GabeCordo/mango/threads"
	"github.com/GabeCordo/mango/utils"
)

var DefaultProcessorFolder = core.DefaultFrameworkFolder + "go/"
var DefaultProcessorConfig = DefaultProcessorFolder + "processor.etl.yaml"
var DefaultModulesFolder = DefaultProcessorFolder + "modules/"

type Processor struct {
	HttpThread  *http.Thread
	Provisioner *provisioner.Thread

	Interrupt chan threads.InterruptEvent
	C1        chan common.ProvisionerRequest
	C2        chan common.ProvisionerResponse

	Logger *utils.Logger
}
