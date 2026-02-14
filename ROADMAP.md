# Backlog App Roadmap

Status: draft

This document outlines the phased plan for a self-hosted, local-only backlog app.
The MVP focuses on games (PC), with a retro 2000s HTMX UI, and a clear path to
additional platforms and media types later.

## Goals
- Local-only, self-hosted, multi-user from day one.
- Simple deployment: single binary + SQLite.
- Import libraries (Steam first) and enrich with third-party metadata (IGDB).
- Manage backlog with statuses, tags, ratings, and notes.
- Provide lightweight recommendations and a read-only chatbot.

## Locked Decisions (MVP)
- Frontend: Go templates + HTMX + minimal JS.
- Metadata: IGDB (OAuth) as primary source.
- Auth: email/password + Steam OpenID.
- Admin: first user becomes admin.
- Manual entry: supported in MVP.
- Chatbot: rules/SQL only (no LLM in MVP).

## MVP Scope
- Authentication (email/password + Steam).
- Library management (manual add + import).
- Backlog statuses: Backlog, Playing, Completed, Dropped, Wishlist, OnHold.
- Filters/search, sorting, tags, ratings, notes.
- Recommendations based on recent play + ratings + genre similarity.
- Chatbot with limited, read-only questions against user data.
- Retro 2000s UI (tables, panels, bevels, fixed width layout).

## Architecture (MVP)
- Go monolith serving HTML + JSON endpoints.
- SQLite database with migrations.
- Background job runner for imports and enrichment.
- Single binary with embedded assets.

## Data Model (Core Tables)
- users
- sessions
- auth_identities
- games
- game_platforms
- game_genres
- external_ids
- library_entries
- tags
- library_entry_tags
- import_jobs
- import_job_items

## Integrations
- Steam Web API: owned games, playtime.
- IGDB: metadata enrichment (title, cover, genres, summary, platforms).
- Optional later: SteamGridDB for artwork.

## Phases

### Phase 0: Product Scope
- Finalize MVP feature list and constraints.
- Define retro UI language and layout rules.

### Phase 1: Architecture and Schema
- Set up Go project, router, template base, asset pipeline.
- Implement SQLite schema + migrations.
- Establish environment config and secrets handling.

### Phase 2: Authentication
- Email/password auth with Argon2id.
- Session cookies.
- Steam OpenID login and user linking.

### Phase 3: Metadata (IGDB)
- OAuth token handling.
- Search and fetch game metadata.
- Cache and normalize metadata locally.

### Phase 4: Steam Import
- Import job pipeline with progress tracking.
- AppID mapping to internal games + IGDB metadata.
- Store playtime and last played data where available.

### Phase 5: HTMX UI
- Retro layout components: tables, panels, banners.
- Library list view with filters + status updates.
- Game details page with metadata and user notes.
- Import UI and status pages.

### Phase 6: Recommendations
- Basic algorithm: recent play + liked + genre similarity.
- Exclude completed/dropped, prefer unplayed/backlog.

### Phase 7: Chatbot (Rules/SQL)
- Intent mapping for common queries.
- Parameterized SQL against read-only views.
- Strict user scoping.

### Phase 8: Packaging
- Single binary + embedded assets.
- Dockerfile (optional).
- Basic docs for self-hosting.

## Out of Scope for MVP
- Federation / ActivityPub.
- Non-game media (books, movies, series).
- Mobile app or native desktop client.
- LLM-based chatbot.

## Future Directions
- Federation between instances (ActivityPub or similar).
- Additional platforms (Epic, GOG, Xbox, PlayStation, Switch).
- Media expansion (books, movies, music, shows).
- Advanced recommendations (collaborative filtering).
- Optional LLM-powered chatbot and semantic search.
