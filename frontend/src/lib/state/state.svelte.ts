type GlobalState = {
	user: {
		name: string | null;
		user_id: string | null;
	};
	menu: {
		lobbies: LobbyPreview[] | null;
		connection: {
			connecting: boolean;
			connected: boolean;
		};
	};
	lobby: Lobby | null;
	session: Session | null;
	game_state: GameState | null;
};

export const connectionState = $state<GlobalState['menu']['connection']>({
	connected: false,
	connecting: false
});

export const globalState = $state<GlobalState>({
	user: {
		name: null,
		user_id: null
	},
	menu: {
		lobbies: null,
		connection: connectionState
	},
	lobby: null,
	session: null,
	game_state: null
});
