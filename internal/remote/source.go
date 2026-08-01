// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package remote

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const remoteBaseURLEnv = "SYNAPSEQ_REMOTE_BASE_URL"

type remoteSource struct {
	baseURL  string
	cacheKey string
	custom   bool
}

func defaultRemoteSource() (remoteSource, error) {
	baseURL := strings.TrimSpace(os.Getenv(remoteBaseURLEnv))
	if baseURL == "" {
		return remoteSource{baseURL: strings.TrimSuffix(t.RemoteIndexURL, "/free/index.json")}, nil
	}

	return customRemoteSource(baseURL)
}

func customRemoteSource(rawURL string) (remoteSource, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return remoteSource{}, fmt.Errorf("invalid remote base URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return remoteSource{}, fmt.Errorf("invalid remote base URL: use http or https")
	}
	if parsed.Host == "" {
		return remoteSource{}, fmt.Errorf("invalid remote base URL: host is required")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.Path != "" && parsed.Path != "/" {
		return remoteSource{}, fmt.Errorf("invalid remote base URL: URL must not include a path, credentials, query, or fragment")
	}

	host := strings.ToLower(parsed.Hostname())
	if host == "" {
		return remoteSource{}, fmt.Errorf("invalid remote base URL: host is required")
	}
	if port := parsed.Port(); port != "" {
		host += "_" + port
	}
	host = cacheSafeName(host)

	return remoteSource{
		baseURL:  (&url.URL{Scheme: parsed.Scheme, Host: parsed.Host}).String(),
		cacheKey: host,
		custom:   true,
	}, nil
}

func cacheSafeName(value string) string {
	var builder strings.Builder
	for _, char := range value {
		if char >= 'a' && char <= 'z' || char >= '0' && char <= '9' || char == '.' || char == '-' {
			builder.WriteRune(char)
			continue
		}
		builder.WriteByte('_')
	}
	return builder.String()
}

func (source remoteSource) indexURL() string {
	if source.custom {
		return source.baseURL + "/index.json"
	}

	return t.RemoteIndexURL
}

func (source remoteSource) sequenceURL(sequencePath string) string {
	return source.baseURL + "/" + strings.TrimPrefix(sequencePath, "/")
}
