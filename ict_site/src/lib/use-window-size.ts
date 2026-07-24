"use client";

import { useSyncExternalStore } from "react";

function subscribe(callback: () => void) {
  window.addEventListener("resize", callback);
  return () => window.removeEventListener("resize", callback);
}

function getSnapshot() {
  return { width: window.innerWidth, height: window.innerHeight };
}

function getServerSnapshot() {
  return { width: 0, height: 0 };
}

export function useWindowSize() {
  return useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);
}
