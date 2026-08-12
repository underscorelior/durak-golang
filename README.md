# Durak Multiplayer

This is an online multiplayer port of the game [Durak](https://en.wikipedia.org/wiki/Durak).

Built on Golang using Gorilla Websockets to communicate between the server and the frontend.

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
```
