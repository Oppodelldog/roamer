# roamer

Demo Project how to automate Windows keyboard/mouse interactions through a web frontend.

* Webservice (HTTP, Websockets) [go]
* Webfrontend [Vue.js]

The webservice uses  **Win32 API** for simulating keyboard and mouse events and **Windows Core Audio API** for sound
settings.

## Install

```cmd
.\install.bat
```

The install script builds `roamer.exe` and writes it to `GOBIN`, or to `GOPATH\bin` when `GOBIN` is not set.
Make sure that directory is on `PATH` if you want to launch Roamer from a terminal or Windows search.
Close any running Roamer instance before installing, otherwise Windows may block replacing `roamer.exe`.

**Attention**

* Read and understand the source code before using this tool
* Some of those tasks require special preparation
* Ensure to execute any of the following actions only with windows focused on the appropriate game.

## Screenshots

### Phone

![Roamer - Phone Overview](Galaxy-J7-1.png)
![Roamer - Phone Macro List](Galaxy-J7-2.png)
![Roamer - Phone Sound Settings](Galaxy-J7.png)
![Roamer - Phone Macro Editor](Galaxy-J7-3.png)

### Tablet landscape

![Roamer - Tablet Landscape Overview](Roamer-Tablet-Landscape-Overview.png)
![Roamer - Tablet Landscape Macro Deck](Roamer-Tablet-Landscape-Macro-Deck.png)
![Roamer - Tablet Landscape Sound Settings](Roamer-Tablet-Landscape-Sound.png)
![Roamer - Tablet Landscape Macro Editor](Roamer-Tablet-Landscape-Macro-Editor.png)

This default configuration brings some samples for some games.

### rust

* run
* paddle kayak
* left-clicking
* mouse pos investigation
* arm first inventory row
* unarm items
* fill up an inventory row
* transfer inventory row
* move inventory row up
* smart breath
* dive tank on/off

### seven days to die

* left-clicking
* walk
* run
* repair slot (1-10)

### altf4

*as you might expect this did not work out really well  
this game beats procedural timed inputs as good as human inputs ;-)*

* attempt for a full run
* several sequences to jump the first hurdles

### valheim

* walk
* run
* grillmaster

---
**Attributions:**  
Background
Image: [aliffian arief](https://unsplash.com/@helip?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyTex)
