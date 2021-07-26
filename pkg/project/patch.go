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

// elasticsearchPatch is a patch for Elasticsearch that fixes the config file so that it doesn't expect a global IP address.
const elasticsearchPatch = `
if [ -f /usr/share/elasticsearch/config/elasticsearch.yml.psh-tmpl ]; then
	sed -i 's/_global_/_site_/g' /usr/share/elasticsearch/config/elasticsearch.yml.psh-tmpl
fi
`

// lispPatch is a patch for Lisp that fixes code that doesn't account for output_dir having a trailing slash.
const lispPatch = `
	sed -i 's/system_name = .*/system_name = asd\[0\]\[len\(builder\.output_dir\)\:-len(\"\.asd\"\)\]/g' /etc/platform/flavor.d/default.py
`

// ensureDirPatch is a patch to fix the import of ensure_dir in platformsh/agent/service.py where it is missing.
const ensureDirPatch = `
	echo "from .util import ensure_dir"|cat - /usr/lib/python2.7/dist-packages/platformsh/agent/service.py > /tmp/out && mv /tmp/out /usr/lib/python2.7/dist-packages/platformsh/agent/service.py
`

// contextualCodePatch is a patch to fix the platformsh matcher bundle to appropiately match platform.cc routes.
const contextualCodePatch = `
if [ -d /app/vendor/contextualcode/platformsh-siteaccess-matcher-bundle ]; then
	cd /app/vendor/contextualcode
	rm -r platformsh-siteaccess-matcher-bundle
	git clone https://gitlab.com/contextualcode/platformsh-siteaccess-matcher-bundle.git -b platform_cc
fi
`

// patchMap is a map of service types to their patch command.
var patchMap = map[string]string{
	"elasticsearch:*": elasticsearchPatch,
	"lisp:*":          lispPatch,
	"memcached:1.4":   ensureDirPatch,
}

// postBuildPatchMap is a map of service types to their post-build patch command.
var postBuildPatchMap = map[string]string{
	"php:*": contextualCodePatch,
}

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

// GetDefinitionPostBuildPatch returns post-build patch command to for given definition.
func (p *Project) GetDefinitionPostBuildPatch(d interface{}) string {
	defType := p.GetDefinitionType(d)
	for m, command := range postBuildPatchMap {
		if wildcardCompare(defType, m) {
			return command
		}
	}
	return ""
}
