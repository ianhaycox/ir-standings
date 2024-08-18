import React, { useState, useEffect } from 'react';
import { SendTelemetry, SendSessionInfo } from '../../../wailsjs/go/main/App';

export const Telemetry = () => {
  const [result, setResult] = useState(null);
  const [worker, setWorker] = useState<Worker | null>(null);

  useEffect(() => {
    // Create a new web worker
    const myWorker = new Worker(new URL('./iracing.js', import.meta.url));

    // Set up event listener for messages from the worker
    myWorker.onmessage = function (event) {
      console.log('Received result from worker:', event.data);
      SendSessionInfo("session");
      SendTelemetry("telementry");

      setResult(event.data);
    };

    // Save the worker instance to state
    setWorker(myWorker);

    myWorker.postMessage("start");

    // Clean up the worker when the component unmounts
    return () => {
      myWorker.terminate();
    };
  }, []); // Run this effect only once when the component mounts

  return null
};
