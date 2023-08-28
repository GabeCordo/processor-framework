package main

import (
	"github.com/GabeCordo/commandline"
	"github.com/GabeCordo/keitt/controllers"
)

func main() {

	cli := commandline.NewCommandLine()

	mc := cli.AddCommand("module", controllers.ModuleCommand{})
	mc.SetCategory("modules").SetDescription("used to install or delete a module from the global space")

	dc := cli.AddCommand("doctor", controllers.DoctorCommand{})
	dc.SetCategory("utils").SetDescription("used to verify the global processor configuration")

	ic := cli.AddCommand("init", controllers.InitCommand{})
	ic.SetCategory("utils").SetDescription("used to initialize the global processor configuration")

	sc := cli.AddCommand("start", controllers.StartCommand{})
	sc.SetCategory("execution").SetDescription("used to start the processor in standalone mode")

	cc := cli.AddCommand("connect", controllers.ConnectCommand{})
	cc.SetCategory("execution").SetDescription("used to connect to a core and run connected mode")

	cli.Run()
}
