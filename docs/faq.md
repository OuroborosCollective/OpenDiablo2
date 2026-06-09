# Ouroboros Collective - FAQ

**ARE-Diablo2-BaalAal**

---

## ❓ Allgemein

### Was ist ARE-Diablo2-BaalAal?
ARE-Diablo2-BaalAal ist ein Open-Source-Projekt von **Ouroboros Collective**, das eine moderne Re-Implementation des klassischen Diablo 2 ARPG-Spiels bietet, optimiert für Mobile und Web.

### Wer steht hinter dem Projekt?
**Ouroboros Collective** - eine Community von Entwicklern und Enthusiasten, die das Diablo 2 Erlebnis für moderne Plattformen zugänglich machen.

---

## 🐍 Axiomatic BaalAal Engine

### Was ist die Axiomatic BaalAal Engine?
Die **Axiomatic BaalAal Engine** ist das Herzstück des Projekts:
- **Deterministic State:** Verwendet Kappa-space Koordinaten für reproduzierbare Spielzustände
- **Ouroboros Cycle:** Rekursives Resonanz-System für Stabilität
- **Security:** Deaktivierter `js` Command, nur sichere `ax-*` Operationen

### Was bedeutet "deterministic"?
Deterministisch bedeutet, dass bei gleichen Eingaben immer der gleiche Zustand erreicht wird. Dies ist wichtig für Multiplayer-Synchronisation und Replay-Funktionalität.

---

## 📱 Mobile & Web

### Wird es eine Mobile App geben?
Ja, die Mobile-Migration ist Teil der aktuellen Roadmap. Geplant sind:
- Progressive Web App (PWA)
- iOS/Android native Wrapper

### Welche Geräte werden unterstützt?
- **Mobile:** iOS 14+, Android 10+
- **Desktop:** Windows, macOS, Linux
- **Web:** Chrome, Firefox, Safari, Edge

### Wie funktioniert das Spiel auf Mobile?
Der Web-Client nutzt Touch-Controls:
- Virtual Joystick für Bewegung
- Tap-to-Attack
- On-screen Hotbar für Fähigkeiten

---

## 🛠️ Technisch

### Welche Programmiersprachen werden verwendet?
- **Backend:** Go (Engine, Server)
- **Frontend:** TypeScript, React
- **3D Rendering:** Babylon.js
- **2D Rendering:** Ebiten (Desktop)

### Kann ich zum Projekt beitragen?
Ja! Wir freuen uns über Beiträge. Siehe [CONTRIBUTING.md](./CONTRIBUTING.md).

### Wo finde ich die Dokumentation?
- [Building](./building.md) - Build-Anleitung
- [Development](./development.md) - Entwickler-Anleitung
- [Status](./status.md) - Aktueller Projekt-Status

---

## 🔒 Sicherheit

### Ist das Spiel sicher?
Ja, die **Axiomatic BaalAal Engine** ist von Grund auf sicher konzipiert:
- Keine arbitrary Code Execution
- WebSocket Origin Validierung
- JS Execution Timeouts

---

## 📜 Lizenz

### Unter welcher Lizenz steht das Projekt?
**GPL v3** - Alle derivative Werke müssen unter der gleichen Lizenz veröffentlicht werden.

### Brauche ich Diablo 2?
Ja, Sie benötigen eine legal erworbene Kopie von Diablo 2 und der Erweiterung Lord of Destruction.

---

## 💬 Kontakt

### Wie kann ich Hilfe bekommen?
- **Discord:** [Join our server](https://discord.gg/pRy8tdc)
- **GitHub Issues:** Für Bug-Reports und Feature-Requests

---

*Ouroboros Collective - Juni 2026*
