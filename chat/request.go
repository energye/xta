//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

package chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/energye/xta/tools/httpclient"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// HttpPost 发送 Http Post 请求
func HttpPost(ai IAI) {
	url := ai.API()
	metaData := ai.MetaData()
	metaData.Messages = *ai.History()
	if Debug {
		log.Println("[Debug] [XTA] - HttpPost URL:", url)
	}
	reqBody, err := json.Marshal(metaData)
	if err != nil {
		log.Println("[XTA] HttpPost Marshal-Body", err)
		return
	}
	if Debug {
		log.Println("[Debug] [XTA] - Request-Body-Size:", len(reqBody))
		log.Println("[Debug] [XTA] - Request-Body-Content:", string(reqBody))
	}
	httpclient.Post(url, reqBody, func(header http.Header) {
		for key, value := range ai.Header() {
			header.Add(key, value[0])
		}
	}, func(resp *http.Response) {
		errorcall := func(code int, type_ string, err error) {
			log.Println("[ERROR] [XTA] HttpPost Response", err)
			if ai.OnFail() != nil {
				ai.OnFail()(&TResponseError{Code: strconv.Itoa(code), Message: err.Error(), Type: type_})
			}
			return
		}
		if resp.StatusCode != 200 {
			if ai.OnFail() != nil {
				data, err := io.ReadAll(resp.Body)
				if err != nil {
					errorcall(resp.StatusCode, "response_read_body", err)
					return
				}
				merr := make(map[string]any)
				err = json.Unmarshal(data, &merr)
				if err != nil {
					errorcall(resp.StatusCode, "response_unmarshal_json", err)
					return
				}
				terr := &TResponseError{Other: make(TResponseErrorOther)}
				for key, value := range merr {
					if key == "message" {
						terr.Message = fmt.Sprintf("%v", value)
					} else if key == "type" {
						terr.Type = fmt.Sprintf("%v", value)
					} else if key == "code" {
						terr.Code = fmt.Sprintf("%v", value)
					} else {
						terr.Other[key] = fmt.Sprintf("%v", value)
					}
				}
				ai.OnFail()(terr)
			}
			return
		}
		respBody := resp.Body
		if metaData.Stream {
			buffer := make([]byte, 1024)
			contents := bytes.Buffer{}
			tmpBuf := bytes.Buffer{}
			callsuccess := func() {
				defer tmpBuf.Reset()
				if tmpBuf.Len() == 0 {
					return
				}
				if Debug {
					log.Println("[Debug] [XTA] - Message:", tmpBuf.String())
				}
				if tmpBuf.Len() > 6 {
					value := tmpBuf.Bytes() // 去除 "data: "
					if string(value[:6]) == "data: " {
						value = value[6:]
					}
					if string(value) == "[DONE]" {
						newHistory := Message{
							Role:    RoleAssistant,
							Content: contents.String(),
						}
						ai.History().Add(newHistory)
						ai.OnReceive()(nil)
						return
					}
					response := &TResponse{}
					err = json.Unmarshal(value, response)
					if err != nil {
						errorcall(200, "response_unmarshal_json", err)
						return
					}
					if response.Error == "" {
						contents.WriteString(response.Choices.ToString())
					}
					if ai.OnReceive() != nil {
						ai.OnReceive()(response)
					}
				} else {
					errorcall(200, "response_data", errors.New(tmpBuf.String()))
				}
			}
			for {
				n, err := respBody.Read(buffer)
				if err != nil {
					if err == io.EOF {
						break
					}
					errorcall(200, "response_read ", err)
					return
				}
				for _, v := range buffer[:n] {
					if v == '\n' {
						callsuccess()
					} else {
						tmpBuf.WriteByte(v)
					}
				}
			}
			if tmpBuf.Len() > 0 {
				callsuccess()
			}
		} else {
			data, err := ioutil.ReadAll(respBody)
			if err != nil {
				errorcall(200, "response_read ", err)
				return
			}
			response := &TResponse{}
			err = json.Unmarshal(data, response)
			if err != nil {
				errorcall(200, "response_unmarshal_json ", err)
				return
			}
			if ai.OnReceive() != nil {
				ai.OnReceive()(response)
			}
		}
	})
}
