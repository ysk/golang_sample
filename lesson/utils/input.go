package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

func Input(msg string) (string, error) {
	fmt.Print(msg + ": ")
	if !scanner.Scan() {
		return "", fmt.Errorf("入力エラー")
	}
	return strings.TrimSpace(scanner.Text()), nil
}
