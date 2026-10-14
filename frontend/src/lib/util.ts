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
