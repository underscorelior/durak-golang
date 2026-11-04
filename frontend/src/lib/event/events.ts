import type { Action, TurnPhase } from '$lib/game';

export enum GameEvent {
	Error = 'error',
	ConnectionEstablished = 'connection_established',
	MenuLobbiesUpdated = 'menu_lobbies_updated', // TODO: Find better name
	UpdateUser = 'update_user',
	UserUpdated = 'user_updated',
	CreateLobby = 'create_lobby',
	LobbyCreated = 'lobby_created',
	JoinLobby = 'join_lobby',
	LobbyJoined = 'lobby_joined',
	MemberJoined = 'member_joined',
	LeaveLobby = 'leave_lobby',
	LobbyLeft = 'lobby_left',
	MemberLeft = 'member_left',
	PlayerDisconnected = 'player_disconnected',
	RejoinLobby = 'rejoin_lobby',
	LobbyRejoined = 'lobby_rejoined',
	PlayerRejoined = 'player_rejoined',
	StartGame = 'start_game',
	GameStarted = 'game_started',
	MakeAction = 'make_action',
	ActionMade = 'action_made' // TODO: Find better name
}

type ConnectionEstablishedPayload = {
	user_id: string;
	name: string;
	lobbies: LobbyPreview[];
	is_rejoining: boolean;
};

type MenuLobbiesUpdatedPayload = {
	lobbies: LobbyPreview[];
	next_update: Date; // TODO: IMPLEMENT A VISUAL TIMER, ALSO IMPLEMENT A FORCE REFRESH
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

type LobbyJoinedPayload = {
	lobby: Lobby;
};

type MemberJoinedPayload = {
	member: Member;
};

type LeaveLobbyPayload = object;

type LobbyLeftPayload = {
	lobbies: LobbyPreview[];
};

type MemberLeftPayload = {
	user_id: string;
	host_id: string;
	position: number;
};

type PlayerDisconnectedPayload = {
	user_id: string;
	disconnected_at: Date;
};

type RejoinLobbyPayload = {
	lobby_code: string;
	session_token: string;
};

type LobbyRejoinedPayload = {
	lobby: Lobby;
};

type PlayerRejoinedPayload = {
	user_id: string;
};

type StartGamePayload = object;

type GameStartedPayload = {
	lobby: Lobby;
	session: Session;
};

type MakeActionPayload = {
	action: Action;
	card?: Card;
	pair_index?: number;
};

type ActionMadePayload = {
	user_id: string;
	action: Action;
	card?: Card;
	pair_index?: number;
	turn_phase?: TurnPhase;
};

type ErrorPayload = {
	code: string;
	message: string;
	error_for: GameEvent;
};

export type EventPayloads = {
	[GameEvent.ConnectionEstablished]: ConnectionEstablishedPayload;
	[GameEvent.MenuLobbiesUpdated]: MenuLobbiesUpdatedPayload;
	[GameEvent.UpdateUser]: UpdateUserPayload;
	[GameEvent.UserUpdated]: UserUpdatedPayload;
	[GameEvent.CreateLobby]: CreateLobbyPayload;
	[GameEvent.LobbyCreated]: LobbyCreatedPayload;
	[GameEvent.JoinLobby]: JoinLobbyPayload;
	[GameEvent.LobbyJoined]: LobbyJoinedPayload;
	[GameEvent.MemberJoined]: MemberJoinedPayload;
	[GameEvent.LeaveLobby]: LeaveLobbyPayload;
	[GameEvent.LobbyLeft]: LobbyLeftPayload;
	[GameEvent.MemberLeft]: MemberLeftPayload;
	[GameEvent.PlayerDisconnected]: PlayerDisconnectedPayload;
	[GameEvent.RejoinLobby]: RejoinLobbyPayload;
	[GameEvent.LobbyRejoined]: LobbyRejoinedPayload;
	[GameEvent.PlayerRejoined]: PlayerRejoinedPayload;
	[GameEvent.StartGame]: StartGamePayload;
	[GameEvent.GameStarted]: GameStartedPayload;
	[GameEvent.MakeAction]: MakeActionPayload;
	[GameEvent.ActionMade]: ActionMadePayload;
	[GameEvent.Error]: ErrorPayload;
};

export type WSEvent = {
	[K in keyof EventPayloads]: {
		type: K;
		payload: EventPayloads[K];
	};
}[keyof EventPayloads];
