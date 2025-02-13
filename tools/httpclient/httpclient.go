//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

package httpclient

import (
	"bytes"
	"crypto/tls"
	"log"
	"net/http"
)

func Post(url string, data []byte, headerFN func(header http.Header), responseFN func(resp *http.Response)) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Println("[XTA] HttpClient NewRequest", err)
		return
	}
	if headerFN != nil {
		headerFN(req.Header)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[XTA] HttpClient Do", err)
		return
	}
	defer resp.Body.Close()
	if responseFN != nil {
		responseFN(resp)
	}
}
