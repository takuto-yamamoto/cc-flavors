# cc-flavors

A scrapbook of Claude Code flavor texts.

## What it does

`cc-flavors` collects "flavor texts" shown by Claude Code (e.g. `Moonwalking...`).

It ships as a single Go binary and a tiny tmux hook.

## Install

### TPM

Add to `.tmux.conf`:

```
set -g @plugin 'takuto-yamamoto/cc-flavors'
run '~/.tmux/plugins/tpm/tpm'
```

Reload tmux, then press `prefix + I` to install.

### Manual

```bash
git clone https://github.com/takuto-yamamoto/cc-flavors.git ~/.tmux/plugins/cc-flavors
```

Add to `.tmux.conf`:

```
run '~/.tmux/plugins/cc-flavors/cc-flavors.tmux'
```

## Usage

### Requirements

- `cc-flavors` binary must be on `PATH`
- `tmux`
- `claude` command (or set `@cc_flavors_cmd`)

### Options

Set in `.tmux.conf` (optional):

```
set -g @cc_flavors_cmd "claude"
set -g @cc_flavors_scan_interval "5"
```

### Export (human readable)

```
cc-flavors export
```

Sample output:

```
    12  Moonwalking
     7  Thinking
     3  Refactoring
```
