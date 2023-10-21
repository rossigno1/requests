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
	"mime/multipart"
	"os"
	"path/filepath"
)

var MAXTIMEOUT = time.Duration(30)

func SetTimeout(max int){
	MAXTIMEOUT = time.Duration(max)
	return
}

/****** POST METHOD ************************/
func Post(uri string, data interface{}, headers map[string]string,cookie []*http.Cookie,verify bool)(resp *http.Response,err error){

	var payload io.Reader 

	client := &http.Client{
		Timeout: MAXTIMEOUT * time.Second,
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
	
    req, err := http.NewRequest("POST", uri, payload )
	if err != nil { return nil,err }

	defer req.Body.Close()
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
		Timeout: MAXTIMEOUT * time.Second,
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
/**************** Head METHOD ******************/
func Head(uri string, params map[string]string,headers map[string]string,cookie []*http.Cookie, verify bool)(resp *http.Response,err error){
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
	req, err := http.NewRequest("HEAD", u,nil)
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
		Timeout: MAXTIMEOUT * time.Second,
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


/****************** PUT Method ******************/
func Put(uri string, data interface{}, headers map[string]string,cookie []*http.Cookie,verify bool)(resp *http.Response,err error){

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
		log.Printf("Error interface format not known : %v\n ",data)
	}

	
    req, err := http.NewRequest("PUT", uri, payload )

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


func Delete(uri string, data interface{}, headers map[string]string,cookie []*http.Cookie,verify bool)(resp *http.Response,err error){

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
	if data != nil {
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
			log.Printf("Error interface format not known : %v ",data)
		}
	}
	
    req, err := http.NewRequest("DELETE", uri, payload )

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


/************ Specific PostFile ***********/
func PostFile(uri string, fname string, headers map[string]string,cookie []*http.Cookie,verify bool)(resp *http.Response,err error){


	file,err := os.Open(fname)
	if err != nil { return }

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil { return }
	io.Copy(part, file)
	writer.Close()

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	if verify == false {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

    req, err := http.NewRequest("POST", uri, payload )

	//Add header
	for k,v := range headers {
		req.Header.Set(k,v)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

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