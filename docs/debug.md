# Ouroboros Collective - Debugging

**ARE-Diablo2-BaalAal**

---

## 🐛 Häufige Fehler

### Server startet nicht
```bash
# Prüfe Port-Verfügbarkeit
lsof -i :8080

# Starte mit Debug-Log
LOG_LEVEL=debug go run .
```

### WebSocket Verbindungsfehler
```bash
# CORS Einstellungen prüfen
# WebSocket Port prüfen (Standard: 8080)
```

---

## 🔧 Debug Commands

### Go Backend
```bash
# Hilfe
go run . --help

# Debug Mode
DEBUG=1 go run .
```

---

*Ouroboros Collective - Juni 2026*
