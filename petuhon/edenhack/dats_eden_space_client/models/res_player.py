from typing import TYPE_CHECKING, Any, Dict, List, Type, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.res_player_universe_item_item import ResPlayerUniverseItemItem
    from ..models.res_ship import ResShip


T = TypeVar("T", bound="ResPlayer")


@_attrs_define
class ResPlayer:
    """
    Attributes:
        name (Union[Unset, str]):  Example: MyTeam.
        ship (Union[Unset, ResShip]):
        universe (Union[Unset, List[List['ResPlayerUniverseItemItem']]]):  Example: [['Earth', 'Reinger', 100],
            ['Reinger', 'Earth', 100], ['Reinger', 'Larkin', 100]].
    """

    name: Union[Unset, str] = UNSET
    ship: Union[Unset, "ResShip"] = UNSET
    universe: Union[Unset, List[List["ResPlayerUniverseItemItem"]]] = UNSET
    additional_properties: Dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> Dict[str, Any]:
        name = self.name

        ship: Union[Unset, Dict[str, Any]] = UNSET
        if not isinstance(self.ship, Unset):
            ship = self.ship.to_dict()

        universe: Union[Unset, List[List[Dict[str, Any]]]] = UNSET
        if not isinstance(self.universe, Unset):
            universe = []
            for universe_item_data in self.universe:
                universe_item = []
                for universe_item_item_data in universe_item_data:
                    universe_item_item = universe_item_item_data.to_dict()
                    universe_item.append(universe_item_item)

                universe.append(universe_item)

        field_dict: Dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if name is not UNSET:
            field_dict["name"] = name
        if ship is not UNSET:
            field_dict["ship"] = ship
        if universe is not UNSET:
            field_dict["universe"] = universe

        return field_dict

    @classmethod
    def from_dict(cls: Type[T], src_dict: Dict[str, Any]) -> T:
        from ..models.res_player_universe_item_item import ResPlayerUniverseItemItem
        from ..models.res_ship import ResShip

        d = src_dict.copy()
        name = d.pop("name", UNSET)

        _ship = d.pop("ship", UNSET)
        ship: Union[Unset, ResShip]
        if isinstance(_ship, Unset):
            ship = UNSET
        else:
            ship = ResShip.from_dict(_ship)

        universe = []
        _universe = d.pop("universe", UNSET)
        for universe_item_data in _universe or []:
            universe_item = []
            _universe_item = universe_item_data
            for universe_item_item_data in _universe_item:
                universe_item_item = universe_item_item_data#ResPlayerUniverseItemItem.from_dict(universe_item_item_data)

                universe_item.append(universe_item_item)

            universe.append(universe_item)

        res_player = cls(
            name=name,
            ship=ship,
            universe=universe,
        )

        res_player.additional_properties = d
        return res_player

    @property
    def additional_keys(self) -> List[str]:
        return list(self.additional_properties.keys())

    def __getitem__(self, key: str) -> Any:
        return self.additional_properties[key]

    def __setitem__(self, key: str, value: Any) -> None:
        self.additional_properties[key] = value

    def __delitem__(self, key: str) -> None:
        del self.additional_properties[key]

    def __contains__(self, key: str) -> bool:
        return key in self.additional_properties
