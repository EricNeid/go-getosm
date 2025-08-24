// SPDX-FileCopyrightText: Copyright (C) 2009-2021 German Aerospace Center (DLR) and others.
// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: EPL-2.0 OR GPL-2.0-or-later
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const apiURL = "https://www.overpass-api.de/api/interpreter"
const maxRetries = 5
const retryDelay = 1 * time.Second

var (
	bbox         = ""
	prefix       = "osm"
	timeout      = 240
	elementLimit = 1073741824
)

var errorInvalidBB = errors.New("invalid bounding box given")
var errorDownload = errors.New("could not retrieve data")

func init() {
	log.SetFlags(0)
	log.SetPrefix("osmGet: ")

	flag.StringVar(&bbox, "b", bbox, "Bounding box: west,south,east,north")
	flag.StringVar(&prefix, "prefix", prefix, "Prefix of output file")

	flag.IntVar(&timeout, "timeout", timeout, "timeout for connection")
	flag.IntVar(&elementLimit, "elementLimit", elementLimit, "Prefix for osm xml")

	flag.Parse()
}

func main() {
	bb, err := readBoundingBox(bbox)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	query := formatQuery(bb)
	result, err := download(query)
	for retry := 1; err != nil && retry <= maxRetries; retry++ {
		log.Printf("error downloading data: %v, attempting retry %d of %d in %s seconds\n", err, retry, maxRetries, retryDelay)
		result, err = download(query)
		time.Sleep(retryDelay)
	}
	if err != nil {
		log.Fatalf("error downloading data: %v, maximum retries reached\n", err)
	}

	output := fmt.Sprintf("%s_bbox.osm.xml", prefix)
	os.WriteFile(output, *result, os.ModeAppend)
}
