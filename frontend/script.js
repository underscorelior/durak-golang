let conn;

class GameState {
    constructor() {
        this.name = "default"
        this.userId = ""
        this.lobbyId = ""
        this.lobbyHost = ""
        this.players = []
        this.position = -1
        this.maxPlayers = -1
    }
}

let state = new GameState()

class WSEvent {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}

class ConnectionEstablishedEvent {
    constructor({ name, userId, lobbies }) {
        this.name = name
        this.userId = userId
        this.lobbies = lobbies
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

class CreateLobbyEvent {
    constructor() {

    }
}

class LobbyCreatedEvent {
    constructor(lobbyId) {
        this.lobbyId = lobbyId
    }
}

class JoinLobbyEvent {
    constructor(lobbyId) {
        this.lobbyId = lobbyId
    }
}

// TODO: In order to avoid passing in lobbyID, maybe a system where I can reference a past request and grab it from there?
class JoinLobbyFailedEvent {
    constructor({ code, message, lobbyId }) {
        this.code = code
        this.message = message
        this.lobbyId = lobbyId
    }
}

class LobbyJoinedEvent {
    constructor(lobby) {
        this.lobby = lobby
    }
}

class PlayerJoinedEvent {
    constructor(player) {
        this.player = player
    }
}


function routeEvent(event) {
    if (event.type === undefined) {
        alert('no type field in the event');
    }

    switch (event.type) {
        case 'connection_established': {
            const connectionEstablishedEvent = new ConnectionEstablishedEvent(event.payload)
            connectionEstablishedHandler(connectionEstablishedEvent)
            break;
        }
        case 'user_updated': {
            const userUpdatedEvent = new UserUpdatedEvent(event.payload.name) // find a better pattern
            userUpdatedHandler(userUpdatedEvent)
            break;
        }
        case 'lobby_created': {
            const lobbyCreatedEvent = new LobbyCreatedEvent(event.payload.lobbyId)
            lobbyCreatedHandler(lobbyCreatedEvent)
            break;
        }
        case 'join_lobby_failed': {
            const joinLobbyFailedEvent = new JoinLobbyFailedEvent(event.payload)
            joinLobbyFailedHandler(joinLobbyFailedEvent)
            break;
        }
        case 'lobby_joined': {
            const lobbyJoinedEvent = new LobbyJoinedEvent(event.payload.lobby)
            lobbyJoinedHandler(lobbyJoinedEvent)
            break;
        }
        case 'player_joined': {
            const playerJoinedEvent = new PlayerJoinedEvent(event.payload.player)
            playerJoinedHandler(playerJoinedEvent)
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
    state.name = updatedEvent.name
    state.userId = updatedEvent.userId
    state.lobbies = updatedEvent.lobbies

    document.getElementById('current-name').innerHTML = state.name;
}

function updateUser() {
    let name = document.getElementById('username').value

    const updateEvent = new UpdateUserEvent(name)
    sendEvent('update_user', updateEvent)

    return false;
}

function userUpdatedHandler(updatedEvent) {
    state.name = updatedEvent.name
    document.getElementById('current-name').innerHTML = state.name;
}

function createLobby() {
    const createEvent = new CreateLobbyEvent()
    sendEvent('create_lobby', createEvent)

    return false;
}

function lobbyCreatedHandler(event) {
    joinLobby(event.lobbyId)
}


function joinLobby(lobbyId = null) {
    let lobbyID = lobbyId ?? document.getElementById('lobby-id').value

    const joinEvent = new JoinLobbyEvent(lobbyID)
    sendEvent('join_lobby', joinEvent)

    return false;
}

function joinLobbyFailedHandler(event) {
    alert(`${event.code}: ${event.message} (${event.lobbyId})`)
}

function lobbyJoinedHandler(event) {
    state.lobbyId = event.lobby.lobbyId
    state.players = event.lobby.players
    state.lobbyHost = event.lobby.host
    state.maxPlayers = event.lobby.maxPlayers
    state.position = event.lobby.position
    document.getElementById('current-lobby').innerHTML = state.lobbyId
    document.getElementById('lobby-players').innerHTML = JSON.stringify(state.players)
    document.getElementById('lobby-host').innerHTML = state.lobbyHost === state.userId
    updatePlayerDisplay(state.position, state.players)
}

function playerJoinedHandler(event) {
    state.players.push(event.player)
    document.getElementById('lobby-players').innerHTML = JSON.stringify(state.players)
    updatePlayerDisplay(state.position, state.players)
}

function updatePlayerDisplay(position, players) {
    let you = players.find((pl) => pl.position == position)
    if (you)
        document.getElementById('bottom-box').innerHTML = `${you.name} (${you.position})`// You

    let left = players.find((pl) => pl.position == (position + 1) % state.maxPlayers)
    if (left)
        document.getElementById('left-box').innerHTML = `${left.name} (${left.position})` // Player to your left (position + 1 % maxPlayers)

    let front = players.find((pl) => pl.position == (position + 2) % state.maxPlayers)
    if (front)
        document.getElementById('front-box').innerHTML = `${front.name} (${front.position})` // Player to your front (position + 2 % maxPlayers)

    let right = players.find((pl) => pl.position == (position + 3) % state.maxPlayers)
    if (right)
        document.getElementById('right-box').innerHTML = `${right.name} (${right.position})` // player to your right (positon + 3 % maxPlayers)
}

window.onload = function () {
    document.getElementById('join-lobby').addEventListener("click", () => joinLobby())
    document.getElementById('set-name').onclick = updateUser
    document.getElementById('create-lobby').onclick = createLobby
    connectWebsocket()
};