package keyword_extract

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/mordeve/stopwords"
)

func stemmer(keyword string) map[string]string {
	requestBody, err := json.Marshal(map[string]string{
		"input": keyword,
	})

	if err != nil {
		log.Fatalln(err)
	}

	timeout := time.Duration(10 * time.Second)
	client := http.Client{Timeout: timeout}

	request, err := http.NewRequest("POST",
		"http://localhost:5000/predict/",
		bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}
	// log.Println(string(body))
	var results map[string]string
	json.Unmarshal([]byte(body), &results)

	return results
}

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

func getStem(s []string) []string {
	var res []string
	for _, word := range s {
		res = append(res, stemmer(word)["stem"])
	}
	return res
}

func Extract(result map[string]interface{},
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

	cleaned_hyli_num := re.ReplaceAllString(sentence_hyli, " ")
	cleaned_hyli_punc := re_punc.ReplaceAllString(cleaned_hyli_num, " ")
	cleaned_hyli := stopwords.CleanString(cleaned_hyli_punc, stopwordMap, true)

	split := strings.Split(cleaned_hyli, " ")

	split_un := Unique(split)
	split_un_clean := delete_empty(split_un)
	split_un_stem := getStem(split_un_clean)

	m := make(map[string]float32)

	for k := range split_un_stem {
		res1 := strings.Count(cleaned_hyli, split_un_stem[k])
		tf := float32(res1) / float32(len(split))
		fmt.Println("word: ", split_un_stem[k])
		fmt.Println("TF: ", tf)
		idf := result[split_un_stem[k]]
		if idf == nil {
			idf = 5.65
		}
		iAreaId := idf.(float64)
		//fmt.Println(idf)
		//fmt.Println(tf)
		//fmt.Println(tf * float32(iAreaId))
		m[split_un_stem[k]] = (tf * float32(iAreaId))
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
