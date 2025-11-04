<!--
SPDX-FileCopyrightText: 2021 Eric Neidhardt
SPDX-License-Identifier: CC-BY-4.0
-->
<!-- markdownlint-disable MD041-->
[![Go Report Card](https://goreportcard.com/badge/github.com/EricNeid/go-getosm?style=flat-square)](https://goreportcard.com/report/github.com/EricNeid/go-getosm)
![Test](https://github.com/EricNeid/go-getosm/actions/workflows/tests.yml/badge.svg)
![Linting](https://github.com/EricNeid/go-getosm/actions/workflows/linting.yml/badge.svg)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/EricNeid/go-getosm)
[![Release](https://img.shields.io/github/release/EricNeid/go-getosm.svg?style=flat-square)](https://github.com/EricNeid/go-getosm/releases/latest)
[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/EricNeid/go-getosm)

# About

A simple downloader for OSM xml data, using <https://www.overpass-api.de>.
Software is based on <https://github.com/eclipse/sumo/blob/main/tools/osmGet.py>

## Quickstart

1. Get go from <https://golang.org/dl/>

2. Download osmget

```bash
go install github.com/EricNeid/go-getosm/cmd/osmget@latest
```

## Usage

Download xml for given bounding box:

```bash
osmget -b 13.168487548828123,52.29189255277229,13.278350830078125,52.35211857272093
osmget -b 13.168487548828123,52.29189255277229,13.278350830078125,52.35211857272093 -t 4 -verbose -continue
```

Options:

```bash
osmget -prefix osm2025q3 # prepend downloaded tiles with given label
osmget -t 2              # split output into multiple tiles
osmget -timeout 10       # connection timeout
osmget -retries 10       # retry failed tile download
osmget -retryDelay 5     # wait this many seconds before retry download
osmget -elementLimit 5   # element limit in osm file
osmget -verbose          # extra debug output
osmget -continue         # skip already downloaded files
osmget -h                # print usage
```

## Question or comments

Please feel free to open a new issue:
<https://github.com/EricNeid/go-getosm/issues>
