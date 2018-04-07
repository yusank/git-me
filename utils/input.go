package utils

import (
	"bufio"
	"fmt"
	"os"
)

var inputReader *bufio.Reader

// read user input
func ReadInput(output string) (op string) {
	fmt.Print(output)
	inputReader = bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Print(err)
		return op
	}

	op = input[:len(input)-1]
	return
}
