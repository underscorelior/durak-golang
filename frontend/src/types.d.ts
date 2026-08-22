type Lobby = {
	lobby_code: string;
	host_id: string;
	players: Player[];
	max_players: number;
	created_at: string;
	game_state: GameState | null;
};

type LobbyPreview = {
	lobby_code: string;
	host_name: string;
	player_count: number;
	max_players: number;
	created_at: string;

	is_open: boolean;
	is_playing: boolean;
};

type GameState = {
	// TODO: Combine both the lobby players and the game state players locally
	players: { user_id: string; hand_size: number };
	hand: Card[];
	trump: Card;
	deck_size: number;
	turn: Turn;
};

type Player = {
	user_id: string;
	name: string;
	position: string;
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
	defense_card: Card;
	is_defended: boolean;
};

type Turn = {
	table_state: CardPair[];
	defender_id: string;
	initial_attacker_id: string;
	// attackerIds: string[];
	phase: TurnPhase;
};

enum TurnPhase {
	Initial,
	Attack,
	Defense,
	Complete
}
