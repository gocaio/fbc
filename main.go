/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at
   http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
*/

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	//"net/url"

	"github.com/fatih/color"
)

var (
	yellow  = color.New(color.Bold, color.FgYellow).SprintFunc()
	red     = color.New(color.Bold, color.FgRed).SprintFunc()
	cyan    = color.New(color.Bold, color.FgCyan).SprintFunc()
	green   = color.New(color.Bold, color.FgGreen).SprintFunc()
	blue    = color.New(color.Bold, color.FgBlue).SprintFunc()
	magenta = color.New(color.Bold, color.FgMagenta).SprintFunc()
	black   = color.New(color.FgBlack, color.BgWhite).SprintFunc()
)

var urlFlag = flag.String("url", "http://www.google.com", "Request URL")
var apiFlag = flag.String("api", "", "Firebase.io API key")
var projectFlag = flag.String("project", "", "https://<project>.firebase.com name")
var postFlag = flag.Bool("post", false, "Print POST data used")

type Response struct {
	Error struct {
		Code    int    `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
		Status  string `json:"status,omitempty"`
		Details []struct {
			Type  string `json:"@type,omitempty"`
			Links []struct {
				Description string `json:"description,omitempty"`
				URL         string `json:"url,omitempty"`
			} `json:"links,omitempty"`
		} `json:"details,omitempty"`
	} `json:"error,omitempty"`
}

type Data struct {
	URL string `json:"url"`
	API string `json:"key"`
}

func main() {
	log.SetFlags(0)

	flag.Parse()

	if *apiFlag == "" && *projectFlag == "" {
		flag.PrintDefaults()
		return
	}

	if *apiFlag != "" {
		fmt.Fprintf(color.Output, "\nChecking restrictions on API key %v \n", magenta(*apiFlag))
		fmt.Fprintf(color.Output, "Tested on: %v \n\n", cyan(time.Now().Format(time.RFC1123)))

		if *postFlag == true {
			//fmt.Fprintf(color.Output, "%v \n  {\n    %v: '%s',\n    %v: '%s'\n  }\n", blue("POST Body:"), red(`"url"`), green(*urlFlag), red(`"key"`), green(*apiFlag))
			fmt.Fprintf(color.Output, "%v      { %v: '%s', %v: '%s' }\n", blue("POST Body:"), red(`"url"`), green(*urlFlag), red(`"key"`), green(*apiFlag))
		}

		constructedURL := "https://searchconsole.googleapis.com/v1/urlTestingTools/mobileFriendlyTest:run?key=" + *apiFlag
		body := MakePOSTRequest(constructedURL)

		var response Response
		err := json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal("Error --> ", err)
		}

		var msgRegex = `^.+?(\.)`
		msgMatch := regexp.MustCompile(msgRegex)
		msg := msgMatch.FindString(response.Error.Message)
		fmt.Fprintf(color.Output, "%v %v\n", blue("Response Code: "), response.Error.Code)
		fmt.Fprintf(color.Output, "%v %v\n", blue("Status:        "), response.Error.Status)
		fmt.Fprintf(color.Output, "%v %v\n\n", blue("Message:       "), msg)
	}

	if *projectFlag != "" {
		projectURL := fmt.Sprintf("https://%s.firebase.com/.json", *projectFlag)
		println(projectURL)

	}

}

// MakePOSTRequest will do the POST request
func MakePOSTRequest(host string) (body []byte) {
	//var proxyUrl, err = url.Parse("http://127.0.0.1:8080")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//Proxy:           http.ProxyURL(proxyUrl),
	}
	client := &http.Client{Transport: tr}

	postData := fmt.Sprintf("{'url': '%s','key': '%s'}", *urlFlag, *apiFlag)
	var jsonStr = []byte(postData)

	e, _ := json.Marshal(string(jsonStr))
	req, err := http.NewRequest("POST", host, bytes.NewBuffer(e))
	req.Header.Set("User-Agent", "FireBase Scanner v1.0.0")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	return body

}

// MakeGETRequest will do the GET request
func MakeGETRequest(host string) (body []byte) {
	//var proxyUrl, err = url.Parse("http://127.0.0.1:8080")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//Proxy:           http.ProxyURL(proxyUrl),
	}
	client := &http.Client{Transport: tr}

	postData := fmt.Sprintf("{'url': '%s','key': '%s'}", *urlFlag, *apiFlag)
	var jsonStr = []byte(postData)

	e, _ := json.Marshal(string(jsonStr))
	req, err := http.NewRequest("POST", host, bytes.NewBuffer(e))
	req.Header.Set("User-Agent", "FireBase Scanner v1.0.0")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	return body

}
