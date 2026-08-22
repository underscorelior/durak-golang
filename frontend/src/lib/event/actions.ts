import { sendEvent } from '$lib/ws/connection';
import { Events, type EventPayloads } from './events';

export function createLobby() {
	sendEvent(Events.CreateLobby, {});
}

export function joinLobby(lobby_code: string) {
	const joinLobby = { lobby_code } as EventPayloads[Events.JoinLobby];

	sendEvent(Events.JoinLobby, joinLobby);
}

export function leaveLobby() {
	sendEvent(Events.LeaveLobby, {});
}
