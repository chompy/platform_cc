package api

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var dirToTarGzBasePath = ".platform"

// dirToTarGz - convert contents of directory to tar.gz
func dirToTarGz(path string) ([]byte, error) {
	path = filepath.Join(dirToTarGzBasePath, path)
	out := bytes.Buffer{}
	gw := gzip.NewWriter(&out)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	if err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
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
			header.Name = strings.TrimLeft(strings.Replace(path, dirToTarGzBasePath, "", 1), "/")
			header.Name = strings.ReplaceAll(header.Name, "\\", "/")
			log.Println(header.Name)
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
		return []byte{}, err
	}
	return out.Bytes(), nil
}

// dirToTarGzB64 - convert contents of directory to tar.gz base64 encoded
func dirToTarGzB64(path string) (string, error) {
	data, err := dirToTarGz(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// unmarshalYamlValue - get value of current node
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
				log.Println(err)
			}
			return out
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

// unmarshalYamlWithCustomTags - unmarshal yaml with custom tags
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
