import { globalState } from './state/state.svelte';
import { cardInHand, isCardEmpty } from './util';

export enum TurnPhase {
	INITIAL = 0,
	GAME,
	TAKE,
	DRAW,
	COMPLETE
}

export enum Action {
	DEFEND = 0,
	ATTACK,
	TAKE,
	PASS
}

/**
 * Checks if user is a defender; e.g. the user is the defender and the phase is GAME
 * @returns {boolean} `true` if the user can defend, `false` if not
 */
export function isDefender(): boolean {
	const user_id = globalState.user.user_id;
	const turn = globalState.lobby?.game_state?.turn;

	if (turn == null || user_id == null) {
		return false;
	}

	return turn.defender_id == user_id && turn.phase == TurnPhase.GAME;
}

/**
 * Checks if user is an attacker, not whether an attack could be valid.
 *  - User is initial attacker during INITIAL phase
 *  - User is an attacker during GAME phase
 * @returns {boolean} `true` if the user can attack, `false` if not
 */
export function isAttacker(): boolean {
	const user_id = globalState.user.user_id;
	const turn = globalState.lobby?.game_state?.turn;

	if (turn == null || user_id == null) {
		return false;
	}

	if (turn.phase == TurnPhase.INITIAL && turn.initial_attacker_id == user_id) {
		return true;
	}

	return turn.attacker_ids.includes(user_id) && turn.phase == TurnPhase.GAME;
}

/**
 * Checks if user can take the cards; e.g. the user is the defender and the phase is GAME
 * @returns {boolean} `true` if the user can take cards, `false` if not
 */
export function canTake(): boolean {
	const user_id = globalState.user.user_id;
	const turn = globalState.lobby?.game_state?.turn;

	if (turn == null || user_id == null) {
		return false;
	}

	return user_id == turn.defender_id && turn.phase == TurnPhase.GAME;
}

/**
 * Checks if user is can pass their turn; e.g. the user is an attacker and the phase is GAME
 * @returns {boolean} `true` if the user can pass their turn, `false` if not
 */
export function canPass(): boolean {
	const user_id = globalState.user.user_id;
	const turn = globalState.lobby?.game_state?.turn;

	if (turn == null || user_id == null) {
		return false;
	}

	return turn.attacker_ids.includes(user_id) && turn.phase == TurnPhase.GAME;
}

/**
 * Checks if the user can attack using the specific card; e.g. user is a valid attacker, the card is in hand, and either the phase is INITIAL or card's rank exists within the table
 * @param card card to be used in attack
 * @returns {boolean} `true` if is a valid attack, `false` if not
 */
export function canAttack(card: Card): boolean {
	const turn = globalState.lobby?.game_state?.turn;

	if (turn == null) {
		return false;
	}

	if (!isAttacker() || !cardInHand(card)) {
		return false;
	}

	if (turn.phase == TurnPhase.INITIAL) {
		return true;
	}

	return canPlace(card);
}

/**
 * Checks if a card can be used as an attack, e.g. has the card's rank already been played
 * @param card attack card to check against
 * @returns {boolean} `true` if the card is placable under current conditions, `false` if not
 */
export function canPlace(card: Card): boolean {
	const turn = globalState.lobby?.game_state?.turn;
	const pairs = globalState.lobby?.game_state?.turn.table_state;

	if (turn == null || pairs == null) {
		return false;
	}

	const validRanks = new Set<number>(
		pairs.flatMap((pair) => [pair.attack_card.rank, pair.defense_card?.rank || 0])
	);

	return card.rank != 0 && validRanks.has(card.rank);
}

/**
 * Checks if a defense is a valid move; e.g. user is the defender, user has the card in hand, and the card is able to beat the attack card.
 * @param pair card pair containing the attack card to be defended
 * @param defense_card card used for defending
 * @returns {boolean} `true` if card is able to defend against attack card, `false` if not
 */
export function canDefend(pair: CardPair, defense_card: Card): boolean {
	const turn = globalState.lobby?.game_state?.turn;
	const trump_suit = globalState.lobby?.game_state?.trump.suit;
	const atk = pair.attack_card;

	if (turn == null && trump_suit == null) {
		return false;
	}

	if (!isDefender() || !cardInHand(defense_card)) {
		return false;
	}

	if (isCardEmpty(atk) || (pair.defense_card != null && !isCardEmpty(pair.defense_card))) {
		return false;
	}

	if (atk.suit == defense_card.suit) {
		return atk.rank < defense_card.rank;
	}

	return defense_card.suit == trump_suit;
}
