package es_driver

import (
	"database/sql/driver"
	"math/rand"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ESConnConfig struct {
	Username string   `json:"username"`
	Password string   `json:"passowrd"`
	AddrList []string `json:"addrlist"`
}

type ESConfSQL struct {
	Method   string
	Database string
	Table    string
	Where    []string
	Order    []string
	Field    []string
	Alias    []string
	Value    []interface{}
	From     int
	Size     int
}

type ESConn struct {
	username   string
	password   string
	httpAddr   string
	httpClient *http.Client
}

func (ec *ESConn) Init(c *ESConnConfig) {
	ec.httpClient = &http.Client{}
	ec.username = c.Username
	ec.password = c.Password
	ec.httpAddr = c.AddrList[rand.Intn(len(c.AddrList))]
}

func (ec *ESConn) Prepare(query string) (driver.Stmt, error) {
	fmt.Printf("ESConn Prepare\n")
	return nil, errors.New("Prepare no support")
}

func (ec *ESConn) Close() error {
	fmt.Printf("ESConn Close\n")
	return nil
}

func (ec *ESConn) Begin() (driver.Tx, error) {
	fmt.Printf("ESConn Begin\n")
	return nil, errors.New("Begin no support")
}

func (ec *ESConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	fmt.Printf("ESConn Exec\n")
	return nil, errors.New("Exec no support")
}

func (ec *ESConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	rows := new(ESRows)

	sql, err := parseESSQL(query)
	if err != nil {
		return nil, err
	}

	var request SearchRequest
	filter, err := parseFilterQuery(sql.Where)
	if err != nil {
		return nil, err
	}

	request.Sort = make([]interface{}, 0)
	request.Filter = filter
	request.From = 0
	request.Size = 10

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s/%s/%s/_search", ec.httpAddr, sql.Database, sql.Table)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	fmt.Printf("url:%s\n", url)
	fmt.Printf("body:%s\n", string(data))
	fmt.Printf("sql:%+v\n", sql)
	resp, err := ec.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &rows.searchResult)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
