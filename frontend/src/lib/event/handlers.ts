import { globalState } from '$lib/state/state.svelte';
import { joinLobby, rejoinLobby } from './actions';
import type { EventPayloads, Events } from './events';
import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import { clearLocalSession, setLocalSession } from '$lib/ws/connection';
import { sortHand } from '$lib/util';
import { Action } from '$lib/game';

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

export function actionMadeHandler(payload: EventPayloads[Events.ActionMade]) {
	const game_state = globalState.lobby?.game_state;

	if (!game_state || !globalState.lobby?.game_state) {
		alert('ERROR LOBBY OR GAME STATE NULL');
		return;
	}
	const playerIndex = game_state.players.findIndex((gp) => gp.user_id == payload.user_id);

	const { user_id, action, turn_phase, card, pair_index } = payload;

	if ([Action.Attack, Action.Defense].includes(action) && pair_index && card) {
		const pair = game_state.turn.table_state[pair_index];

		if (!pair) {
			alert("Couldn't find card pair for actionMade");
		}

		if (action == Action.Attack) {
			pair.attack_card = card;
		}

		if (action == Action.Defense && pair_index != null && card) {
			pair.defense_card = card;
			pair.is_defended = true;
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

export function actionFailedHandler(payload: EventPayloads[Events.ActionFailed]) {
	alert(`${payload.code}: ${payload.message}`);
}
