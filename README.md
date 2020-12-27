# Dank Grinder
The most feature-rich, advanced Dank Memer automation tool.

![](https://img.shields.io/github/last-commit/dankgrinder/dankgrinder) ![](https://img.shields.io/github/v/release/dankgrinder/dankgrinder) ![](https://img.shields.io/github/workflow/status/dankgrinder/dankgrinder/Go)

Need support or have a suggestion? [Join our Discord](https://discord.gg/Fjzpq8YPSn). Encountered a bug? [Report a bug](https://github.com/dankgrinder/dankgrinder/issues/new?assignees=&labels=bug&template=bug-report.md&title=)

![](https://i.imgur.com/IHXrnxC.png)
## Features
* Cycles through currency commands to earn approximately 150,000 coins per hour
* Responds to fishing and hunting events (i.e. captures the dragon, etc.)
* Responds to global events, search, and  postmeme
* Automation of highlow
* Reports an average income every cycle
* Automatically buys a fishing pole, hunting rifle or laptop if they have broken
* High configurability

## Getting started

### Installation
1. Download the latest release for your operating system [here](https://github.com/dankgrinder/dankgrinder/releases/latest), or [build from source](#building-from-source). If you build from source you can delete everything besides the compiled binaries and config.json.

#### Windows and MacOS
2. Extract the archive
3. [Configure](#configuration), make sure `token` and `channel_id` are set, it will not run otherwise
4. Double click dankgrinder.exe for windows or the dankgrinder binary for MacOS to run

#### Linux
2. Extract the tarball:  
   `$ tar -xzf dankgrinder-0.1.0-linux-amd64.tar.gz`
3. [Configure](#configuration), make sure `token` and `channel_id` are set, it will not run otherwise
4. Run the binary:  
   `$ ./dankgrinder`
   
You can use the candy binary (or candy.exe on Windows) to automatically consume all candy. Simply double click it. Running this also requires `token`, `channel_id` and `user_id` to be set in config.json.
   
### Getting an authorization token
1. Go to Discord, logged into the account you want to use
2. Open the inspector by pressing `ctrl` + `shift` + `i`
3. Click the "network" tab
4. Refresh the page
5. Look for any entry called "science" and click on it
6. Expand the "request headers" and look for the value "authorization", that is your token

### Enabling Discord's developer mode
To obtain a channel or user ID you will need to have developer mode enabled.
1. Go to your user settings on Discord
2. Click "appearance" under "app settings"
3. Scroll down to "advanced" and enable developer mode

You can now right click any user or channel to obtain the ID.

## Building from source
If you use an architecture other than amd64, or you want to build from source for another reason, doing so is quite simple.

### Prerequisites
A working Go installation, see https://golang.org/doc/install

### Building
1. Clone the repository:  
   `$ git clone https://github.com/dankgrinder/dankgrinder.git`
2. Build:  
   `$ make build`
   
An executable should then be in the repository's main directory.

## Configuration
All configuration can be done by editing config.yml. The comments in the file itself provide extra explanation where necessary. [Join our Discord](https://discord.gg/Fjzpq8YPSn) if you have any extra questions.

## Disclaimer
This is a self-bot. Such bots are against Discord's terms of service. Automation of Dank Memer commands also breaks Dank Memer's rules. By using this software you acknowledge that we take no responsibility whatsoever for any action taken against your account, whether by Discord or Dank Memer, for not abiding by their respective rules.

Despite this, we believe the chance of detection by either Discord or Dank Memer to be low provided that you take appropriate action to ensure this. This includes but is not limited to running the bot only in private channels and not being open about the fact that you use it.
