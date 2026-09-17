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

const (
	DebianBenchmarkURL = "dists/bookworm/main/binary-amd64/Release"
)

var DebianHostPattern = regexp.MustCompile(`/debian(-security)?/(.+)$`)

// DebianSecurityHostPattern matches requests addressed to the Debian security
// archive host. The classic sources.list form
//
//	deb http://security.debian.org <suite>-security main
//
// requests paths like /dists/trixie-security/InRelease, with no
// /debian-security/ prefix for DebianHostPattern to latch onto, so these
// requests are recognised by Host instead. An optional :port is tolerated
// because Host carries whatever the client dialled.
//
// Host patterns are matched against a lower-cased host (DNS names are
// case-insensitive, so Security.Debian.Org is the same archive), hence the
// lower-case literal here.
var DebianSecurityHostPattern = regexp.MustCompile(`^security\.debian\.org(:\d+)?$`)

// BuiltinHostPattern returns the compile-time Host matcher for a distro type,
// or nil when it has none. LoadFromConfig falls back to it so a
// distributions.yaml entry that omits host_pattern keeps the built-in
// behaviour, the same way omitting mirrors keeps the built-in mirror list.
func BuiltinHostPattern(distType int) *regexp.Regexp {
	if distType == TypeDebian {
		return DebianSecurityHostPattern
	}
	return nil
}

// Official Debian primary and US regional mirrors.
// Sites that contain protocol headers restrict access to resources using that protocol.
var DebianOfficialMirrors = []string{
	"ftp.us.debian.org/debian/",
	"deb.debian.org/debian/",
	"security.debian.org/debian-security/",
	"mirrors.kernel.org/debian/",
	"mirrors.mit.edu/debian/",
	"mirror.us.leaseweb.net/debian/",
	"mirror.clarkson.edu/debian/",
	"mirrors.rit.edu/debian/",
	"mirror.math.princeton.edu/pub/debian/",
	"mirrors.ocf.berkeley.edu/debian/",
	"mirror.cs.vt.edu/pub/debian/",
}

var DebianCustomMirrors = []string{
	"mirror.steadfast.net/debian/",
	"mirror.cc.columbia.edu/pub/linux/debian/",
}

var BuiltinDebianMirrors = GenerateBuildInList(DebianOfficialMirrors, DebianCustomMirrors)

var DebianDefaultCacheRules = newDebStyleRules(TypeDebian)
