import { GameEvent, type WSEvent } from './events';
import {
	actionMadeHandler,
	connectionEstablishedHandler,
	errorHandler,
	gameStartedHandler,
	lobbyCreatedHandler,
	lobbyJoinedHandler,
	lobbyLeftHandler,
	lobbyRejoinedHandler,
	memberJoinedHandler,
	memberLeftHandler,
	menuLobbiesUpdatedHandler,
	playerDisconnectedHandler,
	playerRejoinedHandler,
	userUpdatedHandler
} from './handlers';

export function routeEvent(event: WSEvent) {
	if (event.type === undefined) {
		// TODO: Error Handling
		throw Error('Undefined type field in event.');
	}

	switch (event.type) {
		case GameEvent.Error: {
			errorHandler(event.payload);
			break;
		}
		case GameEvent.ConnectionEstablished: {
			connectionEstablishedHandler(event.payload);
			break;
		}
		case GameEvent.MenuLobbiesUpdated: {
			menuLobbiesUpdatedHandler(event.payload);
			break;
		}
		case GameEvent.UserUpdated: {
			userUpdatedHandler(event.payload);
			break;
		}
		case GameEvent.LobbyCreated: {
			lobbyCreatedHandler(event.payload);
			break;
		}
		case GameEvent.LobbyJoined: {
			lobbyJoinedHandler(event.payload);
			break;
		}
		case GameEvent.MemberJoined: {
			memberJoinedHandler(event.payload);
			break;
		}
		case GameEvent.LobbyLeft: {
			lobbyLeftHandler(event.payload);
			break;
		}
		case GameEvent.MemberLeft: {
			memberLeftHandler(event.payload);
			break;
		}
		case GameEvent.PlayerDisconnected: {
			playerDisconnectedHandler(event.payload);
			break;
		}
		case GameEvent.LobbyRejoined: {
			lobbyRejoinedHandler(event.payload);
			break;
		}
		case GameEvent.PlayerRejoined: {
			playerRejoinedHandler(event.payload);
			break;
		}
		case GameEvent.GameStarted: {
			gameStartedHandler(event.payload);
			break;
		}
		case GameEvent.ActionMade: {
			actionMadeHandler(event.payload);
			break;
		}
	}
}
