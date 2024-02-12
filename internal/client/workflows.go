package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"
)

// GetWorkflows retreives all workflows
func (c *Client) GetWorkflows() ([]Workflow, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/workflows/export", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
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

func (c *Client) PostWorkflow(workflowYaml string) ([]byte, error) {
	//rb, err := json.Marshal(workflowYaml)
	//if err != nil {
	//	return nil, err
	//}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/workflows", c.HostURL), strings.NewReader(workflowYaml))
	//req, err := http.NewRequest("POST", fmt.Sprintf("%s/workflows", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) DeleteWorkflow(workflowYamlString string) (string, error) {
	var workflowYaml map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(workflowYamlString), &workflowYaml)
	if err != nil {
		return "", err
	}

	id := workflowYaml["id"].(string)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/workflows/%s", c.HostURL, id), nil)
	if err != nil {
		return "", err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/workflows/%s", c.HostURL, id), nil
}