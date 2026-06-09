# Ouroboros Collective - Web Client Agent Guide

**ARE-Diablo2-BaalAal - Mobile/Web Frontend**

---

## 🎯 Projekt-Ziel

Der Web-Client ist das mobile-optimierte Frontend für die **Axiomatic BaalAal Engine**. Er ermöglicht das Spielen von Diablo 2 über den Browser.

---

## 🐍 Axiomatic Integration

### WebSocket Kommunikation
- Echtzeit-Sync mit Go Backend
- Kappa-space Koordinaten für deterministische State
- Resonanz-System für Stabilität

### Mobile-Optimierung
- Touch-Controls (Virtual Joystick, Hotbar)
- Responsive UI
- PWA Support

---

## 📁 Struktur

```
web-client/
├── app/           # Next.js App Router
├── components/    # React Komponenten
├── public/        # Statische Assets
└── utils/         # Hilfsfunktionen
```

---

## 🔧 Entwicklung

```bash
npm install
npm run dev   # Development
npm run build # Production
```

---

*Ouroboros Collective - Juni 2026*
