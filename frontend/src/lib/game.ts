import { globalState } from './state/state.svelte';

export enum Action {
	Defense = 0,
	Attack,
	Pass,
	Take
}

export enum TurnPhase {
	Initial = 0,
	Attack,
	Defense,
	Complete
}

export function validateMove() {}

export function canDefend() {
	return true;
}

export function canAttack() {
	return true;
}

export function isAttacking(): boolean {
	const game_state = globalState.lobby?.game_state;
	if (!game_state) {
		return false;
	}

	if (
		game_state.turn.phase == TurnPhase.Initial &&
		game_state.turn.initial_attacker_id == globalState.user.user_id &&
		canDefend()
	) {
		return true;
	}

	// TODO: FIX THE || ''
	if (
		game_state.turn.attacker_ids.includes(globalState.user.user_id || '') &&
		game_state.turn.phase == TurnPhase.Attack &&
		canAttack()
	) {
		return true;
	}

	return false;
}

export function isDefending(): boolean {
	const game_state = globalState.lobby?.game_state;
	if (!game_state) {
		return false;
	}

	if (
		game_state.turn.defender_id == globalState.user.user_id &&
		game_state.turn.phase == TurnPhase.Defense
	) {
		return true;
	}
	return false;
}
