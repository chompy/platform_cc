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
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

const dockerNetworkName = "pcc"

// createNetwork creates a global network for use with all PCC containers.
func (d Docker) createNetwork() error {
	output.LogDebug("Create Docker network.", dockerNetworkName)
	if _, err := d.client.NetworkCreate(
		context.Background(),
		dockerNetworkName,
		types.NetworkCreate{
			CheckDuplicate: true,
		},
	); err != nil {
		if !strings.Contains(err.Error(), "exists") {
			return tracerr.Wrap(err)
		}
	}
	return nil
}

// deleteNetwork deletes the global network.
func (d Docker) deleteNetwork() error {
	done := output.Duration("Delete network.")
	err := d.client.NetworkRemove(
		context.Background(),
		dockerNetworkName,
	)
	if !client.IsErrNetworkNotFound(err) {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// getNetworkHostIP gets the host IP address for the given project's network.
func (d Docker) getNetworkHostIP() (string, error) {
	net, err := d.client.NetworkInspect(
		context.Background(),
		dockerNetworkName,
	)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	return net.IPAM.Config[0].Gateway, nil
}
