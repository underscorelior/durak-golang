type Lobby = {
	lobbyId: string;
	hostId: string;
	players: Player[];
	maxPlayers: number;
	gameState: GameState;
};

type MenuLobby = {
	lobbyId: string;
	hostName: string;
	playerCount: number;
	maxPlayers: number;

	isOpen: boolean;
	isPlaying: boolean;
};

type GameState = {
	// TODO: Combine both the lobby players and the game state players locally
	players: { userId: string; handSize: number };
	hand: Card[];
	trump: Card;
	deckSize: number;
	turn: Turn;
};

type Player = {
	userId: string;
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
	attackCard: Card;
	defenseCard: Card;
	isDefended: boolean;
};

type Turn = {
	tableState: CardPair[];
	defenderId: string;
	initialAttackerId: string;
	// attackerIds: string[];
	phase: TurnPhase;
};

enum TurnPhase {
	Initial,
	Attack,
	Defense,
	Complete
}
