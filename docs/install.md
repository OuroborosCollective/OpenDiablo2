# Ouroboros Collective - Installationsanleitung

**ARE-Diablo2-BaalAal**

---

## 📋 Voraussetzungen

- Go 1.21+
- Node.js 18+
- MPQ-Dateien von Diablo 2

---

## 🔧 Installation

### 1. Repository klonen
```bash
git clone https://github.com/OuroborosCollective/OpenDiablo2.git
cd OpenDiablo2
```

### 2. Go Dependencies installieren
```bash
go mod download
```

### 3. Web-Client installieren
```bash
cd web-client
npm install
cd ..
```

---

## 🎮 Spiel starten

### Backend (Go Server)
```bash
go run .
# Server startet auf http://localhost:8080
```

### Web-Client (separate Terminal)
```bash
cd web-client
npm run dev
# Client startet auf http://localhost:3000
```

---

## 📱 Mobile Installation

### Web App (PWA)
1. Öffne `http://localhost:3000` im Mobile Browser
2. Tippe auf "Zum Home-Bildschirm hinzufügen"
3. Die App erscheint auf dem Home-Screen

### Native Apps (geplant)
- iOS App Store (Q4 2026)
- Google Play Store (Q4 2026)

---

## 🐛 Fehlerbehebung

### Go Modulfehler
```bash
go clean -cache
go mod tidy
```

### Node Module Fehler
```bash
cd web-client
rm -rf node_modules package-lock.json
npm install
```

---

## 📁 MPQ-Dateien

Das Spiel benötigt die originalen Diablo 2 MPQ-Dateien:
- `d2exp.mpq`
- `d2music.mpq`
- `d2sfx.mpq`
- `d2data.mpq`
- `d2char.mpq`
- `d2vid.mpq`

Siehe [MPQ](./mpq.md) für weitere Details.

---

*Ouroboros Collective - Juni 2026*
