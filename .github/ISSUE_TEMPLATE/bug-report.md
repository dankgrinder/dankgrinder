---
name: Bug report
about: Report an issue or a bug
title: ''
labels: bug
assignees: ''

---

**Expected behavior**
The behavior you expected

**Actual behavior**
What actually happened

**Steps to reproduce**
How the problem can be reproduced from a clean installation

**Logs**
If relevant, paste any errors or other logs here.

**Environment**
OS: e.g. MacOS 10.14.2, Windows, Ubuntu 20.04
Dank Grinder version: e.g. 1.4.12
Config: make sure to remove your tokens! E.g.
```yaml
instances:
  - token: ""
    channel_id: ""
    shifts:
      - state: "active"
        duration:
          base: 0
          variance: 0

features:
  commands:
    fish: true
    hunt: true
  custom_commands:
  auto_buy:
    fishing_pole: true
    hunting_rifle: true
    laptop: true
  auto_sell:
    enable: false
    interval: 0
    items:
      - "boar"
      - "dragon"
      - "duck"
      - "fish"
      - "deer"
      - "rabbit"
      - "skunk"
  auto_gift:
    enable: false
    to: ""
    interval: 0
    items:
      - "kn"
      - "zz"
  balance_check: true
  log_to_file: true
  debug: false

compatibility:
  postmeme:
    - "f"
    - "r"
    - "i"
    - "c"
    - "k"
  allowed_searches:
    - "bus"
    - "coat"
    - "dresser"
    - "grass"
    - "laundromat"
    - "mailbox"
    - "pantry"
    - "pocket"
    - "shoe"
    - "sink"
    - "car"
  search_cancel:
    - "no"
    - "bad options"
    - "i don't want to die"
    - "stop"
    - "."
  cooldown:
    beg: 45
    search: 30
    highlow: 20
    postmeme: 60
    fish: 60
    hunt: 60
    margin: 3
  await_response_timeout: 4

suspicion_avoidance:
  typing:
    base: 0
    variance: 250
    speed: 450
  message_delay:
    base: 100
    variance: 400
```
