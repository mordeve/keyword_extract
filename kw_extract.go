package main

import (
	"log"
	"regexp"
	"sort"
	"strings"
)

func Unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}

func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func kw_extract(result map[string]interface{},
	stopwordMap map[string]interface{},
	sentence_hyli string) []string {

	// result 		 -> words vs. idf scores
	//  stopwordMap  -> a map for the stopwords
	// sentence_hyli -> input string

	re, err := regexp.Compile("[0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	re_punc, err2 := regexp.Compile("['\"!#$%&()*+,-./:;<=>?@[\\]^_`{|}~']+")
	if err2 != nil {
		log.Fatal(err2)
	}

	// sentence_hyli := "derin öğrenme aynı zamanda derin yapılandırılmış öğrenme Yer Açtı öğrenme ya da derin makine öğrenmesi bir veya daha fazla gizli katman içeren yapay sinir ağları ve benzeri makine öğrenme algoritmaları kapsayan çalışma alanıdır yani En az 1 adet yapay sinir ağı kullanıldığı ve birçok Algoritma ile insan eldeki verilerden yeni veriler elde etmesidir derin öğrenme gözetimli gözetimi ve gözetim Sokak gerçekleştirilebilir yapay sinir ağları pekiştirmeli öğrenme yaklaşımıyla da başarılı sonuç vermiştir yapay sinir ağları biyolojik sistemlerde ki bilgi işleme ve dağıtım iletişim düğünlerinin esinlenilmiştir farkları vardır özellikle sinir ağları statik ve sembolik olma elimdeyken çoğu canlı organizmanın biyolojik beyni dinamik plastik ve"
	cleaned_hyli_num := re.ReplaceAllString(sentence_hyli, " ")
	cleaned_hyli_punc := re_punc.ReplaceAllString(cleaned_hyli_num, " ")
	cleaned_hyli := CleanString(cleaned_hyli_punc, stopwordMap, true)

	split := strings.Split(cleaned_hyli, " ")

	split_un := Unique(split)
	split_un_clean := delete_empty(split_un)

	m := make(map[string]float32)

	for k := range split_un_clean {
		res1 := strings.Count(cleaned_hyli, split_un_clean[k])
		tf := float32(res1) / float32(len(split))
		idf := result[split_un_clean[k]]
		if idf == nil {
			idf = 5.65
		}
		iAreaId := idf.(float64)
		//fmt.Println(idf)
		//fmt.Println(tf)
		//fmt.Println(tf * float32(iAreaId))
		m[split_un_clean[k]] = (tf * float32(iAreaId))
	}

	type kv struct {
		Key   string
		Value float32
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	var hyli_list []string

	for _, kv := range ss {
		hyli_list = append(hyli_list, kv.Key)
		//fmt.Println(kv.Key, kv.Value)
	}
	return hyli_list
}
