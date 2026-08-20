let conn;

class GameState {
    constructor() {
        this.name = "default"
        this.lobby = {}
        this.userId = ""
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

class GameStartedEvent {
    constructor(lobby) {
        this.lobby = lobby
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
        case 'game_started': {
            const gameStartedEvent = new GameStartedEvent(event.payload.lobby)
            gameStartedHandler(gameStartedEvent)
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
    sendEvent('create_lobby', {})
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
    state.lobby = event.lobby
    document.getElementById('current-lobby').innerHTML = state.lobby.lobbyId
    document.getElementById('lobby-players').innerHTML = JSON.stringify(state.lobby.players)
    document.getElementById('lobby-host').innerHTML = state.lobby.hostId === state.userId
    updatePlayerDisplay(state.lobby.position, state.lobby.players)

    if (state.lobby.hostId === state.userId && state.lobby.players.length > 1) {
        document.getElementById('start-game').disabled = false
    }
}

function playerJoinedHandler(event) {
    let isHost = state.lobby.hostId === state.userId
    state.lobby.players.push(event.player)
    document.getElementById('lobby-players').innerHTML = JSON.stringify(state.lobby.players)
    updatePlayerDisplay(state.lobby.position, state.lobby.players)
    console.log(isHost, state.lobby.players.length, state)
    if (isHost && state.lobby.players.length > 1) {
        console.log("Is host and len > 1")
        document.getElementById('start-game').disabled = false
    }
}

function startGame() {
    sendEvent('start_game', {})
    return false
}

function gameStartedHandler(event) {
    state.lobby = event.lobby
    updatePlayerDisplay(state.lobby.position, state.lobby.players)
    document.getElementById("hand").innerHTML = JSON.stringify(state.lobby.gameState.hand)
    document.getElementById("trump").innerHTML = JSON.stringify(state.lobby.gameState.trump)
    document.getElementById("deck-size").innerHTML = state.lobby.gameState.deckSize
}

function updatePlayerDisplay(position, players) {
    let you = players.find((pl) => pl.position == position)
    if (you)
        document.getElementById('bottom-box').innerHTML = `${you.name} (${you.position})`// You

    let left = players.find((pl) => pl.position == (position + 1) % state.lobby.maxPlayers)
    if (left)
        document.getElementById('left-box').innerHTML = `${left.name} (${left.position})` // Player to your left (position + 1 % maxPlayers)

    let front = players.find((pl) => pl.position == (position + 2) % state.lobby.maxPlayers)
    if (front)
        document.getElementById('front-box').innerHTML = `${front.name} (${front.position})` // Player to your front (position + 2 % maxPlayers)

    let right = players.find((pl) => pl.position == (position + 3) % state.lobby.maxPlayers)
    if (right)
        document.getElementById('right-box').innerHTML = `${right.name} (${right.position})` // player to your right (positon + 3 % maxPlayers)
}

window.onload = function () {
    document.getElementById('join-lobby').addEventListener("click", () => joinLobby())
    document.getElementById('set-name').onclick = updateUser
    document.getElementById('create-lobby').onclick = createLobby
    document.getElementById('start-game').onclick = startGame
    connectWebsocket()
};