import { globalState } from './state/state.svelte';

export const rankToString = [
	'',
	'A',
	'2',
	'3',
	'4',
	'5',
	'6',
	'7',
	'8',
	'9',
	'10',
	'J',
	'Q',
	'K',
	'A'
];

export const suitToColor = ['#000', '#c00', '#c00', '#000'];
export const suitToSymbol = ['♣', '♦', '♥', '♠'];

export type SortingPreference = 'rank-suit' | 'suit-rank' | 'trump-rank-suit' | 'trump-suit-rank';

const rankSort = (c1: Card, c2: Card) => c2.rank - c1.rank;
const suitSort = (c1: Card, c2: Card) => c2.suit - c1.suit;

export function sortHand(hand?: Card[], sort: SortingPreference = 'rank-suit'): Card[] {
	if (!hand) return [];

	if (sort == 'rank-suit') {
		hand = hand.sort(suitSort);
		hand = hand.sort(rankSort);
	} else if (sort == 'suit-rank') {
		hand = hand.sort(rankSort);
		hand = hand.sort(suitSort);
	}

	return hand;
}

export function shID(uuid: string | null = '') {
	return (uuid || '').split('-')[0];
}

export function fanDistance(index: number, length: number) {
	let dist = index - (length - 1) / 2;

	if (length % 2 == 0) {
		dist = dist < 0 ? dist - 1 : dist + 1;
	}

	return dist < 0 ? Math.ceil(dist) : Math.floor(dist);
}

export function findMemberByID(user_id: string): Member | undefined {
	const players = globalState.lobby?.members;

	if (players == null) {
		return undefined;
	}
	return players.find((p) => p.user_id === user_id);
}

export function findPlayerByID(user_id: string): Player | undefined {
	return globalState.lobby?.game_state?.players.find((p) => p.user_id === user_id);
}

export function cardInHand(card: Card): boolean {
	const hand = globalState.lobby?.game_state?.hand;

	if (hand == null) {
		return false;
	}

	return hand.some((c) => card.rank === c.rank && card.suit === c.suit);
}

export function isCardEmpty(card: Card): boolean {
	return card.rank == 0 && card.suit == 0;
}

export const EMPTY_CARD: Card = { suit: 0, rank: 0 };
