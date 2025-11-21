# QuadTree

Drzewo grafowe majace zastosowanie w przeszukiwaniu przestrzeni dwuwymiarowej.
Formalnie zdefiniowane w ten spoób.
$$
Q =
\begin{cases}
\text{leaf}(S), & \text{jeśli } |S| \le L \\
(Q_{NE}, Q_{NW}, Q_{SE}, Q_{SW}), & \text{w przeciwnym razie}
\end{cases}
$$

Co mozna przeczytac jako:
Jezeli wezel zawiera wiecej niz L elementow, to stworz 4 węzły potomne i dystybuj elementy do wezlow potomnych.

## Przeszukiwanie punktow na plaszczyznie
Wyobrazmy sobie, ze na pewnej plaszczyznie mamy N punktow.
Chcielibysmy sprawdzic czy punkty koliduja, ze soba lub sa w zadanej odleglosci `d` od siebie. Wówczas dla kazdej pary ${(X_1, Y_1), (X_2, Y_2)}$ musielibysmy wykonac nastepujace sprawdzenie:

$$
\sqrt{(X_1-X_2)^2+(Y_1-Y_2)^2} < d
$$

Aby to zrobic musimy porownac
wszystkie mozliwe pary uporządkowane, gdzie powtórzenia sa dozwolone. Ujmując to matematcznie: "szukamy liczby wariancji 2 elementowcyh na zbiorze N elementowym".
Liczba waraincje z powtorzeniami zdana jest wzorem:

$$
N^k
$$

gdzie `N` liczba dostepnych elementow, `k` dlugosc sekwencji.

Reasumujac:
Naiwne podejscie do problemu wymagalo by sprawdzenia kazdej mozliwej pary obiektow, czyliL
$N^2$ aplikacji sprawdzenia odległości.

Optymalniej byloby porownywac kazda pare tylko raz - bo kolejnosc nie ma znacznia  („kolizja A z B to to samo co B z A”). Czyli szukamy liczbe kombinacji bez poworzen:

$$
\binom{N}{k}
$$

W rozwazanym przypadku szukamy "liczy mozliwych sposobow wybrania 2 elementów ze zbioru N różnych elementów,
bez możliwości powtarzania elementów i bez uwzględniania kolejności."

$$
\binom{N}{2} = \frac{N!}{2!(N-2)!} = \frac{N(N-1)}{2}
$$

Otrzymaliśmy wynik o połowę mniejszy, ale nadal złzonosc algorytmiczna jest rzedu $O(N^2)$

## Jak Quadtree usprawnia przeszukiwanie?

### Quadtree przechowujace punkty plaszczyzny
Zalozmy, ze korzen drzewa reprezentuje AABB (Axis-align bounding box) plaszczyzny.
Wprowadzmy reguly dla podzialu wezla:
- AABB wezla $Q$ dzieli sie na 4 rowne AABB reprezentujace fragmenty rodzica, tzn $\{Q_{NE}, Q_{NW}, Q_{SE}, Q_{SW}\}$
- kazdy z elementow jest przenoszony do odpowiedniego nowo powstalego liscia, w taki sposob ze zawieraja sie w jego AABB

Ponizsze ilustracje prezentuja przykladowe punkty na plaszczyznie umieszczone w Quadtree. W miare dodawania punktow obszar jest dzielony na mniejsze fragmenty.
<div style="display: flex; justify-content: center; gap: 4%; align-items: stretch; margin: 0 auto;">
  <figure style="flex: 1; margin: 0; text-align: left; display: flex; flex-direction: column;">
    <img src="quadtree_plane.svg" style="width: 100%; margin: 0; padding: 0;" />
    <figcaption style="flex: 1; margin-top: 0.5rem;">Podział płaszczyzny na fragmenty i ich kolejne podziały; kolorowane ramki pokazują granice fragmentów, a punkty leżą w ich AABB. To tu dzieje się „geometria” quadtree.
    <br>Tam gdzie wieksze zageszczenie punktow tam silniejsza fragmentacja obszaru, np obszar SE zawiera 1 punkt - nie jest wymagana dalsza fragmentacja. Natomiast w obszarze NW nastapila fragmentacja,
a nastepnie kolejna fragmentacja NW.SE.</figcaption>
  </figure>
  <figure style="flex: 1; margin: 0; text-align: left; display: flex; flex-direction: column;">
    <img src="quadtree_graph.svg" style="width: 100%; margin: 0; padding: 0;" />
    <figcaption style="flex: 1; margin-top: 0.5rem;">Reprezentacja grafowa tego samego podziału: węzły z etykietami ilustrują fragmenty(np. NW.SE), a linie ilustruja historie fragmentacji = strukture dziedziczenia wezlow drzewa</figcaption>
  </figure>
</div>

### Przeszukiwanie Quadtree
Jezeli wszystkie punkty nalezace do plaszczyzny umiescimy w Quadtree, to mozemy skorzystac z wlasnosci struktury drzewa w celu znacznego usprawnienia przeszukiwania.

<div style="display: flex; justify-content: center; gap: 4%; align-items: stretch; margin: 0 auto;">
  <figure style="flex: 1; margin: 0; text-align: left; display: flex; flex-direction: column;">
    <img src="quadtree_search_plane.svg" style="width: 100%; margin: 0; padding: 0;" />
    <figcaption style="flex: 1; margin-top: 0.5rem;">Dla wybranego punktu, zaznaczonego kolorem pomaranczowym, szukamy sasiadow znajdujacych sie w otoczeniu oznaczonym pomaranczowym kwadaratem. Szukamy punktow wewnatrz zadanego AABB, o centrum w wybranym punkcie x.</figcaption>
  </figure>
  <figure style="flex: 1; margin: 0; text-align: left; display: flex; flex-direction: column;">
    <img src="quadtree_search_graph.svg" style="width: 100%; margin: 0; padding: 0;" />
    <figcaption style="flex: 1; margin-top: 0.5rem;">Przeszukiwanie drzewa polega na sprawdzeniu czy zadany AABB przecina sie z AABB wezla.
    <br>Dla omawianego przypadku:
    <ul>
      <li>dla korzenia uzyskujemy wynik pozytywny</li>
      <li>sprawdzamy jego dzieci, te dll ktorych uzyksano wynik pozytywny obrysowano pomaranczowym kwadratem</li>
      <li>rekurencyjnie kontynuujemy sprawdzanie dzieci wezlow, dla ktorych uzykalismy wynik pozytywny, az dojdziemy do liscia
    </ul>
    </figcaption>
  </figure>
</div>

#### Algorytm przeszukiwania Quadtree

```
Search(Q, queryAabb):
  if Q.aabb nie przecina queryAabb:
    return ∅

  if Q jest liściem:
    zwróć wszystkie punkty p ∈ Q.points, dla których p ∈ queryAabb

  // węzeł wewnętrzny – sprawdzamy tylko dzieci, których AABB przecina zapytanie
  result ← ∅
  dla dziecka C w {Q_NE, Q_NW, Q_SE, Q_SW}:
    if C.aabb przecina queryAabb:
      result ← result ∪ Search(C, queryAabb)
  return result
```

W skrócie: odrzucamy całe poddrzewa, których AABB nie przecina obszaru zapytania; w węzłach liści filtrujemy punkty testem należenia do AABB. Dzięki temu liczba sprawdzeń maleje do poddrzew faktycznie przecinanych przez obszar wyszukiwania.

**Złożoność przeszukiwania.**

Jezeli $N$ to liczba punktow.
Maksymalna liczba node’ów przecinanych przez AABB na jednym poziomie to 4.

Wysokosc drzewa to $h = \log_{4}{N}$

Liczb nodow odwiedzanych w calej operacji $\leq 4*h = 4\log_{4}{N}$

Czyli czas na przejscie calego drzewa to:
$$
\leq O(4\log_{4}N) 
$$

Uwzgledniajac, ze trzeba wykonac testy odleglosci dla kazdego z potencjalnych 4 punktow przechowywanych w kazdym lisciu:

$$
O(4\log_{4}{N} + 4)
$$

Odrzucajac stale (czyli liczbe 4) ostatecznie zlozonosc algorytmiczna przeszukiwania drzewa wynosi:

$$
\log{N}
$$

Co stanowi znakomicie lepszy wynik niz bezmyslne porownywanie kazdej pary punkow, kotre przypomnijmy mialo zlozonosc $O(N^2)$