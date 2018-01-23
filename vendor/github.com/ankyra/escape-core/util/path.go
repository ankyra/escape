/*
Copyright 2017, 2018 Ankyra

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

package util

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsDir(path string) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	return st.IsDir()
}

func RemoveTree(path string) error {
	return os.RemoveAll(path)
}

func MkdirRecursively(path string) error {
	return os.MkdirAll(path, 0755)
}

func CopyFile(src, dst string) error {
	st, err := os.Stat(src)
	if err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	tmp, err := ioutil.TempFile(filepath.Dir(dst), "")
	if err != nil {
		return err
	}

	if _, err = io.Copy(tmp, in); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}

	if err = tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}

	if err = os.Chmod(tmp.Name(), st.Mode()); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), dst)
}
