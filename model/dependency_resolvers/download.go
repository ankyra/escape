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

package dependency_resolvers

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ankyra/escape-client/util"
	core "github.com/ankyra/escape-core"
)

var downloadLogger util.Logger

func DoDownloads(downloads []*core.DownloadConfig, logger util.Logger) error {
	downloadLogger = logger
	for _, download := range downloads {
		if err := DoDownload(download); err != nil {
			return err
		}
	}
	return nil
}

func DoDownload(download *core.DownloadConfig) error {
	if !download.OverwriteExistingDest {
		if util.PathExists(download.Dest) {
			downloadLogger.Log("download.skip_overwrite", map[string]string{
				"URL":  download.URL,
				"dest": download.Dest,
			})
			return nil
		}
	}
	if download.Platform != "" && download.Platform != runtime.GOOS {
		downloadLogger.Log("download.skip_platform", map[string]string{
			"URL":      download.URL,
			"platform": download.Platform,
			"actual":   runtime.GOOS,
		})
		return nil
	}
	if download.Arch != "" && download.Arch != runtime.GOARCH {
		downloadLogger.Log("download.skip_arch", map[string]string{
			"URL":    download.URL,
			"arch":   download.Arch,
			"actual": runtime.GOARCH,
		})
		return nil
	}
	dir, _ := filepath.Split(download.Dest)
	if dir != "" {
		if err := util.MkdirRecursively(dir); err != nil {
			return fmt.Errorf("Failed to make directory '%s' for zip file '%s': %s", dir, download.Dest, err.Error())
		}
	}
	out, err := os.Create(download.Dest)
	if err != nil {
		return fmt.Errorf("Couldn't open download destination '%s': %s", download.Dest, err.Error())
	}
	defer out.Close()

	downloadLogger.Log("download.start", map[string]string{
		"URL":  download.URL,
		"dest": download.Dest,
	})
	resp, err := http.Get(download.URL)
	if err != nil {
		return fmt.Errorf("Couldn't download '%s': %s", download.URL, err.Error())
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("Couldn't write '%s' to '%s': %s", download.URL, download.Dest, err.Error())
	}

	downloadLogger.Log("download.finished", map[string]string{
		"URL": download.URL,
	})
	return DoUnpack(download, dir)
}

func DoUnpack(download *core.DownloadConfig, targetDir string) error {
	if !download.Unpack {
		return nil
	}
	downloadLogger.Log("download.unpack", map[string]string{
		"dest": download.Dest,
	})
	if strings.HasSuffix(download.Dest, ".zip") {
		return UnpackZipFile(download.Dest, targetDir)
	} else if strings.HasSuffix(download.Dest, ".tgz") {
		return UnpackTgzFile(download.Dest, targetDir)
	} else if strings.HasSuffix(download.Dest, ".tar.gz") {
		return UnpackTgzFile(download.Dest, targetDir)
	} else if strings.HasSuffix(download.Dest, ".tar") {
		return UnpackTarFile(download.Dest, targetDir)
	}
	return fmt.Errorf("Can't unpack: unknown archive extension '%s'", download.Dest)
}

func UnpackZipFile(zipFile, targetDir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("Couldn't open '%s': %s", zipFile, err.Error())
	}
	defer r.Close()
	for _, f := range r.File {
		name := filepath.Join(targetDir, f.Name)
		dir, _ := filepath.Split(name)
		if err := util.MkdirRecursively(dir); err != nil {
			return fmt.Errorf("Failed to make directory '%s' for zip file '%s': %s", dir, name, err.Error())
		}
		out, err := os.Create(name)
		if err != nil {
			return fmt.Errorf("Couldn't create file '%s' for zip file '%s': %s", name, zipFile, err.Error())
		}
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("Couldn't open '%s' in zip file '%s': %s", name, zipFile, err.Error())
		}
		_, err = io.Copy(out, rc)
		if err != nil {
			return fmt.Errorf("Couldn't copy '%s' from zip file '%s': %s", name, zipFile, err.Error())
		}
		rc.Close()
		os.Chmod(name, f.Mode())
	}
	downloadLogger.Log("download.unpack_finished", map[string]string{
		"file": zipFile,
	})
	return nil
}

func UnpackTgzFile(file, targetDir string) error {
	fp, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Couldn't open archive '%s': %s", file, err.Error())
	}
	defer fp.Close()

	gzf, err := gzip.NewReader(fp)
	if err != nil {
		return fmt.Errorf("Couldn't read gzip archive '%s': %s", file, err.Error())
	}
	defer gzf.Close()
	return UnpackTar(file, targetDir, gzf)
}

func UnpackTarFile(file, targetDir string) error {
	fp, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Couldn't open archive '%s': %s", file, err.Error())
	}
	defer fp.Close()
	return UnpackTar(file, targetDir, fp)
}

func UnpackTar(file, targetDir string, reader io.Reader) error {
	tarReader := tar.NewReader(reader)
	if err := UnpackTarReader(tarReader, targetDir); err != nil {
		return fmt.Errorf("Failed to unpack '%s': %s", file, err.Error())
	}
	downloadLogger.Log("download.unpack_finished", map[string]string{
		"file": file,
	})
	return nil
}
