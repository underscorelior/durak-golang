# Durak Multiplayer

This is an online multiplayer port of the game [Durak](https://en.wikipedia.org/wiki/Durak).

Built on Golang using Gorilla Websockets to communicate between the server and the frontend.

# UML Diagrams

### Sequence Diagram of Websocket Events

```mermaid
sequenceDiagram
    actor Pl as Player
    participant Man as Manager
    participant Lob as Lobby
    participant Oth as Other Players

    Pl->>Man: WebSocketUpgrade
    Man-->>Pl: ConnectionEstablished

    Pl->>Man: UpdateUser
    Man-->>Pl: UserUpdated

    rect rgba(245, 245, 245, 0.5)
        Note over Pl,Oth: Create Lobby
        Pl->>Man: CreateLobby
        Man-->>Lob: NewLobby
        Lob-->>Man: lobby
        Man->>Pl: LobbyCreated
    end

    rect rgba(245, 245, 245, 0.5)
        Note over Pl,Oth: Join Lobby

        Pl->>Man: JoinLobby

        alt Lobby exists
            Man->>Lob: AddClient()
            Lob-->>Man: success
            Man-->>Pl: LobbyJoined
            Lob->>Oth: PlayerJoined
        else Lobby does not exist
            Man-->>Pl: JoinLobbyFailed
        else Lobby full
            Man-->>Pl: JoinLobbyFailed
        end
    end

    rect rgba(245, 245, 245, 0.5)
        Note over Pl,Oth: Rejoin Lobby

    end

    rect rgba(245, 245, 245, 0.5)
        Note over Pl,Oth: Start Game
    end
```

### Activiy Diagram of Turn Phases

```mermaid
%%{init: {"flowchart": {"markdownAutoWrap": false}}}%%
flowchart TB
    INITIAL -->| First attack card placed | GAME

    GAME["GAME<br/><span style='font-size:12px;white-space:nowrap'>Attackers + Defender place cards asynchronously</span>"]
    GAME -->| Defender takes cards | TAKE
    GAME -->|No valid attacks/ \nattackers PASS| DRAW
    TAKE --> DRAW
    DRAW -->|Everyone draws| COMPLETE
    COMPLETE -.->|Next round| INITIAL
```
