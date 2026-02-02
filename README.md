# Claudia

Claudia is a custom Waybar module that provides real-time visibility into
claude-code activity. It indicates whether claude-code is currently Batch,
Edit, Read, Write, Glob, Grep, Task, Web, or Idle.

## TODO

- Support multiple concurrent sessions.
- Start the listener when at least one Claude session is active.
- Stop the listener when all Claude sessions have exited.


## Usage

Add (or merge) into `~/.config/waybar/config`:
```json
{
	"modules-right": ["custom/claudia", "..."],
	"custom/claudia": {
		"exec": "~/.config/waybar/scripts/claudia",
		"return-type": "json",
        "exec-on-event": true
	}
}
```

Add (or merge) into `~/.claude/settings.json`
```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "socat - UNIX-CONNECT:/tmp/claudia.sock"
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "socat - UNIX-CONNECT:/tmp/claudia.sock"
          }
        ]
      }
    ],
    "Notification": [
      {
        "matcher": "idle_prompt",
        "hooks": [
          {
            "type": "command",
            "command": "socat - UNIX-CONNECT:/tmp/claudia.sock"
          }
        ]
      }
    ]
  }
}
```
