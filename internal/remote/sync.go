// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package remote

// RemoteSync updates the local Remote index cache.
func RemoteSync() error {
	source, err := defaultRemoteSource()
	if err != nil {
		return err
	}

	return remoteSync(source)
}

// RemoteSyncURL updates the cache for a custom Remote base URL.
func RemoteSyncURL(baseURL string) error {
	source, err := customRemoteSource(baseURL)
	if err != nil {
		return err
	}

	return remoteSync(source)
}

func remoteSync(source remoteSource) error {
	cache, err := openRemoteCacheForSource(source)
	if err != nil {
		return err
	}

	data, response, err := downloadURL(source.indexURL())
	if err != nil {
		return err
	}
	if err := validateJSONContentType(response); err != nil {
		return err
	}
	if _, err := parseRemoteIndex(data); err != nil {
		return err
	}

	if err = cache.index().write(data); err != nil {
		return err
	}

	return nil
}
