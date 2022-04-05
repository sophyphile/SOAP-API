package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

var url = "http://www.dneonline.com/calculator.asmx"
var method = "POST"

func main() {
	callSOAPClientSteps()
}

func callSOAPClientSteps() {

	req := populateRequest()

	httpReq, err := generateSOAPRequest(req)
	if err != nil {
		fmt.Println("Some problem occurred in request generation")
	}

	response, err := soapCall(httpReq)
	response = response
	if err != nil {
		fmt.Println("Problem occurred in making a SOAP call")
	}
}

var getTemplate = (`<?xml version="1.0" encoding="utf-8"?>
	<soapenv:Envelope
	xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
	xmlns:xsd="http://www.w3.org/2001/XMLSchema">
	<soapenv:Header/>
	<soapenv:Body>
		<Subtract xmlns="http://tempuri.org/">
		<intA>{{.IntA}}</intA>
		<intB>{{.IntB}}</intB>
		</Subtract>
	</soapenv:Body>
	</soapenv:Envelope>`)

type Request struct {
	//Values are set in below fields as per the request
	IntA int
	IntB int
}

func populateRequest() *Request {
	req := Request{}
	req.IntA = 17
	req.IntB = 3
	return &req
}

type Response struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	SoapBody *SOAPBodyResponse
}

type SOAPBodyResponse struct {
	XMLName xml.Name `xml:"Body"`
	Resp    *SubtractResponse
}

type SubtractResponse struct {
	XMLName        xml.Name `xml:"http://tempuri.org/ SubtractResponse"`
	SubtractResult string   `xml:"SubtractResult"`
}

func generateSOAPRequest(req *Request) (*http.Request, error) {
	// Using the var getTemplate to construct request
	template, err := template.New("InputRequest").Parse(getTemplate)
	if err != nil {
		fmt.Println("Error while marshling object. %s ", err.Error())
		return nil, err
	}

	doc := &bytes.Buffer{}
	// Replacing the doc from template with actual req values
	err = template.Execute(doc, req)
	if err != nil {
		fmt.Println("template.Execute error. %s ", err.Error())
		return nil, err
	}

	buffer := &bytes.Buffer{}
	encoder := xml.NewEncoder(buffer)
	err = encoder.Encode(doc.String())
	if err != nil {
		fmt.Println("encoder.Encode error. %s ", err.Error())
		return nil, err
	}

	r, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(doc.String())))
	if err != nil {
		fmt.Println("Error making a request. %s ", err.Error())
		return nil, err
	}

	return r, nil
}

func soapCall(req *http.Request) (*Response, error) {
	client := &http.Client{}
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("SOAPAction", "http://tempuri.org/Subtract")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("The error is: ", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("The error is: ", err)
		return nil, err
	}

	fmt.Println(string(body))
	defer resp.Body.Close()

	r := &Response{}
	err = xml.Unmarshal(body, &r)

	if err != nil {
		fmt.Println("The error is: ", err)
		return nil, err
	}

	str := fmt.Sprintf("%#v", r.SoapBody)
	fmt.Println(r)
	fmt.Println(str)

	fmt.Println(r.SoapBody)
	fmt.Println(*r.SoapBody)
	fmt.Println(*r.SoapBody.Resp)
	fmt.Println(&r.SoapBody.Resp.SubtractResult)
	fmt.Println(r.SoapBody.Resp.SubtractResult)

	if r.SoapBody.Resp.SubtractResult != "14" {
		fmt.Println("The error is: ", err)
		return nil, err
	}

	return r, nil
}
