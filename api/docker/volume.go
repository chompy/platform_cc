package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/ztrue/tracerr"
)

// containerVolumeNameFormat is the format for mount volume names.
const containerVolumeNameFormat = dockerNamingPrefix + "v-%s"

// CreateNFSVolume creates a NFS Docker volume.
func (d *Client) CreateNFSVolume(pid string, name string) error {
	pathString := fmt.Sprintf(":/System/Volumes/Data/%s", GetVolumeName(pid, name))
	_, err := d.cli.VolumeCreate(
		context.Background(),
		volume.VolumesCreateBody{
			Name:   GetVolumeName(pid, name),
			Driver: "local",
			DriverOpts: map[string]string{
				"type":   "nfs",
				"o":      "addr=host.docker.internal,rw,nolock,hard,nointr,nfsvers=3",
				"device": pathString,
			},
		},
	)
	return tracerr.Wrap(err)
}

// GetProjectVolumes gets a list of all volumes for given project.
func (d *Client) GetProjectVolumes(pid string) (volume.VolumesListOKBody, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", fmt.Sprintf(dockerNamingPrefix+"*", pid))
	return d.cli.VolumeList(
		context.Background(),
		filterArgs,
	)
}

// GetAllVolumes gets a list of all volumes used by PCC.
func (d *Client) GetAllVolumes() (volume.VolumesListOKBody, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", "pcc-*")
	return d.cli.VolumeList(
		context.Background(),
		filterArgs,
	)
}

// DeleteProjectVolumes deletes all volumes for given project.
func (d *Client) DeleteProjectVolumes(pid string) error {
	volList, err := d.GetProjectVolumes(pid)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return d.deleteVolumes(volList)
}

// DeleteAllVolumes deletes all volumes related to PCC.
func (d *Client) DeleteAllVolumes() error {
	volList, err := d.GetAllVolumes()
	if err != nil {
		return tracerr.Wrap(err)
	}
	return d.deleteVolumes(volList)
}

func (d *Client) deleteVolumes(volList volume.VolumesListOKBody) error {
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
			return tracerr.Wrap(err)
		}
	}
	return nil
}

// GetVolumeName generates a volume name for given project id and container name.
func GetVolumeName(pid string, name string) string {
	return fmt.Sprintf(dockerNamingPrefix+"%s-v", pid, name)
}
