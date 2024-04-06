from typing import TYPE_CHECKING, Any, Dict, List, Type, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.res_planet_diff import ResPlanetDiff
    from ..models.res_travel_planet_garbage import ResTravelPlanetGarbage
    from ..models.res_travel_ship_garbage import ResTravelShipGarbage


T = TypeVar("T", bound="ResTravel")


@_attrs_define
class ResTravel:
    """
    Attributes:
        fuel_diff (Union[Unset, int]):  Example: 1000.
        planet_diffs (Union[Unset, List['ResPlanetDiff']]):
        planet_garbage (Union[Unset, ResTravelPlanetGarbage]):  Example: {'6fTzQid': [[0, 0], [0, 1], [1, 1]],
            'RVnTkM59': [[0, 0], [0, 1], [1, 1], [2, 1], [1, 2]]}.
        ship_garbage (Union[Unset, ResTravelShipGarbage]):  Example: {'71B2XMi': [[2, 10], [2, 9], [2, 8], [3, 8]]}.
    """

    fuel_diff: Union[Unset, int] = UNSET
    planet_diffs: Union[Unset, List["ResPlanetDiff"]] = UNSET
    planet_garbage: Union[Unset, "ResTravelPlanetGarbage"] = UNSET
    ship_garbage: Union[Unset, "ResTravelShipGarbage"] = UNSET
    additional_properties: Dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> Dict[str, Any]:
        fuel_diff = self.fuel_diff

        planet_diffs: Union[Unset, List[Dict[str, Any]]] = UNSET
        if not isinstance(self.planet_diffs, Unset):
            planet_diffs = []
            for planet_diffs_item_data in self.planet_diffs:
                planet_diffs_item = planet_diffs_item_data.to_dict()
                planet_diffs.append(planet_diffs_item)

        planet_garbage: Union[Unset, Dict[str, Any]] = UNSET
        if not isinstance(self.planet_garbage, Unset):
            planet_garbage = self.planet_garbage.to_dict()

        ship_garbage: Union[Unset, Dict[str, Any]] = UNSET
        if not isinstance(self.ship_garbage, Unset):
            ship_garbage = self.ship_garbage.to_dict()

        field_dict: Dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if fuel_diff is not UNSET:
            field_dict["fuelDiff"] = fuel_diff
        if planet_diffs is not UNSET:
            field_dict["planetDiffs"] = planet_diffs
        if planet_garbage is not UNSET:
            field_dict["planetGarbage"] = planet_garbage
        if ship_garbage is not UNSET:
            field_dict["shipGarbage"] = ship_garbage

        return field_dict

    @classmethod
    def from_dict(cls: Type[T], src_dict: Dict[str, Any]) -> T:
        if src_dict['shipGarbage'] is None:
            src_dict['shipGarbage'] = {}

        from ..models.res_planet_diff import ResPlanetDiff
        from ..models.res_travel_planet_garbage import ResTravelPlanetGarbage
        from ..models.res_travel_ship_garbage import ResTravelShipGarbage

        d = src_dict.copy()
        fuel_diff = d.pop("fuelDiff", UNSET)

        planet_diffs = []
        _planet_diffs = d.pop("planetDiffs", UNSET)
        for planet_diffs_item_data in _planet_diffs or []:
            planet_diffs_item = ResPlanetDiff.from_dict(planet_diffs_item_data)

            planet_diffs.append(planet_diffs_item)

        _planet_garbage = d.pop("planetGarbage", UNSET)
        planet_garbage: Union[Unset, ResTravelPlanetGarbage]
        if isinstance(_planet_garbage, Unset):
            planet_garbage = UNSET
        else:
            planet_garbage = ResTravelPlanetGarbage.from_dict(_planet_garbage)

        _ship_garbage = d.pop("shipGarbage", UNSET)
        ship_garbage: Union[Unset, ResTravelShipGarbage]
        if isinstance(_ship_garbage, Unset):
            ship_garbage = UNSET
        else:
            ship_garbage = ResTravelShipGarbage.from_dict(_ship_garbage)

        res_travel = cls(
            fuel_diff=fuel_diff,
            planet_diffs=planet_diffs,
            planet_garbage=planet_garbage,
            ship_garbage=ship_garbage,
        )

        res_travel.additional_properties = d
        return res_travel

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
