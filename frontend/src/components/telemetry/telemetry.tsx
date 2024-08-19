import { useState, useEffect } from 'react';
import worker_func from './irsdk';

export const Telemetry = () => {
  const [result, setResult] = useState(null);
  const [worker, setWorker] = useState<Worker | null>(null);

  useEffect(() => {
    // Create a new web worker
    const myWorker = new Worker(URL.createObjectURL(new Blob(["("+worker_func.toString()+")()"], {type: 'text/javascript'})));

    // Save the worker instance to state
    setWorker(myWorker);

    // Clean up the worker when the component unmounts
    return () => {
      myWorker.terminate();
    };
  }, []); // Run this effect only once when the component mounts

  return null
};
