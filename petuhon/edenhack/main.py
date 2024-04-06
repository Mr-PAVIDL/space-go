import time

from tqdm.auto import tqdm

from collect import Loader
from dats_eden_space_client import AuthenticatedClient
from mapper import Mapper

# client = AuthenticatedClient(base_url="https://datsedenspace.datsteam.dev", token="660c35366abee660c35366abf1",
client = AuthenticatedClient(base_url="http://localhost:3333", token="660c35366abee660c35366abf1",
                             auth_header_name='X-Auth-Token', prefix=None)

from dats_eden_space_client.models import *
from dats_eden_space_client.api.ship import get_player_universe, post_player_travel, post_player_collect
from dats_eden_space_client.api.round_ import get_player_rounds, delete_player_reset
from dats_eden_space_client.types import Response

res: ResRoundList = get_player_rounds.sync_detailed(client=client).parsed
print(res)


# uni.ship
DELAY = 1/4


# mapper.print_graph()

def fly_to(mapper: Mapper, orig, dest, grab_garbage=False):
    cost, path = mapper.find_path(orig, dest)
    print(f'Trip {orig} -> {dest} will cost {cost}')
    input('confirm?')
    travel: ResTravel = post_player_travel.sync_detailed(client=client, body=ReqTravel(planets=[
        *path[1:],
    ])).parsed
    mapper.update(travel.planet_diffs)
    time.sleep(DELAY)
    if grab_garbage:
        # print(travel)
        loader = Loader(player.ship.capacity_x, player.ship.capacity_y, travel.ship_garbage.additional_properties)
        loaded = loader.load(travel.planet_garbage.additional_properties)
        if loaded == 0:
            return False

        loader.print()

        # input('confirm?')
        collect_resp = post_player_collect.sync_detailed(client=client,
                                                         body=ReqCollect(
                                                             ReqCollectGarbage.from_dict(loader.to_dict()))).parsed
        time.sleep(DELAY)
        print(collect_resp)
        # check if no garbage is left
        return loaded != len(travel.planet_garbage.additional_properties)
    return False


# loader.load(travel.planet_garbage)

# for d in travel.planet_diffs:


if __name__ == '__main__':
    # delete_player_reset.sync_detailed(client=client)
    # 0/0
    player: ResPlayer = get_player_universe.sync_detailed(client=client).parsed
    mapper = Mapper(player)
    mapper.print_graph()
    print(player.ship.planet.name)

    # 0/0

    home = 'Eden'
    if player.ship.planet.name != home:
        fly_to(mapper, player.ship.planet.name, home)

    planets = sorted(mapper.planets)
    while planets:
        for planet in tqdm(planets[:]):
            if planet == home:
                continue

            if not fly_to(mapper, home, planet, True):
                planets.remove(planet)
                print('remove planet ', planet)
            fly_to(mapper, planet, home, False)

    # fly_to('Conn', True)
    # fly_to('Eden', False)
