# Ouroboros Collective - Entwickler-Anleitung

**ARE-Diablo2-BaalAal**

---

## 🚀 Development Setup

### 1. Repository klonen
```bash
git clone https://github.com/OuroborosCollective/OpenDiablo2.git
cd OpenDiablo2
```

### 2. Backend starten
```bash
go run .
# Server läuft auf http://localhost:8080
```

### 3. Web-Client starten (separate Terminal)
```bash
cd web-client
npm run dev
# Client läuft auf http://localhost:3000
```

---

## 🐍 Axiomatic BaalAal Engine

### Engine Struktur
```
d2script/
├── axiomatic.go      # Haupt-Engine
├── axiomatic_test.go # Tests
├── kappa_system.go    # Kappa-Koordinaten
└── BaalAal/          # Ouroboros Cycle
```

### Engine initialisieren
```go
engine := NewBaalAalEngine()
engine.ProcessCycle(tick)
```

### Axiomatic Commands
- `ax-state` - Zeige aktuellen State
- `ax-kappa` - Kappa-space Koordinaten
- `ax-resonance` - Resonanz-Level

---

## 🌐 Web-Client Entwicklung

### Komponenten erstellen
```bash
cd web-client/components
touch NewComponent.tsx
```

### Babylon.js Integration
```typescript
// components/BabylonScene.tsx
import { Scene, Engine } from '@babylonjs/core';

export function BabylonScene() {
  // 3D Scene Setup
}
```

### Mobile-spezifische Entwicklung
```typescript
// utils/mobile.ts
export const isMobile = /Android|iPhone|iPad/.test(navigator.userAgent);
```

---

## 🧪 Testing

### Go Tests
```bash
go test ./d2script/...
go test ./d2common/...
```

### Web-Client Tests
```bash
cd web-client
npm test
```

---

## 📝 Code-Stil

- Go: `gofmt` und `golangci-lint`
- TypeScript: ESLint + Prettier
- Commits: Conventional Commits (`feat:`, `fix:`, `docs:`)

---

## 🔧 Nützliche Commands

```bash
# Linting
golangci-lint run

# Formatierung
go fmt ./...
npm run lint -- --fix

# Profiling
go tool pprof http://localhost:6060/debug/pprof/
```

---

## 🐛 Debugging

### Backend Debugging
```bash
# Mit Debug-Output
LOG_LEVEL=debug go run .

# JSON-RPC Console
curl -X POST http://localhost:8080/rpc -d '{"method":"debug.state"}'
```

### Web-Client Debugging
```bash
# React DevTools
npm install @babylonjs/react

# Network Tab für WebSocket
```

---

*Ouroboros Collective - Juni 2026*
