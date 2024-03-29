package menandmice

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	ResourceNotFound           = 16544
	ObjectNotFoundForReference = 2049
)

type Mmclient struct {
	serverLocation *time.Location
	resty.Client
}

// Cfg config to construct client
type Cfg struct {
	MMEndpoint string
	MMUsername string
	MMPassword string
	MMTimezone string
	TLSVerify  bool
	Timeout    int
	Version    string
	Debug      bool
}

func init() {
	// Remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// ClientInit establishes default settings on the REST client.
func ClientInit(c *Cfg) (*Mmclient, error) {
	client := Mmclient{
		// default current timezone
		Client: *resty.New(),
	}
	client.SetDebug(c.Debug)
	if c.MMEndpoint == "" {
		return nil, errors.New("REST API endpoint must be configured")
		// TODO check if it resolaves
	}

	if match, _ := regexp.MatchString("^(http|https)://", c.MMEndpoint); !match {

		return nil, fmt.Errorf("REST API endpoint: %s must start with \"http://\" or \"https://\"", c.MMEndpoint)
	}

	if c.MMUsername == "" {
		return nil, errors.New("Invalid username")
	}
	if c.MMPassword == "" {
		return nil, errors.New("Invalid password")
	}
	if c.MMTimezone != "" {

		if location, err := time.LoadLocation(c.MMTimezone); err == nil {
			client.serverLocation = location
		} else {
			return nil, err
		}
	} else {
		client.serverLocation = time.Now().Location()
	}

	if !c.TLSVerify {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	} else {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: false})
	}

	client.SetBasicAuth(c.MMUsername, c.MMPassword)
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("User-Agen", "terraform-provider-menandmice "+c.Version)
	client.SetTimeout(time.Duration(c.Timeout) * time.Second)
	client.SetBaseURL(c.MMEndpoint + "/mmws/api")

	// work around for micetro not understanding + in query string
	client.SetPreRequestHook(func(_ *resty.Client, rawreq *http.Request) error {
		rawreq.URL.RawQuery = strings.ReplaceAll(rawreq.URL.RawQuery, "+", "%20")
		return nil
	})

	// Test if we can make a connection

	// TODO use request that need authentication
	_, err := client.R().Get("")
	if err != nil {
		return nil, fmt.Errorf("Could not connect with endpoint: %s\n\t%s", c.MMEndpoint, err)
	}

	return &client, err
}

type DeleteRequest struct {
	SaveComment  string `json:"saveComment"`
	ForceRemoval bool   `json:"forceRemoval"`
	ObjType      string `json:"objType,omitempty"`
}

func deleteRequest(objType string) DeleteRequest {
	return DeleteRequest{
		ForceRemoval: true,
		SaveComment:  "deleted by terraform",
		ObjType:      objType,
	}
}

type RefResponse struct {
	Result struct {
		Ref string `json:"ref"`
	} `json:"result"`
}
type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type RequestError struct {
	Method     string
	URL        string
	HTTPCode   int
	StatusCode int
	ErrMessage string
}

func (r *RequestError) Error() string {

	var operation string
	switch r.Method {
	case "GET":
		operation = "reading"

	case "PUT":
		operation = "updating"

	case "POST":
		operation = "creating"

	case "DELETE":
		operation = "deleting"
	default:
		operation = "accesing"
	}

	url, err := url.Parse(r.URL)
	if err != nil {
		log.Fatal(err)
	}
	resource := url.RequestURI()

	return fmt.Sprintf("Failed with %v %v\n\tHTTP code:%v: %v", operation, resource, r.HTTPCode, r.ErrMessage)
}

func ResponseError(response *resty.Response, errorResponse ErrorResponse) error {

	if !response.IsSuccess() {
		return &RequestError{

			Method:     response.Request.Method,
			URL:        response.Request.URL,
			HTTPCode:   response.StatusCode(),
			StatusCode: errorResponse.Error.Code,
			ErrMessage: errorResponse.Error.Message,
		}
	}
	return nil
}
func map2filter(filter map[string]interface{}) string {
	if filter == nil {
		return ""
	}
	var condition string
	conditions := make([]string, 0, len(filter))
	for key, val := range filter {
		condition = fmt.Sprintf("%s=%v", key, val)
		conditions = append(conditions, condition)
	}
	return strings.Join(conditions, "&")
}

// TODO add context
func (c *Mmclient) Get(result interface{}, path string, query map[string]interface{}) error {

	//TODO better error message
	var errorResponse ErrorResponse

	request := c.R().
		SetError(&errorResponse).
		SetResult(&result)

	for key, val := range query {

		request = request.SetQueryParam(key, fmt.Sprintf("%v", val))
	}

	r, err := request.Get(path)

	if err != nil {
		return err
	}

	return ResponseError(r, errorResponse)
}

// TODO add context
func (c *Mmclient) Post(data interface{}, result interface{}, path string) error {

	//TODO better error message
	var errorResponse ErrorResponse
	r, err := c.R().
		SetBody(data).
		SetError(&errorResponse).
		SetResult(&result).
		Post(path)

	if err != nil {
		return err
	}

	return ResponseError(r, errorResponse)
}

// TODO add context
func (c *Mmclient) Delete(data interface{}, path string) error {

	var err error
	var errorResponse ErrorResponse
	r, err := c.R().
		SetBody(data).
		SetError(&errorResponse).
		Delete(path)

	if err != nil {
		return err
	}

	return ResponseError(r, errorResponse)
}

// TODO add context
func (c *Mmclient) Put(data interface{}, path string) error {
	var errorResponse ErrorResponse
	response, err := c.R().
		SetBody(data).
		SetError(&errorResponse).
		Put(path)

	if err != nil {
		return err
	}
	return ResponseError(response, errorResponse)
}
