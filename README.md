# Dank Grinder [![](https://goreportcard.com/badge/github.com/dankgrinder/dankgrinder)](https://goreportcard.com/report/github.com/dankgrinder/dankgrinder) [![](https://img.shields.io/github/workflow/status/dankgrinder/dankgrinder/Go)](https://github.com/dankgrinder/dankgrinder/actions) ![](https://img.shields.io/github/license/dankgrinder/dankgrinder) [![](https://img.shields.io/github/v/release/dankgrinder/dankgrinder)](https://github.com/dankgrinder/dankgrinder/releases/latest) ![](https://img.shields.io/github/downloads/dankgrinder/dankgrinder/total)
The most feature-rich, advanced Dank Memer automation tool.

Want to join the community or need support? [Join our Discord](https://discord.gg/Fjzpq8YPSn). Have a question? [Create a question issue](https://github.com/dankgrinder/dankgrinder/issues/new?assignees=&labels=question&template=question.md&title=). Want to suggest a feature? [Create a suggestion issue](https://github.com/dankgrinder/dankgrinder/issues/new?assignees=&labels=suggestion&template=suggestion.md&title=). Encountered a bug? [Report a bug](https://github.com/dankgrinder/dankgrinder/issues/new?assignees=&labels=bug&template=bug-report.md&title=). Want to contribute? [Read our contribution guidelines](https://github.com/dankgrinder/dankgrinder/blob/master/CONTRIBUTING.md).

<p align="center">
<img src="https://i.imgur.com/3AQk7eh.png" alt="logo" />
</p>

## Features
* Can run many instances at once
* Uses very little system resources
* High configurability; configure custom commands and much more
* Automation of most currency commands and blackjack
* Responds to fishing and hunting events (i.e. captures the dragon and special fish)
* Responds to global events, search, and postmeme
* Automatically uses tidepods and buys lifesavers when dying from them
* Automatically buys a fishing pole, hunting rifle or laptop if they have broken

## Getting started

### Installation
1. Download the latest release for your operating system [here](https://github.com/dankgrinder/dankgrinder/releases/latest) (darwin is for macOS), or [build from source](#building-from-source). If you build from source you can delete everything besides the compiled binaries and config.yml.
2. Extract the archive
3. [Configure](#configuration). Make sure `token` and `channel_id` fields of the instance are set, it will not run otherwise

#### Windows
4. Double click dankgrinder.exe. If the program closes immediately after opening it, open a command prompt and drag the executable to it, then hit enter. You will now be able to see the error

#### Linux
4. Run the binary:  
   `$ ./dankgrinder`

#### macOS
4. Run by double clicking the dankgrinder binary

### Getting an authorization token
1. Go to Discord, logged into the account you want to use
2. Open the inspector by pressing `ctrl` + `shift` + `i`
3. Click the "network" tab
4. Refresh the page
5. Look for any entry called "science" and click on it
6. Expand the "request headers" and look for the value "authorization", that is your token

### Enabling Discord's developer mode
To obtain a channel id for config.yml, you will need to have developer mode enabled.
1. Go to your user settings on Discord
2. Click "appearance" under "app settings"
3. Scroll down to "advanced" and enable developer mode

You can now right click any user or channel to obtain the id.

## Building from source
If you use an architecture other than amd64, or you want to build from source for another reason, doing so is quite simple.

### Prerequisites
A working Go installation, see https://golang.org/doc/install.

### Building
1. Clone the repository:  
   `$ git clone https://github.com/dankgrinder/dankgrinder.git`
2. Build:  
   `$ make build`

The executables should then be in the `/build` folder.

## Configuration
All configuration can be done by editing config.yml with your editor of choice (e.g. Visual Studio Code, Sublime Text, Notepad++). The comments in the file itself provide extra explanation where necessary. For the bot to run, you must at least enter an [authorization token](#getting-an-authorization-token) and a [channel id](#enabling-discords-developer-mode). If you are running into issues, you can [join our Discord](https://discord.gg/Fjzpq8YPSn).

If you do not know how yaml works and are getting fatal errors, use [this guide](https://www.tutorialspoint.com/yaml/yaml_basics.htm) to learn the basics of yaml. Configuration errors are characterized by a near-instant fatal error when starting the program. If the program opens and then closes immediately on Windows, open a command prompt first, drag the executable onto it and hit enter. You should now be able to see the error.

A question mark after a field name means this field is optional.

Name | Type | Description
---- | ---- | ----
`instances` | array of [instance objects](#instance-object) | The instances that are run simultaneously by the program
`shifts` | array of [shift objects](#shift-object) | One or more shifts which the instances use to schedule switching between the active and dormant state. [Read more about shifts](#shifts)
`features` | [features object](#features-object) | Several feature configurations which apply to all instances
`compatibility` | [compatibility object](#compatibility-object) | Several compatibility options which apply to all instances
`suspicion_avoidance` | [suspicion avoidance object](#suspicion-avoidance-object) | Several techniques to avoid suspicion which apply to all instances

### Instance object
Name | Type | Description
---- | ---- | ----
`token` | string | The Discord [authorization token](#getting-an-authorization-token) of the instance
`channel_id` | string | The channel id this instance sends and receives messages in, you must have [Discord developer mode](#enabling-discords-developer-mode) enabled to obtain one
`is_master` | bool |
`features?` | [features object](#features-object) | Override the default features object of the config only for this specific instance, any fields left out will not be overridden and vice-versa, see [default values and when you can leave out fields](#default-values-and-when-you-can-leave-out-fields)
`suspicion_avoidance?` | [suspicion avoidance object](#suspicion-avoidance-object) | Override the default suspicion avoidance object of the config only for this specific instance, any fields left out will not be overridden and vice-versa, see [default values and when you can leave out fields](#default-values-and-when-you-can-leave-out-fields)
`shifts?` | array of [shift objects](#shift-object) | Override the default shifts array of the config only for this specific instance, see [default values and when you can leave out fields](#default-values-and-when-you-can-leave-out-fields)

### Shift object
Name | Type | Description
---- | ---- | ----
`state` | string | The state of the program for this shift, either `active` or `dormant`
`duration.base` | integer | The base duration of this shift in seconds. [Read more about base and variation](#base-and-variation)
`duration.variation` | integer | The random variation of this shift in seconds. [Read more about base and variation](#base-and-variation)

### Features object
Name | Type | Description
---- | ---- | ----
`commands` | [commands object](#commands-object) | Enable or disable certain commands
`custom_commands` | array of [custom command object](#custom-command-object) | Configure your own, custom commands for the program to use
`auto_buy` | [auto-buy object](#auto-buy-object) | Options for the automatic buying of certain items if it is detected that they are not available
`auto_sell` | [auto-sell object](#auto-sell-object) | Options for the automatic, periodic selling of certain items
`auto_gift` | [auto-gift object](#auto-gift-object) | Options for the automatic, periodic gifting of certain items to the master instance
`auto_blackjack` | [auto-blackjack object](#auto-blackjack-object) | Options for automatically using the blackjack command
`auto_share` | [auto-share object](#auto-share-object) | Options for automatically sharing money with the master instance
`auto_tidepod` |  [auto-tidepod](#auto-tidepod-object) | Options for automatically using tidepods
`balance_check` | [balance check object](#balance-check-object) | Options for checking balance
`log_to_file` | boolean | Whether or not to log errors and information to a file. If running multiple instances a log file for each instance will always be created, regardless of the value of this property
`debug` | boolean | Enable logging debug level information. Currently has no effect

### Commands object
Name | Type | Description
---- | ---- | ----
`beg` | boolean | Enable the `pls beg` command
`postmeme` | boolean | Enable the `pls postmeme` command
`search` | boolean | Enable the `pls search` command
`highlow` | boolean | Enable the `pls highlow` command
`fish` | boolean | Enable the `pls fish` command
`hunt` | boolean | Enable the `pls hunt` command

### Custom command object
Name | Type | Description
---- | ---- | ----
`value` | string | The value of the command, for example: `pls dep max`
`interval` | integer | The interval at which this command will be re-sent in seconds. Time may vary depending on other commands and responses. If `0` the command will only run once in the beginning of every active shift
`amount` | integer | The amount of times this command will be run in total every active shift. Set to `0` for no limit
`pause_below_balance` | integer | A wallet balance value below which this command will not be sent. The balance is read from the balance check functionality. Consider having the interval of this quite low, to make sure the balance the program thinks you have is as up-to-date as possible

### Auto-buy object
Name | Type | Description
---- | ---- | ----
`fishing_pole` | boolean | Enable the automatic purchase of a fishing pole when it is detected that one is not available
`hunting_rifle` | boolean | Enable the automatic purchase of a hunting rifle when it is detected that one is not available
`laptop` | boolean | Enable the automatic purchase of a laptop when it is detected that one is not available

### Auto-sell object
Name | Type | Description
---- | ---- | ----
`enable` | boolean | Whether or not to enable automatic selling
`interval` | integer | The interval at which items will be sold during an active shift. If set to 0, items will only be sold once at the beginning of every active shift
`items` | array of strings | The Dank Memer item ids of the items to sell

### Auto-gift object
Name | Type | Description
---- | ---- | ----
`enable` | boolean | Whether or not to enable automatic gifting to the master instance
`interval` | integer | The interval at which items will be gifted during an active shift. If set to 0, items will only be gifted once at the beginning of every active shift
`items` | array of strings | The Dank Memer item ids of the items to gift

### Auto-blackjack object
Name | Type | Description
---- | ---- | ----
`enable` | boolean | Whether or not to enable automatic blackjack
`priority` | boolean | Whether or not to give the command priority over other, regular commands if there are commands queued
`amount` | integer | The amount to bet every time, set to `0` to bet the maximum amount of coins
`pause_below_balance` | integer | The balance below which the program should stop betting. The balance is read from the balance check functionality. Consider having the interval of this quite low, to make sure the balance the program thinks you have is as up-to-date as possible
`logic_table` | dictionary[string]dictionary[string]string | What to do for every possible blackjack hand. The string values are the exact response that will be triggered

### Auto-share object
Name | Type | Description
---- | ---- | ----
`enable` | boolean | Whether or not to enable automatically giving money to the master 
`fund` | boolean | Whether or not instances can request money from the master
`maximum_balance` | integer | The amount of money the instance may have before giving them to the master instance
`minimum_balance` | integer | The amount of money the instance should keep after giving money to the master instance and the amount the master will fund it to if an instance requests it

### Auto-tidepod object
Name | Type | Description
---- | ---- | ----
`enable` | boolean | Whether or not to enable automatic usage of tidepods
`interval` | integer | The interval in seconds at which the program attempts to use tidepods. If set to 0, a tidepod will only be used once at the beginning of every active shift
`buy_lifesaver_on_death` | bool | Whether or not to buy a lifesaver after dying from tidepod usage

### Balance check object
Name | Type | Description
---- | ---- | ----
`enable` | boolean | Whether or not to enable balance checks
`interval` | integer | The interval in seconds at which the program checks the balance. This is the same as the auto-share interval, if enabled


### Compatibility object
Name | Type | Description
---- | ---- | ----
`postmeme` | array of strings | What options can be chosen for the postmeme command. The program will pick one randomly
`allowed_searches` | array of strings | The searches the application is allowed to pick. Items higher/earlier in the list have higher priority
`search_cancel` | array of strings | List of things the program will say to cancel a search when no allowed searches are provided. It will pick one randomly
`cooldown` | [cooldown object](#cooldown-object) | Cooldowns of commands (not custom commands)
`await_response_timeout` | integer | The time that the program will wait for a response when it is expecting one. Set to a higher value when Dank Memer is slow to respond and this causes issues. Values below `3` are not recommended

### Cooldown object
Name | Type | Description
---- | ---- | ----
`beg` | integer | The cooldown of the beg command in seconds, set a few seconds higher to account for network delay
`search` | integer | The cooldown of the search command in seconds, set a few seconds higher to account for network delay
`highlow` | integer | The cooldown of the highlow command in seconds, set a few seconds higher to account for network delay
`postmeme` | integer | The cooldown of the postmeme command in seconds, set a few seconds higher to account for network delay
`fish` | integer | The cooldown of the fish command in seconds, set a few seconds higher to account for network delay
`hunt` | integer | The cooldown of the hunt command in seconds, set a few seconds higher to account for network delay
`blackjack` | integer | The cooldown of the blackjack command in seconds, set a few seconds higher to account for network delay
`sell` | integer | The cooldown of the sell command in seconds, set a few seconds higher to account for network delay
`gift` | integer | The cooldown of the gift command in seconds, set a few seconds higher to account for network delay

### Suspicion avoidance object
Name | Type | Description
---- | ---- | ----
`typing` | [typing object](#typing-object) | Options for the use of typing when sending messages
`message_delay` | [message delay object](#message-delay-object) | Delay between receiving a message and starting to type and send a response

### Typing object
Name | Type | Description
---- | ---- | ----
`base` | integer | The base duration of typing in milliseconds. [Read more about base and variation](#base-and-variation)
`variation` | integer | The random variation of typing in milliseconds. [Read more about base and variation](#base-and-variation)
`speed` | integer | The typing speed based on message length, in characters per minute

### Message delay object
Name | Type | Description
---- | ---- | ----
`base` | integer | The base delay in milliseconds. [Read more about base and variation](#base-and-variation)
`variation` | integer | The random variation of the delay in milliseconds. [Read more about base and variation](#base-and-variation)

### Base and variation
A base is a value that forms the base for a final result; it is the value that the program starts with.

A variation is a random value, from 0 up to but excluding n, added to a starting value.

Say a base value of 100, and a variation of 50 are used. The final result will be a number from 100 up to and including 149 because of the random value added by the variation.

### Shifts
You can use shifts to make sure the bot is not suspicious because it is too active. We highly recommended using this option, and not running this program 24/7. An uptime of 50% per instance or less is advisable.

For example, if you would like to have the bot be active for 6 hours per day, then inactive for the remaining 18 you would do something like this (note that it will loop this, so after completing the last shift it will go back to the first):
```yaml
shifts:
  - state: "active"
    duration:
      base: 21600
      variation: 0
  - state: "dormant"
    duration:
      base: 64800
      variation: 0
```

### Custom commands
```yaml
custom_commands:
  - value: "pls command1"
    interval: 60
  - value: "pls command2"
    interval: 300
    amount: 5
  - value: "pls command3"
  - value: "pls buy zz 20"
    pause_below_balance: 9000000
```
In example custom command 1, the value is sent every 60 seconds for infinite times, until the program enters the dormant state.

In example custom command 2, the value is sent every 5 minutes for a total of 5 times per active shift.

In example custom command 3, the value is sent once in the beginning of every active shift.

In example custom command 4, 20 zz will be bought whenever the balance is above 9,000,000.

### Instances
Example if you would like to run two instances simultaneously and 24/7 (this shift configuration is not recommended):
```yaml
instances:
  - token: "bmljZSB0cnkgYnV0IHRoaXMgaXM.bm90IGE.cmVhbCB0b2tlbg"
    channel_id: "791694339116892202"
    is_master: true
  - token: "b2YgY291cnNlIHRoaXM.aXNuJ3QgYQ.cmVhbCB0b2tlbiBlaXRoZXIsIHNpbGx5"
    channel_id: "791694383098495047"
```

### Default values and when you can leave out fields
Because of the way Go structs work, most times, when you leave out a field in your config, it will default to a value such as `0` or `false`. This is useful to avoid clutter. In the default config many fields are left out because their values would have been set to `0` or `false`, but they are still available for use, of course. Simply lookup what configuration is possible for your object of choice in the [configuration documentation](#configuration) above.

The only exception is the fields which are currently marked as optional, the `features`, `suspicion_avoidance` and `shifts` fields on every [instance object](#instance-object). These fields are used to override the values you have specified in the regular `features`, `suspicion_avoidance` and `shifts` objects in the config. If you leave out one of these fields or a child field of one of these fields, the default configuration is not overridden and is used instead.

```yaml
instances:
  - token: "bmljZSB0cnkgYnV0IHRoaXMgaXM.bm90IGE.cmVhbCB0b2tlbg"
    channel_id: "791694339116892202"
    is_master: true
  - token: "b2YgY291cnNlIHRoaXM.aXNuJ3QgYQ.cmVhbCB0b2tlbiBlaXRoZXIsIHNpbGx5"
    channel_id: "791694383098495047"
    shifts:
      - state: "active"
  - token: "MTI3OTgzNDcyMTkzNDM4Mg.ZmRzdg.dGhpcyBpcyBub3QgYW4gYWN0dWFsIH"
    channel_id: "791691923098486933"
    features:
      auto_tidepod:
        enable: false

shifts:
  - state: "active"
    base: 21600
  - state: "dormant"
    base: 32400
```

### Commonly asked questions
## How do I disable certain features?
You can remove a specific feature by removing it from the config, setting its `enable` value to `false`.

## How do I disable global events?
This is not possible as of now.

## Can I run it on other architectures?
As mentioned above, it is possible to run this program on other architecures than amd64. However, to do this, you will need to [build from source](#building-from-source). This is not difficult to do.

In the example above, all instances use a shift configuration of 6 hours active, 9 hours dormant, except for the second instance. This instance overrides the shift configuration defined below with a shift configuration that is always active. 

The third instance in this example, overrides the normal configuration with one where auto-tidepod is disabled (the rest of the config is left out for simplicity), the rest of the instances would still use auto-tidepod if it was enabled. 

## Disclaimer
This is a self-bot. Such bots are against Discord's terms of service. Automation of Dank Memer commands also breaks Dank Memer's rules. By using this software you acknowledge that we take no responsibility whatsoever for any action taken against your account, whether by Discord or Dank Memer, for not abiding by their respective rules.

Despite this, we believe the chance of detection by either Discord or Dank Memer to be low provided that you take appropriate measures to ensure this. This includes but is not limited to running the bot only in private channels, not being open about the fact that you use it and not running so much as to raise suspicion. 
