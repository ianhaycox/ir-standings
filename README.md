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

## Overlays

### *Standings*

Shows the top ten standings of each class. Click on the 'Class Name' box bottom left to toggle between classes.

![standings](https://github.com/ianhaycox/ir-standings/blob/develop/images/live-standings.png?raw=true)

---

## Installing & Running

The app does not require installation. Just copy the executable to a folder of your choice.

To use it, simply run the executable. It doesn't matter whether you do this before or after launching iRacing.

---

## Development

Install https://wails.io/

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
