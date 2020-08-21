package menandmice

import (
	"crypto/tls"
	"errors"
	"log"
	"time"

	"github.com/go-resty/resty/v2"

	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Cfg config to construct client
type Cfg struct {
	MMWeb      string
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
func ClientInit(c *Cfg) (*resty.Client, error) {
	client := resty.New()

	if c.MMWeb == "" {
		return nil, errors.New("REST API endpoint must be configured as mmWeb")
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
	client.SetHostURL("http://" + c.MMWeb + "/mmws/api") // FIXME https

	return client, nil
}

// TODO make this a method call for new object mmClient that is wrapper
func MmGet(c *resty.Client, result interface{}, path string) diag.Diagnostics {

	var diags diag.Diagnostics

	//TODO better error Message
	r, err := c.R().Get(path)
	if err != nil {
		return diag.FromErr(err)
	}

	if !r.IsSuccess() {
		return diag.Errorf("HTTP error code: %v", r.StatusCode())
	}
	if err != nil {
		return diag.FromErr(err)
	}
	err = json.Unmarshal(r.Body(), &result)

	return diags
}

// TODO make this a method call for new object mmClient that is wrapper
func MmPost(c *resty.Client, data interface{}, result interface{}, path string) diag.Diagnostics {

	var diags diag.Diagnostics
	//TODO better error Message
	r, err := c.R().SetBody(data).Post(path)

	if err != nil {
		return diag.FromErr(err)
	}

	if !r.IsSuccess() {
		return diag.Errorf("HTTP error code:%v", r.StatusCode())
	}
	if err != nil {
		return diag.FromErr(err)
	}
	err = json.Unmarshal(r.Body(), &result)

	return diags
}
