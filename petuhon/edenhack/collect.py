from collections import defaultdict

import colorama

from dats_eden_space_client.api.ship import post_player_collect
from dats_eden_space_client.models import ResTravelPlanetGarbage, ReqCollect, ReqCollectGarbage
from colorama import Fore, Back, Style

from colorama import just_fix_windows_console

just_fix_windows_console()


class Loader:
    def __init__(self, width, height, garbage):
        self.width = width
        self.height = height
        self.trunk = [[None] * width for _ in range(height)]
        for name, cells in garbage.items():
            for cell in cells:
                self.trunk[cell[1]][cell[0]] = name
        self.print()

    def print(self):

        chunk_ids = set([x for xs in self.trunk for x in xs])
        colors = list(vars(colorama.Fore).values()) * 10
        chunk2color = {chunk: color for chunk, color in zip(chunk_ids, colors)}
        from string import ascii_letters
        chunk2letter = {chunk: letter for chunk, letter in zip(chunk_ids, ascii_letters + ascii_letters.upper())}
        print('┌' + '-' * (self.width * 2 + 1) + '┐')
        for row in self.trunk:
            print('| ', end='')
            for cell in row:
                if cell is None:
                    print('.', end=' ')
                else:
                    print(chunk2color[cell] + chunk2letter[cell] + colorama.Fore.RESET, end=' ')
            print('|')
        print('└' + '-' * (self.width * 2 + 1) + '┘')

    def load(self, garbage):
        garbage = sorted(garbage.items(), key=lambda pair: len(pair[1]))
        cnt = 0
        for name, cells in garbage:
            if self.try_fit(name, cells):
                cnt += 1
        return cnt

    def to_dict(self):
        data = defaultdict(lambda: [])
        for y, row in enumerate(self.trunk):
            for x, cell in enumerate(row):
                if cell is not None:
                    data[cell].append((x, y))
        return data

    def try_fit(self, name, cells):
        for y in range(self.height):
            for x in range(self.width):
                if self.can_fit(cells, x, y):
                    for cell in cells:
                        self.trunk[y + cell[1]][x + cell[0]] = name
                    print(f'can fit {name}!')
                    return True
        return False

    def can_fit(self, cells, x, y):
        for cell in cells:
            if y + cell[1] >= self.height or x + cell[0] >= self.width:
                return False
            if self.trunk[y + cell[1]][x + cell[0]] is not None:
                return False

        return True


if __name__ == '__main__':
    loader = Loader(width=5, height=6, garbage={'x': [(0, 0)]})
    loader.print()
    for i in range(3):
        loader.trunk[1][1 + i] = 'abfadassd'
    for i in range(3):
        loader.trunk[2 + i][2] = 'sfdadf'
    loader.print()
    print(loader.can_fit([(0, 0)], 0, 0))
    print(loader.can_fit([(0, 0)], 1, 1))
    print(loader.can_fit([(0, 0), (1, 0)], 0, 1))

    loader.load({'1': [(0, 0), (1, 0), (2, 0)], '2': [(0, 0), (0, 1), (1, 0), (1, 1)]})
    loader.print()
