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

import (
	"regexp"
)

const (
	UbuntuGeoMirrorAPI = "http://mirrors.ubuntu.com/mirrors.txt"
	UbuntuBenchmarkURL = "dists/noble/main/binary-amd64/Release"
)

var UbuntuHostPattern = regexp.MustCompile(`/ubuntu/(.+)$`)

// US-based official Canonical and primary regional mirror endpoints.
// Sites that contain protocol headers restrict access to resources using that protocol.
var UbuntuOfficialMirrors = []string{
	"us.archive.ubuntu.com/ubuntu/",
	"security.ubuntu.com/ubuntu/",
	"archive.ubuntu.com/ubuntu/",
	"releases.ubuntu.com/ubuntu/",
	"ports.ubuntu.com/ubuntu-ports/",
	"mirrors.kernel.org/ubuntu/",
	"mirror.us.leaseweb.net/ubuntu/",
	"mirror.clarkson.edu/ubuntu/",
	"mirrors.rit.edu/ubuntu/",
	"mirror.math.princeton.edu/pub/ubuntu/",
	"mirrors.mit.edu/ubuntu/",
	"mirrors.ocf.berkeley.edu/ubuntu/",
	"mirrors.edge.kernel.org/ubuntu/",
	"mirror.cs.vt.edu/pub/ubuntu/",
	"mirrors.tripole.ru/ubuntu/",
}

var UbuntuCustomMirrors = []string{
	"mirror.steadfast.net/ubuntu/",
	"mirror.cc.columbia.edu/pub/linux/ubuntu/archive/",
}

var BuiltinUbuntuMirrors = GenerateBuildInList(UbuntuOfficialMirrors, UbuntuCustomMirrors)

var UbuntuDefaultCacheRules = newDebStyleRules(TypeUbuntu)
