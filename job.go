package render

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const JobsUrl string = "/jobs"

type Job struct {
	Id           string `json:"id"`
	ServiceId    string `json:"serviceId"`
	StartCommand string `json:"startCommand"`
	PlanId       string `json:"planId"`
	CreatedAt    string `json:"createdAt"`
	StartedAt    string `json:"startedAt"`
	FinishedAt   string `json:"finishedAt"`
	Status       string `json:"status"`
}

type NewJob struct {
	ServiceId    string `json:"serviceId"`    // required
	StartCommand string `json:"startCommand"` // required
	PlanId       string `json:"planId,omitempty"`
}

func (c *Client) GetJob(id string, serviceId string) (*Job, error) {

	url := c.Host + c.ServicesBase + serviceId + "/jobs/" + id

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	job := Job{}
	err = json.Unmarshal(res, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *Client) CreateJob(nj NewJob) (*Job, error) {
	j, err := json.Marshal(nj)
	if err != nil {
		return nil, err
	}

	url := c.Host + c.ServicesBase + nj.ServiceId + "/jobs"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	job := Job{}
	err = json.Unmarshal(res, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
