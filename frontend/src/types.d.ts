type Lobby = {
	lobby_code: string;
	host_id: string;
	is_private: boolean;
	is_started: boolean;
	max_players: number;
	created_at: string;
	players: Player[];
	position: number;
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

type Session = {
	lobby_code: string;
	session_token: string;
};

type GamePlayer = {
	user_id: string;
	hand_size: number;
	position: number;
};

type GameState = {
	// TODO: Combine both the lobby players and the game state players locally
	players: GamePlayer[];
	hand: Card[];
	trump: Card;
	deck_size: number;
	turn: Turn;
};

type Player = {
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
