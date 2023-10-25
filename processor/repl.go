package processor

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/GabeCordo/keitt/processor/components/provisioner"
	"github.com/GabeCordo/keitt/processor/threads/common"
	"github.com/GabeCordo/toolchain/logging"
	"github.com/GabeCordo/toolchain/multithreaded"
	"math/rand"
	"os"
	"strings"
)

// we can execute clusters in a batch mode with the
// [module].[cluster] [key]:[value] [key]:[value] binding

var BadSyntax = errors.New("input must match the format [module].[cluster] [key]:[value] [key]:[value]")

func (processor *Processor) repl() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(logging.Orange + "@" + logging.Green + "keitt " + logging.Reset)

		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		if input == "" {
			// do nothing
		} else if input == "help" {
			processor.help()
		} else if input == "info" {
			processor.info()
		} else if input == "supervisors" {
			processor.supervisors()
		} else if input == "stop" {
			processor.Interrupt <- common.Shutdown
		} else {
			module, cluster, metadata, err := processor.parseInput(input)
			if err == nil {
				processor.executeCluster(module, cluster, metadata)
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}

func (processor *Processor) help() {
	fmt.Printf("%s[module].[cluster] [key]:[value] [key]:[value] binding%s\n", logging.Gray, logging.Reset)
}

func (processor *Processor) info() {

	request := common.ProvisionerRequest{
		Action: common.ProvisionerModuleGet,
		Nonce:  rand.Uint32(),
	}
	processor.C1 <- request

	rsp, didTimeout := multithreaded.SendAndWait(processor.HttpThread.ProvisionerResponseTable, request.Nonce, processor.Config.Timeout)

	if didTimeout {
		fmt.Println("failed to load attached modules")
		return
	}

	response := (rsp).(common.ProvisionerResponse)
	modules := (response.Data).([]*provisioner.ModuleWrapper)

	for _, module := range modules {
		fmt.Printf("├─ %s\n", module.Identifier)

		for _, cluster := range module.GetClusters() {

			fmt.Printf("|  ├─%s\n", cluster.Identifier)
		}
	}
}

func (processor *Processor) supervisors() {

	request := common.ProvisionerRequest{
		Action: common.ProvisionerModuleGet,
		Nonce:  rand.Uint32(),
	}
	processor.C1 <- request

	rsp, didTimeout := multithreaded.SendAndWait(processor.HttpThread.ProvisionerResponseTable, request.Nonce, processor.Config.Timeout)

	if didTimeout {
		fmt.Println("failed to load supervisor data")
		return
	}

	response := (rsp).(common.ProvisionerResponse)
	modules := (response.Data).([]*provisioner.ModuleWrapper)

	for _, module := range modules {

		fmt.Printf("├─ %s\n", module.Identifier)
		for _, cluster := range module.GetClusters() {

			fmt.Printf("|   ├─ %s\n", cluster.Identifier)

			for _, instance := range cluster.FindSupervisors() {

				fmt.Printf("|   |   ├─ %d (state: %s)\n", instance.Id, instance.State.ToString())
			}
		}
	}
}

func (processor *Processor) parseInput(input string) (module, cluster string, metadata map[string]string, err error) {

	metadata = make(map[string]string)

	i := strings.SplitN(input, " ", 2)

	switch numOfParams := len(i); {
	case numOfParams >= 2:
		key := ""
		value := ""
		inKey := true
		inQuotation := false
		for _, c := range []rune(i[1]) {
			if (c == ':') && inKey {
				inKey = false
			} else if (c == ':') && !inKey {
				err = BadSyntax
				break
			} else if (c == '"') && !inKey {
				inQuotation = !inQuotation
			} else if (c == '"') && inKey {
				err = BadSyntax
				break
			} else if ((c == ' ') && !inQuotation) || (c == '\n') {
				inKey = true
				metadata[key] = value
				key = ""
				value = ""
			} else if (c == ' ') && inKey {
				err = BadSyntax
				break
			} else if inKey {
				key += string(c)
			} else {
				value += string(c)
			}
		}
		if len(key) > 0 {
			metadata[key] = value
		}
		fallthrough
	case numOfParams >= 1:
		j := strings.Split(i[0], ".")
		if len(j) == 2 {
			module = j[0]
			cluster = j[1]
		} else {
			err = BadSyntax
		}
	case numOfParams == 0:
		err = BadSyntax
	}

	return module, cluster, metadata, err
}

func (processor *Processor) executeCluster(module, cluster string, metadata map[string]string) {

	metadataStr := formatMap(metadata)
	fmt.Printf("%smodule: %s, cluster: %s metadata: %s%s\n",
		logging.Gray, module, cluster, metadataStr, logging.Reset)

	request := common.ProvisionerRequest{
		Action:   common.ProvisionerSupervisorCreate,
		Source:   common.User,
		Module:   module,
		Cluster:  cluster,
		Metadata: metadata,
		Nonce:    rand.Uint32(),
	}
	processor.C1 <- request
}

func formatMap(data map[string]string) string {
	output := "{"
	idx := 0
	for key, value := range data {
		postfix := ","
		if (len(data) - 1) == idx {
			postfix = ""
		}
		output += fmt.Sprintf("%s:%s%s", key, value, postfix)
		idx++
	}
	output += "}"
	return output
}
