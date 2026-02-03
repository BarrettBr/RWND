## Ideas for Down the line

- TUI / GUI visual of records
  - Left pane: List of requests ; Right pane: Request / Response & Diff if replayed?
- Dockerfile / Release Files w/CI integration
- Logging:
  - Batch logging
  - Backpressure: If logging is experiencing issues use something similar to CWND in networking to adjust the buffer / speed or alert
  - Multiple workers: Keep one writer but multiple preprocessors
- Invariant Rules: Be able to set flags and if they show we list them under a seperate file or something i.e: Response w/401 unauthorized gets logged to seperate file
  - How to make look clean and not _bloated_?
  - Tag map? / Side-channel gets different alerts / YAML file or rules?
- Replay Engine:
  - Editing / Replay to backend & comparison
  - Mutators for Headers / Body / etc, edit them replay it and side by side compare responses

## Performance & General Robustness Ideas

- Body capture: Size limit / content-type allow or blacklist (Prevent zip / image)
  - Skip GET bodies
  - `ReadAll` -> `io.LimitReader()` & `rec.Request.Truncated = true`
- Allow for selection of sampling vs all logging and maybe dynamically choose this for when at load?
- Limit headers we get and allow capturing of all via an option
  - Host, User-Agent, Content-Type / Length, Authorization and skip rest?
- Optimize logging writes:
  - Swap from JSON -> gob/other binary format
- Double check I/O fmt.Println on hot paths to ensure we aren't spamming the console
- Make capturing an option so it passes through otherwise so you can _flick on the switch_ when you want not shutting down the whole app whenever
- Double check proxy settings:
  - Timeouts / idle conns allowed / etc
