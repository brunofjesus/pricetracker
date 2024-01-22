package categories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Parser() {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://www.pcdiga.com/api/menu/vertical",
		nil,
	)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-GB,en;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("TE", "trailers")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/117.0")
	req.Header.Add("Accept-Encoding", "deflate, br")
	req.Header.Add("Cookie", "__cf_bm=lmBYE6R6E4ONewyu6E2ezdvKjLKbPRdhlV_TmA2BmJA-1696546480-0-Admm/o8VkCpkEMfrjlwvt9MK5MnV9iEB3D9/PHcQCYi2oTxP+DaQLCF88xfcFQw49iKlKlIbnrOQYlcoKwiHnPo=")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	fmt.Printf("%s", b)

	var response []PcDigaMenuItem
	err = json.NewDecoder(res.Body).Decode(&response)

	fmt.Println(response)
}
