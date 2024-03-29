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

import "errors"

var (
	// ErrContainerNotFound is an error returned when a container is not found.
	ErrContainerNotFound = errors.New("container not found")
	// ErrContainerNotRunning is an error returned when a container is not running.
	ErrContainerNotRunning = errors.New("container is not running")
	// ErrCannotDeleteCommit is an error returned when a container commit cannot be deleted.
	ErrCannotDeleteCommit = errors.New("cannot delete commit")
	// ErrImageNotFound is an error returned when a container image is not found.
	ErrImageNotFound = errors.New("image not found")
	// ErrInvalidSlot is an error returned when an invalid slot is specified.
	ErrInvalidSlot = errors.New("invalid slot")
	// ErrCommandExited is an error returned when a container command exits with a non zero exit code.
	ErrCommandExited = errors.New("command exited with error")
)
