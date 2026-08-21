import { sendEvent } from '$lib/ws/connection';
import { Events } from './events';

export function createLobby() {
	sendEvent(Events.CreateLobby, {});
}
