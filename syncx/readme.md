# a simple graph goroutine scheduler

a problem during the interview at Bytedance

## description

give you a graph, for example

```mermaid
flowchart TD
    A --> B
    A --> C
    B --> D
    C --> D
```

the `arrow` meaning the goroutine target end stands for is run after the source end's goroutine finished.

It requires that all the goroutine start at once, and be executed accoring to the input graph.
