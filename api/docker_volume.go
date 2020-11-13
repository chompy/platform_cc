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
func (d *DockerClient) GetProjectVolumes(pid string) (volume.VolumesListOKBody, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", fmt.Sprintf(dockerNamingPrefix+"*", pid))
	return d.cli.VolumeList(
		context.Background(),
		filterArgs,
	)
}

// GetAllVolumes - get list of all volumes used by PCC
func (d *DockerClient) GetAllVolumes() (volume.VolumesListOKBody, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", "pcc-*")
	return d.cli.VolumeList(
		context.Background(),
		filterArgs,
	)
}

// DeleteProjectVolumes - delete all volumes for given project
func (d *DockerClient) DeleteProjectVolumes(pid string) error {
	volList, err := d.GetProjectVolumes(pid)
	if err != nil {
		return err
	}
	return d.deleteVolumes(volList)
}

// DeleteAllVolumes - delete all volumes related to PCC
func (d *DockerClient) DeleteAllVolumes() error {
	volList, err := d.GetAllVolumes()
	if err != nil {
		return err
	}
	return d.deleteVolumes(volList)
}

func (d *DockerClient) deleteVolumes(volList volume.VolumesListOKBody) error {
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
		err := <-ch
		if err != nil {
			return err
		}
	}
	return nil
}
