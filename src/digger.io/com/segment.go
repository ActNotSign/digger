package com

import(
    "net/http"
    "io/ioutil"
    "net/url"
    "encoding/json"
)

type SegmentConfig struct {
    Host    string
}

type SegmentCom struct {
    config  *SegmentConfig
}

type SegmentResult struct {
    Ok      bool
    Words   []string
    Error   string
}

func (s *SegmentCom) InitConn(config *SegmentConfig) {
    s.config = config
}

func (s *SegmentCom) Segment(text string, postagging bool) (SegmentResult) {
    tagging := "false"
    if postagging {
        tagging = "true"
    }
    resp, _ := http.PostForm(
        s.config.Host + "/segment",
        url.Values{"content": {text}, "postagging": {tagging}})
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    //format data
    var result SegmentResult
    json.Unmarshal(body, &result)
    return result
}
