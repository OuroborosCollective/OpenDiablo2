export interface MobileConfig {
  isMobile: boolean;
  isTablet: boolean;
  isLandscape: boolean;
  displayRatio: string;
  scaleFactor: number;
  touchEnabled: boolean;
}

export const getMobileConfig = (): MobileConfig => {
  if (typeof window === "undefined") {
    return {
      isMobile: false,
      isTablet: false,
      isLandscape: true,
      displayRatio: "16:9",
      scaleFactor: 1,
      touchEnabled: false,
    };
  }

  const width = window.innerWidth;
  const height = window.innerHeight;
  const userAgent = navigator.userAgent;
  
  const isMobile = /Android|iPhone|iPad|iPod/i.test(userAgent) || width < 768;
  const isTablet = /iPad|Android/i.test(userAgent) || (width >= 768 && width < 1024);
  const isLandscape = width > height;
  
  const ratio = width / height;
  let displayRatio = "16:9";
  if (ratio >= 2.0) displayRatio = "21:9";
  else if (ratio >= 1.7) displayRatio = "16:9";
  else if (ratio >= 1.5) displayRatio = "18:9";
  
  const scaleFactor = Math.min(width / 400, height / 800);
  const clampedScale = Math.max(0.8, Math.min(1.2, scaleFactor));
  
  const touchEnabled = "ontouchstart" in window || navigator.maxTouchPoints > 0;

  return {
    isMobile,
    isTablet,
    isLandscape,
    displayRatio,
    scaleFactor: clampedScale,
    touchEnabled,
  };
};

export const MOBILE_CONFIG = {
  joystickSize: 130,
  hotbarSlots: 6,
  targetFPS: 60,
  touchThreshold: 0.1,
  movementSpeed: 0.15,
  animationDuration: 150,
};
