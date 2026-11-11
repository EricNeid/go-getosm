// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt
// Inspired by SUMO's osmGet.py (no code copied).
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	getosm "github.com/EricNeid/go-getosm"
	"github.com/op/go-logging"
)

const (
	version = "0.10.0"
)

type bboxesFlag []string

func (b *bboxesFlag) String() string {
	return strings.Join(*b, "; ")
}

func (b *bboxesFlag) Set(value string) error {
	*b = append(*b, value)
	return nil
}

var (
	url                = "https://www.overpass-api.de/api/interpreter"
	bboxes             bboxesFlag
	tilesX             = 1
	tilesY             = 1
	prefix             = "osm"
	timeout            = 240
	elementLimit       = 500000000
	verbose            = false
	retries            = 5
	retryDelaySec      = 10
	continueLastFailed = false
	userAgent          = fmt.Sprintf("OsmGet-Tool/%s", version)
)

func parseArguments() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Printf("Version: %s\n", version)
		flag.PrintDefaults()

		fmt.Printf("  Usage  : %s -b w,s,e,n\n", os.Args[0])
		fmt.Printf("  Example: %s -b 13.10,52.20,13.20,52.30 -tx 2 -ty 2 -verbose\n", os.Args[0])
	}

	flag.Var(&bboxes, "b", "Bounding box: west,south,east,north")
	flag.StringVar(&prefix, "prefix", prefix, "Prefix of output file")
	flag.StringVar(&userAgent, "userAgent", userAgent, "Set custom user agent to be used in the header")
	flag.IntVar(&tilesX, "tx", tilesX, "Number of columns tiles to split the bounding box")
	flag.IntVar(&tilesY, "ty", tilesY, "Number of rows to split the bounding box")
	flag.IntVar(&timeout, "timeout", timeout, "Timeout for connection")
	flag.IntVar(&retries, "retries", retries, "How often to retry the download of a failed tile")
	flag.IntVar(&retryDelaySec, "retryDelay", retryDelaySec, "Delay between retries in seconds")
	flag.BoolVar(&verbose, "verbose", verbose, "Verbose output")
	flag.BoolVar(&continueLastFailed, "continue", continueLastFailed, "Do not download already present tiles")
	flag.IntVar(&elementLimit, "elementLimit", elementLimit, "Element limit in osm file")
	flag.StringVar(&url, "url", url, "Overpass API URL")

	flag.Parse()

	if len(bboxes) == 0 {
		log.Errorf("no bounding box given, please specify at least one\n")
		flag.Usage()
		os.Exit(1)
	}

	if tilesX*tilesY < 1 {
		log.Infof("invalid number of tiles given, must be >= 1\n")
		flag.Usage()
		os.Exit(1)
	}
}

var log = getosm.Log

func main() {
	parseArguments()

	if verbose {
		getosm.SetLogLevel(logging.DEBUG)
	} else {
		getosm.SetLogLevel(logging.INFO)
	}

	var outputFiles []string
	for bbIndex, bb := range bboxes {
		bbTiles, err := getosm.ReadBoundingBox(bb, tilesX, tilesY)
		if err != nil {
			flag.Usage()
			os.Exit(1)
		}

		log.Infof("processing bounding box %d of %d\n", bbIndex+1, len(bboxes))
		for tileIndex, bb := range bbTiles {
			log.Infof("downloading tiles %d of %d\n", tileIndex+1, len(bbTiles))

			// output file name
			outputFile := getFileName(prefix, bbIndex, len(bboxes), tileIndex, len(bbTiles))

			// check if file already exists
			if continueLastFailed && doesFileExist(outputFile) {
				log.Infof("skipping already present tile %s\n", outputFile)
				continue
			}

			// download tile until retries are exhausted
			request := getosm.Request{
				Query:           getosm.FormatQuery(bb, timeout, elementLimit),
				CustomUserAgent: userAgent,
				URL:             url,
			}

			result, err := request.Download()
			retryDelay := time.Duration(retryDelaySec) * time.Second
			for retry := 1; err != nil && retry <= retries; retry++ {
				log.Warningf("error downloading data: %v, attempting retry %d of %d in %s\n", err, retry, retries, retryDelay)
				result, err = request.Download()
				time.Sleep(retryDelay)
			}
			if err != nil {
				log.Fatalf("error downloading data: %v, maximum retries reached\n", err)
			}

			// create output
			err = os.WriteFile(outputFile, *result, os.ModePerm)
			if err != nil {
				log.Fatalf("error writing data: %v, could not write file %s\n", err, outputFile)
			}
			outputFiles = append(outputFiles, outputFile)
		}
	}

	log.Infof("all done\n")
	log.Infof("files created: %s\n", strings.Join(outputFiles, ","))
	log.Infof("bounding box: %s\n", getosm.BoundingBoxesToWkt(bboxes))
}

func getFileName(prefix string, bbIndex, bbCount, tileIndex, tileCount int) string {
	var outputFile string
	if tileCount > 1 {
		outputFile = fmt.Sprintf("%s_%d_%d.osm.xml", prefix, tileIndex+1, tileCount)
	} else {
		outputFile = fmt.Sprintf("%s_bbox.osm.xml", prefix)
	}
	if bbCount > 1 {
		outputFile = fmt.Sprintf("bb_%d_%d_%s", bbIndex, bbCount, outputFile) // add bb index if required
	}
	return outputFile
}

func doesFileExist(file string) bool {
	log.Debugf("check if file is already present %s\n", file)
	stats, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	if stats.Size() == 0 {
		return false
	}
	return true
}
