# Buchhaltungsbutler Belegupload

Mit diesem Kommandozeilenprogramm lassen Belege an [Buchhaltungsbutler](https://www.buchhaltungsbutler.de/) übertragen.

Die API-Zugangsdaten müssen in einer .env-Datei hinterlegt werden.

Belege werden nach dem Upload automatisch in den Ordner `.backup` verschoben und nicht gelöscht.

## Beispiel

```bash
git pull https://github.com/rocramer/bhb-receipt-uploader.git
go build

./bhb-receipt-uploader -p=/pfad/zu/den/Ausgangsrechnungen/ -d=outbound
./bhb-receipt-uploader -p=/pfad/zu/den/Eingangsrechnungen/ -d=inbound
``