// SPDX-FileCopyrightText: Copyright (C) 2009-2021 German Aerospace Center (DLR) and others.
// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: EPL-2.0 OR GPL-2.0-or-later
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	getosm "github.com/EricNeid/go-getosm"
	"github.com/op/go-logging"
)

const apiURL = "https://www.overpass-api.de/api/interpreter"
const maxRetries = 5
const retryDelay = 1 * time.Second

var (
	bbox         = ""
	tiles        = 1
	prefix       = "osm"
	timeout      = 240
	elementLimit = 1073741824
	verbose      = false
	queryOutput  = false
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("osmGet: ")

	flag.StringVar(&bbox, "b", bbox, "Bounding box: west,south,east,north")
	flag.StringVar(&prefix, "prefix", prefix, "Prefix of output file")
	flag.IntVar(&tiles, "t", tiles, "Number of tiles to split the bounding box into")
	flag.IntVar(&timeout, "timeout", timeout, "timeout for connection")
	flag.IntVar(&elementLimit, "elementLimit", elementLimit, "Element limit in osm file")
	flag.BoolVar(&verbose, "verbose", verbose, "Verbose output")
	flag.BoolVar(&queryOutput, "q", queryOutput, "Query output")

	flag.Parse()
}

func main() {
	if tiles < 1 {
		log.Println("invalid number of tiles given, must be >= 1")
		flag.Usage()
		os.Exit(1)
	}

	if verbose {
		getosm.SetLogLevel(logging.DEBUG)
	} else {
		getosm.SetLogLevel(logging.INFO)
	}

	bbs, err := getosm.ReadBoundingBox(bbox, tiles)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	for i, bb := range bbs {
		log.Printf("downloading tiles %d of %d\n", i+1, len(bbs))
		query := getosm.FormatQuery(bb, timeout, elementLimit)
		if queryOutput {
			log.Printf("download using query %s\n", query)
		}
		result, err := getosm.Download(apiURL, query)
		for retry := 1; err != nil && retry <= maxRetries; retry++ {
			log.Printf("error downloading data: %v, attempting retry %d of %d in %s seconds\n", err, retry, maxRetries, retryDelay)
			result, err = getosm.Download(apiURL, query)
			time.Sleep(retryDelay)
		}
		if err != nil {
			log.Fatalf("error downloading data: %v, maximum retries reached\n", err)
		}
		if len(bbs) > 1 {
			output := fmt.Sprintf("%s%d_%d.osm.xml", prefix, i+1, len(bbs))
			os.WriteFile(output, *result, os.ModePerm)
		} else {
			output := fmt.Sprintf("%s_bbox.osm.xml", prefix)
			os.WriteFile(output, *result, os.ModePerm)
		}
	}

	log.Println("all done")
}
