/*
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
*/

package def

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/pkg/output"
	"gopkg.in/yaml.v3"
)

var projectPlatformDir = ".platform"

// dirToTarGz converts contents of directory to tar.gz.
func dirToTarGz(pathTo string) (bytes.Buffer, error) {
	pathTo = filepath.Join(projectPlatformDir, pathTo)
	out := bytes.Buffer{}
	gw := gzip.NewWriter(&out)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	if err := filepath.Walk(pathTo, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			header, err := tar.FileInfoHeader(f, f.Name())
			if err != nil {
				return err
			}
			header.Name = strings.TrimLeft(strings.Replace(path, pathTo, "", 1), "/")
			header.Name = strings.ReplaceAll(header.Name, "\\", "/")
			err = tw.WriteHeader(header)
			if err != nil {
				return err
			}
			_, err = io.Copy(tw, file)
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	}); err != nil {
		return out, err
	}
	tw.Flush()
	gw.Flush()
	return out, nil
}

// dirToTarGzB64 converts contents of directory to tar.gz base64 encoded.
func dirToTarGzB64(path string) (string, error) {
	data, err := dirToTarGz(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data.Bytes()), nil
}

// unmarshalYamlValue gets value of current node.
func unmarshalYamlValue(value *yaml.Node) interface{} {
	switch value.Tag {
	case "!!map", "!!seq":
		{
			return unmarshalYamlWithCustomTags(value)
		}
	case "!archive":
		{
			out, err := dirToTarGzB64(value.Value)
			if err != nil {
				output.Warn(err.Error())
			}
			return out
		}
	case "!include":
		{
			path := ""
			if value.Value == "" {
				for i := range value.Content {
					// don't know what to do if type isn't string
					if value.Content[i].Value == "type" && value.Content[i+1].Value != "string" {
						return ""
					}
					if value.Content[i].Value == "path" {
						path = value.Content[i+1].Value
						break
					}
				}
			}
			if path == "" {
				return ""
			}
			path = filepath.Join(projectPlatformDir, path)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				output.Warn(err.Error())
				return ""
			}
			return string(data)
		}
	default:
		{
			var out interface{}
			if err := value.Decode(&out); err != nil {
				return nil
			}
			return out
		}
	}
}

// unmarshalYamlWithCustomTags unmarshals yaml with custom tags.
func unmarshalYamlWithCustomTags(value *yaml.Node) interface{} {
	switch value.Tag {
	case "!!map":
		{
			out := make(map[string]interface{})
			for i := range value.Content {
				if i%2 == 0 {
					out[value.Content[i].Value] = unmarshalYamlValue(value.Content[i+1])
				}
			}
			return out
		}
	case "!!seq":
		{
			out := make([]interface{}, 0)
			for _, child := range value.Content {
				out = append(out, unmarshalYamlValue(child))
			}
			return out
		}
	}
	return nil
}

// YamlMerge provides interface for creating yaml that is mergable with support for custom tags.
type YamlMerge map[string]interface{}

// UnmarshalYAML unmarshals YAML for app def.
func (m *YamlMerge) UnmarshalYAML(value *yaml.Node) error {
	*m = unmarshalYamlWithCustomTags(value).(map[string]interface{})
	return nil
}

// mergeYamls takes multiple yaml byte arrays and merge them and unmarshals in to given interface.
func mergeYamls(data [][]byte, def interface{}) error {
	// unmarshal yaml in to maps and merge the maps
	mapData := map[string]interface{}{}
	for _, raw := range data {
		newData := YamlMerge{}
		if err := yaml.Unmarshal(raw, &newData); err != nil {
			return errors.WithStack(err)
		}
		mergeMaps(mapData, newData)
	}
	// marshal merged maps back in to yaml
	defBytes, err := yaml.Marshal(mapData)
	if err != nil {
		return errors.WithStack(err)
	}
	// unmarshal yaml back in to interface
	return errors.WithStack(yaml.Unmarshal(defBytes, def))
}
