# Ouroboros Collective - Migrations-Guide

**ARE-Diablo2-BaalAal**

---

## 🔄 Migration zur Axiomatic BaalAal Engine

### Von Otto Engine zu Axiomatic

Die **Axiomatic BaalAal Engine** ersetzt die alte Otto Engine:

#### 1. Script-Änderungen
```javascript
// ALT (Otto)
js("player.attack()")

// NEU (Axiomatic)
ax-attack
```

#### 2. State Management
```go
// ALT
engine.ExecuteScript("state.player.health = 100")

// NEU
engine.SetKappaState("player", "health", 100)
```

#### 3. Event System
```go
// ALT
engine.On("player.attack", handler)

// NEU
engine.Subscribe("ax-attack", handler)
```

---

## 📱 Mobile Migration

### Web-Client Setup
```bash
cd web-client
npm install
npm run dev
```

### Touch-Controls
```typescript
// Virtual Joystick
import { Joystick } from './components/Joystick';

// Hotbar
import { Hotbar } from './components/Hotbar';
```

---

## 🐍 Ouroboros Cycle

Das Ouroboros Cycle System für Resonanz-Stabilität:

1. **Feed:** Aktion ausführen
2. **Grow:** Resonanz erhöhen
3. **Shed:** Alte Daten bereinigen
4. **Repeat:** Zyklus fortsetzen

---

## ✅ Checkliste

- [ ] Otto Engine Referenzen entfernt
- [ ] Axiomatic Commands implementiert
- [ ] Kappa-space Koordinaten konfiguriert
- [ ] Web-Client Touch-Controls
- [ ] Mobile Build getestet

---

*Ouroboros Collective - Juni 2026*
