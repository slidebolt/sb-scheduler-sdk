// Package scheduler is the client contract for the SlideBolt schedule service.
package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	domain "github.com/slidebolt/sb-domain"
	messenger "github.com/slidebolt/sb-messenger-sdk"
)

const (
	SubjectCreate  = "scheduler.schedules.create"
	SubjectUpdate  = "scheduler.schedules.update"
	SubjectDelete  = "scheduler.schedules.delete"
	requestTimeout = 5 * time.Second
)

type CreateScheduleRequest struct {
	ID       string          `json:"id,omitempty"`
	Name     string          `json:"name"`
	Schedule domain.Schedule `json:"schedule"`
}
type UpdateScheduleRequest struct {
	Key      domain.EntityKey `json:"key"`
	Name     string           `json:"name,omitempty"`
	Schedule domain.Schedule  `json:"schedule"`
}
type DeleteScheduleRequest struct {
	Key domain.EntityKey `json:"key"`
}
type ScheduleReply struct {
	Key   domain.EntityKey `json:"key"`
	Error string           `json:"error,omitempty"`
}

type Client struct{ msg messenger.Messenger }

func NewClient(msg messenger.Messenger) *Client { return &Client{msg: msg} }
func (c *Client) CreateSchedule(req CreateScheduleRequest) (domain.EntityKey, error) {
	var reply ScheduleReply
	if err := c.request(SubjectCreate, req, &reply); err != nil {
		return domain.EntityKey{}, err
	}
	if reply.Error != "" {
		return domain.EntityKey{}, errors.New(reply.Error)
	}
	return reply.Key, nil
}
func (c *Client) UpdateSchedule(req UpdateScheduleRequest) error {
	var reply ScheduleReply
	if err := c.request(SubjectUpdate, req, &reply); err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}
	return nil
}
func (c *Client) DeleteSchedule(key domain.EntityKey) error {
	var reply ScheduleReply
	if err := c.request(SubjectDelete, DeleteScheduleRequest{Key: key}, &reply); err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}
	return nil
}
func (c *Client) EnableSchedule(key domain.EntityKey) error {
	return messenger.NewCommands(c.msg, domain.LookupCommand).Send(key, domain.ScheduleEnable{})
}
func (c *Client) DisableSchedule(key domain.EntityKey) error {
	return messenger.NewCommands(c.msg, domain.LookupCommand).Send(key, domain.ScheduleDisable{})
}
func (c *Client) request(subject string, req, reply any) error {
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("scheduler: marshal: %w", err)
	}
	msg, err := c.msg.Request(subject, data, requestTimeout)
	if err != nil {
		return fmt.Errorf("scheduler: %s: %w", subject, err)
	}
	if err := json.Unmarshal(msg.Data, reply); err != nil {
		return fmt.Errorf("scheduler: decode reply: %w", err)
	}
	return nil
}
