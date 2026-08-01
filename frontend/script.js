let conn;

class WSEvent {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}

/** Lobby Events:
 * CreateLobby
   * 
 * JoinLobby
   * url
   * players
 * LobbyUpdate (player joins, etc.)
   * Subroute (player update (?))
   * Way to differentiate between add and remove (Join and Leave lobby broadcass this)
 * LeaveLobby
 */

class JoinLobbyEvent {
    constructor(players) {
        self.players = players
    }
}

class LeaveLobbyEvent {
    constructor() {

    }
}

class UpdateLobbyEvent {
    constructor() {

    }
}

/** Game Events:
 * GameStart
   * game object
 */

function routeEvent(event) {
    if (event.type === undefined) {
        alert('no type field in the event');
    }

    switch (event.type) {
        case 'join_lobby': {
            // const joinEvent = new JoinLobbyEvent(...event.payload)
            // joinLobby(joinEvent)
            break;
        }
        case 'leave_lobby': {
            // const leaveEvent = new LeaveLobbyEvent(...event.payload)
            // leaveLobby(leaveEvent)
            break;
        }
        case 'update_lobby': {
            // const updateEvent = new UpdateLobbyEvent(...event.payload)
            // joinLobby(joinEvent)
            break;
        }
        default: {
            alert('unsupported message type');
            break;
        }
    }
}

function sendEvent(eventName, payload) {
    const event = new Event(eventName, payload);

    conn.send(JSON.stringify(event));
}



window.onload = function () {
    if (window['WebSocket']) {
        console.log('supports websocksets');

        conn = new WebSocket(
            'ws://' + document.location.host + '/ws',
        );

        conn.onmessage = function (e) {
            const eventData = JSON.parse(e.data);

            const event = new WSEvent(eventData.type, eventData.payload);

            routeEvent(event);
        };
    } else {
        alert('Browser doesnt support websockets');
    }
};

