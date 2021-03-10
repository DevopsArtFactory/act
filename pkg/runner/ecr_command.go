package runner

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/DevopsArtFactory/act/pkg/constants"
)

// createEcrLoginCommand creates ecr login command
func createEcrLoginCommand(token, endpoint string) (string, error) {
	decodedToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return constants.EmptyString, err
	}

	username, password, err := getLoginInformationFromToken(string(decodedToken))
	if err != nil {
		return constants.EmptyString, err
	}

	return makeDockerLoginCommand(username, password, endpoint), nil
}

// getLoginInformationFromToken retrieves login information
func getLoginInformationFromToken(token string) (string, string, error) {
	splited := strings.Split(token, ":")
	if len(splited) != 2 {
		return constants.EmptyString, constants.EmptyString, fmt.Errorf("token is wrong")
	}

	return splited[0], splited[1], nil
}

// makeDockerLoginCommand makes docker login command
func makeDockerLoginCommand(username, password, endpoint string) string {
	return fmt.Sprintf("docker login -u %s -p %s %s", username, password, endpoint)
}
