package functions

import (
	"fmt"
	"strings"
)

func ParseCommand(cmdLine string) (cmd, param string) {
	parts := strings.Split(cmdLine, " ")
	if len(parts) < 1 {
		return "", ""
	}

	if len(parts) == 1 {
		cmd = strings.TrimSpace(parts[0])
		return cmd, ""
	}

	cmd = strings.TrimSpace(parts[0])
	param = strings.Join(parts[1:], " ")

	fmt.Println("PARSE FUNC command: ", cmd)
	// fmt.Println("PARSE FUNC parameters: ", param)
	return cmd, param
}
