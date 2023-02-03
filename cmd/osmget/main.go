// SPDX-FileCopyrightText: Copyright (C) 2009-2021 German Aerospace Center (DLR) and others.
// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: EPL-2.0 OR GPL-2.0-or-later
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const apiURL = "https://www.overpass-api.de/api/interpreter"

var (
	bbox   = ""
	prefix = "osm"
)

var errorInvalidBB = errors.New("invalid bounding box given")
var errorDownload = errors.New("could not retrieve data")

func init() {
	log.SetFlags(0)
	log.SetPrefix("osmGet: ")

	flag.StringVar(&bbox, "b", bbox, "Bounding box: west,south,east,north")
	flag.StringVar(&prefix, "prefix", prefix, "Prefix of output file")
	flag.Parse()
}

func main() {
	w, s, e, n, err := readBoundingBox()
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	query := formatQuery(w, s, e, n)
	res, err := download(query)
	if err != nil {
		log.Fatal(err)
	}

	output := fmt.Sprintf("%s_bbox.osm.xml", prefix)
	os.WriteFile(output, *res, os.ModeAppend)
}

func readBoundingBox() (w, s, e, n float64, err error) {
	bb := strings.Split(bbox, ",")
	if len(bb) != 4 {
		log.Println("invalid bounding box given: expecting w,s,e,n")
		return w, s, e, n, errorInvalidBB
	}
	w, err = strconv.ParseFloat(bb[0], 64)
	if err != nil {
		log.Printf("could not parse west %s\n", bb[0])
		return w, s, e, n, errorInvalidBB
	}
	s, err = strconv.ParseFloat(bb[1], 64)
	if err != nil {
		log.Printf("could not parse south %s\n", bb[1])
		return w, s, e, n, errorInvalidBB
	}
	e, err = strconv.ParseFloat(bb[2], 64)
	if err != nil {
		log.Printf("could not parse east %s\n", bb[2])
		return w, s, e, n, errorInvalidBB
	}
	n, err = strconv.ParseFloat(bb[3], 64)
	if err != nil {
		log.Printf("could not parse north %s\n", bb[3])
		return w, s, e, n, errorInvalidBB
	}
	return w, s, e, n, nil
}

func formatQuery(w, s, e, n float64) string {
	return fmt.Sprintf(`
	<osm-script timeout="240" element-limit="1073741824">
	<union>
		<bbox-query n="%f" s="%f" w="%f" e="%f"/>
		<recurse type="node-relation" into="rels"/>
		<recurse type="node-way"/>
		<recurse type="way-relation"/>
	</union>
	<union>
		<item/>
		<recurse type="way-node"/>
	</union>
	<print mode="body"/>
	</osm-script>
	`, n, s, w, e)
}

func download(query string) (*[]byte, error) {
	log.Printf("download using query: %s\n", query)

	client := http.Client{
		Timeout: 0, // no timeout
	}
	resp, err := client.Post(apiURL, "text/xml", strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Println("reading response body")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		log.Println("download complete")
	} else {
		log.Printf("download failed with status %s\n", resp.Status)
		log.Println("response is:")
		log.Println(string(body))

		return nil, errorDownload
	}

	return &body, err
}
