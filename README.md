# Buchhaltungsbutler Belegupload

Mit diesem Kommandozeilenprogramm lassen Belege eines bestimmten Ordners an [Buchhaltungsbutler](https://www.buchhaltungsbutler.de/) übertragen.

Die API-Zugangsdaten werden in der .env-Datei hinterlegt.

Belege werden nach dem Upload automatisch in den Ordner `.backup` (im selben Verzeichnis) verschoben und nicht gelöscht.

## Beispiel

```bash
# Download und build
git pull https://github.com/rocramer/bhb-receipt-uploader.git
go build

# Ausführen
./bhb-receipt-uploader -p /pfad/zu/den/Ausgangsrechnungen/ -d outbound
./bhb-receipt-uploader -p /pfad/zu/den/Eingangsrechnungen/ -d inbound
```

## Cronjob 
Zur täglichen (oder häufigeren) Synchronisation kann ein Cronjob eingerichtet werden. Folgender Eintrag überträgt die Belege täglich um 02:00 Uhr (Pfadangaben anpassen!):

```
0 2 * * * ./bhb-receipt-uploader -p /pfad/zu/den/Ausgangsrechnungen/ -d outbound
0 2 * * * ./bhb-receipt-uploader -p /pfad/zu/den/Eingangsrechnungen/ -d inbound
``