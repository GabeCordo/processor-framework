package processor

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/GabeCordo/mango-go/processor/threads/common"
	"github.com/GabeCordo/mango/threads"
	"github.com/GabeCordo/mango/utils"
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
		fmt.Print(utils.Orange + "@" + utils.Green + "mango " + utils.Reset)
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		if input == "help" {
			processor.help()
		} else if input == "stop" {
			processor.Interrupt <- threads.Shutdown
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
	fmt.Printf("%s[module].[cluster] [key]:[value] [key]:[value] binding%s\n", utils.Gray, utils.Reset)
}

func (processor *Processor) parseInput(input string) (module, cluster string, metadata map[string]string, err error) {

	metadata = make(map[string]string)

	input = strings.Replace(input, "\n", "", -1)

	i := strings.Split(input, " ")

	switch numOfParams := len(i); {
	case numOfParams >= 2:
		i := i[1:]
		for _, pair := range i {
			k := strings.Split(pair, ":")
			if len(k) != 2 {
				err = BadSyntax
				break
			} else {
				metadata[k[0]] = k[1]
			}
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
		utils.Gray, module, cluster, metadataStr, utils.Reset)

	request := common.ProvisionerRequest{
		Action:   common.ProvisionerCreateSupervisor,
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
		var postfix string
		if idx != (len(data) - 1) {
			postfix = fmt.Sprint(",")
		} else {
			postfix = fmt.Sprint("}")
		}
		output += fmt.Sprintf("%s:%s%s", key, value, postfix)
		idx++
	}
	return output
}
