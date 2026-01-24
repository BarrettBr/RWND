## Ideas for Down the line

- TUI / GUI visual of records
- Logging:
  - Batch logging
  - Backpressure: If logging is experiencing issues use something similar to CWND in networking to adjust the buffer / speed or alert
  - Multiple workers: Increase write throughput by parallelizing work (Better for SQLite later)
- Graceful shutdown / flushing of buffer: When stopping the program flush the buffers logs either to disk / log so to prevent dropping
  - `Close()` on Logger / File and call upon shutdown
