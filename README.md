# ir Standings - iRacing Live Championship overlay

This project provides a lightweight overlay for iRacing to display live Championship Standings.



The windows overlay code is based on [iRon](https://github.com/lespalt/iRon) with the backend `go` code from myself. A big thank you to L. E. Spalt
for making iRon Open Source.

I don't currently plan to extend it further. That said, I'm making it available in the hope it might be useful to others in the iRacing community,
either for direct use or as a starting point for other apps.

# Contents

- [Where to Download](#where-to-download)
- [Overlays](#overlays)
  - [*Standings*](#standings)
- [Installing & Running](#installing--running)
- [Configuration](#configuration)
- [Building from source](#building-from-source)
- [Dependencies](#dependencies)
- [Bug reports and feature requests](#bug-reports-and-feature-requests)

---

## Where to Download

The latest binary release can be found [here](https://github.com/ianhaycox/ir-standings/releases/latest).

## Overlays

### *Standings*

Shows the standings of the entire field, including safety rating, iRating, and number of laps since the last pit stop ("pit age"). I usually leave this off by default and switch it on during cautions. Or glimpse at it pre-race to get a sense of the competition level.

This will highlight buddies in green (Dale Jr. in the example below).

![standings](https://github.com/ianhaycox/ir-standings/blob/develop/live-standings.png?raw=true)

---

## Installing & Running

The app does not require installation. Just copy the executable to a folder of your choice. Make sure the folder is not write protected, as iRon will attempt to save its configuration file in the working directory.

To use it, simply run the executable. It doesn't matter whether you do this before or after launching iRacing. A console window will pop up, indicating that iRon is running. Once you're in the car in iRacing, the overlays should show up, and you can configure things to your liking. I recommend running iRacing in borderless window mode. Overlays *might* work in other modes as well, but I haven't tested it.

---

## Configuration

To place and resize the overlays, press ALT-j. This will enter a mode in which you can move overlays around with the mouse and resize them by dragging their bottom-right corner. Press ALT-j again to go back to normal mode.

Overlays can be switched on and off at runtime using the hotkeys displayed during startup. All hotkeys are configurable.

Certain aspects of the overlays, such as colors, font types, sizes etc. can be customized. To do that, open the file **config.json** that iRon created and experiment by editing the (hopefully mostly self-explanatory) parameters. You can do that while the app is running -- the changes will take effect immediately whenever the file is saved.

_Note that currently, the config file will be created only after the overlays have been "touched" for the first time, usually by dragging or resizing them._

---

## Building from source

This app is built with Visual Studio 2022. The free version should suffice, though I haven't verified it. The project/solution files should work out of the box. Depending on your Visual Studio setup, you may need to install additional prerequisites (static libs) needed to build DirectX applications.

---

## Dependencies

There are no runtime dependencies other than standard Windows components like DirectX.  Those should already be present on most if not all systems that can run iRacing.

Build dependencies (most notably the iRacing SDK and picojson) are kept to a minimum and are included in the repository.

---

## Bug reports and feature requests

If you encounter a problem, please file a github issue and I'll do my best to address it. Pull requests with fixes are welcome too, of course.

If you'd like to see a specific feature added, feel free to file a github issue as well. If it's something small, I may actually get to it :-) No promises though, as unfortunately the time I can spend on this project is quite limited.

---

## Donations

If you like this project enough to wonder whether you can contribute financially: first of all, thank you! I'm not looking for donations, but **please consider giving to Ukraine-related charities instead**.
