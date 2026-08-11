package server

import (
	"encoding/json"
	"fmt"
)

func UpdateUser(event Event, c *Client) error {
	var updateUserEvent UpdateUserEvent
	if err := json.Unmarshal(event.Payload, &updateUserEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	c.Name = updateUserEvent.Name

	var userUpdatedMsg UserUpdatedEvent

	userUpdatedMsg.Name = c.Name

	data, err := json.Marshal(userUpdatedMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal connection established message: %v\n", err)
	}

	userUpdated := Event{
		Payload: data,
		Type:    EventUserUpdated,
	}

	c.egress <- userUpdated

	return nil
}
