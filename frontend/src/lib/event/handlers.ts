import { globalState } from '$lib/state/state.svelte';
import type { EventPayloads, Events } from './events';

export function connectionEstablishedHandler(payload: EventPayloads[Events.ConnectionEstablished]) {
	globalState.user.name = payload.name;
	globalState.user.user_id = payload.user_id;
	globalState.menu.lobbies = payload.lobbies;
}

export function userUpdatedHandler(payload: EventPayloads[Events.UserUpdated]) {
	globalState.user.name = payload.name;
}

export function lobbyCreatedHandler(payload: EventPayloads[Events.LobbyCreated]) {}

export function joinLobbyFailedHandler(payload: EventPayloads[Events.JoinLobbyFailed]) {}

export function lobbyJoinedHandler(payload: EventPayloads[Events.LobbyJoined]) {}
