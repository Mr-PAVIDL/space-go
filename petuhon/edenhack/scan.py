import time

from tqdm.auto import tqdm

from collect import Loader
import json
from dats_eden_space_client import AuthenticatedClient
from mapper import Mapper

client = AuthenticatedClient(base_url="https://datsedenspace.datsteam.dev", token="660c35366abee660c35366abf1",
                             auth_header_name='X-Auth-Token', prefix=None)

from dats_eden_space_client.models import *
from dats_eden_space_client.api.ship import get_player_universe, post_player_travel, post_player_collect
from dats_eden_space_client.api.round_ import get_player_rounds, delete_player_reset
from dats_eden_space_client.types import Response

# res: ResRoundList = get_player_rounds.sync_detailed(client=client).parsed
# print(res)



player: ResPlayer = get_player_universe.sync_detailed(client=client).parsed
mapper = Mapper(player)

with open(f'graph_dump_2.json', 'w') as f:
    json.dump(mapper.graph, f)

def fly_to(orig: str, dest: str) -> ResTravel:
    cost, path = mapper.find_path(orig, dest)
    print(f'Trip {orig} -> {dest} will cost {cost}')
    # input('confirm?')
    travel: ResTravel = post_player_travel.sync_detailed(client=client, body=ReqTravel(planets=[
        *path[1:],
    ])).parsed
    return travel


data = {}

if player.ship.planet.name != 'Earth':
    fly_to(player.ship.planet.name, 'Earth')

prev = 'Earth'
for planet in tqdm(sorted(mapper.planets)[:]):
    resp = fly_to(prev, planet)
    data[planet] = resp.planet_garbage.additional_properties
    prev = planet
    time.sleep(1/4)

print(data)

with open('planet_dump_2.json', 'w') as f:
    json.dump(data, f)
