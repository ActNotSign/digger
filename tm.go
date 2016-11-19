package main

import(
    "log"
    "errors"
    "strconv"
    "strings"
    "reflect"
    "time"
    "regexp"
)

const (
    FORMAT_DATE = "2006-01-02"
    FORMAT_ALL = "2006-01-02 15:04:05"
    FORMAT_EMPTY = "0000-00-00"

    REPLY_MODEL_SIMILARITY = "similarity"
    REPLY_MODEL_ABSOLUTE = "absolute"

    MSG_NO_FOUND = "没有找到相关信息，暂时无法回答您"
    MSG_END = "我还有什么可以帮您的吗？"

    STORY_TYPE_ADDITION = "addition"
    STORY_TYPE_OUTPUT = "output"
    STORY_TYPE_ASSIGN = "assign"

    QUESTION_TYPE_QUESTION = "question"
    QUESTION_TYPE_ENV = "env"
    QUESTION_TYPE_TIME = "time"
    QUESTION_PERIOD = "。"

    OPERATOR_IN = "in"
    OPERATOR_NOTIN = "!in"
    OPERATOR_GREATER_THAN = ">"
    OPERATOR_LESS_THAN = "<"
    OPERATOR_TEMPLATE = "{*}"

    M_NAME = "\033[31m 机器人：\033[0m \033[32m"
    U_NAME = "\033[0m 您："

    SIMILARITY_BEFOREHAND_REPLY = .4
    SIMILARITY_BEFOREHAND_INTENT = .4
)

func (t *TriggerMessage) numeric(text string, r Reply) (bool) {
    log.Println("TirggerMessage.numeric.{param}.text", text)
    switch r.Action {
        case OPERATOR_IN:
            if len(r.Answer) == 2 {
                min, _ := strconv.ParseFloat(r.Answer[0], 64)
                max, _ := strconv.ParseFloat(r.Answer[1], 64)
                f, err := (*t.kb).numExt.Number(text)
                if min <= f && f <= max && err == nil{
                    return true
                }
            }
            break

        case OPERATOR_NOTIN:
            if len(r.Answer) == 2 {
                min, _ := strconv.ParseFloat(r.Answer[0], 64)
                max, _ := strconv.ParseFloat(r.Answer[1], 64)
                f, err := (*t.kb).numExt.Number(text)
                if min > f || f > max && err == nil{
                    return true
                }
            }
            break

        case OPERATOR_GREATER_THAN:
            if len(r.Answer) == 1 {
                min, _ := strconv.ParseFloat(r.Answer[0], 64)
                f, err := (*t.kb).numExt.Number(text)
                if min <= f && err == nil {
                    return true
                }
            }
            break

        case OPERATOR_LESS_THAN:
            if len(r.Answer) == 1 {
                max, _ := strconv.ParseFloat(r.Answer[0], 64)
                f, err := (*t.kb).numExt.Number(text)
                if max >= f && err == nil {
                    return true
                }
            }
            break
    }
    return false
}

func (t *TriggerMessage) time(time_str string, r Reply) (bool) {
    today_date := time.Unix(time.Now().Unix(), 0).Format(FORMAT_DATE)
    tm, err := time.Parse(FORMAT_ALL, time_str)
    if err != nil {
        return false
    }
    switch r.Action {
        case OPERATOR_IN:
            if len(r.Answer) == 2 {
               min := strings.Replace(r.Answer[0], FORMAT_EMPTY, today_date, -1)
               max := strings.Replace(r.Answer[1], FORMAT_EMPTY, today_date, -1)
               m1, err := time.Parse(FORMAT_ALL, strings.TrimSpace(min))
               m2, err := time.Parse(FORMAT_ALL, strings.TrimSpace(max))
               if tm.After(m1) && tm.Before(m2) && err == nil{
                    return true
               }
            }
            break

        case OPERATOR_GREATER_THAN:
            if len(r.Answer) == 1 {
               min := strings.Replace(r.Answer[0], FORMAT_EMPTY, today_date, -1)
               m1, err := time.Parse(FORMAT_ALL, strings.TrimSpace(min))
               if tm.After(m1) && err == nil {
                    return true
               }
            }
            break

        case OPERATOR_LESS_THAN:
            if len(r.Answer) == 1 {
               max := strings.Replace(r.Answer[1], FORMAT_EMPTY, today_date, -1)
               m2, err := time.Parse(FORMAT_ALL, strings.TrimSpace(max))
               if tm.Before(m2) && err == nil {
                    return true
               }
            }
            break
    }
    return false
}

func (t *TriggerMessage) setc(text string, r Reply) (bool) {
    switch r.Action {
        case OPERATOR_IN:
            if (*t.kb).set[r.Answer[0]].Contains(text) {
                return true
            }
            break

        case OPERATOR_NOTIN:
            if !(*t.kb).set[r.Answer[0]].Contains(text) {
                return true
            }
            break
    }
    return false
}

func (t *TriggerMessage) getStartQuestion() (Question, error) {
    if _, ok := (*t.kb).question["0"]; ok {
        return (*t.kb).question["0"], nil
    }
    return Question{}, errors.New(" no found star(*t.kb).question")
}

func (t *TriggerMessage) isIntent(valueId string) (bool) {
    for _, v := range ((*t.kb).type2value["3"]).Values {
        if v == valueId {
            return true
        }
    }
    return false
}

func (t *TriggerMessage) NewSession() {
    token, err := GenerateRandomString(32)
    if err == nil {
        t.session.Id = token
        t.session.ValueId = ""
        t.session.AdditionValueId = ""
        t.session.TypeMap = make(map[string]string)
        t.session.AdditionTypeMap = make(map[string]string)
        t.session.IsAddition = false
        t.session.Pointer = Question{}
    } else {
        log.Println("TriggerMessage.NewSession->", "new session", err)
    }
}

func (t *TriggerMessage) matchTemplate(input string, temp string) (string, []string) {
    exp := ""
    for _, v := range (*t.kb).template {
        exp = strings.Replace(temp, v, `(.*?)`, -1)
    }
    reg := regexp.MustCompile(exp)
    result := reg.FindStringSubmatch(input)

    for k := 1; k < len(result); k++ {
        temp = strings.Replace(exp, `(.*?)`, result[k], -1)
    }
    if len(result) > 1 {
        return temp, result[1:]
    } else {
        return temp, result
    }
}

func (t *TriggerMessage) matchIntent(text string) (string, []string) {
    var (
        maxSimi = 0.0
        vid = ""
        entitys []string
    )
    iwords := segmentCom.Segment(text, false)
    for _, v := range (*t.kb).reply {
        if t.isIntent(v.ValueId) {
            for _, s := range v.Answer {
                var tmpEntitys []string
                if v.Action == OPERATOR_TEMPLATE {
                    s, tmpEntitys = t.matchTemplate(text, s)
                }
                owords := segmentCom.Segment(s, false)
                tmpSimi := (*t.kb).similarity.Cosine(iwords.Words, owords.Words, v.Id)
                if tmpSimi > maxSimi {
                    maxSimi = tmpSimi
                    vid = v.ValueId
                    entitys = tmpEntitys
                    log.Println("TriggerMessage.matchIntent->", text, s, tmpSimi)
                }
            }
        }
    }

    if maxSimi > SIMILARITY_BEFOREHAND_INTENT {
        t.NewSession()
        t.session.ValueId = vid
        for _, tid := range ((*t.kb).intent[vid]).Types {
            if strings.Trim(tid, " ") != "" {
                t.session.TypeMap[tid] = ""
            }
        }
    }

    return vid, entitys
}

func (t *TriggerMessage) isValidValueId(vid string) (bool) {
    if t.session.Pointer.TypeId == "" {
        return true
    }
    log.Println("TriggerMessage.isValidValueId->", (*t.kb).type2value[t.session.Pointer.TypeId].Values) 
    for _, v := range (*t.kb).type2value[t.session.Pointer.TypeId].Values {
        if v == vid {
            return true
        }
    }
    return false
}

func (t *TriggerMessage) matchSimilarityReply(text string) (Reply, error) {
    maxsimi := 0.0
    id := ""
    iwords := segmentCom.Segment(text, false)
    for _, v := range (*t.kb).reply {
        if !t.isValidValueId(v.ValueId) {
            continue
        }
        if !t.isIntent(v.ValueId) && v.TypeId == "1"{
            for _, s := range v.Answer {
                owords := segmentCom.Segment(s, false)
                tmpsimi := (*t.kb).similarity.Cosine(iwords.Words, owords.Words, v.Id)
                if tmpsimi > maxsimi {
                    maxsimi = tmpsimi
                    id = v.Id
                    log.Println("TriggerMessage.matchSimilarityReply->", maxsimi, v.Answer)
                }
            }
        } else if v.TypeId == "2" && t.numeric(text, v) {
            return v, nil
        }
    }
    if maxsimi > SIMILARITY_BEFOREHAND_REPLY {
        return (*t.kb).reply[id], nil
    }

    return Reply{}, errors.New("no match reply")
}

func (t *TriggerMessage) matchAbsoluteReply(text string) (Reply, error) {
    for _, v := range (*t.kb).reply {
        if !t.isValidValueId(v.ValueId) {
            continue
        }
        if v.TypeId == "12" && t.time(text, v) {
            return v, nil
        } else if v.TypeId == "2" && t.numeric(text, v) {
            return v, nil
        } else if v.TypeId == "23" && t.setc(text, v) {
            return v, nil
        }
    }
    return Reply{}, errors.New("no match reply")
}

func (t *TriggerMessage) matchReply(text string, mode string) (Reply, error) {
    var (
        r   Reply
        err error
    )

    switch mode {
        case REPLY_MODEL_SIMILARITY:
            r, err = t.matchSimilarityReply(text)
            break
        case REPLY_MODEL_ABSOLUTE:
            r, err = t.matchAbsoluteReply(text)
            break
    }

    return r, err
}

func (t *TriggerMessage) updateSessionTypeMap(r *Reply) {
    if r.ValueId == "" {
        if !t.session.IsAddition {
            t.session.TypeMap[t.session.Pointer.TypeId] = "0"
        } else {
            t.session.AdditionTypeMap[t.session.Pointer.TypeId] = "0"
        }
    } else if t.session.Pointer.TypeId != "" {
        if !t.session.IsAddition {
            t.session.TypeMap[t.session.Pointer.TypeId] = r.ValueId
        } else {
            t.session.AdditionTypeMap[t.session.Pointer.TypeId] = r.ValueId
        }
    } else {
        for _, v := range (*t.kb).intent[r.ValueId].Types {
            if _, ok := t.session.TypeMap[v]; ok {
                t.session.TypeMap[v] = r.ValueId
                break
            }
        }
    }
}


func (t *TriggerMessage) core(text string) (Question, Story, error, error) {
    var (
        story   Story
        reply   Reply
    )
    storyErr := errors.New("init error")
    replyErr := errors.New("init error")
    log.Println(*t)
    if t.session.ValueId == "" {
        // intent
        _, ent := t.matchIntent(text)
        for _, v := range ent {
            reply, replyErr = t.matchReply(v, REPLY_MODEL_ABSOLUTE)
            if replyErr == nil {
                // update type map
                t.updateSessionTypeMap(&reply)
                log.Println("TriggerMessage.core.{intent}.reply->", reply, replyErr)
                // story
                story, storyErr = t.matchStory()
                log.Println("TrggerMessage.core.{intent}.story->", story, storyErr)
            }
        }
    } else {
        // reply
        reply, replyErr = t.matchReply(text, REPLY_MODEL_SIMILARITY)
        log.Println("TriggerMessage.core.reply->", reply, replyErr)
        if replyErr == nil {
            // update type map
            t.updateSessionTypeMap(&reply)
            // story
            story, storyErr = t.matchStory()
            log.Println("TrggerMessage.core.story->", story, storyErr)
        } else {
            t.matchIntent(text)
        }
    }
    // reset pointer
    t.session.Pointer = Question{}
    tm := &t.session.TypeMap
    if t.session.IsAddition {
        tm = &t.session.AdditionTypeMap
    }
    log.Println("TriggerMessage.core.{param}.tm->", *tm)
    for k, v := range (*tm) {
        if v != "" {
            continue
        }
        if _, ok := (*t.kb).question[k]; ok {
            switch (*t.kb).question[k].QuestionType {
                case QUESTION_TYPE_QUESTION:
                    t.session.Pointer = (*t.kb).question[k]
                    return t.session.Pointer, story, nil, storyErr

                case QUESTION_TYPE_ENV:
                    if !t.session.IsAddition {
                        t.session.TypeMap[(*t.kb).question[k].TypeId] = (*t.kb).question[k].Text
                    } else {
                        t.session.AdditionTypeMap[(*t.kb).question[k].TypeId] = (*t.kb).question[k].Text
                    }
                    break

                case QUESTION_TYPE_TIME:
                    t.session.Pointer = (*t.kb).question[k]
                    reply, replyErr = t.matchReply(time.Unix(time.Now().Unix(), 0).Format(FORMAT_ALL), REPLY_MODEL_ABSOLUTE)
                    t.updateSessionTypeMap(&reply)
                    break
            }
        } else {
            return t.session.Pointer, story, errors.New("no found the question: " + v), storyErr
        }
    }
    // match story
    if storyErr != nil {
        story, storyErr = t.matchStory()
    }
    return t.session.Pointer, story, errors.New("no found question"), storyErr
}

func (t *TriggerMessage) resetAddtionSession() {
    t.session.AdditionValueId = ""
    t.session.AdditionTypeMap = make(map[string]string)
    t.session.IsAddition = false
}

func (t *TriggerMessage) matchStory() (Story, error) {
    // get typemap
    tm := &t.session.TypeMap
    vid := t.session.ValueId
    if t.session.IsAddition {
        tm = &t.session.AdditionTypeMap
        vid = t.session.AdditionValueId
    }
    // match
    for _, v := range (*t.kb).story[vid] {
        // init match
        mtm := make(map[string]string)
        for k, c := range *tm {
            if _, ok := v.Factor[k]; ok {
                mtm[k] = c
            }
        }
        // match
        log.Println("TriggerMessage.matchStory->", mtm, v.Factor)
        if reflect.DeepEqual(mtm, v.Factor) {
            switch v.Type {
                case STORY_TYPE_ADDITION:
                    t.session.AdditionValueId = v.Text
                    for _, tid := range ((*t.kb).intent[v.Text]).Types {
                        if strings.Trim(tid, " ") != "" {
                            t.session.AdditionTypeMap[tid] = ""
                        }
                    }
                    t.session.IsAddition = true
                    break

                case STORY_TYPE_OUTPUT:
                    return v, nil

                case STORY_TYPE_ASSIGN:
                    datArr := strings.Split(v.Text, ":")
                    if _, ok := t.session.TypeMap[datArr[0]]; ok {
                        t.session.TypeMap[datArr[0]] = datArr[1]
                    }
                    t.resetAddtionSession()
                    break
            }
        }
    }
    return Story{}, errors.New("no match story")
}

func (t *TriggerMessage) Trigger (text string) (Question, Story, error, error) {
    q, s, qr, sr := t.core(text)
    if qr == nil {
        //fmt.Println(M_NAME, q.Text)
    } else {
        t.NewSession()
        if sr == nil {
            s.Text = s.Text + QUESTION_PERIOD  + MSG_END
        } else {
            s.Text = MSG_NO_FOUND + QUESTION_PERIOD + MSG_END
        }
        log.Println("TriggerMessage.Trigger->", "not found question", t.session)
    }
    if q.QuestionType !=  QUESTION_TYPE_QUESTION {
        return Question{}, s, qr, sr
    }
    return q, s, qr, sr
}

func (t *TriggerMessage) Init(kb *KnowledgeBase) (*TriggerMessage){
    // init
    t.kb = kb
    t.NewSession()
    return t
}
