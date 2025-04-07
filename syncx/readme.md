# a simple graph goroutine scheduler

a problem during the interview at Bytedance

## description

give you a graph, for example

```mermaid
flowchart TD
    A --> B
    A --> C
    B --> D
    C-->D
```

the `arrow` meaning the goroutine target end stand for need to be run after the source end's goroutine.

It requires that all the goroutine need to start at once, and be executed accoring to the input graph.
