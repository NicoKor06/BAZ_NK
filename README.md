<h1>Blog-Applikation Z</h1><h4>(…weil X schon von Twitter belegt ist)</h4>
Die verrückte Bloggerin Berta braucht dringend deine Hilfe! Sie hat einen neuen Blog gestartet 
und braucht nun API-Zauberer, die ihr helfen, ihre Artikel zu veröffentlichen, zu bearbeiten 
und zu löschen. Aber Vorsicht, Berta ist sehr wählerisch und möchte, dass alles nach REST API 
Standard läuft, mit einer Datenbank für die Speicherung der Blog-Einträge, Caching für eine schnelle 
Performance und einer Authentifizierung, damit Benutzer nur Zugriff auf ihre eigenen Beiträge zum Bearbeiten haben
und ein Admin Zugriff auf alles bekommt.
Kannst du die Herausforderung annehmen und Bertas Blog zum Erfolg führen? Zeig uns deine 
GO-Kenntnisse, deine Fähigkeit, komplexe Systeme zu entwickeln und deine Lust auf eine 
spaßige Zusammenarbeit mit der verrückten Bloggerin Berta. Und vergiss nicht, die 
Open API Spezifikationen zu beachten, damit Berta immer weiß, was du gerade treibst. 
Viel Erfolg!

<h2>Tasks</h2>
- Beschäftige dich mit der Programmiersprache GO
- Überlege dir, wie du die Dateien in deinem Projekt strukturierst
- Erstelle die API’s nach der vorgegebenen OpenAPI Spezifikation (openapi.json)
    - Wie ist ein OpenAPI File aufgebaut?
    - Ändere die Routen nach REST-Standard
    - Ergänze die fehlenden Methoden in der Spezifikation (aktuell sind alle Routen auf "GET" gemapped, ist das korrekt? Schau auf die summary, was die Route machen soll.)
    - Es gibt Status-Codes die aktuell noch mit "XXX" bezeichnet sind. Wie wären die korrekten Codes?
    - Ergänze notwendige Status-Codes, die eventuell sinnvoll wären für deinen API-Workflow
    - Schaue dir an, mit welchen GO-Hilfsmitteln (Packages oder natives GO) du das HTTP-Handling umsetzt
- Speichere im 1. Schritt (weil's vorerst einfacher ist) die Blog-Einträge z.B. in einer Runtime-Variable oder in einem JSON-File (und wenn du dich bereit fühlst, in einer Datenbank.. hier kommt Docker ins Spiel)
- Schaue dir an, welche Möglichkeiten es gibt, API-Endpunkte zu authorisieren/authentifizieren
- Informiere dich, wie man ein Caching sinnvoll einsetzen könnte (z.B. mit Redis), um die API -Abfragen zu beschleunigen
- Was ist eine gute Applikation ohne Unit-Tests? Informiere dich über Testing in Go und versuche eine möglich hohe Test-Abdeckung für deinen Code zu generieren
- Stelle Fragen 😊
- Deine Backend-Applikation ist fertig und du willst noch ein schönes Frontend dazu bauen? Dann findest du hier alle Infos:
  https://vue-blueprint.schwarz/

<h2>Nützliche Links</h2>
Nur eine kleine Auswahl - ansonsten ist Google dein Freund

<h4>GO</h4>
- https://app.pluralsight.com/id/signin?redirectTo=https%3A%2F%2Fapp.pluralsight.com%2Flibrary%2Fcourses%2Fgo-fundamentals%2Ftable-of-contents
- https://gobyexample.com/
- https://go.dev/tour

<h4>REST</h4>
- https://www.restapitutorial.com/

<h4>Redis</h4>
- https://github.com/redis/go-redis

<h4>Unit-Testing</h4>
- https://www.digitalocean.com/community/tutorials/how-to-write-unit-tests-in-go-using-go-test-and-the-testing-package

<h4>Datenbank</h4>
- https://tutorialedge.net/golang/golang-mysql-tutorial/

<h4>Docker</h4>
- https://www.ibm.com/topics/docker

<h4>OpenAPI</h4>
- https://swagger.io/specification/
- https://www.ionos.de/digitalguide/websites/web-entwicklung/was-ist-openapi/