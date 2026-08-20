import { connectionState } from '$lib/ws/connection.svelte';

type GlobalState = {
	user: {
		name: string | null;
		user_id: string | null;
	};
	menu: {
		lobbies: MenuLobby[] | null;
		connection: {
			connecting: boolean;
			connected: boolean;
		};
	};
	game: GameState | null;
	lobby: Lobby | null;
};

export const state = $state<GlobalState>({
	user: {
		name: null,
		user_id: null
	},
	menu: {
		lobbies: null,
		connection: connectionState
	},
	game: null,
	lobby: null
});
