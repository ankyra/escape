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
	"time"

	. "github.com/ankyra/escape/model/interfaces"
	"github.com/ankyra/escape/model/inventory/types"
)

type LoginController struct{}

func (LoginController) Login(context Context, url, username, password string, insecureSkipVerify bool) error {
	context.GetEscapeConfig().GetCurrentTarget().SetInsecureSkipVerify(insecureSkipVerify)
	authMethods, err := context.GetInventory().GetAuthMethods(url)
	if err != nil {
		return err
	}
	if authMethods == nil {
		fmt.Printf("Authentication not required.\n\nSuccessfully logged in to %s\n", url)
		context.GetEscapeConfig().GetCurrentTarget().SetAuthToken("")
		context.GetEscapeConfig().GetCurrentTarget().SetApiServer(url)
		return context.GetEscapeConfig().Save()
	}

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
	reader := bufio.NewReader(os.Stdin)

	ix := -1
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
	method := methods[ix-1]
	if method.Type == "oauth" {
		openBrowser(method.URL)
		return getEscapeTokenWithRedeemToken(context, url, method.RedeemToken, method.RedeemURL)
	} else if method.Type == "secret-token" {
		if username == "" {
			fmt.Printf("Username: ")
			username, err = reader.ReadString('\n')
			if err != nil {
				return err
			}
			username = strings.TrimSpace(username)
		}
		if password == "" {
			fmt.Printf("Password: ")
			password, _ = reader.ReadString('\n')
			if err != nil {
				return err
			}
			password = strings.TrimSpace(password)
		}
		authToken, err := context.GetInventory().LoginWithSecretToken(method.URL, username, password)
		if err != nil {
			return err
		}
		context.GetEscapeConfig().GetCurrentTarget().SetAuthToken(authToken)
		context.GetEscapeConfig().GetCurrentTarget().SetApiServer(url)
		context.GetEscapeConfig().Save()
		fmt.Printf("\nSuccessfully retrieved and stored auth token %s\n", authToken)
	} else {
		return fmt.Errorf("Unknown auth method.")
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
			context.GetEscapeConfig().GetCurrentTarget().SetAuthToken(string(authToken))
			context.GetEscapeConfig().GetCurrentTarget().SetApiServer(url)
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
