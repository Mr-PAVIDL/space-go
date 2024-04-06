from collections import defaultdict

from dats_eden_space_client.models import ResPlanetDiff, ResPlayer


class Mapper:
    def __init__(self, universe: ResPlayer):
        self.graph = defaultdict(dict)
        self.planets = set()
        for x, y, fuel in universe.universe:
            self.planets.add(x)
            self.planets.add(y)
            self.graph[x][y] = fuel

    def update(self, diffs: [ResPlanetDiff]):
        for d in diffs:
            self.graph[d.from_][d.to] = d.fuel

    def print_graph(self, cut_far_from='Earth', cut_farther_then=250):
        print('----------graph-----------')
        print()

        bad_planets = set()
        for p in self.planets:
            cost, _ = self.find_path(cut_far_from, p)
            if self.cost(cut_far_from, p) > cut_farther_then:
                bad_planets.add(p)

        for planet_from, others in self.graph.items():
            if planet_from in bad_planets:
                continue
            for planet_to, fuel in others.items():
                if planet_to in bad_planets:
                    continue
                print(planet_from, planet_to, fuel)

        print()
        print('----------graph-----------')

    def cost(self, planet_from, planet_to, seen: set = None):
        if planet_from == planet_to:
            return 0

        if seen is None:
            seen = {planet_from}

        best = 1e99
        for other in self.graph[planet_from]:
            if other in seen:
                continue
            if self.graph[planet_from] is None or self.graph[planet_from][other] is None:
                continue

            seen.add(other)
            from_other = self.cost(other, planet_to, seen) + self.graph[planet_from][other]
            seen.remove(other)
            if from_other < best:
                best = from_other
        return best

    def find_path(self, planet_from, planet_to):
        import heapq

        came_from = {}
        queue = []
        heapq.heappush(queue, (0, planet_from, None))
        while queue:
            cost, node, frm = heapq.heappop(queue)
            if node in came_from:
                continue
            else:
                came_from[node] = frm

            if node == planet_to:
                break

            for other, fuel in self.graph[node].items():
                heapq.heappush(queue, (cost + fuel, other, node))
        if planet_to not in came_from:
            return None

        total_cost = 0
        path = [planet_to]
        while path[-1] != planet_from:
            total_cost += self.graph[came_from[path[-1]]][path[-1]]
            path.append(came_from[path[-1]])
        return total_cost, list(reversed(path))