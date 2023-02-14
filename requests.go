package requests

import (
	"net/http"
	"net/url"
	"crypto/tls"
	"strings"
	"bytes"
	"log"
	"fmt"
	"time"
	"io"
	"errors"

)


/****** POST METHOD ************************/
func Post(uri string, data interface{}, headers map[string]string,cookie []*http.Cookie,verify bool)(resp *http.Response,err error){

	var payload io.Reader 

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	if verify == false {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}
	
	switch v := data.(type){
	case map[string]string :
		
		params := url.Values{}
		if len(v) > 0  {
			for key,val := range v {
				params.Add(key,val)
			}
		}
		payload = strings.NewReader(params.Encode())

	case []byte : 
		tmp := []byte(fmt.Sprintf("%s", v))
		payload = bytes.NewBuffer(tmp)

	default :
		err = errors.New(fmt.Sprintf("Error interface format not known : %v\n ",data))
		return 
	}

	log.Printf("%s",payload)
	
    req, err := http.NewRequest("POST", uri, payload )

	//Add header
	for k,v := range headers {
		req.Header.Set(k,v)
	}
	
	//Add cookie
	if cookie != nil {
		for _,c := range cookie {
			req.AddCookie(c)
		}
	}

	resp, err = client.Do(req)
	if err != nil { return nil,err }

	return resp,nil
}


func CheckResponse(msg string)bool{
	if strings.Contains(msg,"Error"){ return false}
	return true
}



/**************** GET METHOD ******************/
func Get(uri string, params map[string]string,headers map[string]string,cookie []*http.Cookie, verify bool)(resp *http.Response,err error){
	var (
		u string
	)
	if len(params)> 0 {
		p := url.Values{}
		for k,v := range params {
			p.Add(k,v)
		}
		if strings.Contains(uri,"?") {
			u = fmt.Sprintf("%s&%s",uri,p.Encode())
		}else{
			u = fmt.Sprintf("%s?%s",uri,p.Encode())
		}
	}else{
		u = uri
	}
	req, err := http.NewRequest("GET", u,nil)
	if err != nil { log.Println(err) ;return nil,err}

    //Add header
	for k,v := range headers {
		req.Header.Set(k,v)
	}
	//Add cookie
	if cookie != nil {
		for _,c := range cookie {
			req.AddCookie(c)
		}
	}
    client := &http.Client{
		Timeout: 15 * time.Second,
	}
	if verify == false {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	resp, err = client.Do(req)
	return 
}

func PatchJson(uri string, param map[string]string, jsonPost []byte, headers map[string]string)(*http.Response,error){

	//Construct URI
	var u string
	params := url.Values{}
	if len(param) > 0  {
		for k,v := range param {
			params.Add(k,v)
		}
		u = fmt.Sprintf("%s?%s",uri,params.Encode())
	}else{
		u = uri
	}

	req, err := http.NewRequest("PATCH", u,bytes.NewBuffer(jsonPost))
	if err != nil { return nil,err}

	//Add header
	for k,v := range headers {
		req.Header.Set(k,v)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil { return nil,err }

	return resp,nil
}
