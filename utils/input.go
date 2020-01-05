package utils

import (
	"bufio"
	"fmt"
	"os"
)

// GetInput Gets input from user until '\n'. It accepts a prompt message as argument
func GetInput(prompt string) (input string, err error) {
	fmt.Printf(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	inputStr := scanner.Text()

	return inputStr, scanner.Err()
}
