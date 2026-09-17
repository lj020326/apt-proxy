// Copyright 2026 LJ Johnson
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package distro

import "regexp"

var CentosHostPattern = regexp.MustCompile(`/centos/(.+)$`)

const CentosBenchmarkURL = "TIME"

// Official CentOS Vault and US regional mirror endpoints.
// Sites that contain protocol headers restrict access to resources using that protocol.
var CentosOfficialMirrors = []string{
	"vault.centos.org/centos/",
	"mirror.centos.org/centos/",
	"mirrors.kernel.org/centos/",
	"mirror.us.leaseweb.net/centos/",
	"mirror.clarkson.edu/centos/",
	"mirrors.rit.edu/centos/",
	"mirror.math.princeton.edu/pub/centos/",
	"mirrors.ocf.berkeley.edu/centos/",
	"mirror.cs.vt.edu/pub/centos/",
}

var CentosCustomMirrors = []string{
	"mirror.steadfast.net/centos/",
	"mirror.cc.columbia.edu/pub/linux/centos/",
}

var BuiltinCentosMirrors = GenerateBuildInList(CentosOfficialMirrors, CentosCustomMirrors)

var CentosDefaultCacheRules = []Rule{
	{Pattern: regexp.MustCompile(`repomd.xml$`), CacheControl: `max-age=3600`, Rewrite: true, OS: TypeCentOS},
	{Pattern: regexp.MustCompile(`filelist.gz$`), CacheControl: `max-age=3600`, Rewrite: true, OS: TypeCentOS},
	{Pattern: regexp.MustCompile(`dir_sizes$`), CacheControl: `max-age=3600`, Rewrite: true, OS: TypeCentOS},
	{Pattern: regexp.MustCompile(`TIME$`), CacheControl: `max-age=3600`, Rewrite: true, OS: TypeCentOS},
	{Pattern: regexp.MustCompile(`timestamp.txt$`), CacheControl: `max-age=3600`, Rewrite: true, OS: TypeCentOS},
	{Pattern: regexp.MustCompile(`.*`), CacheControl: `max-age=100000`, Rewrite: true, OS: TypeCentOS},
}
