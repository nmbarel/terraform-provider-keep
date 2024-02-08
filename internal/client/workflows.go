package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetWorkflows retreives all workflows
func (c *Client) GetWorkflows() ([]Workflow, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/workflows/export", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	fmt.Print(body)
	if err != nil {
		return nil, err
	}

	//workflows := []Workflow{}
	yamls := []string{}

	err = json.Unmarshal(body, &yamls)
	if err != nil {
		return nil, err
	}

	workflows := make([]Workflow, len(yamls))

	for index, yaml := range yamls {
		workflows[index].Yaml = yaml
	} 

	return workflows, nil
}