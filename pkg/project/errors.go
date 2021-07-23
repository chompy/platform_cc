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

package project

import "github.com/pkg/errors"

var (
	// ErrContainerNoIP is an error returned when a container has no IP address.
	ErrContainerNoIP = errors.New("container has no ip address")
	// ErrInvalidRelationship is an error returned when Platform.sh relationships are invalid.
	ErrInvalidRelationship = errors.New("one or more relationships are invalid")
	// ErrNoApplicationFound is an error returned when a project has no applications.
	ErrNoApplicationFound = errors.New("project should have at least one application")
	// ErrInvalidDefinition is an error returned when a service definition is invalid.
	ErrInvalidDefinition = errors.New("invalid definition")
	// ErrContainerRunning is returned when attempting to start an already running container.
	ErrContainerRunning = errors.New("container already running")
)
