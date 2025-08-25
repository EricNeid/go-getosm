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

func Download(apiURL, query string) (*[]byte, error) {
	log.Debugf("download using query: %s\n", query)
	client := http.Client{
		Timeout: 0, // no timeout
	}
	resp, err := client.Post(apiURL, "text/xml", strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Debugf("reading response body")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		log.Debugf("download complete")
	} else {
		log.Errorf("download failed with status %s\n", resp.Status)
		log.Errorf("response is: \n%s\n", string(body))
		return nil, ErrorDownload
	}

	return &body, err
}
