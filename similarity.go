package main

import(
    "github.com/deckarep/golang-set"
    "math"
)

type Similarity struct {
    docWord DocumentWords
}

//
type DocumentWords struct {
    tf      map[string]map[string]float64
    idf     map[string]float64
    tfidf   map[string]map[string]float64
    sumWord float64
    sumDoc  float64
}

// computing words weight
func (d *DocumentWords) Computing(docs *map[string][]string) {
    d.sumDoc = float64(len(*docs))
    d.tf = make(map[string]map[string]float64)
    d.idf = make(map[string]float64)
    d.tfidf = make(map[string]map[string]float64)
    // computing tf
    for k, v := range (*docs) {
        if _, ok := d.tf[k]; !ok {
            d.tf[k] = make(map[string]float64)
        }
        for _, w := range v {
            if _, ok := d.tf[k][w]; !ok {
                d.tf[k][w] = 1.
            } else {
                d.tf[k][w] += 1.
            }
            d.sumWord++
        }
    }
    for k, v := range d.tf {
        for w, i := range v {
            d.tf[k][w] = i / d.sumWord
            d.idf[w] = 0.
        }
    }
    // get idf
    for w, _ := range d.idf {
        for _, v := range d.tf {
            if _, ok := v[w]; ok {
                d.idf[w] += 1.
            }
        }
    }
    for w, v := range d.idf {
        d.idf[w] = math.Log(d.sumDoc / v)
    }
    //get tfidf
    for k, v := range d.tf {
        if _, ok := d.tfidf[k]; !ok {
            d.tfidf[k] = make(map[string]float64)
        }
        for w, f := range v {
            d.tfidf[k][w] = f * d.idf[w]
        }
    }
}

func (d *DocumentWords) Weight(doc string, word string) (float64) {
    if _, ok := d.tfidf[doc][word]; ok && word != "哪里"{
        return d.tfidf[doc][word]
    }
    return .0000000001
}

// computing similarity
func (s *Similarity) Weight(doc string, word string) (float64) {
    return s.docWord.Weight(doc, word)
}

func (s *Similarity) ComputingWeight(docs *map[string][]string) {
    s.docWord.Computing(docs)
}

func (s *Similarity) Jaccard(farr []string, sarr []string) (float64) {
    fm := mapset.NewSet()
    sm := mapset.NewSet()
    for _, o := range farr {
        fm.Add(o)
    }
    for _, o := range sarr {
        sm.Add(o)
    }
    return float64(len(fm.Intersect(sm).ToSlice()))/float64(len(fm.Union(sm).ToSlice()))
}

func (s *Similarity) Cosine(farr []string, sarr []string, doc string) (float64) {
    var (
        // computing params
        x = .0
        y = .0
        z = .0
    )
    // to map set
    fset := mapset.NewSet()
    sset := mapset.NewSet()
    for _, o := range farr {
        fset.Add(o)
    }
    for _, o := range sarr {
        sset.Add(o)
    }
    // to word vector
    fmap := make(map[string]float64)
    smap := make(map[string]float64)
    for _, v := range farr {
        if _, ok := fmap[v]; ok {
            fmap[v] += s.Weight(doc, v)
        } else {
            fmap[v] = s.Weight(doc, v)
        }
    }
    for _, v := range sarr {
        if _, ok := smap[v]; ok {
            smap[v] += s.Weight(doc, v)
        } else {
            smap[v] = s.Weight(doc, v)
        }
    }
    // vector value
    fv := make([]float64, len(fset.Union(sset).ToSlice()))
    sv := make([]float64, len(fset.Union(sset).ToSlice()))
    var k = 0
    for v := range fset.Union(sset).Iter() {
        // interface to string
        tmp := v.(string)
        // update vector
        if _, ok := fmap[tmp]; ok {
            fv[k] = fmap[tmp]
        } else {
            fv[k] = .0
        }
        if _, ok := smap[tmp]; ok {
            sv[k] = smap[tmp]
        } else {
            sv[k] = .0
        }
        // index
        k++
    }
    for k, _ = range fv {
        x += math.Pow(fv[k], 2)
        y += math.Pow(sv[k], 2)
        z += fv[k] * sv[k]
    }
    // cosine
    return z / (math.Sqrt(x) * math.Sqrt(y))
}
