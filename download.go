// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt
// Inspired by SUMO's osmGet.py (no code copied).
package gogetosm

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var ErrorDownload = errors.New("could not retrieve data")

func FormatQuery(bb BoundingBox, timeout, elementLimit int) string {
	return fmt.Sprintf(`
	<osm-script timeout="%d" element-limit="%d">
	<union>
		<bbox-query n="%g" s="%g" w="%g" e="%g"/>
		<recurse type="node-relation" into="rels"/>
		<recurse type="node-way"/>
		<recurse type="way-relation"/>
	</union>
	<union>
		<item/>
		<recurse type="way-node"/>
	</union>
	<print/>
	</osm-script>
	`, timeout, elementLimit, bb.North, bb.South, bb.West, bb.East)
}

func Download(apiURL, query string) (*[]byte, error) {
	Log.Debugf("download using query: %s\n", query)
	client := http.Client{
		Timeout: 0, // no timeout
	}
	resp, err := client.Post(apiURL, "text/xml", strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Log.Errorf("error while reading response body %v\n", err)
		return nil, err
	}

	if resp.StatusCode == 200 {
		Log.Debugf("download complete")
	} else {
		Log.Errorf("download failed with status %s\n", resp.Status)
		Log.Errorf("response is: \n%s\n", string(body))
		return nil, ErrorDownload
	}

	return &body, err
}
