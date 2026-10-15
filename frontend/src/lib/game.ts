import { globalState } from './state/state.svelte';
import { isCardEmpty } from './util';

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

export function canTake(): boolean {
	const user_id = globalState.user.user_id;
	const turn = globalState.lobby?.game_state?.turn;

	return turn != null && user_id == turn.defender_id && turn.phase == TurnPhase.Defense;
}

export function canPass(): boolean {
	const user_id = globalState.user.user_id;
	const turn = globalState.lobby?.game_state?.turn;

	return (
		turn != null && turn.attacker_ids.includes(user_id || '') && turn.phase == TurnPhase.Attack
	); // TODO: Need a way to figure this out
}

export function canDefend(): boolean {
	const user_id = globalState.user.user_id;
	const turn = globalState.lobby?.game_state?.turn;

	return turn != null && turn.defender_id == user_id && turn.phase == TurnPhase.Defense;
}

export function canAttack(): boolean {
	const user_id = globalState.user.user_id;
	const turn = globalState.lobby?.game_state?.turn;

	return (
		turn != null &&
		((turn.attacker_ids.includes(user_id || '') && turn.phase == TurnPhase.Attack) ||
			(turn.initial_attacker_id == user_id && turn.phase == TurnPhase.Initial))
	);
}

export function cardDefends(pair: CardPair, defense_card: Card): boolean {
	const trump_suit = globalState.lobby?.game_state?.trump.suit;
	const atk = pair.attack_card;

	if (
		isCardEmpty(atk) ||
		(pair.defense_card && !isCardEmpty(pair.defense_card)) ||
		pair.is_defended
	) {
		return false;
	}

	if (atk.suit == defense_card.suit) {
		return atk.rank < defense_card.rank;
	}

	return defense_card.suit == trump_suit;
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
