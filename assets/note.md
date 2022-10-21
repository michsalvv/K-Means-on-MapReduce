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

- spawnare 1 o più worker aggiuntivi che resteranno in idle. Nel caso in cui un mappere non finisce il job entro un tot, viene considerato failed e il master assegnerà il suo chunk ad un nuovo worker. Nel frattempo avvierà una procedura di ping per vedere se il mapper crashato è tornato online. Nel caso lo utilizzerà come worker idle. 


# TODO
- Quando crasha il master comunicarlo ai worker collegati (non è richiesto)
- go env -w GO111MODULE=off
- Nella relazione scrivere che è disponibile il log sul master
- Mettere anche screen dell'esecuzione nella relazione
- Nelle slide quando spieghi una funzione metti il diagramma di flusso in cui viene usata
- fare versione migliorata con combiner
- nelle slide mettere anche una prova che con solo 2 iterazioni non funziona
- Troppo poca randomica la selezione dei centroidi iniziali, vedi traccia ci sta qualche info utile
- Mettere tutto più generale, provare con k>3
- Migliorare stampe nei nodi
- Deployare su AWS 
  - Forse su AWS ogni eseguibile deve essere un servizio
- Nella relazione scrivere come il client comunica il path del dataset al master e come il master lo cerca nella directory
  - scrivere anche che comunque è un aspetto secondario poichè il servizio è offerto in modo distribuito solamente in locale

- usare sklearn per trovare i veri cluster per confrontare risultati
- nei risultati escludere gli outlier
  - salvare i risultati in csv [tempo esecuzione; numero iterazioni]

MIGLIORARE CONVERGENZA


## Discussione su random initiliazionation
When initializing the centroids, it is important that the initially selected points are fairly apart. If the points are too close together, there is a good chance the points will find a cluster in their local region and the actual cluster will be blended with another cluster.

When randomly initializing centroids, we have no control over where their initial position will be.

Unfortunately, Lloyd does suffer greatly from a problem many clustering algorithms face. That is, that the forming of clusters is heavily influenced by the choosing of initial centroids. Or in more mathematical terms. Lloyd might get stuck in a local optima, depending on its initialization.

L'inizializzazione tramite kMeans++ migliora ma non troppo, fai vedere differenze con inizializzione randomica. 
Secondo me funziona poco bene sui dataset piccoli perchè se la selezione randomica è particolarmente sfortunata, l'algoritmo avrà pochi punti per convergere e quindi tenterà a oscillare vicino ad una falsa convergenza. 