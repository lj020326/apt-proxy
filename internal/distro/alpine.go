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

var AlpineHostPattern = regexp.MustCompile(`/alpine/(.+)$`)

const AlpineBenchmarkURL = "MIRRORS.txt"

// Official Alpine Linux primary and US regional mirror endpoints.
// Sites that contain protocol headers restrict access to resources using that protocol.
var AlpineOfficialMirrors = []string{
	"dl-cdn.alpinelinux.org/alpine/",
	"uk.alpinelinux.org/alpine/",
	"mirrors.kernel.org/alpine/",
	"mirror.us.leaseweb.net/alpine/",
	"mirror.clarkson.edu/alpine/",
	"mirrors.rit.edu/alpine/",
	"mirrors.ocf.berkeley.edu/alpine/",
}

var AlpineCustomMirrors = []string{
	"mirror.steadfast.net/alpine/",
}

var BuiltinAlpineMirrors = GenerateBuildInList(AlpineOfficialMirrors, AlpineCustomMirrors)

var AlpineDefaultCacheRules = []Rule{
	{Pattern: regexp.MustCompile(`APKINDEX.tar.gz$`), CacheControl: `max-age=3600`, Rewrite: true, OS: TypeAlpine},
	{Pattern: regexp.MustCompile(`tar.gz$`), CacheControl: `max-age=3600`, Rewrite: true, OS: TypeAlpine},
	{Pattern: regexp.MustCompile(`apk$`), CacheControl: `max-age=3600`, Rewrite: true, OS: TypeAlpine},
	{Pattern: regexp.MustCompile(`.*`), CacheControl: `max-age=100000`, Rewrite: true, OS: TypeAlpine},
}
