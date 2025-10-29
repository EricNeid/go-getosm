// SPDX-FileCopyrightText: Copyright (C) 2009-2021 German Aerospace Center (DLR) and others.
// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: EPL-2.0 OR GPL-2.0-or-later
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	app "github.com/EricNeid/go-getosm"
	"github.com/op/go-logging"
)

const apiURL = "https://www.overpass-api.de/api/interpreter"

var (
	bbox          = ""
	tiles         = 1
	prefix        = "osm"
	timeout       = 240
	elementLimit  = 1073741824
	verbose       = false
	retries       = 5
	retryDelaySec = 2
)

func init() {
	flag.StringVar(&bbox, "b", bbox, "Bounding box: west,south,east,north")
	flag.StringVar(&prefix, "prefix", prefix, "Prefix of output file")
	flag.IntVar(&tiles, "t", tiles, "Number of tiles to split the bounding box into")
	flag.IntVar(&timeout, "timeout", timeout, "timeout for connection")
	flag.IntVar(&retries, "retries", retries, "how often to retry the download of a failed tile")
	flag.IntVar(&retryDelaySec, "retryDelay", retryDelaySec, "delay between retries in seconds")
	flag.IntVar(&elementLimit, "elementLimit", elementLimit, "Element limit in osm file")
	flag.BoolVar(&verbose, "verbose", verbose, "Verbose output")

	flag.Parse()
}

var log = app.Log

func main() {
	if tiles < 1 {
		log.Infof("invalid number of tiles given, must be >= 1\n")
		flag.Usage()
		os.Exit(1)
	}

	if verbose {
		app.SetLogLevel(logging.DEBUG)
	} else {
		app.SetLogLevel(logging.INFO)
	}

	bbs, err := app.ReadBoundingBox(bbox, tiles)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	var output []string
	for i, bb := range bbs {
		log.Infof("downloading tiles %d of %d\n", i+1, len(bbs))
		query := app.FormatQuery(bb, timeout, elementLimit)
		result, err := app.Download(apiURL, query)

		retryDelay := time.Duration(retryDelaySec) * time.Second
		for retry := 1; err != nil && retry <= retries; retry++ {
			log.Warningf("error downloading data: %v, attempting retry %d of %d in %s seconds\n", err, retry, retries, retryDelay)
			result, err = app.Download(apiURL, query)
			time.Sleep(retryDelay)
		}
		if err != nil {
			log.Fatalf("error downloading data: %v, maximum retries reached\n", err)
		}

		var outputFile string
		if len(bbs) > 1 {
			outputFile = fmt.Sprintf("%s%d_%d.osm.xml", prefix, i+1, len(bbs))
		} else {
			outputFile = fmt.Sprintf("%s_bbox.osm.xml", prefix)
		}
		os.WriteFile(outputFile, *result, os.ModePerm)
		output = append(output, outputFile)
	}

	log.Infof("all done")
	log.Infof("files created: %s\n", strings.Join(output, ","))
}
