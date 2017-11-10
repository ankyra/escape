/*
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	. "github.com/ankyra/escape/model/interfaces"
	"github.com/ankyra/escape/model/inventory/types"
	"golang.org/x/crypto/ssh/terminal"
)

type LoginController struct{}

func (LoginController) Login(context Context, url, authMethodRequested, username, password string, insecureSkipVerify bool) error {
	context.GetEscapeConfig().GetCurrentProfile().SetInsecureSkipVerify(insecureSkipVerify)
	authMethods, err := context.GetInventory().GetAuthMethods(url)
	if err != nil {
		return err
	}

	if authMethods == nil {
		fmt.Printf("Authentication not required.\n\nSuccessfully logged in to %s\n", url)
		context.GetEscapeConfig().GetCurrentProfile().SetAuthToken("")
		context.GetEscapeConfig().GetCurrentProfile().SetApiServer(url)
		return context.GetEscapeConfig().Save()
	}

	reader := bufio.NewReader(os.Stdin)
	if username != "" && authMethods["service-account"] != nil {
		return secretTokenAuth(reader, context, authMethods["service-account"].URL, username, password)
	}

	var authMethod *types.AuthMethod

	if authMethodRequested != "" {
		authMethod = authMethods[authMethodRequested]
		if authMethod == nil {
			for _, availableAuthMethod := range authMethods {
				if availableAuthMethod.Type == authMethodRequested {
					authMethod = availableAuthMethod
				}
			}
		}
	}

	if authMethod == nil {
		authMethod = authUserSelection(reader, authMethods)
	}
	if authMethod.Type == "oauth" {
		openBrowser(authMethod.URL)
		return getEscapeTokenWithRedeemToken(context, url, authMethod.RedeemToken, authMethod.RedeemURL)
	} else if authMethod.Type == "secret-token" {
		return secretTokenAuth(reader, context, authMethod.URL, username, password)
	} else {
		return fmt.Errorf("Unknown auth method.")
	}
	return nil
}

func authUserSelection(reader *bufio.Reader, authMethods map[string]*types.AuthMethod) *types.AuthMethod {
	sortedKeys := sortAuthMethodMapKeys(authMethods)

	sortedAuthMethods := []*types.AuthMethod{}
	for _, key := range sortedKeys {
		sortedAuthMethods = append(sortedAuthMethods, authMethods[key])
	}

	fmt.Println("Available authentication methods:\n")
	i := 1
	methods := []*types.AuthMethod{}
	for key, authMethod := range sortedAuthMethods {
		fmt.Printf(" %d. %s [%s]\n", i, sortedKeys[key], authMethod.Type)
		methods = append(methods, authMethod)
		i += 1
	}

	ix := -1
	var err error

	for {
		fmt.Printf("\nPlease select an authentication method (1-%d): ", i-1)
		requestedMethod, _ := reader.ReadString('\n')
		ix, err = strconv.Atoi(strings.TrimSpace(requestedMethod))
		if err != nil {
			fmt.Println("Not a number.")
			continue
		}
		if ix < 1 || ix > i-1 {
			fmt.Println("Number out of range.")
			continue
		} else {
			break
		}
	}

	return methods[ix-1]
}

func secretTokenAuth(reader *bufio.Reader, context Context, url, username, password string) error {
	err := credentialsUserInput(reader, &username, &password)
	if err != nil {
		return err
	}
	authToken, err := context.GetInventory().LoginWithSecretToken(url, username, password)
	if err != nil {
		return err
	}
	context.GetEscapeConfig().GetCurrentProfile().SetAuthToken(authToken)
	context.GetEscapeConfig().GetCurrentProfile().SetApiServer(url)
	context.GetEscapeConfig().Save()
	fmt.Printf("\nSuccessfully retrieved and stored auth token %s\n", authToken)
	return nil
}

func credentialsUserInput(reader *bufio.Reader, username, password *string) error {
	var err error
	if *username == "" {
		fmt.Printf("Username: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		*username = strings.TrimSpace(input)
	}
	if *password == "" {
		fmt.Printf("Password: ")
		passwordBytes, _ := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		*password = strings.TrimSpace(string(passwordBytes))
	}

	return nil
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Unsupported platform. Please login manually at %s", url)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}

func getEscapeTokenWithRedeemToken(context Context, url, redeemToken, redeemURL string) error {

	currentTry := 0
	tries := 25
	timeOut := time.Duration(1)
	client := &http.Client{}
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	redeemURL += "?redeem-token=" + redeemToken
	for currentTry < tries {

		req, err := http.NewRequest("GET", redeemURL, nil)
		if err != nil {
			return fmt.Errorf("Couldn't retrieve token from server: %s", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("Couldn't retrieve token from server: %s", err)
		}
		if resp.StatusCode == 200 {
			authToken, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("Couldn't read response from server '%s': %s", redeemURL, resp.Status)
			}
			context.GetEscapeConfig().GetCurrentProfile().SetAuthToken(string(authToken))
			context.GetEscapeConfig().GetCurrentProfile().SetApiServer(url)
			context.GetEscapeConfig().Save()
			fmt.Printf("\nSuccessfully retrieved and stored auth token %s\n", authToken)
			return nil
		}
		if resp.StatusCode != 404 {
			return fmt.Errorf("Couldn't retrieve token from server. Got status code %d", resp.Status)
		}
		time.Sleep(timeOut * time.Second)
		currentTry++
		if currentTry == 5 {
			timeOut *= 2
		}
		if currentTry == 10 {
			timeOut *= 2
		}
	}
	return nil
}

func sortAuthMethodMapKeys(authMethods map[string]*types.AuthMethod) []string {
	sortedKeys := make([]string, len(authMethods))
	i := 0
	for k, _ := range authMethods {
		sortedKeys[i] = k
		i++
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}
