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
    constructor(name) {
        this.name = name
    }
}

class LobbyJoinedEvent {
    constructor(state) {
        this.state = state
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
        case 'lobby_joined': {
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
    const event = new WSEvent(eventName, payload);

    conn.send(JSON.stringify(event));
}

function connectWebsocket(name) {
    if (window['WebSocket']) {
        console.log('supports websocksets');

        conn = new WebSocket(
            'ws://' + document.location.host + '/ws',
        );

        conn.onopen = function () {
            document.getElementById('connection-status').innerHTML =
                'True';
            const joinEvent = new JoinLobbyEvent(name)
            sendEvent('join_lobby', joinEvent)
        };
        conn.onclose = function () {
            document.getElementById('connection-status').innerHTML =
                'False';
            // add some sort of reconnection/jitter protection
        };

        conn.onmessage = function (e) {
            const eventData = JSON.parse(e.data);

            const event = WSEvent(...eventData);
            routeEvent(event);
        };
    } else {
        alert('Browser doesnt support websockets');
    }
}


function joinLobby() {
    let name = document.getElementById('username').value

    connectWebsocket(name)

    // fetch('login', {
    //     method: 'POST',
    //     body: JSON.stringify(formData),
    // })
    //     .then((response) => {
    //         if (response.ok) {
    //             return response.json();
    //         } else {
    //             throw `unauthorized`;
    //         }
    //     })
    //     .then((data) => {
    //         connectWebsocket(data.otp);
    //     })
    //     .catch((e) => {
    //         alert(e);
    //     });

    // const joinEvent = new JoinLobbyEvent(name)
    // console.log("test")
    // sendEvent('join_lobby', joinEvent)

    return false;
}


window.onload = function () {
    document.getElementById('join-lobby').onclick = joinLobby
};