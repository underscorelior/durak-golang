import { globalState } from '$lib/state/state.svelte';
import { joinLobby, rejoinLobby } from './actions';
import type { EventPayloads, Events } from './events';
import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import { clearLocalSession, setLocalSession } from '$lib/ws/connection';

export function connectionEstablishedHandler(payload: EventPayloads[Events.ConnectionEstablished]) {
	globalState.user.name = payload.name;
	globalState.user.user_id = payload.user_id;
	globalState.menu.lobbies = payload.lobbies;

	if (payload.is_rejoining && globalState.session) {
		rejoinLobby(globalState.session);
	}

	if (!payload.is_rejoining) {
		console.log('NOT REJOINING');
		clearLocalSession();
	}
}

export function menuLobbiesUpdatedHandler(payload: EventPayloads[Events.MenuLobbiesUpdated]) {
	globalState.menu.lobbies = payload.lobbies;
}

export function userUpdatedHandler(payload: EventPayloads[Events.UserUpdated]) {
	globalState.user.name = payload.name;
}

export function lobbyCreatedHandler(payload: EventPayloads[Events.LobbyCreated]) {
	joinLobby(payload.lobby_code);
}

export function lobbyJoinedHandler(payload: EventPayloads[Events.LobbyJoined]) {
	globalState.lobby = payload.lobby;

	goto(resolve(`/${payload.lobby.lobby_code}`));
}

// TODO: Proper err handling
export function joinLobbyFailedHandler(payload: EventPayloads[Events.JoinLobbyFailed]) {
	clearLocalSession();
	alert(`${payload.code}: ${payload.message} (${payload.lobby_code})`);

	goto(resolve('/'));
}

export function playerJoinedHandler(payload: EventPayloads[Events.PlayerJoined]) {
	const players = [...(globalState.lobby?.players ?? []), payload.player].sort(
		(a, b) => a.position - b.position
	);
	globalState.lobby = { ...globalState.lobby, players } as Lobby;
}

export function lobbyLeftHandler(payload: EventPayloads[Events.LobbyLeft]) {
	clearLocalSession();
	goto(resolve('/'));

	globalState.lobby = null;
	globalState.menu.lobbies = payload.lobbies;
}

export function playerLeftHandler(payload: EventPayloads[Events.PlayerLeft]) {
	const players = globalState.lobby?.players.filter((p) => p.user_id != payload.user_id);
	globalState.lobby = { ...globalState.lobby, host_id: payload.host_id, players } as Lobby;
}

export function playerDisconnectedHandler(payload: EventPayloads[Events.PlayerDisconnected]) {
	const disconnected = globalState.lobby?.players.find((p) => p.user_id === payload.user_id);

	if (disconnected) {
		disconnected.disconnected_at = payload.disconnected_at;
		disconnected.is_connected = false;
	}
}

export function rejoinLobbyFailedHandler(payload: EventPayloads[Events.RejoinLobbyFailed]) {
	clearLocalSession();
	alert(`${payload.code}: ${payload.message} (${payload.lobby_code})`);

	goto(resolve('/'));
}

export function lobbyRejoinedHandler(payload: EventPayloads[Events.LobbyRejoined]) {
	globalState.lobby = payload.lobby;
}

export function playerRejoinedHandler(payload: EventPayloads[Events.PlayerRejoined]) {
	const rejoined = globalState.lobby?.players.find((p) => p.user_id === payload.user_id);

	if (rejoined) {
		rejoined.is_connected = true;
	}
}

export function gameStartedHandler(payload: EventPayloads[Events.GameStarted]) {
	const session = payload.session;
	setLocalSession(session);

	globalState.lobby = payload.lobby;
}
