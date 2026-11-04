import { globalState } from '$lib/state/state.svelte';
import { joinLobby, rejoinLobby } from './actions';
import { GameEvent, type EventPayloads } from './events';
import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import { clearLocalSession, setLocalSession } from '$lib/ws/connection';
import { findMemberByID, sortHand } from '$lib/util';
import { Action } from '$lib/game';

// TODO: Proper err handling
export function errorHandler(payload: EventPayloads[GameEvent.Error]) {
	alert(`${payload.code}: ${payload.message} (${payload.error_for})`);

	if (payload.error_for == GameEvent.JoinLobby || payload.error_for == GameEvent.RejoinLobby) {
		clearLocalSession();
		goto(resolve('/'));
	}
}

export function connectionEstablishedHandler(
	payload: EventPayloads[GameEvent.ConnectionEstablished]
) {
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

export function menuLobbiesUpdatedHandler(payload: EventPayloads[GameEvent.MenuLobbiesUpdated]) {
	globalState.menu.lobbies = payload.lobbies;
}

export function userUpdatedHandler(payload: EventPayloads[GameEvent.UserUpdated]) {
	globalState.user.name = payload.name;
}

export function lobbyCreatedHandler(payload: EventPayloads[GameEvent.LobbyCreated]) {
	joinLobby(payload.lobby_code);
}

export function lobbyJoinedHandler(payload: EventPayloads[GameEvent.LobbyJoined]) {
	globalState.lobby = payload.lobby;

	goto(resolve(`/${payload.lobby.lobby_code}`));
}

export function memberJoinedHandler(payload: EventPayloads[GameEvent.MemberJoined]) {
	const members = [...(globalState.lobby?.members ?? []), payload.member].sort(
		(a, b) => a.position - b.position
	);
	globalState.lobby = { ...globalState.lobby, members } as Lobby;
}

export function lobbyLeftHandler(payload: EventPayloads[GameEvent.LobbyLeft]) {
	clearLocalSession();
	goto(resolve('/'));

	globalState.lobby = null;
	globalState.menu.lobbies = payload.lobbies;
}

export function memberLeftHandler(payload: EventPayloads[GameEvent.MemberLeft]) {
	const members = globalState.lobby?.members.filter((p) => p.user_id != payload.user_id);
	members?.forEach((m) => (m.position -= Number(m.position >= payload.position)));
	globalState.lobby = { ...globalState.lobby, host_id: payload.host_id, members } as Lobby;
}

export function playerDisconnectedHandler(payload: EventPayloads[GameEvent.PlayerDisconnected]) {
	const disconnected = findMemberByID(payload.user_id);

	if (disconnected) {
		disconnected.disconnected_at = payload.disconnected_at;
		disconnected.is_connected = false;
	}
}

export function lobbyRejoinedHandler(payload: EventPayloads[GameEvent.LobbyRejoined]) {
	globalState.lobby = payload.lobby;

	globalState.lobby.game_state = {
		...globalState.lobby.game_state,
		hand: sortHand(globalState.lobby.game_state?.hand),
		turn: {
			...payload.lobby.game_state?.turn,
			attacker_ids: Object.keys(
				payload.lobby.game_state?.turn.attacker_ids as never as { [key: string]: boolean } // TODO: Fix this up so that it makes better sense, maybe a different way within golang itself
			)
		}
	} as GameState; // TODO: Fix thsi?
}

export function playerRejoinedHandler(payload: EventPayloads[GameEvent.PlayerRejoined]) {
	const rejoined = findMemberByID(payload.user_id);

	if (rejoined) {
		rejoined.is_connected = true;
	}
}

export function gameStartedHandler(payload: EventPayloads[GameEvent.GameStarted]) {
	const session = payload.session;

	setLocalSession(session);

	globalState.lobby = payload.lobby;

	globalState.lobby.game_state = {
		...globalState.lobby.game_state,
		hand: sortHand(globalState.lobby.game_state?.hand),
		turn: {
			...payload.lobby.game_state?.turn,
			attacker_ids: Object.keys(
				payload.lobby.game_state?.turn.attacker_ids as never as { [key: string]: boolean } // TODO: Fix this up so that it makes better sense, maybe a different way within golang itself
			)
		}
	} as GameState; // TODO: Fix thsi?
}

export function actionMadeHandler(payload: EventPayloads[GameEvent.ActionMade]) {
	const game_state = globalState.lobby?.game_state;

	if (!game_state || !globalState.lobby?.game_state) {
		alert('ERROR LOBBY OR GAME STATE NULL');
		return;
	}
	const playerIndex = game_state.players.findIndex((gp) => gp.user_id == payload.user_id);

	const { user_id, action, turn_phase, card, pair_index } = payload;

	if ([Action.ATTACK, Action.DEFEND].includes(action) && pair_index && card) {
		const pair = game_state.turn.table_state[pair_index];

		if (!pair) {
			alert("Couldn't find card pair for actionMade");
		}

		if (action == Action.ATTACK) {
			pair.attack_card = card;
		}

		if (action == Action.DEFEND) {
			pair.defense_card = card;
		}

		game_state.players[playerIndex].hand_size -= 1;
		if (user_id == globalState.user.user_id) {
			globalState.lobby.game_state.hand = game_state.hand.filter(
				(c) => c.rank != card?.rank || c.suit != card?.suit
			);
		}
	}

	game_state.turn.phase = turn_phase;
}
