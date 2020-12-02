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

package docker

import (
	"context"
	"strings"

	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ztrue/tracerr"
)

// globalNetworkName is the name of the global network.
const globalNetworkName = "pcc"

// CreateNetwork creates a global network for use with all PCC containers.
func (d *Client) CreateNetwork() error {
	done := output.Duration("Create network.")
	if _, err := d.cli.NetworkCreate(
		context.Background(),
		globalNetworkName,
		types.NetworkCreate{
			CheckDuplicate: true,
		},
	); err != nil {
		if !strings.Contains(err.Error(), "exists") {
			return tracerr.Wrap(err)
		}
	}
	done()
	return nil
}

// DeleteNetwork deletes the global network.
func (d *Client) DeleteNetwork() error {
	done := output.Duration("Delete network.")
	err := d.cli.NetworkRemove(
		context.Background(),
		globalNetworkName,
	)
	if !client.IsErrNetworkNotFound(err) {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// GetNetworkHostIP gets the host IP address for the given project's network.
func (d *Client) GetNetworkHostIP() (string, error) {
	net, err := d.cli.NetworkInspect(
		context.Background(),
		globalNetworkName,
	)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	return net.IPAM.Config[0].Gateway, nil
}
