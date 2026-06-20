# Purgatorio

To run the backend with file logging, use the command inside of the backend directory
```bash
LOG_FILE="logs/purg.$(date +"%F_%T").log" go run ./cmd/server/main.go
```

To run the frontend, run `npm run dev` inside of the frontend directory

---

## Setup

Run the `setup.sh` file to create the necessary keys for JWT Tokens and environment variables. Then just run `docker compose up`. You can also supply the `LOG_FILE` environment variable to enable JSON-based file logging.

---

## Backend

### Structure

```
backend/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── api/https/              # HTTP handlers + WS
│   │   ├── auth/               # Login, register, token refresh
│   │   ├── battle/             # Battle lifecycle, WebSocket stream
│   │   ├── economy/            # Resource endpoints
│   │   ├── army/               # Troop management
│   │   ├── user/               # Profile, stats
│   │   ├── search/             # Player search
│   │   ├── admin/              # Admin endpoints
│   │   └── middleware/         # Auth middleware (JWT verification)
│   ├── domain/
│   │   ├── service/            # Business logic
│   │   ├── model/              # Domain types
│   │   └── repository/         # Repository interfaces
│   ├── engine/                 # Battle simulation engine
│   └── infrastructure/
│       └── postgres/           # Repository implementations
└── pkg/                        # Shared utilities
    ├── ctxkeys/                # Contains the context keys for objects
    ├── purgerr/                # Multi-errors and stack trace interfaces
    └── response/               # Utility functions for HTTP Responses
```

### WebSocket Protocol (`/api/battles/{id}/ws`)

The battle uses a single WebSocket connection through three phases:

**Phase 1 — Deployment (60s window)**

| Direction | Type | Payload |
|-----------|------|---------|
| Server → Client | `deployment_start` | `{time_left: 60}` |
| Client → Server | `deploy` | `{troops: [{troop_type, position{x,y}, count}]}` |
| Server → Client | `deploy_ack` | `{message: "ok"}` |
| Client → Server | `done` | _(empty)_ |

**Phase 2 — Streaming**

| Direction | Type | Payload |
|-----------|------|---------|
| Server → Client | `tick_batch` | `{ticks: TickResult[], batch_start}` |
| Server → Client | `battle_end` | `{outcome, destruction, loot, sin_meter, duration_ticks}` |
| Client → Server | `skip` | _(empty)_ — ends battle early |
| Server → Client | `error` | `{message}` |

Ticks are streamed in batches of 10 with a 500ms delay between batches. Each `TickResult` contains HP changes, position updates, and whether the simulation has ended.

---

## Frontend

### Structure

```
frontend/
├── src/
│   ├── api/                    # WebSocket client, fetch wrappers
│   ├── app/                    # Providers (Auth, Game), routing
│   ├── components/             # Shared UI components
│   ├── game/
│   │   └── phaser/             # Phaser 3 scenes, sprites, isometric math
│   ├── hooks/                  # React hooks (useBattleSocket, useGame, useAuth)
│   ├── types/                  # Shared TypeScript types
│   └── ui/
│       ├── battle/             # Battle overlay (deployment, playback, HUD)
│       ├── economy/            # Resource views
│       ├── army/               # Army management (train, detrain, barracks)
└       ├── panels/             # Reusable panels (DeploymentScreen, BattleResultScreen)
```

---

##### Sin Meter
- What is it : A percentage bar on your profile, separate from all other resources. Every attack you launch increases the sin meter. Your next attack must achieve at least the sin meter's current percentage in destruction to be considered "successful", otherwise it is marked as "failed to meet threshold" and the loot is not granted.
- Thought behind it : Prevents players from spam attacking and also goes to show that such fighting for power and escape is a sin in itself. It drains over time (you repent and hence it decreases), making the players able to attack again comfortably :)

| Property | Value |
|----------|-------|
| Drain rate | 10% per hour (1% per 6 minutes) |
| On successful attack | Sin meter increases by 10% |
| On failed attack | Sin meter resets to 0% |
| Threshold check | Attack must achieve destruction ≥ current sin % to be a victory |
| Display | Hover tooltip shows time to drain |

The sin meter prevents spam attacking — players must wait for their sin to drain before launching attacks that can succeed.

##### Shield
On a failed defense, the defender gets a 2hr shield to protect and recover the base

##### Battle Mechanism
- The base gets an extra 6 tile padding where the attacker can select and mark troops for deployment
- The player gets a 60s deploy window and then the battle star

---

### Context and Glossary
This game draws heavy inspiration from "Purgatorio" and the basic storyline (which was supposed to be conveyed through cutscenes, but :( nevermind) follows as

```
You were not sent here. You escaped. Minos's judgment was yours to bear in the circles below.
But something in you refused and now you stand at the foot of the mountain, carrying Hell's stain
- on a soul that still wants to be clean.
Virgil found you here. He will guide you upward, though he cannot follow you to the end.
The mountain should be a place of quiet suffering and slow grace.
But others like you discovered what Hell's taint can do, and they turned it on each other.
They are called The Striving. They wanted Paradise so badly they made war of the path to it.
You'll have to go through all of them to get out.
```

##### Buildings
- **Bastion**: alternative name to a wall
- **Sanctum**: The town hall. shoots up a beam of light from the center
- **Angel Spire**: something like archer tower. looks like a high tower with an angel with light on top
- **Barracks**: Self explanatory. looks like a long hut
- **Lament Basin**: the resource collector. contains 'penitent souls' within the chambers to extract out their 'penitence' as a resource

##### Troops
- **Ashwalker**: a burning soul and an all-rounder troop
- **Coveter**: a soul plagued with envy and hence looks drained out. fast but less damaging
- **Hoarder**: a soul plagued with greed, always collecting loot
- **Ravager**: symbolises gluttony and hence is easily identified by the large belly. a tank troop
- **Stone Bearer**: the proud ones. their penitence is carrying massive slabs. has most hp (shiels behind the slab)

##### Resources
- **Penitence**: the primary resource. resource generation as of now is really fast (~2 mins to fill the entire collector)
- **Grace**: like gems. awarded for completing achievements (not in the game for now)

---
