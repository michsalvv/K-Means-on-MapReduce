## Note

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