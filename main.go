package main

import (
	"github.com/GabeCordo/commandline"
	"github.com/GabeCordo/processor-framework/controllers"
)

func main() {

	cli := commandline.NewCommandLine()

	dc := cli.AddCommand("doctor", controllers.DoctorCommand{})
	dc.SetCategory("utils").SetDescription("used to verify the global processor configuration")

	ic := cli.AddCommand("init", controllers.InitCommand{})
	ic.SetCategory("utils").SetDescription("used to initialize the global processor configuration")

	sc := cli.AddCommand("start", controllers.StartCommand{})
	sc.SetCategory("runtime").SetDescription("used to start the processor in standalone mode")

	cc := cli.AddCommand("connect", controllers.ConnectCommand{})
	cc.SetCategory("runtime").SetDescription("used to connect to a core and run connected mode")

	cli.Run()
}
