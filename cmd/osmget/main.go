// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt
// Inspired by SUMO's osmGet.py (no code copied).
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"time"

	getosm "github.com/EricNeid/go-getosm"
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

var log = getosm.Log

func main() {
	parseArguments()

	if verbose {
		getosm.SetLogLevel(slog.LevelDebug)
	} else {
		getosm.SetLogLevel(slog.LevelInfo)
	}

	if err := run(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

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
		log.Error("no bounding box given, please specify at least one")
		flag.Usage()
		os.Exit(1)
	}

	if tilesX*tilesY < 1 {
		log.Info("invalid number of tiles given, must be >= 1")
		flag.Usage()
		os.Exit(1)
	}
}

// run downloads the tiles for all given bounding boxes.
// An error is returned if a download or a file operation failed.
func run() error {
	// setup context, so downloads can be aborted by the user (ctrl+c)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var outputFiles []string
	for bbIndex, bb := range bboxes {
		bbTiles, err := getosm.ReadBoundingBox(bb, tilesX, tilesY)
		if err != nil {
			flag.Usage()
			return fmt.Errorf("invalid bounding box %s: %w", bb, err)
		}

		log.Info("processing bounding box", "current", bbIndex+1, "total", len(bboxes))
		for tileIndex, tile := range bbTiles {
			log.Info("downloading tiles", "current", tileIndex+1, "total", len(bbTiles))

			// output file name
			outputFile := getFileName(prefix, bbIndex, len(bboxes), tileIndex, len(bbTiles))

			// check if file already exists
			if continueLastFailed && doesFileExist(outputFile) {
				log.Info("skipping already present tile", "file", outputFile)
				continue
			}

			// download tile until retries are exhausted
			request := getosm.Request{
				Query:           getosm.FormatQuery(tile, timeout, elementLimit),
				CustomUserAgent: userAgent,
				URL:             url,
			}

			result, err := request.Download(ctx)
			retryDelay := time.Duration(retryDelaySec) * time.Second
			for retry := 1; err != nil && retry <= retries && ctx.Err() == nil; retry++ {
				log.Warn("error downloading data, attempting retry", "error", err, "attempt", retry, "retries", retries, "delay", retryDelay)
				select {
				case <-time.After(retryDelay):
				case <-ctx.Done():
					return fmt.Errorf("interrupted: %w", ctx.Err())
				}
				result, err = request.Download(ctx)
			}
			if err != nil {
				if ctx.Err() != nil {
					return fmt.Errorf("interrupted: %w", ctx.Err())
				}
				return fmt.Errorf("error downloading data: %w, maximum retries reached", err)
			}

			// create output
			err = os.WriteFile(outputFile, result, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error writing data: %w, could not write file %s", err, outputFile)
			}
			outputFiles = append(outputFiles, outputFile)
		}
	}

	log.Info("all done")
	log.Info("files created", "files", strings.Join(outputFiles, ","))
	log.Info("bounding box", "wkt", getosm.BoundingBoxesToWkt(bboxes))
	return nil
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
	log.Debug("check if file is already present", "file", file)
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
