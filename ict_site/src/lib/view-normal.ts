import * as React from "react";

const NORMAL_BREAKPOINT = 1024;

export function useIsNormal() {
  const [isNormal, setIsNormal] = React.useState<boolean | undefined>(undefined);
  React.useEffect(() => {
    const mql = window.matchMedia(`(min-width: ${NORMAL_BREAKPOINT}px)`);
    const onChange = () => {
      setIsNormal(window.innerWidth >= NORMAL_BREAKPOINT);
    };
    mql.addEventListener("change", onChange);
    setIsNormal(window.innerWidth >= NORMAL_BREAKPOINT);
    return () => mql.removeEventListener("change", onChange);
  }, []);
  return !!isNormal;
}
