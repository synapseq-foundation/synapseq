/*
 * SynapSeq - Text-Driven Audio Sequencer for Brainwave Entrainment
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package hub

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var hubHTTPGet = http.Get

func downloadURL(url string) ([]byte, *http.Response, error) {
	response, err := hubHTTPGet(url)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	if err := validateHTTPResponse(response, url); err != nil {
		return nil, response, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, response, err
	}

	return data, response, nil
}

func validateHTTPResponse(response *http.Response, url string) error {
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("unexpected response downloading %s: %s", url, response.Status)
	}

	return nil
}

func validateJSONContentType(response *http.Response) error {
	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return fmt.Errorf("invalid content-type for manifest file: %s", contentType)
	}

	return nil
}

func downloadFile(url, path string) error {
	data, _, err := downloadURL(url)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
