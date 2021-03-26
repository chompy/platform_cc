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

package container

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

type imagePullProgress struct {
	Status string `json:"status"`
	Detail struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"progressDetail"`
	ID string `json:"id"`
}

// imagePullSingle pulls the latest image for the given container.
func (d Docker) imagePullSingle(c Config, progress func(p imagePullProgress)) error {
	done := output.Duration(
		fmt.Sprintf("Pull Docker container image for '%s.'", c.GetContainerName()),
	)
	output.LogDebug(
		"Pull container image.",
		map[string]interface{}{
			"container_id": c.GetContainerName(),
			"image":        c.Image,
		},
	)
	r, err := d.client.ImagePull(
		context.Background(),
		c.Image,
		types.ImagePullOptions{},
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer r.Close()
	for {
		b := make([]byte, 2048)
		n, _ := r.Read(b)
		if n == 0 {
			break
		}
		// report on progress
		if progress != nil {
			prog := imagePullProgress{}
			if err := json.Unmarshal(bytes.Trim(b, "\x00"), &prog); err == nil {
				progress(prog)
			}
		}
	}
	done()
	return nil
}

// ImagePull pulls one or more Docker images.
func (d Docker) ImagePull(c []Config) error {
	// get list of images and prepare progress output
	images := make([]string, 0)
	msgs := make([]string, 0)
	for _, cc := range c {
		hasImage := false
		for _, image := range images {
			if image == cc.Image {
				hasImage = true
				break
			}
		}
		if !hasImage {
			images = append(images, cc.Image)
			msgs = append(msgs, cc.GetHumanName())
		}
	}
	prog := output.Progress(msgs)
	// simultaneously pull images
	outputEnabled := output.Enable
	var wg sync.WaitGroup
	for i, c := range c {
		wg.Add(1)
		go func(i int, c Config) {
			defer wg.Done()
			defer func() { output.Enable = outputEnabled }()
			output.Enable = false
			imagePullProg := func(p imagePullProgress) {
				prog(i, output.ProgressMessageWait, &p.Detail.Current, &p.Detail.Total)
			}
			if err := d.imagePullSingle(c, imagePullProg); err != nil {
				prog(i, output.ProgressMessageError, nil, nil)
				output.LogError(err)
				return
			}
			prog(i, output.ProgressMessageDone, nil, nil)
		}(i, c)
	}
	wg.Wait()
	return nil
}

// hasImage returns true if given image exists.
func (d Docker) hasImage(name string) bool {
	_, _, err := d.client.ImageInspectWithRaw(
		context.Background(), name,
	)
	return err == nil
}

// listProjectImages returns list of all project Docker images.
func (d Docker) listProjectImages(pid string) ([]types.ImageSummary, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(
		"reference",
		fmt.Sprintf(
			"%s%s*",
			dockerCommitTagPrefix,
			fmt.Sprintf(containerNamingPrefix, pid),
		),
	)
	return d.client.ImageList(
		context.Background(),
		types.ImageListOptions{
			Filters: filterArgs,
		},
	)
}

// listAllImages returns list of all Platform.CC Docker images.
func (d Docker) listAllImages() ([]types.ImageSummary, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("reference", dockerCommitTagPrefix+"*")
	return d.client.ImageList(
		context.Background(),
		types.ImageListOptions{
			Filters: filterArgs,
		},
	)
}

// deleteImages deletes given Docker images.
func (d Docker) deleteImages(images []types.ImageSummary) error {
	// prepare progress output
	output.LogDebug("Delete Docker image.", images)
	msgs := make([]string, len(images))
	for i, img := range images {
		msgs[i] = img.ID
	}
	done := output.Duration("Delete images.")
	prog := output.Progress(msgs)
	// delete volumes
	var wg sync.WaitGroup
	for i, img := range images {
		wg.Add(1)
		go func(name string, i int) {
			defer wg.Done()
			if _, err := d.client.ImageRemove(
				context.Background(),
				name,
				types.ImageRemoveOptions{Force: true},
			); err != nil {
				prog(i, output.ProgressMessageError, nil, nil)
				output.Warn(err.Error())
				return
			}
			prog(i, output.ProgressMessageDone, nil, nil)
		}(img.ID, i)
	}
	wg.Wait()
	done()
	return nil
}

// serviceTypeFromImage extract the service type name from the Docker image.
func (d Docker) serviceTypeFromImage(name string) string {
	// extract type name from Platform.sh image name
	sType := typeFromImageName(name)
	if sType != "" {
		return sType
	}
	// type name or raw image id was provided, inspect the image
	inspectImage, _, err := d.client.ImageInspectWithRaw(
		context.Background(),
		name,
	)
	if err != nil {
		output.LogError(err)
		return ""
	}
	if len(inspectImage.RepoTags) > 0 {
		sType := typeFromImageName(inspectImage.RepoTags[0])
		if sType != "" {
			return sType
		}
	}
	// inspect parent
	parentInspectImage, _, err := d.client.ImageInspectWithRaw(
		context.Background(),
		inspectImage.Parent,
	)
	if err != nil {
		output.LogError(err)
		return ""
	}
	if len(parentInspectImage.RepoTags) > 0 {
		return typeFromImageName(parentInspectImage.RepoTags[0])
	}
	return ""
}
