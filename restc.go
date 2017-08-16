package restc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eleme/jsonpath"
	"github.com/gregjones/httpcache"
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
	if result, ok := eval.Next(); ok {
		json.Unmarshal(result.Value, data)
	} else {
		return fmt.Errorf("failed to get data from eval")
	}
	if eval.Error != nil {
		return errors.Wrap(eval.Error, "eval next")
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

func NewClient() *Client {
	return &Client{http.Client{}}
}
