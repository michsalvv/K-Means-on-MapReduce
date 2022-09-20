## Note
Script bash con input nome dataset e numero cluster. Istanzia K reducer, 1 master e X mapper.

Nella relazione mettere la foto del flusso presente sulle slide

Realizzare un’architettura master-worker, in cui il master distribuisce il carico di lavoro tra i nodi
worker, che implementano i mapper ed i reducer.
- Quindi non è necessario un DFS sottostante, è il master che distribuisce i chunk come nell'esercizio del grep distribuito
#### Map Phase
Ciascun mapper riceve un **chunk** di punti dell'insieme di partenza e per ciasucno di questi punti calcola la distanza euclidea tra il punt ed i *k* centroidi, identificando cosi il centroide che minimizza la distanza e al cui cluster il punto viene assegnato. 

#### Reduce Phase
Ciascun reduer riceve in parallelo in input tutti i punti assegnati ad un determinato cluster e calcola il valore aggiornato del centroide per quel cluster. 

- Quindi se vogliamo trovare k=4 cluster, dovremmo avere solamente 4 reducer. 

I reducer devono effettuare la remote read dai mapper. I risultati dei mapper sono locali e non tornano al master obv. 
Quando tutti i mapper hanno finito, il master notifica ad ogni reducer da quali mapper devono leggere. 
    - In realtà penso che devono leggere da tutti i mapper per prendersi i punti che ogni mapper ha assegnato al cluster in gestione dal singolo reducer
    - L'importante è che il master non invia dati ai reducer

#### Convergenza algoritmo
Il numero di iterazioni può essere determinato a priori oppure dipendere dalla convergenza, utilizzando come criterio di convergenza la minimizzazione della somma dei quadrati all’interno del cluster


## Architettura cluster distribuito
Per il deployment dell’applicazione, si richiede di usare container Docker, di fornire i relativi file per la creazione delle immagini e di effettuare il deployment dell’applicazione su una istanza EC2 utilizzando il grant AWS a disposizione.

Tanti container quanti sono mapper+worker+master+client

Scrivere che nel docker-compose è stato utilizzato il sistema di environment variables tramite il file .env
    - oppure utilizzare ```bash docker-compose up --scale nome_servizio=numero replica```

## Possibili migliorie
Scale-out automatico del numero di worker:
	- ES: se il dataset ha 500 istanze, tramite un tasso di splitting deciso dall'utente verranno avviati 5 mapper ed il dataset splittato in 5 chunk


## Gestione failure
https://yunuskilicdev.medium.com/distributed-mapreduce-algorithm-and-its-go-implementation-12273720ff2f


# TODO
- Quando crasha il master comunicarlo ai worker collegati (non è richiesto)
- go env -w GO111MODULE=off
- Ogni worker quando si connette deve comunicare IP e porta esposta dal container. Può essere fatto facendo settare delle variabili d'ambiente dal docker-compose e lette poi successivamente in Go.
	- https://stackoverflow.com/questions/64717467/how-to-read-linux-environment-variables-in-go
	- Forse anche nell'esempio dell'implementazione del map reduce fa qualcosa del genere ma senza virtualizzazione
