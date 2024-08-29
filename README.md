# ir Standings - iRacing Live Championship overlay

This project provides a lightweight browser overlay for iRacing to display live Championship Standings.

# Contents

- [Where to Download](#where-to-download)
- [Overlays](#overlays)
  - [*Standings*](#standings)
- [Installing & Running](#installing--running)
- [Development](#development)
- [Bug reports and feature requests](#bug-reports-and-feature-requests)

---

## Where to Download

The latest binary release can be found [here](https://github.com/ianhaycox/ir-standings/releases/latest).

### Installing & Running

The app does not require installation. Just copy the executable to a folder of your choice.

To use it, simply run the executable. It doesn't matter whether you do this before or after launching iRacing.

If you prefer to use the installer, that is available as well.

---

## Overlays

### *Standings*

Shows the top ten standings of each class. Click on the 'Class Name' box bottom left to toggle between classes.

![standingsgtp](https://github.com/ianhaycox/ir-standings/blob/develop/images/live-standings-gtp.png?raw=true)

![standingsgto](https://github.com/ianhaycox/ir-standings/blob/develop/images/live-standings-gto.png?raw=true)

Login with your iRacing email and password. The details are not saved, but are required to download the results for previous broadcast races.

The VCR Championship rules are used to calculate championship points from the previous races and the current race.
The current race does not have to be the Saturday broadcast race, the current positions are used along with the prior broadcast results to work out the live standings.

Drivers greyed out in the table are not present in the current session.

The overlay has very low resource usage only updating every 3 seconds to determine the new race positions.

NOTE: The window title bar is invisible but is there on screen. To re-size or close the overlay you'll have to guess where the borders are and the close icon is.

---

## Development

Install https://wails.io/

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

## TODO

- Allow the window to be used as an OBS Browser source
- Sort out closing the app/window resizing and title bar.
- Cache previous results to avoid fetching each time.
- Configurable settings - refresh rate, best of, points system, top 10, colors etc.
- 