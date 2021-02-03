# Dank Grinder [![](https://goreportcard.com/badge/github.com/dankgrinder/dankgrinder)](https://goreportcard.com/report/github.com/dankgrinder/dankgrinder) [![](https://img.shields.io/github/workflow/status/dankgrinder/dankgrinder/Go)](https://github.com/dankgrinder/dankgrinder/actions) [![](https://img.shields.io/github/v/release/dankgrinder/dankgrinder)](https://github.com/dankgrinder/dankgrinder/releases/latest)
The most feature-rich, advanced Dank Memer automation tool.

Want to join the community or need support? [Join our Discord](https://discord.gg/Fjzpq8YPSn). Have a question? [Create a question issue](https://github.com/dankgrinder/dankgrinder/issues/new?assignees=&labels=question&template=question.md&title=). Want to suggest a feature? [Create a suggestion issue](https://github.com/dankgrinder/dankgrinder/issues/new?assignees=&labels=suggestion&template=suggestion.md&title=). Encountered a bug? [Report a bug](https://github.com/dankgrinder/dankgrinder/issues/new?assignees=&labels=bug&template=bug-report.md&title=). Want to contribute? [Read our contribution guidelines](https://github.com/dankgrinder/dankgrinder/blob/master/CONTRIBUTING.md).

<center>
<img src="https://i.imgur.com/VSYIKaX.png" alt="logo" style="margin-top:10px" />
</center>

## Features
* Cycles through currency commands to earn approximately 150,000 coins per hour
* Responds to fishing and hunting events (i.e. captures the dragon, etc.)
* Responds to global events, search, and postmeme
* Automation of highlow
* Reports an average income every 2 minutes
* Automatically buys a fishing pole, hunting rifle or laptop if they have broken
* High configurability
* Can run many instances at once
* Uses very little system resources

## Getting started

### Installation
1. Download the latest release for your operating system [here](https://github.com/dankgrinder/dankgrinder/releases/latest) (darwin is for macOS), or [build from source](#building-from-source). If you build from source you can delete everything besides the compiled binaries and config.yml.
2. Extract the archive
3. [Configure](#configuration), make sure `token` and `channel_id` fields of the instance are set, it will not run otherwise

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
   
The executables should then be in the `/build` folder.

## Configuration
All configuration can be done by editing config.yml with your editor of choice (e.g. Visual Studio Code, Sublime Text, Notepad++). The comments in the file itself provide extra explanation where necessary. For the bot to run, you must at least enter an [authorization token](#getting-an-authorization-token) and a [channel id](#enabling-discords-developer-mode). If you are running into issues, you can [join our Discord](https://discord.gg/Fjzpq8YPSn). Example configuration for one instance:  
```yaml
instances:
  - token: "bmljZSB0cnkgYnV0IHRoaXMgaXM.bm90IGE.cmVhbCB0b2tlbg"
    channel_id: "791694339116892202"
    shifts:
      - state: "active"
        duration:
          base: -1
          variance: 0
```

Add more instances to run more bots at once, for example:  
```yaml
instances:
  - token: "bmljZSB0cnkgYnV0IHRoaXMgaXM.bm90IGE.cmVhbCB0b2tlbg"
    channel_id: "791694339116892202"
    shifts:
      - state: "active"
        duration:
          base: -1
          variance: 0
  - token: "b2YgY291cnNlIHRoaXM.aXNuJ3QgYQ.cmVhbCB0b2tlbiBlaXRoZXIsIHNpbGx5"
       channel_id: "791694383098495047"
       shifts:
          - state: "active"
            duration:
               base: -1
               variance: 0
```

## Disclaimer
This is a self-bot. Such bots are against Discord's terms of service. Automation of Dank Memer commands also breaks Dank Memer's rules. By using this software you acknowledge that we take no responsibility whatsoever for any action taken against your account, whether by Discord or Dank Memer, for not abiding by their respective rules.

Despite this, we believe the chance of detection by either Discord or Dank Memer to be low provided that you take appropriate measures to ensure this. This includes but is not limited to running the bot only in private channels, not being open about the fact that you use it and not running so much as to raise suspicion. 
