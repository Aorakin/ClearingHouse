```mermaid
---
title: GET / - Welcome
---
sequenceDiagram
    actor User
    participant ClearingHouse

    User->>ClearingHouse: GET /
    ClearingHouse-->>User: welcome response
```
