package server

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"
)

func stringToReply(message string) string {
	return fmt.Sprintf("+%s\r\n", message)
}

func makeErrorResponse(message string) string {
	return fmt.Sprintf("-%s\r\n", message)
}

func makeNullValueResponse() string {
	return "$-1\r\n"
}

func isReplArray(command string) bool {
	return strings.HasPrefix(command, "*")
}

func extractArgumentsFromReplArray(command string) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(command))
	scanner.Scan()
	numArgs, err := stringToInt(scanner.Text()[1:])

	if err != nil {
		log.Println("Failed to convert numArgs to int")
		return nil, errors.New("failed to convert numArgs to int")
	}

	arguments := make([]string, numArgs)

	for i := 0; i < numArgs; i++ {
		scanner.Scan()
		argLength, err := stringToInt(scanner.Text()[1:])

		if err != nil {
			log.Println("Failed to convert numArgs to int")
			return nil, errors.New("failed to convert numArgs to int")
		}

		scanner.Scan()
		arguments[i] = scanner.Text()[:argLength]
	}

	return arguments, nil
}
