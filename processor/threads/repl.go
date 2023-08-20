package threads

import (
	"bufio"
	"fmt"
	"os"
)

func repl() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)
	}
}
