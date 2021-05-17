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

package platformsh

import "github.com/pkg/errors"

var (
	ErrProjectNotFound     = errors.New("platform.sh project not found")
	ErrBadAPIResponse      = errors.New("platform.sh api returned unexpected response")
	ErrMissingAPIToken     = errors.New("platform.sh api token not found, please use the platformsh:login command to generate it")
	ErrInvalidEnvironment  = errors.New("invalid platform.sh environment")
	ErrEnvironmentNotFound = errors.New("platform.sh environment not found")
)
