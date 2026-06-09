# Ouroboros Collective - Beitragsrichtlinien

**ARE-Diablo2-BaalAal**

---

## 🎯 Willkommen!

Wir freuen uns über Beiträge zur **ARE-Diablo2-BaalAal** Engine! Dieses Projekt ist ein Open-Source-Vorhaben von **Ouroboros Collective**, das darauf abzielt, das klassische Diablo 2 Erlebnis auf modernen Plattformen zugänglich zu machen.

---

## 🚀 Erste Schritte

### Voraussetzungen
- Go 1.21+
- Node.js 18+
- Git

### Setup
```bash
# Repository klonen
git clone https://github.com/OuroborosCollective/OpenDiablo2.git
cd OpenDiablo2

# Dependencies installieren
go mod download

# Lokalen Server starten
go run .
```

---

## 🐍 Axiomatic BaalAal Engine

Das Projekt nutzt die **Axiomatic BaalAal Engine** - eine deterministische, logikbasierte Engine:

- **Deterministic State:** Kappa-space Koordinaten
- **Recursive Logic:** Ouroboros Cycle System
- **Security:** Deaktivierter `js` Command, nur `ax-*` Operationen

### Neue Features entwickeln
```bash
# Feature-Branch erstellen
git checkout -b feature/axiomatic-improvement

# Änderungen vornehmen
# ...

# Commit und Push
git add .
git commit -m "feat: add axiomatic improvement"
git push origin feature/axiomatic-improvement
```

---

## 🧪 Testing

Wir erwarten, dass alle neuen Features mit Tests abgedeckt werden:

```bash
# Alle Tests ausführen
go test ./...

# Spezifische Tests
go test ./d2script/...
go test ./d2common/...
```

---

## 📝 Code-Stil

- **Go:** Wir verwenden `golangci-lint` für Linting
- **TypeScript/React:** ESLint + Prettier
- **Commits:** Clear, descriptive messages

### Linting aktivieren
```bash
# Go Linter installieren
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Linting ausführen
golangci-lint run
```

---

## 🔄 Pull Request Prozess

1. Fork erstellen
2. Feature-Branch erstellen (`git checkout -b feature/xyz`)
3. Änderungen vornehmen und testen
4. Commit mit klarer Beschreibung
5. PR erstellen auf GitHub

### PR Requirements
- ✅ Tests für neue Features
- ✅ Keine Lint-Fehler
- ✅ Klare Commit-Beschreibung
- ✅ Referenz zu relevanten Issues

---

## 🐛 Issues melden

- Verwende GitHub Issues für Bug-Reports
- Beschreibe das Problem detailliert
- Füge Stack-Traces und Reproduktionsschritte bei
- Setze Labels (bug, enhancement, security)

---

## 📜 Lizenz

Mit dem Beitrag zum Projekt stimmen Sie zu, dass Ihr Code unter **GPL v3** lizenziert wird.

---

## 💬 Kontakt

- **Discord:** [Join our server](https://discord.gg/pRy8tdc)
- **GitHub:** [OuroborosCollective](https://github.com/OuroborosCollective)

---

*Vielen Dank für Ihre Unterstützung!*
*Ouroboros Collective - Juni 2026*
