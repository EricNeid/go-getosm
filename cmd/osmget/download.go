package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func formatQuery(bb boundingBox) string {
	return fmt.Sprintf(`
	<osm-script timeout="%d" element-limit="%d">
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
	`, timeout, elementLimit, bb.n, bb.s, bb.w, bb.e)
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
