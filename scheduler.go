// Package scheduler is the client contract for the SlideBolt timer service.
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
	SubjectCreate  = "scheduler.timers.create"
	SubjectUpdate  = "scheduler.timers.update"
	SubjectDelete  = "scheduler.timers.delete"
	requestTimeout = 5 * time.Second
)

type CreateTimerRequest struct {
	ID    string       `json:"id,omitempty"`
	Name  string       `json:"name"`
	Timer domain.Timer `json:"timer"`
}

type UpdateTimerRequest struct {
	Key   domain.EntityKey `json:"key"`
	Name  string           `json:"name,omitempty"`
	Timer domain.Timer     `json:"timer"`
}

type DeleteTimerRequest struct {
	Key domain.EntityKey `json:"key"`
}

type TimerReply struct {
	Key   domain.EntityKey `json:"key"`
	Error string           `json:"error,omitempty"`
}

type Client struct{ msg messenger.Messenger }

func NewClient(msg messenger.Messenger) *Client { return &Client{msg: msg} }

func (c *Client) CreateTimer(req CreateTimerRequest) (domain.EntityKey, error) {
	var reply TimerReply
	if err := c.request(SubjectCreate, req, &reply); err != nil {
		return domain.EntityKey{}, err
	}
	if reply.Error != "" {
		return domain.EntityKey{}, errors.New(reply.Error)
	}
	return reply.Key, nil
}

func (c *Client) UpdateTimer(req UpdateTimerRequest) error {
	var reply TimerReply
	if err := c.request(SubjectUpdate, req, &reply); err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}
	return nil
}

func (c *Client) DeleteTimer(key domain.EntityKey) error {
	var reply TimerReply
	if err := c.request(SubjectDelete, DeleteTimerRequest{Key: key}, &reply); err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}
	return nil
}

func (c *Client) EnableTimer(key domain.EntityKey) error {
	return messenger.NewCommands(c.msg, domain.LookupCommand).Send(key, domain.TimerEnable{})
}

func (c *Client) DisableTimer(key domain.EntityKey) error {
	return messenger.NewCommands(c.msg, domain.LookupCommand).Send(key, domain.TimerDisable{})
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
