# 001 tmux scan only

## Context
We want tmux-based collection without wrappers. We accept periodic scanning instead of focus events.

## Decision
- Start a scan loop on `client-attached`.
- Every 5 seconds, list all panes and detect `#{pane_current_command}` matching a configurable command name (default `claude`).
- If a pane matches, attach `pipe-pane -o` to stream output.
- If a pane does not match, detach any existing `pipe-pane`.
- Do not use `pane-focus-in`.

## Reason
- A 5-second loop is fast enough and keeps the hook logic simple.
- `-o` avoids clobbering other `pipe-pane` usage.
- A configurable command name supports different CLI entry points.
