import { IRacingSDK } from "irsdk-node";
import { SendTelemetry, SendSessionInfo } from '../../../wailsjs/go/main/App';

const TIMEOUT = Math.floor((1 / 60) * 1000); // 60fps

// Create a loop to use if the service is running.
function loop(sdk: IRacingSDK): void {
    if (sdk.waitForData(TIMEOUT)) {
        const session = sdk.getSessionData();
        const telemetry = sdk.getTelemetry();

        // Do something with the data

        SendTelemetry(JSON.stringify(telemetry, null, 2))
        SendSessionInfo(JSON.stringify(session, null, 2))

        // Wait for new data
        loop(sdk);
    } else {
        // Disconnected.
    }
}

// Check the iRacing service is running
if (await IRacingSDK.IsSimRunning()) {
    const sdk = new IRacingSDK();

    // Start the SDK
    if (await sdk.ready()) {
        // If you want the SDK to auto-enable telemetry
        sdk.autoEnableTelemetry = true;

        // Create a loop
        loop(sdk);
    }
}
