package main


import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"bufio"
	"net/http/cookiejar"
	"crypto/tls"
	"strings"
	"strconv"
	"time"
	"log"
)
var (
//the path of file saving the url that need to visite
 filePath = "./url.ini"
//the path fo file saving the config that customize the program
confPath = "./default.conf"
//max times of sending require, when the work is complete the program will exit
maxtimes = 0
//the interval of time between two require 
interval = 0
//max days is the longest time program run,unite is day
maxdays = 0
)

var urlArry []string
var confMap = make(map[string]int)


func main(){
	readConf()
	readUrl()
	if interval == 0 || maxtimes ==0  {
		log.Fatal("the config might not right!")
	}
	ticker := time.NewTicker( time.Minute * time.Duration(interval))
    go func() {
		counter := 0
        for t := range ticker.C {
			counter++
			fmt.Println()
			fmt.Println(counter,"     ", t)
			for _,url := range urlArry {
				visitUrl(url)
				//fmt.Println(url)
			}
			if counter >= maxtimes {
				ticker.Stop()
				os.Exit(0)
			}
        }
	}()
	time.Sleep( 24 * time.Hour * time.Duration(maxdays))
}

//read config from default.conf and save them in map
func readConf(){
	urlFile,err:= os.Open(confPath)	
	if err!=nil{
		fmt.Println(err)
		os.Exit(2)
	}
	defer urlFile.Close()
	buf := bufio.NewReader(urlFile)
	for {
		byteLine, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		strLine := string(byteLine)
		if strLine == "" ||  strLine[0] == '#' {
			continue
		}
		tmpStr := strings.TrimSpace(strLine)
		index := strings.Index(tmpStr, "=")
		if index <= 0 {
			continue
		}
		key := strings.TrimSpace(tmpStr[:index])
		value := strings.TrimSpace(tmpStr[index+1:])
		if len(key) == 0 || len(value) == 0 {
			continue
		}
		int_value,_:= strconv.Atoi(value)
		confMap[key] = int_value
	}
	maxtimes = confMap["maxtimes"]
	interval = confMap["interval"]
	maxdays = confMap["maxdays"]
}

//read url from the filepath and save them into array
func readUrl(){
	urlFile,err:= os.Open(filePath)	
	if err!=nil{
		fmt.Println(err)
		os.Exit(2)
	}
	defer urlFile.Close()
	buf := bufio.NewReader(urlFile)
	for {
		byteLine, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		strLine := string(byteLine)
		if strLine == "" {
			continue
		}
		urlArry = append(urlArry,strLine)
	}
}

//send a get require to the url and get reponse,
//after visit a url , it will printf the url and the byte of reponsed
func visitUrl(url string){
	getReq, _ := http.NewRequest("GET", url, nil)
	CurCookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	clen:=http.Client{
		Jar: CurCookieJar,
		Transport: tr,
	}
	resp, _ := clen.Do(getReq)
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(url,"     ",len(data))
}