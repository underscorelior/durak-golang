import { sendEvent } from '$lib/ws/connection';
import { Events, type EventPayloads } from './events';

export function createLobby() {
	sendEvent(Events.CreateLobby, {});
}

export function joinLobby(lobby_code: string) {
	const joinLobby = { lobby_code } as EventPayloads[Events.JoinLobby];

	sendEvent(Events.JoinLobby, joinLobby);
}

export function rejoinLobby(session: Session) {
	const rejoinLobby = { ...session } as EventPayloads[Events.RejoinLobby];

	sendEvent(Events.RejoinLobby, rejoinLobby);
}

export function leaveLobby() {
	sendEvent(Events.LeaveLobby, {});
}
