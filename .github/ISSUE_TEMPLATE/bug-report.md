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
Config: e.g.
```yaml
# See https://github.com/dankgrinder/dankgrinder#getting-an-authorization-token
# for instructions on how to get a token.
token: ""

# See https://github.com/dankgrinder/dankgrinder#enabling-discords-developer-mode
# for instructions on how to get a channel ID.
channel_id: ""

features:
  commands:
    fish: true
    hunt: true
  auto_buy:
    fishing_pole: true
    hunting_rifle: true
    laptop: true
  balance_check: true

compatibility:
  postmeme:
    - "f"
    - "r"
    - "i"
    - "c"
    - "k"

  # Search options the bot will pick if provided to it. It will prefer options
  # which are higher in the list.
  search:
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

  # Cooldown values in seconds for each command. Default is non-donator.
  cooldown:
    beg: 45
    search: 30
    highlow: 20
    postmeme: 60
    fish: 60
    hunt: 60

    # This value in seconds in added to each cooldown to accommodate timing errors.
    margin: 3


suspicion_avoidance:
  typing:
    # This is a base typing duration in milliseconds that is added to the TOTAL
    # typing time.
    base: 0

    # The typing speed in keys per minute, approx. 350 is fast. Can be set to 0
    # to use only the base and variance.
    speed: 450

    # A random value between 0 and this will be added to the TOTAL typing time
    # in milliseconds.
    variance: 250

  # The delay between deciding to send a command and starting to type. Base and
  # variance work just like they do in typing.
  message_delay:
    base: 100
    variance: 400

```

**Additional comments**
Any other information that might be helpful but does not fit in the above categories.
