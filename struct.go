package main

import(
    "github.com/deckarep/golang-set"
)

type Question struct {
    TypeId          string
    QuestionType    string
    Text            string
}

type Reply  struct {
    Id      string
    ValueId string
    TypeId  string
    Action  string
    Answer  []string
}

type Type struct {
    Id      string
    Text    string
}

type Value struct {
    Id      string
    Text    string
}

type Type2Value struct {
    TypeId  string
    Values  []string
}

type Intent struct {
    ValueId string
    Types   []string
}

type Story struct {
    ValueId string
    Factor  map[string]string
    Type    string
    Text    string
}

type Session struct {
    Id              string
    ValueId         string
    AdditionValueId string
    TypeMap         map[string]string
    AdditionTypeMap map[string]string
    Pointer         Question
    IsAddition      bool
}

type KnowledgeBase struct {
    question        map[string]Question
    reply           map[string]Reply
    types           map[string]Type
    value           map[string]Value
    type2value      map[string]Type2Value
    intent          map[string]Intent
    story           map[string][]Story
    set             map[string]mapset.Set
    template        map[string]string
    similarity      Similarity
    numExt          NumericalExtract
}

type TriggerMessage struct {
    kb          *KnowledgeBase
    session     Session
}

