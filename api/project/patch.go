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

import (
	"regexp"
	"strings"
)

// elasticsearchPatch is a command for Elasticsearch that fixes the config file so that it doesn't expect a global IP address.
const elasticsearchPatch = `
if [ -f /usr/share/elasticsearch/config/elasticsearch.yml.psh-tmpl ]; then
	sed -i 's/_global_/_site_/g' /usr/share/elasticsearch/config/elasticsearch.yml.psh-tmpl
fi
`

// lispPatch is a command for Lisp that fixes code that doesn't account for output_dir having a trailing slash.
const lispPatch = `
	sed -i 's/system_name = .*/system_name = asd\[0\]\[len\(builder\.output_dir\)\:-len(\"\.asd\"\)\]/g' /etc/platform/flavor.d/default.py
`

// patchMap is a map of service types to their patch command.
var patchMap = map[string]string{
	"elasticsearch:*": elasticsearchPatch,
	"lisp:*":          lispPatch,
}

// wildcardMatchCharacter is the character to use as the wildcard character.
const wildcardMatchCharacter = "*"

// WildcardCompare tests if original string matches test string with wildcard.
func wildcardCompare(original string, test string) bool {
	test = regexp.QuoteMeta(test)
	test = strings.Replace(test, "\\*", ".*", -1)
	regex, err := regexp.Compile(test + "$")
	if err != nil {
		return false
	}
	return regex.MatchString(original)
}

// GetDefinitionPatch returns patch command to for given definition.
func (p *Project) GetDefinitionPatch(d interface{}) string {
	defType := p.GetDefinitionType(d)
	for m, command := range patchMap {
		if wildcardCompare(defType, m) {
			return command
		}
	}
	return ""
}
