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
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/registry/types"
)

type LoginController struct{}

func (LoginController) Login(context Context, url, username, password string) error {
	authMethods, err := context.GetRegistry().GetAuthMethods(url)
	if err != nil {
		return err
	}
	if authMethods == nil {
		fmt.Printf("Registry at %s does not implement authentication.", url)
		context.GetEscapeConfig().GetCurrentTarget().SetAuthToken("")
		context.GetEscapeConfig().GetCurrentTarget().SetApiServer(url)
		return context.GetEscapeConfig().Save()
	}
	fmt.Println("Available authentication methods:\n")
	i := 1
	methods := []*types.AuthMethod{}
	for key, authMethod := range authMethods {
		fmt.Printf(" %d. %s [%s]\n", i, key, authMethod.Type)
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
		return nil
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
		authToken, err := context.GetRegistry().LoginWithSecretToken(method.URL, username, password)
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
