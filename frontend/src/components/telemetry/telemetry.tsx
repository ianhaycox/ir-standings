import React, { useState, useEffect } from 'react';
import { SendTelemetry, SendSessionInfo } from '../../../wailsjs/go/main/App';

export const Telemetry = () => {
  const [result, setResult] = useState(null);
  const [worker, setWorker] = useState<Worker | null>(null);

  useEffect(() => {
    // Create a new web worker
    const myWorker = new Worker(new URL('./irsdk-node.ts', import.meta.url));

    // Save the worker instance to state
    setWorker(myWorker);

    // Clean up the worker when the component unmounts
    return () => {
      myWorker.terminate();
    };
  }, []); // Run this effect only once when the component mounts

  return null
};
