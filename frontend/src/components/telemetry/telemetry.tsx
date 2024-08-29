import { useEffect, useRef } from 'react';
import { useAppDispatch } from '../../app/hooks';
import { getLatestStandings } from './telemetrySlice';

type Delay = number | null;
type TimerHandler = (...args: any[]) => void;

export const Telemetry = () => {
  const useInterval = (callback: TimerHandler, delay: Delay) => {
    const savedCallbackRef = useRef<TimerHandler>();

    useEffect(() => {
      savedCallbackRef.current = callback;
    }, [callback]);

    useEffect(() => {
      const handler = (...args: any[]) => savedCallbackRef.current!(...args);

      if (delay !== null) {
        const intervalId = setInterval(handler, delay);
        return () => clearInterval(intervalId);
      }
    }, [delay]);
  };

  const dispatch = useAppDispatch();

  useInterval(() => {
    dispatch(getLatestStandings())
  }, 3000);

  return null
};
