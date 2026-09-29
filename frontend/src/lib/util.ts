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

export type SortingPreference = 'suit-rank' | 'rank-suit' | 'trump-suit-rank' | 'trump-rank-suit';

export function sortHand(hand: Card[], sort: SortingPreference = 'rank-suit') {}

export function shID(uuid: string | null = '') {
	return (uuid || '').split('-')[0];
}
