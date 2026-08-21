import { Events, type WSEvent } from './events';
import {
	connectionEstablishedHandler,
	joinLobbyFailedHandler,
	lobbyCreatedHandler,
	lobbyJoinedHandler,
	userUpdatedHandler
} from './handlers';

export function routeEvent(event: WSEvent) {
	if (event.type === undefined) {
		// TODO: Error Handling
		throw Error('Undefined type field in event.');
	}

	switch (event.type) {
		case Events.ConnectionEstablished: {
			connectionEstablishedHandler(event.payload);
			break;
		}
		case Events.UserUpdated: {
			userUpdatedHandler(event.payload);
			break;
		}
		case Events.LobbyCreated: {
			lobbyCreatedHandler(event.payload);
			break;
		}
		case Events.JoinLobbyFailed: {
			joinLobbyFailedHandler(event.payload);
			break;
		}
		case Events.LobbyJoined: {
			lobbyJoinedHandler(event.payload);
			break;
		}
	}
}
