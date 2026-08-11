let conn;

class WSEvent {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}

class ConnectionEstablishedEvent {
    constructor(name) {
        this.name = name
    }
}

class UpdateUserEvent {
    constructor(name) {
        this.name = name
    }
}

class UserUpdatedEvent {
    constructor(name) {
        this.name = name
    }
}

class JoinLobbyEvent {
    constructor(lobbyID) {
        this.lobbyID = lobbyID
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

class LobbyUpdatedEvent {
    constructor(delta_name, delta_operation, delta_payload) {
        self.delta_name = delta_name
        self.delta_operation = delta_operation
        self.delta_payload = delta_payload
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
        case 'connection_established': {
            const connectionEstablishedEvent = new ConnectionEstablishedEvent(event.payload.name)
            connectionEstablishedHandler(connectionEstablishedEvent)
            break;
        }
        case 'update_user': {
            break;
        }
        case 'user_updated': {
            const userUpdatedEvent = new UserUpdatedEvent(event.payload.name) // find a better pattern
            userUpdatedHandler(userUpdatedEvent)
            break;
        }
        case 'join_lobby': {
            // const joinEvent = new JoinLobbyEvent(...event.payload)
            // joinLobby(joinEvent)
            break;
        }
        case 'lobby_joined': {
            const lobbyJoinedEvent = new LobbyJoinedEvent(...event.payload)
            lobbyJoinedHandler(lobbyJoinedEvent)
            break;
        }
        case 'leave_lobby': {
            // const leaveEvent = new LeaveLobbyEvent(...event.payload)
            // leaveLobby(leaveEvent)
            break;
        }
        case 'lobby_updated': {
            // const updateEvent = new LobbyUpdatedEvent(...event.payload)
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

function connectWebsocket() {
    if (window['WebSocket']) {
        console.log('supports websocksets');

        conn = new WebSocket(
            'ws://' + document.location.host + '/ws?username=default',
        );

        conn.onopen = function () {
            document.getElementById('connection-status').innerHTML =
                'True';
            document.getElementById('current-lobby').innerHTML = 'Not Connected'
        };
        conn.onclose = function () {
            document.getElementById('connection-status').innerHTML =
                'False';
            document.getElementById('current-lobby').innerHTML = 'Disconnected'
            document.getElementById('current-name').innerHTML = 'N/A'
            // add some sort of reconnection/jitter protection
        };

        conn.onmessage = function (e) {
            const eventData = JSON.parse(e.data);

            const event = new WSEvent(eventData.type, eventData.payload);
            routeEvent(event);
        };
    } else {
        alert('Browser doesnt support websockets');
    }
}

function connectionEstablishedHandler(updatedEvent) {
    console.log(updatedEvent)
    document.getElementById('current-name').innerHTML = updatedEvent.name;
}

function updateUser() {
    let name = document.getElementById('username').value

    const updateEvent = new UpdateUserEvent(name)
    sendEvent('update_user', updateEvent)

    return false;
}

function userUpdatedHandler(updatedEvent) {
    console.log(updatedEvent)
    document.getElementById('current-name').innerHTML = updatedEvent.name;
}


function joinLobby() {
    let lobbyID = document.getElementById('lobby-id').value

    const joinEvent = new JoinLobbyEvent(lobbyID)
    sendEvent('join_lobby', joinEvent)


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

    return false;
}

function lobbyJoinedHandler(event) {

}


window.onload = function () {
    document.getElementById('join-lobby').onclick = joinLobby
    document.getElementById('set-name').onclick = updateUser
    connectWebsocket()
};