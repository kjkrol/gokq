# Lista rzeczy do poprawy

*przekreslone zrobione

1. Dla quadtree zostal wprowadzony parametr `maxDepth`. Nie jest on uzywany. Pomysl jest taki, by okreslal on maksymalna ilosc poziomow drzewa. Innymi slowy od pewnego pozioomu elementy maja byc skladowane z pominieciem CAPACITY dla tego drzewa. Trzeba dodac test sprawdzajacy dzialanie.
2. ~~Po wprodzadzeniu w zmian do gokg - np usunieciu DistanctTo ze Spatial, nalezy zakutalizowac pakiet quadtree.~~
3. ~~poprawic `findIntersectingNodesUnique` - sotawilem przy niej notatkie~~
4. ~~Poprawic `FindNeighbors`~~
5. ~~Finder przechowuje wskaźnik do korzenia (…_finder.go:12). Dopóki struktura drzewa nie podmienia root na nowy obiekt, wszystko gra; jeżeli kiedyś wprowadzisz rekonstrukcję root-a (np. przy kompresji), trzeba będzie zsynchronizować finder.root~~
6. ~~W wariancie z wieloma probe’ami nadal przechodzimy te same gałęzie, jeśli fragmenty się pokrywają; mapa usuwa duplikaty elementów, ale nie skraca rekurencji. Można rozważyć memoizację odwiedzonych node’ów, jeśli profilowanie pokaże, że wrapy mocno dublują pracę.~~