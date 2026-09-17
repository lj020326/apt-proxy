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
	UbuntuPortsGeoMirrorAPI = "http://mirrors.ubuntu.com/mirrors.txt"
	// Ubuntu Ports targets non-amd64 architectures. Use the arm64 Release
	// file which mirrors universally publish.
	UbuntuPortsBenchmarkURL = "dists/noble/main/binary-arm64/Release"
)

var UbuntuPortsHostPattern = regexp.MustCompile(`/ubuntu-ports/(.+)$`)

// Official Ubuntu Ports primary and US regional mirror endpoints.
// Sites that contain protocol headers restrict access to resources using that protocol.
var UbuntuPortsOfficialMirrors = []string{
	"ports.ubuntu.com/ubuntu-ports/",
	"mirrors.kernel.org/ubuntu-ports/",
	"mirror.us.leaseweb.net/ubuntu-ports/",
	"mirror.clarkson.edu/ubuntu-ports/",
	"mirrors.rit.edu/ubuntu-ports/",
	"mirrors.ocf.berkeley.edu/ubuntu-ports/",
	"mirror.cs.vt.edu/pub/ubuntu-ports/",
}

var UbuntuPortsCustomMirrors = []string{
	"mirror.steadfast.net/ubuntu-ports/",
}

var BuiltinUbuntuPortsMirrors = GenerateBuildInList(UbuntuPortsOfficialMirrors, UbuntuPortsCustomMirrors)

var UbuntuPortsDefaultCacheRules = newDebStyleRules(TypeUbuntuPorts)
