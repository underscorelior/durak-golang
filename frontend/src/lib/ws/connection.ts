import type { EventPayloads, WSEvent } from '$lib/event/events';
import { routeEvent } from '$lib/event/routing';
import { PUBLIC_WEBSOCKET_URL } from '$env/static/public';
import { connectionState, globalState } from '$lib/state/state.svelte';

let socket: WebSocket | null = null;

export default function connectWebsocket(): Error | null {
	if (socket) return null;
	if (window['WebSocket']) {
		const session = getLocalSession();

		if (session !== null) {
			globalState.session = JSON.parse(session);
		}

		connectionState.connecting = true;

		socket = new WebSocket(PUBLIC_WEBSOCKET_URL + (hasLocalSession() ? `?session=${session}` : ''));

		if (socket === null) {
			connectionState.connecting = false;
			return Error('Failed to connect to WebSocket server.');
		}

		socket.onopen = () => {
			connectionState.connected = true;
			connectionState.connecting = false;
		};

		// TODO: Add some sort of jitter protection
		socket.onclose = () => {
			connectionState.connected = false;
			connectionState.connecting = false;
			socket = null;
		};
		// TODO: Support conn.onerror?
		socket.onmessage = onConnectionMessage;

		return null;
	}

	return Error('WebSockets are not supported by your browser.');
}

function onConnectionMessage(event: MessageEvent<string>) {
	const message: WSEvent = JSON.parse(event.data);

	routeEvent(message);
}

export function sendEvent<E extends keyof EventPayloads>(eventName: E, payload: EventPayloads[E]) {
	const event = { type: eventName, payload } as WSEvent;

	const message = JSON.stringify(event);

	socket?.send(message);
}

export function getLocalSession() {
	if (hasLocalSession()) {
		return localStorage.getItem('durak_session') as string;
	}

	return null;
}

export function hasLocalSession() {
	const session = localStorage.getItem('durak_session');
	console.log(session);
	if (session == null) {
		return false;
	}

	const parsedSession = JSON.parse(session) as Session;
	if (parsedSession.lobby_code && parsedSession.session_token) {
		return true;
	}
	return false;
}

export function clearLocalSession() {
	localStorage.setItem('durak_session', '{}');
}
