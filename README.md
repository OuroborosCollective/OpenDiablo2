# Ouroboros Collective - ARE-Diablo2-BaalAal

![Logo](d2logo.png)

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Discord](https://img.shields.io/discord/515518620034662421?label=Discord&style=social)](https://discord.gg/pRy8tdc)

---

## 🎮 Projektübersicht

**ARE-Diablo2-BaalAal** ist ein Open-Source-Re-Implementation des klassischen ARPG-Spiels Diablo 2, entwickelt von **Ouroboros Collective**.

Das Projekt nutzt die **Axiomatic BaalAal Engine** - eine deterministische, logikbasierte Engine, die für Mobile und Web optimiert ist.

> **Hinweis:** Das Projekt benötigt die originalen Diablo 2 Spieldateien. Sie benötigen eine legal erworbene Kopie von [Diablo 2](https://us.shop.battle.net/en-us/product/diablo-ii) und der Erweiterung [Lord of Destruction](https://us.shop.battle.net/en-us/product/diablo-ii-lord-of-destruction).

---

## 🐍 Axiomatic BaalAal Engine

Die **Axiomatic BaalAal Engine** ist das Herzstück unseres Projekts:

- **Deterministic State:** Kappa-space Koordinaten (1000-base integer scaling)
- **Recursive Logic:** Ouroboros "selbst-fressende Schlange" Zyklus für Resonanz-Stabilität
- **Security:** Arbitrary Code Execution via `js` ist deaktiviert - nur axiomatic operations (`ax-*`)
- **Mobile-First:** Optimiert für Touch-Controls und mobile Displays

---

## 📋 Projekt-Status (Juni 2026)

### ✅ Erfolgreich integrierte Änderungen:

**🔒 Sicherheit:**
- WebSocket Origin Validierung implementiert
- JS Execution Timeouts konfiguriert  
- Lokaler WebSocket-Server als Standard

**⚡ Performance:**
- Tile Lookup und Caching optimiert
- Escape Menu UI State Management verbessert
- Ebiten Renderer Verbesserungen

**🧪 Tests:**
- Unit Tests für stringutils Utility Funktionen
- Tests für Ext2SourceType und Ext2AssetType

**🧹 Code-Qualität:**
- NewGameControls Refactoring
- Gamma und Contrast Configuration in Ebiten Renderer
- 3D Audio Bias Implementation

**🌐 Mobile Migration:**
- Web-Client mit Next.js Framework
- Mobile-freundliche Asset-Metadaten-Anzeige
- Touch-Controls vorbereitet

---

## 🚀 Schnellstart

### Voraussetzungen
- Go 1.21+
- Node.js 18+
- MPQ-Dateien von Diablo 2

### Backend (Go Server)
```bash
go build -o server .
./server
```

### Web-Client (Next.js)
```bash
cd web-client
npm install
npm run dev
```

---

## 📁 Projektstruktur

```
ouroboros-collective/
├── d2app/           # Hauptanwendung
├── d2common/        # Gemeinsame Utilities
├── d2core/          # Core Engine
├── d2game/          # Game Logik
├── d2networking/    # Netzwerk & WebSocket
├── d2script/        # Axiomatic Engine (BaalAal)
├── d2thread/        # Threading
├── web-client/      # Mobile/Web Frontend (Next.js)
└── docs/            # Dokumentation
```

---

## 🔥 Für Entwickler

* [Building](./docs/building.md) - Build-Anleitung
* [Development](./docs/development.md) - Entwickler-Anleitung
* [Contributing](./docs/CONTRIBUTING.md) - Beitragsrichtlinien
* [Roadmap](./docs/roadmap.md) - Projekt-Roadmap

---

## 📜 Lizenz

Dieses Projekt ist lizenziert unter **GPL v3**.

Diablo 2 und seine Inhalte sind ©2000 Blizzard Entertainment, Inc. Alle Rechte vorbehalten.

---

*Ein Projekt von [Ouroboros Collective](./docs/CONTRIBUTING.md)*
