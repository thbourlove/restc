package restc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eleme/jsonpath"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/pkg/errors"
)

type Client struct {
	http.Client
}

func (client *Client) FetchJsonDataWithPath(request *http.Request, data interface{}, path string) error {
	response, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "client do request")
	}
	defer response.Body.Close()

	p, err := jsonpath.ParsePath(path)
	if err != nil {
		return errors.Wrap(err, "parse paths")
	}
	eval, err := jsonpath.EvalPathInReader(response.Body, p)
	if err != nil {
		return errors.Wrap(err, "eval paths")
	}

	results := [][]byte{}
	for {
		result, ok := eval.Next()
		if result != nil {
			results = append(results, result.Value)
		}
		if !ok {
			break
		}
	}

	if eval.Error != nil {
		return errors.Wrap(eval.Error, "eval next")
	}

	if len(results) <= 0 {
		return fmt.Errorf("failed to get data from eval")
	} else if len(results) == 1 {
		json.Unmarshal(results[0], data)
	} else {
		resultBytes := append([]byte("["), bytes.Join(results, []byte(","))...)
		resultBytes = append(resultBytes, ']')
		json.Unmarshal(resultBytes, data)
	}

	return nil
}

func (client *Client) FetchJsonData(request *http.Request, data interface{}) error {
	response, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "client do request")
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(data); err != nil {
		return errors.Wrap(err, "decode body")
	}
	return nil
}

func NewDebugClient(cache httpcache.Cache) *Client {
	return &Client{
		http.Client{
			Transport: NewDebugTransport(cache),
		},
	}
}

func NewDebugClientWithDiskCache(dir string) *Client {
	return &Client{
		http.Client{
			Transport: NewDebugTransport(diskcache.New(dir)),
		},
	}
}

func NewClient() *Client {
	return &Client{http.Client{}}
}
