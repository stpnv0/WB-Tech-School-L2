package urlutil

import (
	"net/url"
	"path"
	"strings"
)

func Normalize(rawURL string) (string, error) {
	if rawURL == "" {
		return "", nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	u.Fragment = ""

	if u.Path != "" {
		hadTrailingSlash := strings.HasSuffix(u.Path, "/")
		u.Path = path.Clean(u.Path)
		if hadTrailingSlash && u.Path != "/" {
			u.Path += "/"
		}
	}

	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	u.Host = removeDefaultPort(u)

	if u.RawQuery == "" {
		u.ForceQuery = false
	}

	return u.String(), nil
}

func removeDefaultPort(u *url.URL) string {
	host := u.Host
	if !containsPort(host) {
		return host
	}

	hostParts := strings.Split(host, ":")
	if len(hostParts) != 2 {
		return host
	}

	hostWithoutPort := hostParts[0]
	port := hostParts[1]

	switch u.Scheme {
	case "http":
		if port == "80" {
			return hostWithoutPort
		}
	case "https":
		if port == "443" {
			return hostWithoutPort
		}
	}

	return host
}

func containsPort(host string) bool {
	return strings.Contains(host, ":") && !strings.HasPrefix(host, "[")
}
