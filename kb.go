package main

import (
    "os"
    "io"
    "bufio"
    "strings"
    "log"
    "github.com/deckarep/golang-set"
    "fmt"
)

const (
    WHERE = 0
    NOTE = "#"
    SPLIT = "\t"
)

func (k *KnowledgeBase) makePath(id string, filename string) (string){
    return fmt.Sprintf("data/%s/%s", id, filename)
}

func (k *KnowledgeBase) loadType2Value(id string) {
    // make
    k.type2value = make (map[string]Type2Value)
    // read
    f, err := os.Open(k.makePath(id, "type2value.txt"))
    if err != nil {
        panic(err)
    }
    defer f.Close()
    rd := bufio.NewReader(f)
    // format
    for {
        line, err := rd.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        } else if (strings.Index(line, NOTE) == WHERE ) {
            continue
        }
        arr := strings.Split(strings.Trim(line, "\n"), SPLIT)
        log.Println(id, "KnowledgeBase.loadType2value->", arr)
        var q Type2Value
        q.TypeId = arr[0]
        q.Values = strings.Split(arr[1], ",")
        k.type2value[q.TypeId] = q
    }
}

func (k *KnowledgeBase) loadIntent(id string) {
    // make
    k.intent = make (map[string]Intent)
    // read
    f, err := os.Open(k.makePath(id, "intent.txt"))
    if err != nil {
        panic(err)
    }
    defer f.Close()
    rd := bufio.NewReader(f)
    // format
    for {
        line, err := rd.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        } else if (strings.Index(line, NOTE) == WHERE ) {
            continue
        }
        arr := strings.Split(strings.Trim(line, "\n"), SPLIT)
        log.Println(id, "KnowledgeBase.loadIntent->", arr)
        var i Intent
        i.ValueId = arr[0]
        i.Types = strings.Split(arr[1], ",")
        k.intent[i.ValueId] = i
    }
}

func (k *KnowledgeBase) loadQuestion(id string) {
    // make
    k.question = make (map[string]Question)
    // read
    f, err := os.Open(k.makePath(id, "question.txt"))
    if err != nil {
        panic(err)
    }
    defer f.Close()
    rd := bufio.NewReader(f)
    // format
    for {
        line, err := rd.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        } else if (strings.Index(line, NOTE) == WHERE ) {
            continue
        }
        arr := strings.Split(strings.Trim(line, "\n"), SPLIT)
        log.Println(id, "KnowledgeBase.loadQuestion->", arr)
        var q Question
        q.TypeId = arr[0]
        q.QuestionType = arr[1]
        q.Text = arr[2]
        k.question[q.TypeId] = q
    }
}

func (k *KnowledgeBase) loadValue(id string) {
    // make
    k.value = make (map[string]Value)
    // read
    f, err := os.Open(k.makePath(id, "value.txt"))
    if err != nil {
        panic(err)
    }
    defer f.Close()
    rd := bufio.NewReader(f)
    // format
    for {
        line, err := rd.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        } else if (strings.Index(line, NOTE) == WHERE ) {
            continue
        }
        arr := strings.Split(strings.Trim(line, "\n"), SPLIT)
        log.Println(id, "KnowledgeBase.loadValue->", arr)
        var r Value
        r.Id = arr[0]
        r.Text = arr[1]
        k.value[r.Id] = r
    }
}

func (k *KnowledgeBase) loadType(id string) {
    // make
    k.types = make (map[string]Type)
    // read
    f, err := os.Open(k.makePath(id, "type.txt"))
    if err != nil {
        panic(err)
    }
    defer f.Close()
    rd := bufio.NewReader(f)
    // format
    for {
        line, err := rd.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        } else if (strings.Index(line, NOTE) == WHERE ) {
            continue
        }
        arr := strings.Split(strings.Trim(line, "\n"), SPLIT)
        log.Println(id, "KnowledgeBase.loadType->", arr)
        var r Type
        r.Id = arr[0]
        r.Text = arr[1]
        k.types[r.Id] = r
    }
}

func (k *KnowledgeBase) loadReply(id string) {
    // make
    k.reply = make (map[string]Reply)
    docs := make(map[string][]string)
    // read
    f, err := os.Open(k.makePath(id, "reply.txt"))
    if err != nil {
        panic(err)
    }
    defer f.Close()
    rd := bufio.NewReader(f)
    // format
    for {
        line, err := rd.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        } else if (strings.Index(line, NOTE) == WHERE ) {
            continue
        }
        arr := strings.Split(strings.Trim(line, "\n"), SPLIT)
        log.Println(id, "KnowledgeBase.loadReply->", arr)
        var r Reply
        r.Id = arr[0]
        r.ValueId = arr[1]
        r.TypeId = arr[2]
        r.Action = arr[3]
        r.Answer = strings.Split(arr[4], "#")
        k.reply[r.Id] = r
        //init similarity
        result := segmentCom.Segment(strings.Replace(arr[4], NOTE, " ", -1), false)
        docs[r.Id] = result.Words
    }
    // weight
    k.similarity.ComputingWeight(&docs)
}

func (k *KnowledgeBase) loadStory(id string) {
    // make
    k.story = make (map[string][]Story)
    // read
    f, err := os.Open(k.makePath(id, "story.txt"))
    if err != nil {
        panic(err)
    }
    defer f.Close()
    rd := bufio.NewReader(f)
    // format
    for {
        line, err := rd.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        } else if (strings.Index(line, NOTE) == WHERE ) {
            continue
        }
        arr := strings.Split(strings.Trim(line, "\n"), SPLIT)
        log.Println(id, "KnowledgeBase.loadStory->", arr)
        var s Story
        s.ValueId = arr[0]
        s.Factor = make(map[string]string)
        for _, a := range strings.Split(arr[1], ",") {
            t := strings.Split(a, ":")
            if len(t) == 2 {
                s.Factor[t[0]] = t[1]
            }
        }
        s.Type = arr[2]
        s.Text = arr[3]
        k.story[s.ValueId] = append(k.story[s.ValueId], s)
    }
}

func (k *KnowledgeBase) loadSet(id string) {
    k.set = make (map[string]mapset.Set)
    // read
    f, err := os.Open(k.makePath(id, "set.txt"))
    if err != nil {
        panic(err)
    }
    defer f.Close()
    rd := bufio.NewReader(f)
    // format
    for {
        line, err := rd.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        } else if (strings.Index(line, NOTE) == WHERE ) {
            continue
        }
        arr := strings.Split(strings.Trim(line, "\n"), SPLIT)
        log.Println(id, "KnowledgeBase.loadSet->", arr)
        if _, ok := k.set[arr[0]]; !ok {
            k.set[arr[0]] = mapset.NewSet()
        }
        (k.set[arr[0]]).Add(arr[1])
    }
}


func (k *KnowledgeBase) loadTemplate(id string) {
    k.template = make (map[string]string)
    // read
    f, err := os.Open(k.makePath(id, "template.txt"))
    if err != nil {
        panic(err)
    }
    defer f.Close()
    rd := bufio.NewReader(f)
    // format
    for {
        line, err := rd.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        } else if (strings.Index(line, NOTE) == WHERE ) {
            continue
        }
        arr := strings.Split(strings.Trim(line, "\n"), SPLIT)
        log.Println(id, "KnowledgeBase.loadTemplate->", arr)
        k.template[arr[0]] = arr[1]
    }
}

func (k *KnowledgeBase) Init(id string) (*KnowledgeBase) {
    k.loadStory(id)
    k.loadQuestion(id)
    k.loadType(id)
    k.loadValue(id)
    k.loadIntent(id)
    k.loadType2Value(id)
    k.loadReply(id)
    k.loadSet(id)
    k.loadTemplate(id)
    return k
}
