package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	// "net/url"
)

type Client struct {
	Suburl  string
	Pushurl string
	Signurl string
	Cname   string
	Token   string
	Seq     float64
}

func (cl *Client) Sub() []map[string]interface{} {
	var err error
	// url := fmt.Sprintf("%s?cb=%s", cl.suburl, "test")
	// params := url.Values{}
	// params.Set("cname", cl.Cname)
	// params.Add("seq", cl.Seq)
	// params.Add("token", cl.Token)

	geturl := fmt.Sprintf("%s?cname=%s&seq=%.f&token=%s", cl.Suburl, cl.Cname, cl.Seq, cl.Token)

	res, err := http.Get(geturl)
	if err != nil {
		fmt.Println("get err", err.Error())
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read err")
		panic(err)
	}

	bodyStr := string(body)
	bodyStr = strings.TrimSpace(bodyStr)
	// fmt.Println(bodyStr)
	bodyStr = strings.TrimPrefix(bodyStr, "(")
	// fmt.Println(bodyStr)
	bodyStr = strings.TrimSuffix(bodyStr, ");")

	// fmt.Println(bodyStr)
	if !strings.HasPrefix(bodyStr, "[") {
		bodyStr = fmt.Sprintf("[%s]", bodyStr)
	}

	var data []map[string]interface{}

	if err = json.Unmarshal([]byte(bodyStr), &data); err != nil {
		fmt.Println("json err")
		panic(err)
	}

	// fmt.Println(data)
	// fmt.Println("data")

	return data
}

func (cl *Client) Sign() {
	res, err := http.Get(fmt.Sprintf("%s?cname=%s", cl.Signurl, cl.Cname))
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}

	if err = json.Unmarshal(body, &data); err != nil {
		fmt.Println("json err")
		panic(err)
	}

	fmt.Println("注册", data)
	cl.Token = data["token"].(string)
	cl.Cname = data["cname"].(string)
}

func (cl *Client) Push(content string) {
	res, err := http.Get(fmt.Sprintf("%s?cname=%s&content=%s", cl.Pushurl, cl.Cname, content))
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}

	if err = json.Unmarshal(body, &data); err != nil {
		fmt.Println("json err")
		panic(err)
	}

	fmt.Println(" 发送消息", data)

}
