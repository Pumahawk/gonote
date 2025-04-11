# GoNote

Voglio creare uno strumento che mi permetta di creare note.

Un requisito è che deve salvare le note in un filesystem in modo da permettermi di ricercare con strumenti **Gnu** come **Grep** o **Git**.

Tenendo conto che posso scrivere le note direttamente con **NeoVim** posso sfruttare le sue funzionalità per la gestione dei path e modifica dei file.

Una nota può essere un file in **Markdown** con intestazione oppure un file **yaml** contenente più note.

Il programma **non** ha il compito di scrivere o gestire il contenuto della nota.
Il suo obbiettivo è quello di permettere una facile ricerca delle note.

Il path della nota non deve influenzare l'identificazione della nota. Le note devono poter essere riorganizzate secondo il path senza compromettere la
loro consultazione.

Sono convinto che per prendere note basa scrivere dei file di testo normali all'interno di una cartella.
Molte funzionalità di scrittura e ricerca si possono affidare a utilities Gnu.
Il programma deve permettere una ricerca mirata. Sui metadati delle note.

Ad esempio posso sfruttare una ricerca ricorsiva di grep per fare una ricerca di un campo testuale ma sarebbe molto più
difficile riscure a filtrare in maniera mirata sui metadati di una nota ignorando il resto.

## Links

Le note si possono collegare tra du loro tramite link. Non intendo link in markdown. Ma la possibilità di aggiungere
link nella intestazione cosi da peter navigare tra le note.

Non è possibile spostarsi tra le note usando **NeoVim** perché i percorsi dei file non possono essere usati come riferimento.
Quindi il programma deve avere una funzionalità che mostra i link tra le note.

## Funzionalità del programma

- Lista        - Elenco note.
- Aggregato    - Numero note, tag usati.
- Ricerca      - Ricerca note per metadata o regex.
- Stampa       - Stampa della nota senza doverla aprire con Nvim.
- Modifica     - Apertura della nota tramite l'editor preferito (come kubectl).

**Es**:
- ls --details --tilte --tags --json --tag=xxx,yyy --tag-or zzz,jjj
- stats --tags --tag=xxx
- show --id
- edit id

## Metadati nota

|   Name      |   Description                                                |
| =========== | ============================================================ |
| id          | Identificativo documento. Generato manualmente.              |
| title       | Titolo nota.                                                 |
| tag         | Lista di stringhe.                                           |
| createAt    | Data creazione. Campo libero da inserire manualmente.        |
| updateAt    | Data modifica. Campo libero da inserire manualmente.         |
| meta        | Campo oggetto libero. Può essere utile per altri script.     |
| note        | Contenuto della nota.                                        |

## Esempio di note

```yaml
notes:
- id: "personal-first-note"
  title: "This id my first note"
  tag: ["todo", "status:working"]
  createAt: 2025-04-11T21:26:05+02:00
  updateAt: 2025-04-11T21:36:29+02:00
  meta: {}
  note: |
    Questa è la mia prima nota
```

```markdown
---
notes:
- id: "personal-second-note"
  title: "This id my second note"
  tag: ["todo", "status:done"]
  createAt: 2025-04-11T21:29:51+02:00
  updateAt: 2025-04-11T21:39:58+02:00
  meta: {}
---
# Title my second note

This is my note project.
```
