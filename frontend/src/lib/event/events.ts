export enum Events {
	EventConnectionEstablished = 'connection_established',
	EventUpdateUser = 'update_user',
	EventUserUpdated = 'user_updated',
	EventCreateLobby = 'create_lobby',
	EventLobbyCreated = 'lobby_created',
	EventJoinLobby = 'join_lobby',
	EventJoinLobbyFailed = 'join_lobby_failed',
	EventLobbyJoined = 'lobby_joined',
	EventPlayerJoined = 'player_joined',
	EventLeaveLobby = 'leave_lobby',
	EventPlayerLeft = 'player_left',
	EventStartGame = 'start_game',
	EventGameStarted = 'game_started'
}

type ConnectionEstablishedPayload = {
	name: string;
	userId: string;
	lobbies: MenuLobby[];
};

export type EventPayloads = {
	[Events.EventConnectionEstablished]: ConnectionEstablishedPayload;
};

export type WSEvent = {
	[K in keyof EventPayloads]: {
		type: K;
		payload: EventPayloads[K];
	};
}[keyof EventPayloads];
