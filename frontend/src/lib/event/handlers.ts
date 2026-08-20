import { state } from '$lib/state/state.svelte';
import type { EventPayloads, Events } from './events';

export function connectionEstablishedHandler(
	payload: EventPayloads[Events.EventConnectionEstablished]
) {
	state.user.name = payload.name;
	state.user.user_id = payload.userId;
	state.menu.lobbies = payload.lobbies;
}
