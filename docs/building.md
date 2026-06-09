# Ouroboros Collective - Build-Anleitung

**ARE-Diablo2-BaalAal**

---

## 🏗️ Voraussetzungen

- Go 1.21+
- Node.js 18+ (für Web-Client)
- Git

---

## 🔨 Backend (Go Server) bauen

### Alle Module herunterladen
```bash
go mod download
```

### Development Build
```bash
go build -o dev-server .
./dev-server
```

### Production Build
```bash
go build -ldflags="-s -w" -o prod-server .
./prod-server
```

### Mit Tags bauen
```bash
go build -tags="mobile" -o mobile-server .
```

---

## 🌐 Web-Client bauen

### Installation
```bash
cd web-client
npm install
```

### Development Server
```bash
npm run dev
# Öffnet http://localhost:3000
```

### Production Build
```bash
npm run build
npm start
```

### Docker Build
```bash
docker build -t ouroboros/web-client .
docker run -p 3000:3000 ouroboros/web-client
```

---

## 📱 Mobile Builds

### iOS (requires MacOS)
```bash
cd web-client
npx next export
# Export für iOS WebView vorbereiten
```

### Android
```bash
cd web-client
npx next export
# Für Android WebView verpacken
```

---

## 🧪 Tests ausführen

```bash
# Go Tests
go test ./...

# Web-Client Tests
cd web-client
npm test
```

---

## 🐛 Troubleshooting

### Go Build Fehler
```bash
# Cache leeren
go clean -cache
go mod tidy
```

### Node Module Fehler
```bash
cd web-client
rm -rf node_modules
npm install
```

---

## 📦 Abhängigkeiten

| Paket | Version | Beschreibung |
|-------|---------|--------------|
| Go | 1.21+ | Backend |
| Next.js | 14+ | Web Framework |
| Ebiten | latest | 2D Rendering |
| Babylon.js | latest | 3D Rendering |

---

*Ouroboros Collective - Juni 2026*
