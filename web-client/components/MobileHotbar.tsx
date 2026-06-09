"use client";

import React, { useState, useEffect } from "react";

interface HotbarProps {
  slots?: number;
  onSkillUse?: (slot: number) => void;
}

interface Skill {
  id: number;
  name: string;
  icon: string;
  cooldown?: number;
}

const DEFAULT_SKILLS: Skill[] = [
  { id: 1, name: "Attack", icon: "⚔️" },
  { id: 2, name: "Skill 1", icon: "🔥" },
  { id: 3, name: "Skill 2", icon: "❄️" },
  { id: 4, name: "Skill 3", icon: "⚡" },
  { id: 5, name: "Skill 4", icon: "🛡️" },
  { id: 6, name: "Ultimate", icon: "💥" },
];

export default function MobileHotbar({
  slots = 6,
  onSkillUse,
}: HotbarProps) {
  const [isMobile, setIsMobile] = useState(false);
  const [activeSkill, setActiveSkill] = useState<number | null>(null);
  const [cooldowns, setCooldowns] = useState<Record<number, number>>({});

  useEffect(() => {
    setIsMobile(/Android|iPhone|iPad|iPod/i.test(navigator.userAgent));
  }, []);

  const handleSkillPress = (skillId: number) => {
    setActiveSkill(skillId);
    onSkillUse?.(skillId);
    
    // Reset after short visual feedback
    setTimeout(() => setActiveSkill(null), 150);
    
    // Start cooldown simulation
    setCooldowns(prev => ({ ...prev, [skillId]: 100 }));
    
    // Cooldown countdown
    const interval = setInterval(() => {
      setCooldowns(prev => {
        const current = prev[skillId] || 0;
        if (current <= 0) {
          clearInterval(interval);
          return prev;
        }
        return { ...prev, [skillId]: current - 10 };
      });
    }, 100);
  };

  if (!isMobile) return null;

  return (
    <div className="absolute bottom-8 right-8 z-30 flex flex-col gap-2">
      {/* Skill slots */}
      <div className="flex flex-wrap gap-2 max-w-[200px] justify-end">
        {DEFAULT_SKILLS.slice(0, slots).map((skill) => {
          const cooldownPercent = cooldowns[skill.id] || 0;
          const isActive = activeSkill === skill.id;
          
          return (
            <button
              key={skill.id}
              onClick={() => handleSkillPress(skill.id)}
              className={`
                relative w-12 h-12 rounded-lg border-2 
                flex items-center justify-center text-2xl
                transition-all duration-150 transform
                ${isActive 
                  ? "bg-blue-500/50 border-blue-300 scale-110" 
                  : "bg-black/60 border-white/30 hover:border-white/50"
                }
              `}
              disabled={cooldownPercent > 0}
              style={{
                boxShadow: isActive 
                  ? "0 0 20px rgba(59,130,246,0.8)" 
                  : "0 2px 8px rgba(0,0,0,0.5)"
              }}
            >
              {skill.icon}
              
              {/* Cooldown overlay */}
              {cooldownPercent > 0 && (
                <div 
                  className="absolute inset-0 bg-black/70 rounded-lg flex items-center justify-center"
                  style={{ clipPath: `inset(${100 - cooldownPercent}% 0 0 0)` }}
                >
                  <span className="text-xs text-white font-bold">
                    {Math.ceil(cooldownPercent / 10)}
                  </span>
                </div>
              )}
              
              {/* Key hint */}
              <span className="absolute -top-1 -right-1 w-4 h-4 bg-neutral-800 rounded-full text-[8px] flex items-center justify-center text-white/70">
                {skill.id}
              </span>
            </button>
          );
        })}
      </div>
      
      {/* Quick actions row */}
      <div className="flex gap-2 justify-end">
        <button 
          className="w-10 h-10 rounded-full bg-red-600/80 border-2 border-red-400/50 flex items-center justify-center text-white text-sm font-bold"
          title="Inventory"
        >
          📦
        </button>
        <button 
          className="w-10 h-10 rounded-full bg-yellow-600/80 border-2 border-yellow-400/50 flex items-center justify-center text-white text-sm font-bold"
          title="Skills"
        >
          ⚡
        </button>
        <button 
          className="w-10 h-10 rounded-full bg-neutral-700/80 border-2 border-neutral-500/50 flex items-center justify-center text-white text-sm font-bold"
          title="Menu"
        >
          ☰
        </button>
      </div>
    </div>
  );
}
