# Ouroboros Collective - Projekt-Roadmap

**ARE-Diablo2-BaalAal - Mobile-First Strategie**

---

## 🎯 Vision

Unser Ziel ist es, das klassische Diablo 2 Erlebnis auf modernen Plattformen zugänglich zu machen - mit Fokus auf **Mobile**, **Web** und **Desktop**. Wir nutzen die **Axiomatic BaalAal Engine** für deterministische, sichere und performante Spielabläufe.

---

## 📊 Phasen-Übersicht

### Phase 1: Foundation ✅ (Abgeschlossen)
- [x] Axiomatic BaalAal Engine Integration
- [x] WebSocket Sicherheit (Origin Validation, Timeouts)
- [x] MPQ Hardening
- [x] Unit Tests für Core-Funktionen

### Phase 2: Mobile Core 🔄 (Aktiv)
- [ ] Web-Client Framework (Next.js)
- [ ] Touch-Controls Implementation
- [ ] Responsive UI für mobile Displays
- [ ] Asset-Streaming Optimierung

### Phase 3: Multiplayer 🔜 (Geplant)
- [ ] WebSocket Server für Mobile
- [ ] Session Management
- [ ] Cross-Platform Matchmaking

### Phase 4: Content 🔜 (Geplant)
- [ ] D2 Data File Parser
- [ ] Sprite/Rendering Engine
- [ ] Audio System

---

## 🐍 Axiomatic BaalAal Engine Features

| Feature | Status | Beschreibung |
|---------|--------|--------------|
| Deterministic State | ✅ | Kappa-space Koordinaten (1000-base) |
| Ouroboros Cycle | ✅ | Selbst-fressende Schlange Resonanz |
| KappaSystem | ✅ | Logik-Koordinaten-System |
| Security Layer | ✅ | Deaktivierter `js` Command |
| Mobile API | 🔄 | Touch-Event Integration |

---

## 🌐 Mobile-Specific Ziele

### Touch-Controls
- Virtual Joystick für Bewegung
- Tap-to-Attack Interface
- Gesture Recognition (swipe, pinch)
- On-screen Hotbar

### Responsive Design
- 16:9, 18:9, 20:9 Display Support
- Dynamic HUD scaling
- Portrait/Landscape mode

### Performance
- 60 FPS target auf Mobile
- Low-latency input handling
- Efficient memory usage

### Connectivity
- Offline-Capability
- Cloud Save Sync
- Cross-device progression

---

## 🔧 Technologie-Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| Backend | Go 1.21+ | Game Server, Engine |
| Engine | Axiomatic BaalAal | Deterministic Logic |
| Web Client | Next.js 14+ | Mobile Web UI |
| Rendering | Babylon.js / Ebiten | 3D/2D Graphics |
| State | Kappa-space | Coordinate System |

---

## 📅 Meilensteine

### Q2 2026 - Mobile Alpha
- Web-Client funktionsfähig
- Basis Touch-Controls
- Single-Player Modus

### Q3 2026 - Mobile Beta
- Multiplayer support
- Account-System
- Cloud saves

### Q4 2026 - Production
- iOS/Android Apps
- App Store Veröffentlichung
- Community Features

---

## 🤝 Beitragen

Siehe [CONTRIBUTING.md](./CONTRIBUTING.md) für Details zur Mitarbeit am Projekt.

---

*Letzte Aktualisierung: Juni 2026*
*Ouroboros Collective*
