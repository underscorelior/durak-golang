type Lobby = {
	lobby_code: string;
	host_id: string;
	is_hidden: boolean;
	in_game: boolean;
	max_members: number;
	created_at: string;
	members: Member[];
	position: number;
	game_state: GameState | null;
};

type LobbyPreview = {
	lobby_code: string;
	host_name: string;
	member_count: number;
	max_members: number;
	created_at: string;

	is_open: boolean;
	in_game: boolean;
};

type Session = {
	lobby_code: string;
	session_token: string;
};

type Player = {
	user_id: string;
	hand_size: number;
	position: number;
	passed: boolean;
};

type GameState = {
	// TODO: Combine both the lobby players and the game state players locally
	players: Player[];
	hand: Card[];
	trump: Card;
	deck_size: number;
	turn: Turn;
	first_defense: boolean;
};

type Member = {
	user_id: string;
	name: string;
	position: number;
	is_connected: boolean;

	joined_at: Date;
	disconnected_at: Date | null;
};

enum Suit {
	Club,
	Diamond,
	Heart,
	Spade
}

type Card = {
	suit: Suit;
	rank: number;
};

type CardPair = {
	attack_card: Card;
	defense_card?: Card;
};

type Turn = {
	table_state: CardPair[];
	defender_id: string;
	initial_attacker_id: string;
	attacker_ids: string[];
	phase: TurnPhase;
};
