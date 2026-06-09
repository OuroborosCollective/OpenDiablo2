"use client";

import React, { useRef, useState, useCallback, useEffect } from "react";

interface JoystickProps {
  size?: number;
  onMove?: (x: number, y: number) => void;
  onDirectionChange?: (direction: string) => void;
  position?: "left" | "right";
  sensitivity?: number;
}

export default function VirtualJoystick({
  size = 120,
  onMove,
  onDirectionChange,
  position = "left",
  sensitivity = 1.0,
}: JoystickProps) {
  const [active, setActive] = useState(false);
  const [positionState, setPositionState] = useState({ x: 0, y: 0 });
  const containerRef = useRef<HTMLDivElement>(null);
  const [isMobile, setIsMobile] = useState(false);
  const touchIdRef = useRef<number | null>(null);

  useEffect(() => {
    setIsMobile(/Android|iPhone|iPad|iPod/i.test(navigator.userAgent));
  }, []);

  const getDirection = (x: number, y: number): string => {
    const threshold = 0.3;
    if (Math.abs(x) < threshold && Math.abs(y) < threshold) return "idle";
    
    const angle = Math.atan2(y, x) * (180 / Math.PI);
    
    if (angle >= -22.5 && angle < 22.5) return "right";
    if (angle >= 22.5 && angle < 67.5) return "down-right";
    if (angle >= 67.5 && angle < 112.5) return "down";
    if (angle >= 112.5 && angle < 157.5) return "down-left";
    if (angle >= 157.5 || angle < -157.5) return "left";
    if (angle >= -157.5 && angle < -112.5) return "up-left";
    if (angle >= -112.5 && angle < -67.5) return "up";
    if (angle >= -67.5 && angle < -22.5) return "up-right";
    
    return "idle";
  };

  const handleTouchStart = useCallback((e: React.TouchEvent) => {
    e.preventDefault();
    if (touchIdRef.current !== null) return;
    
    const touch = e.changedTouches[0];
    touchIdRef.current = touch.identifier;
    setActive(true);
  }, []);

  const handleTouchMove = useCallback((e: React.TouchEvent) => {
    if (!containerRef.current || touchIdRef.current === null) return;
    
    const touch = Array.from(e.changedTouches).find(t => t.identifier === touchIdRef.current);
    if (!touch) return;
    
    const rect = containerRef.current.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    
    let deltaX = (touch.clientX - centerX) * sensitivity;
    let deltaY = (touch.clientY - centerY) * sensitivity;
    
    const maxDistance = rect.width / 2 - 20;
    const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);
    
    if (distance > maxDistance) {
      deltaX = (deltaX / distance) * maxDistance;
      deltaY = (deltaY / distance) * maxDistance;
    }
    
    setPositionState({ x: deltaX, y: deltaY });
    
    const normalizedX = deltaX / maxDistance;
    const normalizedY = deltaY / maxDistance;
    
    onMove?.(normalizedX, normalizedY);
    
    const direction = getDirection(normalizedX, normalizedY);
    onDirectionChange?.(direction);
  }, [onMove, onDirectionChange, sensitivity]);

  const handleTouchEnd = useCallback((e: React.TouchEvent) => {
    const touch = Array.from(e.changedTouches).find(t => t.identifier === touchIdRef.current);
    if (!touch) return;
    
    touchIdRef.current = null;
    setActive(false);
    setPositionState({ x: 0, y: 0 });
    onMove?.(0, 0);
    onDirectionChange?.("idle");
  }, [onMove, onDirectionChange]);

  if (!isMobile) return null;

  const normalizedMagnitude = Math.sqrt(positionState.x ** 2 + positionState.y ** 2) / (size / 2 - 20);
  const opacity = 0.3 + normalizedMagnitude * 0.7;

  return (
    <div
      ref={containerRef}
      className={`absolute ${position === "left" ? "bottom-12 left-6" : "bottom-12 right-6"} z-30 select-none`}
      style={{ width: size, height: size }}
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleTouchEnd}
    >
      {/* Outer ring with gradient */}
      <div
        className="absolute inset-0 rounded-full border-2 border-white/40 backdrop-blur-md transition-all duration-75"
        style={{
          background: `radial-gradient(circle, rgba(59,130,246,${opacity * 0.2}) 0%, rgba(0,0,0,${opacity * 0.5}) 100%)`,
          boxShadow: active 
            ? `0 0 ${30 + normalizedMagnitude * 20}px rgba(59,130,246,${opacity})` 
            : "0 0 15px rgba(59,130,246,0.3)"
        }}
      >
        {/* Inner glow ring */}
        <div className="absolute inset-2 rounded-full border border-white/20" />
      </div>
      
      {/* Inner stick with glow */}
      <div
        className="absolute rounded-full transition-transform duration-75"
        style={{
          width: size * 0.45,
          height: size * 0.45,
          left: "50%",
          top: "50%",
          transform: `translate(calc(-50% + ${positionState.x}px), calc(-50% + ${positionState.y}px))`,
          background: active
            ? `linear-gradient(135deg, #3b82f6, #8b5cf6, #ec4899)`
            : "linear-gradient(135deg, #1e40af, #4f46e5, #7c3aed)",
          border: `2px solid rgba(255,255,255,${0.3 + normalizedMagnitude * 0.5})`,
          boxShadow: active
            ? `0 0 ${20 + normalizedMagnitude * 15}px rgba(59,130,246,${0.8}), inset 0 0 20px rgba(255,255,255,0.3)`
            : "0 4px 15px rgba(0,0,0,0.5), inset 0 0 10px rgba(255,255,255,0.1)"
        }}
      >
        {/* Center dot */}
        <div className="absolute inset-0 flex items-center justify-center">
          <div 
            className="w-3 h-3 rounded-full bg-white/60"
            style={{
              boxShadow: active ? "0 0 8px rgba(255,255,255,0.8)" : "none"
            }}
          />
        </div>
      </div>
      
      {/* Direction indicators with active state */}
      <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
        {[
          { dir: "up", top: 8, left: "50%", transform: "translateX(-50%)" },
          { dir: "down", bottom: 8, left: "50%", transform: "translateX(-50%)" },
          { dir: "left", left: 8, top: "50%", transform: "translateY(-50%)" },
          { dir: "right", right: 8, top: "50%", transform: "translateY(-50%)" },
        ].map(({ dir, ...pos }) => {
          const isActiveDir = active && getDirection(
            positionState.x / (size / 2 - 20),
            positionState.y / (size / 2 - 20)
          ).includes(dir);
          
          return (
            <div
              key={dir}
              {...pos}
              className="absolute w-2 h-2 rounded-full transition-all duration-100"
              style={{
                background: isActiveDir ? "rgba(59,130,246,0.9)" : "rgba(255,255,255,0.2)",
                boxShadow: isActiveDir ? "0 0 10px rgba(59,130,246,0.8)" : "none"
              }}
            />
          );
        })}
      </div>
      
      {/* Magnitude indicator */}
      {active && (
        <div className="absolute -bottom-6 left-1/2 transform -translate-x-1/2 text-[10px] text-white/60 font-mono">
          {Math.round(normalizedMagnitude * 100)}%
        </div>
      )}
    </div>
  );
}
