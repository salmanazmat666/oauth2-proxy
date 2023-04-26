package postgres

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/options"
)

type APIClient struct {
	BaseURL string
	Client  *resty.Client
}

func NewAPIClient(baseurl string) *APIClient {

	ac := new(APIClient)
	ac.BaseURL = baseurl
	ac.Client = resty.New()

	return ac
}

func (c *APIClient) CreateProviderConfig(ctx context.Context, opts *options.Provider) error {

	path := c.BaseURL + "/provider"
	resp, err := c.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(opts).
		SetContext(ctx).
		SetError(&ErrorResponse{}).
		Post(path)

	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	if resp.StatusCode() == http.StatusConflict {
		return NewAlreadyExistError("config with this id already exists")
	}
	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("status code returned : %d and err: %v", resp.StatusCode(), resp.Error())
	}

	return nil
}

func (c *APIClient) UpdateProviderConfig(ctx context.Context, opts *options.Provider) error {

	path := c.BaseURL + "/provider"
	resp, err := c.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(opts).
		SetError(&ErrorResponse{}).
		SetContext(ctx).
		Put(path)

	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	if resp.StatusCode() == http.StatusNotFound {
		return NewNotFoundError("config entry not found")
	}
	if resp.StatusCode() != http.StatusAccepted {
		return fmt.Errorf("status code returned : %d and err: %v", resp.StatusCode(), resp.Error())
	}

	return nil
}

func (c *APIClient) GetProviderConfig(ctx context.Context, id string) (*options.Provider, error) {
	path := c.BaseURL + "/provider" + "/" + id
	providerConf := &options.Provider{}
	resp, err := c.Client.R().
		SetResult(providerConf).
		ForceContentType("application/json").
		SetContext(ctx).
		SetError(&ErrorResponse{}).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	if resp.StatusCode() == http.StatusNotFound {
		return nil, NewNotFoundError("config entry not found")
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status code returned : %d and err: %v ", resp.StatusCode(), resp.Error())
	}

	return providerConf, nil

}

func (c *APIClient) DeleteProviderConfig(ctx context.Context, id string) error {
	path := c.BaseURL + "/provider" + "/" + id
	resp, err := c.Client.R().
		SetContext(ctx).
		SetError(&ErrorResponse{}).
		Delete(path)

	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	if resp.StatusCode() == http.StatusNotFound {
		return NewNotFoundError("entry to delete not found")
	}
	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("status code returned : %d and err: %v", resp.StatusCode(), resp.Error())
	}
	return nil
}