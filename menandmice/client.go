package menandmice

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"

	"encoding/json"
)

// Cfg config to construct client
type Mmclient struct{ resty.Client }
type Cfg struct {
	MMEndpoint string
	MMUsername string
	MMPassword string
	TLSVerify  bool
	Timeout    int
}

func init() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// ClientInit establishes default settings on the REST client
func ClientInit(c *Cfg) (*Mmclient, error) {
	client := Mmclient{Client: *resty.New()}

	if c.MMEndpoint == "" {
		return nil, errors.New("REST API endpoint must be configured")
		//TODO check if it resolaves
	}
	if c.MMUsername == "" {
		return nil, errors.New("Invalid Username setting")
	}
	if c.MMPassword == "" {
		return nil, errors.New("Invalid Password setting")
	}

	if c.TLSVerify == false {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	} else {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: false})
	}

	client.SetBasicAuth(c.MMUsername, c.MMPassword)
	client.SetHeader("Content-Type", "application/json")
	client.SetTimeout(time.Duration(c.Timeout) * time.Second)
	client.SetHostURL("https://" + c.MMEndpoint + "/mmws/api")

	return &client, nil
}

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Mmclient) Get(result interface{}, path string) error {

	//TODO better error Message
	var errorResponse ErrorResponse
	r, err := c.R().
		SetError(&errorResponse).
		Get(path)
	if err != nil {
		return err
	}

	if !r.IsSuccess() {
		jsonError := r.Error().(*ErrorResponse)
		return errors.New(fmt.Sprintf("HTTP error code:%v\n%v",
			r.StatusCode(),
			jsonError.Error.Message))
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(r.Body(), &result)

	return err
}

func (c *Mmclient) Post(data interface{}, result interface{}, path string) error {

	//TODO better error Message
	var errorResponse ErrorResponse
	r, err := c.R().
		SetBody(data).
		SetError(&errorResponse).
		Post(path)

	if err != nil {
		return err
	}

	if !r.IsSuccess() {
		return fmt.Errorf("HTTP error code:%v\n%v",
			r.StatusCode(),
			errorResponse.Error.Message)
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(r.Body(), &result)

	return err
}

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

	if !r.IsSuccess() {
		return fmt.Errorf("HTTP error code:%v\n%v",
			r.StatusCode(),
			errorResponse.Error.Message)
	}

	return err
}

func (c *Mmclient) Put(data interface{}, path string) error {
	var errorResponse ErrorResponse
	r, err := c.R().
		SetBody(data).
		SetError(&errorResponse).
		Put(path)

	if err != nil {
		return err
	}

	if !r.IsSuccess() {
		return fmt.Errorf("HTTP error code:%v\n%v",
			r.StatusCode(),
			errorResponse.Error.Message)
	}

	return err
}
