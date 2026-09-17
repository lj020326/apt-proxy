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

package proxy

import (
	"net"
	"net/http"
	"strings"

	"github.com/lj020326/apt-proxy/internal/distro"
)

// tlsRewriteMarker is apt-cacher-ng's "tell-me-what-you-need" marker.
// Clients write http://HTTPS///<host>/...; the proxy fetches https://<host>/...
const tlsRewriteMarker = "HTTPS//"

// hasTLSRewriteMarker reports whether r carries the ACNG HTTPS/// marker.
func hasTLSRewriteMarker(r *http.Request) bool {
	if r == nil || r.URL == nil {
		return false
	}
	for _, authority := range [2]string{r.URL.Host, r.Host} {
		if authority == "" {
			continue
		}
		host := authority
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		if strings.EqualFold(host, "HTTPS") {
			return true
		}
	}
	path := strings.ToUpper(r.URL.EscapedPath())
	return strings.HasPrefix(path, "/"+tlsRewriteMarker) ||
		strings.Contains(path, "/"+tlsRewriteMarker+"/")
}

// parseTLSRewrite extracts the real upstream host and path from an
// HTTPS///-marked request. Returns ok=false when the marker is present but
// the remainder is not a usable host/path pair.
//
// Supported shapes:
//
//	Host=HTTPS, Path=///get.docker.com/ubuntu/dists/...
//	Host=HTTPS:3142, Path=///get.docker.com/...
//	Path=/HTTPS///get.docker.com/ubuntu/...
//	Path=/https///get.docker.com/ubuntu/...  (case-insensitive marker)
func parseTLSRewrite(r *http.Request) (host, path string, ok bool) {
	if r == nil || r.URL == nil {
		return "", "", false
	}

	escaped := r.URL.EscapedPath()
	upper := strings.ToUpper(escaped)

	// Path form: /HTTPS///<host>/<path...>
	if idx := strings.Index(upper, "/"+tlsRewriteMarker+"/"); idx >= 0 {
		rest := escaped[idx+len("/"+tlsRewriteMarker+"/"):]
		return splitHostPath(rest)
	}
	// Path form with exactly three slashes after marker: /HTTPS///<host>...
	if strings.HasPrefix(upper, "/"+tlsRewriteMarker) {
		rest := escaped[len("/"+tlsRewriteMarker):]
		rest = strings.TrimPrefix(rest, "/")
		return splitHostPath(rest)
	}

	// Authority form: Host is "HTTPS" (optionally with port), path is ///host/...
	for _, authority := range [2]string{r.URL.Host, r.Host} {
		if authority == "" {
			continue
		}
		h := authority
		if name, _, err := net.SplitHostPort(h); err == nil {
			h = name
		}
		if !strings.EqualFold(h, "HTTPS") {
			continue
		}
		return splitHostPath(escaped)
	}
	return "", "", false
}

func splitHostPath(rest string) (host, path string, ok bool) {
	// Authority form yields paths like "///host/..." (three leading slashes).
	// Strip all of them so the first segment is the real upstream host.
	rest = strings.TrimLeft(rest, "/")
	if rest == "" {
		return "", "", false
	}
	slash := strings.IndexByte(rest, '/')
	if slash < 0 {
		host = rest
		path = "/"
	} else {
		host = rest[:slash]
		path = rest[slash:]
	}
	if host == "" || strings.Contains(host, " ") {
		return "", "", false
	}
	// Reject obvious non-hosts (e.g. scheme leftovers).
	if strings.Contains(host, "://") {
		return "", "", false
	}
	return host, path, true
}

// applyTLSRewrite mutates r so the reverse proxy fetches via HTTPS from the
// host encoded in the marker. Query string is preserved.
func applyTLSRewrite(r *http.Request) bool {
	host, path, ok := parseTLSRewrite(r)
	if !ok {
		return false
	}
	r.URL.Scheme = "https"
	r.URL.Host = host
	r.URL.Path = path
	r.URL.RawPath = ""
	r.Host = host
	return true
}

// defaultTLSCacheRule picks Cache-Control for HTTPS/// requests that do not
// match a registered distribution. Mirrors apt-cacher-ng's practical split:
// indexes/metadata stay short-lived; immutable package blobs live longer.
func defaultTLSCacheRule(path string) *distro.Rule {
	lower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lower, ".deb"),
		strings.HasSuffix(lower, ".rpm"),
		strings.HasSuffix(lower, ".apk"),
		strings.HasSuffix(lower, ".tar.gz"),
		strings.HasSuffix(lower, ".tar.xz"),
		strings.HasSuffix(lower, ".whl"),
		strings.Contains(lower, "/pool/"):
		return &distro.Rule{CacheControl: "public, max-age=31536000, immutable"}
	case strings.Contains(lower, "/dists/"),
		strings.HasSuffix(lower, "inrelease"),
		strings.HasSuffix(lower, "release"),
		strings.HasSuffix(lower, "packages"),
		strings.HasSuffix(lower, "sources"),
		strings.HasSuffix(lower, "apkindex.tar.gz"):
		return &distro.Rule{CacheControl: "public, max-age=300"}
	default:
		return &distro.Rule{CacheControl: "public, max-age=3600"}
	}
}
