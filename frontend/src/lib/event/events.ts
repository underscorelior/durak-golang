export enum Events {
	ConnectionEstablished = 'connection_established',
	MenuLobbiesUpdated = 'menu_lobbies_updated', // TODO: Find better name
	UpdateUser = 'update_user',
	UserUpdated = 'user_updated',
	CreateLobby = 'create_lobby',
	LobbyCreated = 'lobby_created',
	JoinLobby = 'join_lobby',
	JoinLobbyFailed = 'join_lobby_failed',
	LobbyJoined = 'lobby_joined',
	PlayerJoined = 'player_joined',
	LeaveLobby = 'leave_lobby',
	LobbyLeft = 'lobby_left',
	PlayerLeft = 'player_left',
	StartGame = 'start_game',
	GameStarted = 'game_started'
}

type ConnectionEstablishedPayload = {
	name: string;
	user_id: string;
	lobbies: LobbyPreview[];
};

type MenuLobbiesUpdatedPayload = {
	lobbies: LobbyPreview[];
};

type UpdateUserPayload = {
	name: string;
};

type UserUpdatedPayload = {
	name: string;
};

type CreateLobbyPayload = object;

type LobbyCreatedPayload = {
	lobby_code: string;
};

type JoinLobbyPayload = {
	lobby_code: string;
};

type JoinLobbyFailedPayload = {
	lobby_code: string;
	code: string;
	message: string;
};

type LobbyJoinedPayload = {
	lobby: Lobby;
};

type PlayerJoinedPayload = {
	player: Player;
};

type LeaveLobbyPayload = object;

type LobbyLeftPayload = {
	lobbies: LobbyPreview[];
};

type PlayerLeftPayload = {
	user_id: string;
};

export type EventPayloads = {
	[Events.ConnectionEstablished]: ConnectionEstablishedPayload;
	[Events.MenuLobbiesUpdated]: MenuLobbiesUpdatedPayload;
	[Events.UpdateUser]: UpdateUserPayload;
	[Events.UserUpdated]: UserUpdatedPayload;
	[Events.CreateLobby]: CreateLobbyPayload;
	[Events.LobbyCreated]: LobbyCreatedPayload;
	[Events.JoinLobby]: JoinLobbyPayload;
	[Events.JoinLobbyFailed]: JoinLobbyFailedPayload;
	[Events.LobbyJoined]: LobbyJoinedPayload;
	[Events.PlayerJoined]: PlayerJoinedPayload;
	[Events.LeaveLobby]: LeaveLobbyPayload;
	[Events.LobbyLeft]: LobbyLeftPayload;
	[Events.PlayerLeft]: PlayerLeftPayload;
};

export type WSEvent = {
	[K in keyof EventPayloads]: {
		type: K;
		payload: EventPayloads[K];
	};
}[keyof EventPayloads];
