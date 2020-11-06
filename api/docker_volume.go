package api

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

// containerVolumeNameFormat - format for mount volume names
const containerVolumeNameFormat = dockerNamingPrefix + "v-%s"

// GetProjectVolumes - get list of all volumes for given project
func (d *dockerClient) GetProjectVolumes(pid string) (volume.VolumesListOKBody, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", fmt.Sprintf(dockerNamingPrefix+"*", pid))
	return d.cli.VolumeList(
		context.Background(),
		filterArgs,
	)
}

// DeleteProjectVolumes - delete all volumes for given project
func (d *dockerClient) DeleteProjectVolumes(pid string) error {
	volList, err := d.GetProjectVolumes(pid)
	if err != nil {
		return err
	}
	ch := make(chan error)
	for _, vol := range volList.Volumes {
		log.Printf("Delete Docker volume '%s.'", vol.Name)
		go func(volName string) {
			if err := d.cli.VolumeRemove(
				context.Background(),
				volName,
				true,
			); err != nil {
				ch <- err
			}
			ch <- nil
		}(vol.Name)
	}
	for range volList.Volumes {
		err = <-ch
		if err != nil {
			return err
		}
	}
	return nil
}
