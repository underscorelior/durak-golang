package server

import (
	"durak/internal/game"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

func UpdateUser(event Event, c *Client) error {
	var updateUserEvent UpdateUserEvent
	if err := json.Unmarshal(event.Payload, &updateUserEvent); err != nil {
		return fmt.Errorf("(UpdateUser) Bad payload in request: %w", err)
	}

	c.Name = updateUserEvent.Name

	var userUpdatedPayload UserUpdatedEvent
	userUpdatedPayload.Name = c.Name

	userUpdated, err := newEvent(EventUserUpdated, userUpdatedPayload)
	if err != nil {
		return err
	}

	c.egress <- userUpdated

	return nil
}

func CreateLobby(event Event, c *Client) error {
	// TODO: Need to implement ratelimit system
	lobby := c.manager.NewLobby(c.UserID)
	c.manager.addLobby(lobby)

	var lobbyCreatedPayload LobbyCreatedEvent
	lobbyCreatedPayload.LobbyCode = lobby.LobbyCode

	lobbyCreated, err := newEvent(EventLobbyCreated, lobbyCreatedPayload)
	if err != nil {
		return err
	}

	c.egress <- lobbyCreated

	return nil
}

func JoinLobby(event Event, c *Client) error {
	var joinLobbyEvent JoinLobbyEvent
	if err := json.Unmarshal(event.Payload, &joinLobbyEvent); err != nil {
		return fmt.Errorf("(JoinLobby) Bad payload in request: %w", err)
	}

	if c.lobby != nil && joinLobbyEvent.LobbyCode == c.lobby.LobbyCode {
		return &ClientError{
			Code:    "ALREADY_IN_LOBBY",
			Message: "Cannot join a lobby you are already a part of",
		}
	}

	if c.lobby != nil {
		// TODO: Or should we handle this another way?
		c.lobby.removeClient(c)
	}

	lobby, ok := c.manager.lobbies[joinLobbyEvent.LobbyCode]

	if !ok {
		return &ClientError{
			Code:    "LOBBY_NOT_FOUND",
			Message: fmt.Sprintf("Lobby (%v) does not exist", joinLobbyEvent.LobbyCode),
		}
	}

	if lobby.InGame {
		return &ClientError{
			Code:    "LOBBY_IN_GAME",
			Message: fmt.Sprintf("Lobby (%v) is in an active game", joinLobbyEvent.LobbyCode),
		}
	}

	if len(lobby.members) >= int(lobby.MaxMembers) {
		return &ClientError{
			Code:    "LOBBY_FULL",
			Message: fmt.Sprintf("Lobby is full (Max %v)", lobby.MaxMembers),
		}
	}

	if _, ok := lobby.clients[c.UserID]; ok {
		return &ClientError{
			Code:    "LOBBY_DUPLICATE_MEMBER",
			Message: "Lobby contains a member with the same UserID as you",
		}
	}

	pos := lobby.nextAvailablePosition()
	if err := lobby.usePosition(pos); err != nil {
		return &ClientError{
			Code:    "LOBBY_NO_VALID_POSITION",
			Message: "There was an issue finding a valid position",
		}
	}

	c.lobby = lobby

	p := &Member{
		UserID:      c.UserID,
		Name:        c.Name,
		Position:    pos,
		JoinedAt:    time.Now(),
		IsConnected: true,
	}

	lobby.addClient(c, p)

	var lobbyJoinedPayload LobbyJoinedEvent
	snapshot := lobby.SnapshotFor(c)
	if snapshot == nil {
		return fmt.Errorf("(JoinLobby) No snapshot for user %s in lobby %s", c.UserID, lobby.LobbyCode)
	}

	lobbyJoinedPayload.Lobby = *snapshot

	lobbyJoined, err := newEvent(EventLobbyJoined, lobbyJoinedPayload)
	if err != nil {
		return err
	}

	c.egress <- lobbyJoined

	var memberJoinedPayload MemberJoinedEvent
	memberJoinedPayload.Member = *p

	memberJoined, err := newEvent(EventMemberJoined, memberJoinedPayload)
	if err != nil {
		return err
	}

	ignored := ClientList{
		c.UserID: c,
	}
	lobby.broadcast(memberJoined, ignored)

	return nil
}

func LeaveLobby(event Event, c *Client) error {
	lobby := c.lobby

	if lobby == nil {
		return &ClientError{
			Code:    "NOT_IN_LOBBY",
			Message: "Cannot leave lobby when not in lobby",
		}
	}

	if err := lobby.removeClient(c); err != nil {
		log.Println("(LeaveLobby)", err)
		return err
	}

	var lobbyLeftPayload LobbyLeftEvent
	lobbyLeftPayload.Lobbies = c.manager.MenuLobbies()

	lobbyLeft, err := newEvent(EventLobbyLeft, lobbyLeftPayload)
	if err != nil {
		return err
	}

	c.egress <- lobbyLeft

	return nil
}

func StartGame(event Event, c *Client) error {
	l := c.lobby

	if l == nil {
		return &ClientError{
			Code:    "NOT_IN_LOBBY",
			Message: "Cannot start game when not in lobby",
		}

	}

	if c.UserID != l.HostID {
		return &ClientError{
			Code:    "MEMBER_NOT_HOST",
			Message: fmt.Sprintf("User (%v) is not lobby host (%v)", c.UserID[:6], l.HostID[:6]),
		}
	}

	if l.InGame {
		return &ClientError{
			Code:    "GAME_ALREADY_STARTED",
			Message: "Game has already started",
		}
	}

	if len(l.members) <= 1 {
		return &ClientError{
			Code:    "NOT_ENOUGH_MEMBERS",
			Message: "Cannot start lobby with less than 2 members",
		}
	}

	gameMembers := make(map[string]struct{ Position int })

	for k, v := range l.members { // TODO: Fix this BS
		gameMembers[k] = struct{ Position int }{Position: v.Position}
	}

	l.game = game.InitializeGame(gameMembers)
	l.InGame = true

	snapshot := l.Snapshot()

	l.GenerateSessions()

	for _, cl := range l.clients {
		p := l.members[cl.UserID]

		lobbySnapshot := snapshot
		lobbySnapshot.Position = p.Position

		gameState := l.game.SnapshotFor(cl.UserID)
		if gameState == nil {
			return fmt.Errorf("(StartGame) No snapshot for user %s in lobby %s", cl.UserID, l.LobbyCode)
		}

		var gameStartedPayload GameStartedEvent
		gameStartedPayload.Lobby = lobbySnapshot
		gameStartedPayload.Session = l.members[cl.UserID].session
		gameStartedPayload.GameState = *gameState // TODO: for gameStarted, I dont think we want to provide the turn state

		gameStarted, err := newEvent(EventGameStarted, gameStartedPayload)
		if err != nil {
			return err
		}

		cl.egress <- gameStarted
	}

	return nil
}

func RejoinLobby(event Event, c *Client) error {
	var rejoinLobbyEvent RejoinLobbyEvent

	if err := json.Unmarshal(event.Payload, &rejoinLobbyEvent); err != nil {
		return fmt.Errorf("(RejoinLobby) Bad payload in request: %v", err)
	}

	if c.lobby != nil && rejoinLobbyEvent.LobbyCode == c.lobby.LobbyCode {
		return nil // How else should I handle this?
	}

	lobby, ok := c.manager.lobbies[rejoinLobbyEvent.LobbyCode]

	if !ok {
		return &ClientError{
			Code:    "LOBBY_NOT_FOUND",
			Message: "Lobby does not exist",
		}
	}

	userID, ok := lobby.sessions[rejoinLobbyEvent.SessionToken]

	if !ok {
		return &ClientError{
			Code:    "INVALID_SESSION_TOKEN",
			Message: fmt.Sprintf("Your session token is not valid for lobby: %v", rejoinLobbyEvent.LobbyCode),
		}
	}

	// FIXME: Should we really give the newer client the older ID? Should Member have its own "lobby-specific ID", then clients have their own ID?
	if userID != c.UserID {
		return &ClientError{
			Code:    "INCORRECT_USER_ID",
			Message: fmt.Sprintf("Your user ID does not match : %v", rejoinLobbyEvent.LobbyCode),
		}
	}

	p := lobby.members[c.UserID]

	if p == nil {
		return errors.New("(RejoinLobby) Failed to locate member")
	}

	c.lobby = lobby
	p.IsConnected = true

	lobby.addClient(c, p)

	var lobbyRejoinedPayload LobbyRejoinedEvent
	snapshot := lobby.SnapshotFor(c)
	if snapshot == nil {
		return fmt.Errorf("(RejoinLobby) No snapshot for user %s in lobby %s", c.UserID, lobby.LobbyCode)
	}

	lobbyRejoinedPayload.Lobby = *snapshot

	lobbyRejoined, err := newEvent(EventLobbyRejoined, lobbyRejoinedPayload)
	if err != nil {
		return err
	}

	c.egress <- lobbyRejoined

	var memberRejoinedPayload MemberRejoinedEvent
	memberRejoinedPayload.UserID = c.UserID

	memberRejoined, err := newEvent(EventMemberRejoined, memberRejoinedPayload)
	if err != nil {
		return err
	}

	ignored := ClientList{
		c.UserID: c,
	}
	lobby.broadcast(memberRejoined, ignored)

	return nil
}

func MakeAction(event Event, c *Client) error {
	if c.lobby == nil {
		return errors.New("(MakeAction) Cant make action when not in lobby")
	}

	if c.lobby.game == nil {
		return errors.New("(MakeAction) Cant make an action when there is no active game")
	}

	var makeAction MakeActionEvent

	if err := json.Unmarshal(event.Payload, &makeAction); err != nil {
		return fmt.Errorf("(MakeAction) Bad payload in request: %w", err)
	}

	if gmErr := c.lobby.game.MakeAction(c.UserID, makeAction.Action, makeAction.Card, makeAction.PairIndex); gmErr != nil {
		return gmErr
	}

	var actionMadePayload ActionMadeEvent
	actionMadePayload = ActionMadeEvent{
		UserID:    c.UserID,
		Action:    makeAction.Action,
		Card:      makeAction.Card,
		PairIndex: makeAction.PairIndex,
		Phase:     &c.lobby.game.Turn.Phase,
	}

	actionMade, err := newEvent(EventActionMade, actionMadePayload)
	if err != nil {
		return err
	}

	c.lobby.broadcast(actionMade, ClientList{})

	return nil
}
