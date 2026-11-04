import type { Action } from '$lib/game';
import { sendEvent } from '$lib/ws/connection';
import { GameEvent, type EventPayloads } from './events';

export function updateUser(name: string) {
	const updateUser = { name } as EventPayloads[GameEvent.UpdateUser];
	sendEvent(GameEvent.CreateLobby, updateUser);
}

export function createLobby() {
	sendEvent(GameEvent.CreateLobby, {});
}

export function joinLobby(lobby_code: string) {
	const joinLobby = { lobby_code } as EventPayloads[GameEvent.JoinLobby];

	sendEvent(GameEvent.JoinLobby, joinLobby);
}

export function rejoinLobby(session: Session) {
	const rejoinLobby = { ...session } as EventPayloads[GameEvent.RejoinLobby];

	sendEvent(GameEvent.RejoinLobby, rejoinLobby);
}

export function leaveLobby() {
	sendEvent(GameEvent.LeaveLobby, {});
}

export function startGame() {
	sendEvent(GameEvent.StartGame, {});
}

export function makeAction(action: Action, card?: Card, pair_index?: number) {
	const makeAction = {
		action,
		card,
		pair_index
	} as EventPayloads[GameEvent.MakeAction];

	sendEvent(GameEvent.MakeAction, makeAction);
}
