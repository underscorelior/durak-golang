import { Events, type WSEvent } from './events';
import { connectionEstablishedHandler } from './handlers';

export function routeEvent(event: WSEvent) {
	if (event.type === undefined) {
		// TODO: Error Handling
		throw Error('Undefined type field in event.');
	}

	switch (event.type) {
		case Events.EventConnectionEstablished: {
			connectionEstablishedHandler(event.payload);
			break;
		}
	}
}
