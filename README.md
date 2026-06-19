# Purgatorio

To run the backend with file logging, use the command inside of the backend directory
```bash
LOG_FILE="logs/purg.$(date +"%F_%T").log" go run ./cmd/server/main.go
```

To run the frontend, run `npm run dev` inside of the frontend directory

---

## Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Backend | Go 1.26, Chi router | HTTP API + WebSocket server |
| Database | PostgreSQL 16 | Persistence |
| Frontend | React 19, TypeScript, Vite | SPA UI |
| Game Rendering | Phaser 3 (custom build) | Isometric battle canvas |
| Auth | JWT (RS256), HTTP-only cookies | Stateless auth |
| Containerization | Docker Compose | Dev environment |

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

### Key Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/auth/register` | Create account |
| `POST` | `/api/auth/login` | Authenticate, set cookie |
| `POST` | `/api/auth/refresh` | Rotate refresh token |
| `GET` | `/api/battles/match` | Find an opponent to attack |
| `GET` | `/api/battles/{id}` | Get battle state |
| `GET` | `/api/battles/{id}/ws` | WebSocket — deployment + tick streaming |
| `GET` | `/api/battles/{id}/replay` | Fetch reconstructed replay ticks |
| `POST` | `/api/battles/protect` | Set shield for your own base |
| `GET` | `/api/economy/me` | Get own resources |
| `GET` | `/api/army/me` | Get own troop inventory |
| `POST` | `/api/army/train` | Train troops (costs penitence) |
| `POST` | `/api/army/detrain` | Refund 50% of troops back to penitence |
| `GET` | `/api/users/me` | Own profile + sin meter + shield |
| `GET` | `/api/users/search?q=` | Search players by username |
| `GET` | `/api/users/{id}` | Other player's profile (economy + stats) |

### Database Schema (Key Tables)

| Table | Purpose |
|-------|---------|
| `auth` | Auth |
| `refresh_tokens` | JWT refresh token store |
| `users` | Profile |
| `user_economy` | Penitence, collector pending |
| `user_combat` | Sin meter, last attack, shield |
| `user_army` | Troop counts per type |
| `user_stats` | Attack/defense win/loss counters |
| `base_layouts` | User placed buildings aka layout |
| `base_snapshots` | Defender building snapshot at battle time |
| `buildings` | Contains information for buildings owned by user |
| `buildings` | Contains shop information for buildings |
| `building_levels` | Contains information for upgrades |
| `building_limits` | Contains limits on how many buildings per terrace level |
| `troops` | Contains information about troops |
| `battles` | Battle records (outcome, destruction, loot) |
| `battle_replays` | Replay data (seed + deployments — not tick data) |

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

### Key Frontend Flows

| Flow | Key Files | Description |
|------|-----------|-------------|
| Battle Lifecycle | `BattleOverlay.tsx`, `useBattleSocket.ts`, `BattleSocket` (ws.ts) | Deployment → playback → result |
| Phaser Canvas | `BattleCanvas.tsx`, `BattleScene.ts` | Isometric rendering of buildings + troops |
| Army | `BarracksPanel.tsx` | Train, detrain, capacity management |
| Economy | `EconomyPanel.tsx` | Resource display with real-time sin drain |

---

## Implementation Details

### Deterministic Battle Simulation

The battle engine runs the **entire simulation upfront** on the server i.e. every tick is computed before the first byte is sent to the client. This is possible because the engine is fully deterministic: given the same seed, troop deployments, and defender building snapshot, it always produces identical results.

```
seed + deployments + building_snapshot  →  BattleInput  →  Simulation  →  []TickResult
```

This enables two key features:
- **Replay by re-simulation**: the database stores only the input parameters (seed, deployments, snapshot ID), not the tick data. Replays reconstruct the battle by creating a new simulation with the same inputs.
- **Instant skip**: the full result is already computed, so early termination just trims the output.

### A* Pathfinding with Building Avoidance

Troops navigate using A* on a 4-directional grid. The pathfinder treats building footprints as impassable and generates waypoints to cells adjacent to the target building. Paths are recalculated when targets change (buildings destroyed, re-targeting). Fallback: if no path exists, troops move directly toward the target's edge.

### Sin Meter

| Property | Value |
|----------|-------|
| Drain rate | 10% per hour (1% per 6 minutes) |
| On successful attack | Sin meter increases by 10% |
| On failed attack | Sin meter resets to 0% |
| Threshold check | Attack must achieve destruction ≥ current sin % to be a victory |
| Frontend sync | Recalculated on fetch, local tick every 10s matches server rate |
| Display | Hover tooltip shows time to drain |

The sin meter prevents spam attacking — players must wait for their sin to drain before launching attacks that can succeed.

### WebSocket Graceful Skip

The skip mechanism uses a goroutine that reads WebSocket messages concurrently with the tick-streaming loop. When a `skip` arrives:
1. The streaming loop breaks (either before the next batch write or during the 500ms inter-batch delay)
2. The current batch index is recorded as `endTick`
3. `ResolveAndStore` runs with the truncated simulation result
4. `battle_end` is sent immediately

This means the result screen appears in one network round-trip instead of waiting for all remaining batches to stream.

---

## Setup

Run the `setup.sh` file to create the necessary keys for JWT Tokens and environment variables. Then just run `docker compose up`. You can also supply the `LOG_FILE` environment variable to enable JSON-based file logging.

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

### Mechanisms
While the game draws most of its inspiration from CoC, it still does some things differently such as the Sin Meter and deploying mechanism

##### Sin Meter
- What is it : A percentage bar on your profile, separate from all other resources. Every attack you launch increases the sin meter. Your next attack must achieve at least the sin meter's current percentage in destruction to be considered "successful", otherwise it is marked as "failed to meet threshold" and the loot is not granted.
- Thought behind it : Prevents players from spam attacking and also goes to show that such fighting for power and escape is a sin in itself. It drains over time (you repent and hence it decreases), making the players able to attack again comfortably :)

##### Shield
On a failed defense, the defender gets a 2hr shield to protect and recover the base

##### Battle Mechanism
- The base gets an extra 6 tile padding where the attacker can select and mark troops for deployment
- The player gets a 60s deploy window and then the battle starts

