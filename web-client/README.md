# Ouroboros Collective - Web Client

**ARE-Diablo2-BaalAal - Mobile/Web Frontend**

---

## 🎮 Überblick

Der **Web Client** ist das mobile-optimierte Frontend für die **Axiomatic BaalAal Engine**. Er ermöglicht das Spielen von Diablo 2 über den Browser auf Desktop und Mobile-Geräten.

### Technologie-Stack
- **Framework:** Next.js 14+
- **Sprache:** TypeScript
- **Styling:** Tailwind CSS
- **3D Rendering:** Babylon.js
- **State:** React hooks

---

## 🐍 Axiomatic Integration

Der Web-Client kommuniziert mit der **Axiomatic BaalAal Engine** über:
- WebSocket für Echtzeit-Kommunikation
- REST API für Asset-Metadaten
- Kappa-space Koordinaten für deterministische State-Synchronisation

---

## 🚀 Getting Started

### Voraussetzungen
- Node.js 18+
- npm oder yarn

### Installation
```bash
cd web-client
npm install
```

### Development
```bash
npm run dev
# Öffnet http://localhost:3000
```

### Production Build
```bash
npm run build
npm start
```

---

## 📱 Mobile-Optimierung

### Touch-Controls
- Virtual Joystick für Bewegung
- Tap-to-Attack Interface
- On-screen Hotbar

### Responsive Design
- Unterstützt 16:9, 18:9, 20:9 Displays
- Portrait und Landscape Mode
- Dynamic HUD scaling

### Performance
- Lazy Loading für Assets
- Progressive Web App (PWA) Support
- Offline-Capability (Service Worker)

---

## 📁 Projektstruktur

```
web-client/
├── app/              # Next.js App Router
│   └── page.tsx      # Hauptseite
├── components/        # React Komponenten
│   └── BabylonScene.tsx  # 3D Rendering
├── public/           # Statische Assets
├── utils/            # Utility Funktionen
└── package.json      # Dependencies
```

---

## 🔧 Konfiguration

### Environment Variables
```env
NEXT_PUBLIC_WS_URL=ws://localhost:8080
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

### Mobile-spezifische Einstellungen
```typescript
// utils/mobile-config.ts
export const MOBILE_CONFIG = {
  joystickSize: 120,
  hotbarSlots: 6,
  targetFPS: 60,
};
```

---

## 🧪 Testing

```bash
# Unit Tests
npm test

# E2E Tests
npm run test:e2e

# Linting
npm run lint
```

---

## 📜 Lizenz

GPL v3 - Siehe Haupt-[README.md](../README.md)

---

*Ouroboros Collective - Juni 2026*
