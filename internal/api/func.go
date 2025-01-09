// Copyright (c) HashiCorp, Inc.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetInstance(ID string) (*Instance, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/instance/%s", c.HostURL, ID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	instance := Instance{}
	err = json.Unmarshal(body, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

func (c *Client) CreateInstance(
	Title string,
	ImageId string,
	InstanceTypeId string,
	Disk int,
) (*Instance, error) {
	rb, err := json.Marshal(map[string]interface{}{
		"title":            Title,
		"image_id":         ImageId,
		"instance_type_id": InstanceTypeId,
		"disk":             Disk,
	})

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/instance", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	instance := Instance{}
	err = json.Unmarshal(body, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

func (c *Client) UpdateInstance(
	ID string,
	Title string,
	Disk int,
) (string, error) {
	updateFields := map[string]interface{}{
		"title": Title,
		"disk":  Disk,
	}

	rb, err := json.Marshal(updateFields)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/instance/%s", c.HostURL, ID), strings.NewReader(string(rb)))
	if err != nil {
		return "", err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	instanceID, ok := response["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	return instanceID, nil
}

func (c *Client) DeleteInstance(ID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/instance/%s", c.HostURL, ID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
