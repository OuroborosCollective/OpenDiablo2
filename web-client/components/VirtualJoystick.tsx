"use client";

import React, { useRef, useState, useCallback, useEffect } from "react";

interface JoystickProps {
  size?: number;
  onMove?: (x: number, y: number) => void;
  position?: "left" | "right";
}

export default function VirtualJoystick({
  size = 120,
  onMove,
  position = "left",
}: JoystickProps) {
  const [active, setActive] = useState(false);
  const [positionState, setPositionState] = useState({ x: 0, y: 0 });
  const containerRef = useRef<HTMLDivElement>(null);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    setIsMobile(/Android|iPhone|iPad|iPod/i.test(navigator.userAgent));
  }, []);

  const handleTouchStart = useCallback((e: React.TouchEvent) => {
    e.preventDefault();
    setActive(true);
  }, []);

  const handleTouchMove = useCallback((e: React.TouchEvent) => {
    if (!containerRef.current) return;
    
    const touch = e.touches[0];
    const rect = containerRef.current.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    
    let deltaX = touch.clientX - centerX;
    let deltaY = touch.clientY - centerY;
    
    const maxDistance = rect.width / 2 - 20;
    const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);
    
    if (distance > maxDistance) {
      deltaX = (deltaX / distance) * maxDistance;
      deltaY = (deltaY / distance) * maxDistance;
    }
    
    setPositionState({ x: deltaX, y: deltaY });
    
    // Normalize to -1 to 1
    const normalizedX = deltaX / maxDistance;
    const normalizedY = deltaY / maxDistance;
    
    onMove?.(normalizedX, normalizedY);
  }, [onMove]);

  const handleTouchEnd = useCallback(() => {
    setActive(false);
    setPositionState({ x: 0, y: 0 });
    onMove?.(0, 0);
  }, [onMove]);

  if (!isMobile) return null;

  return (
    <div
      ref={containerRef}
      className={`absolute ${position === "left" ? "bottom-8 left-8" : "bottom-8 right-8"} z-30`}
      style={{ width: size, height: size }}
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleTouchEnd}
    >
      {/* Outer ring */}
      <div
        className="absolute inset-0 rounded-full bg-white/10 border-2 border-white/30 backdrop-blur-sm"
        style={{
          boxShadow: active 
            ? "0 0 20px rgba(255,255,255,0.3)" 
            : "0 0 10px rgba(255,255,255,0.1)"
        }}
      />
      
      {/* Inner stick */}
      <div
        className="absolute rounded-full bg-gradient-to-br from-blue-500 to-purple-600 border-2 border-white/50"
        style={{
          width: size * 0.4,
          height: size * 0.4,
          left: "50%",
          top: "50%",
          transform: `translate(calc(-50% + ${positionState.x}px), calc(-50% + ${positionState.y}px))`,
          boxShadow: active 
            ? "0 0 15px rgba(59,130,246,0.6)" 
            : "0 2px 8px rgba(0,0,0,0.3)"
        }}
      />
      
      {/* Direction indicators */}
      <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
        <div className="absolute top-2 w-0.5 h-2 bg-white/20 rounded-full" />
        <div className="absolute bottom-2 w-0.5 h-2 bg-white/20 rounded-full" />
        <div className="absolute left-2 h-0.5 w-2 bg-white/20 rounded-full" />
        <div className="absolute right-2 h-0.5 w-2 bg-white/20 rounded-full" />
      </div>
    </div>
  );
}
