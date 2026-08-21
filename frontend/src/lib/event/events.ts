export enum Events {
	ConnectionEstablished = 'connection_established',
	UpdateUser = 'update_user',
	UserUpdated = 'user_updated',
	CreateLobby = 'create_lobby',
	LobbyCreated = 'lobby_created',
	JoinLobby = 'join_lobby',
	JoinLobbyFailed = 'join_lobby_failed',
	LobbyJoined = 'lobby_joined',
	PlayerJoined = 'player_joined',
	LeaveLobby = 'leave_lobby',
	PlayerLeft = 'player_left',
	StartGame = 'start_game',
	GameStarted = 'game_started'
}

type ConnectionEstablishedPayload = {
	name: string;
	user_id: string;
	lobbies: MenuLobby[];
};

type UpdateUserPayload = {
	name: string;
};

type UserUpdatedPayload = {
	name: string;
};

type CreateLobbyPayload = object;

type LobbyCreatedPayload = {
	lobby_id: string;
};

type JoinLobbyPayload = {
	lobby_id: string;
};

type JoinLobbyFailedPayload = {
	lobby_id: string;
	code: string;
	message: string;
};

type LobbyJoinedEvent = {
	lobby: Lobby;
};

export type EventPayloads = {
	[Events.ConnectionEstablished]: ConnectionEstablishedPayload;
	[Events.UpdateUser]: UpdateUserPayload;
	[Events.UserUpdated]: UserUpdatedPayload;
	[Events.CreateLobby]: CreateLobbyPayload;
	[Events.LobbyCreated]: LobbyCreatedPayload;
	[Events.JoinLobby]: JoinLobbyPayload;
	[Events.JoinLobbyFailed]: JoinLobbyFailedPayload;
	[Events.LobbyJoined]: LobbyJoinedEvent;
};

export type WSEvent = {
	[K in keyof EventPayloads]: {
		type: K;
		payload: EventPayloads[K];
	};
}[keyof EventPayloads];
