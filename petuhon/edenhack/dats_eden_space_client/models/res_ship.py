from typing import TYPE_CHECKING, Any, Dict, List, Type, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.res_planet import ResPlanet
    from ..models.res_ship_garbage import ResShipGarbage


T = TypeVar("T", bound="ResShip")


@_attrs_define
class ResShip:
    """
    Attributes:
        capacity_x (Union[Unset, int]):  Example: 8.
        capacity_y (Union[Unset, int]):  Example: 11.
        fuel_used (Union[Unset, int]):  Example: 1000.
        garbage (Union[Unset, ResShipGarbage]):  Example: {'6fTzQid': [[0, 0], [0, 1], [1, 1]], 'RVnTkM59': [[0, 0], [0,
            1], [1, 1], [2, 1], [1, 2]]}.
        planet (Union[Unset, ResPlanet]):
    """

    capacity_x: Union[Unset, int] = UNSET
    capacity_y: Union[Unset, int] = UNSET
    fuel_used: Union[Unset, int] = UNSET
    garbage: Union[Unset, "ResShipGarbage"] = UNSET
    planet: Union[Unset, "ResPlanet"] = UNSET
    additional_properties: Dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> Dict[str, Any]:
        capacity_x = self.capacity_x

        capacity_y = self.capacity_y

        fuel_used = self.fuel_used

        garbage: Union[Unset, Dict[str, Any]] = UNSET
        if not isinstance(self.garbage, Unset):
            garbage = self.garbage.to_dict()

        planet: Union[Unset, Dict[str, Any]] = UNSET
        if not isinstance(self.planet, Unset):
            planet = self.planet.to_dict()

        field_dict: Dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if capacity_x is not UNSET:
            field_dict["capacityX"] = capacity_x
        if capacity_y is not UNSET:
            field_dict["capacityY"] = capacity_y
        if fuel_used is not UNSET:
            field_dict["fuelUsed"] = fuel_used
        if garbage is not UNSET:
            field_dict["garbage"] = garbage
        if planet is not UNSET:
            field_dict["planet"] = planet

        return field_dict

    @classmethod
    def from_dict(cls: Type[T], src_dict: Dict[str, Any]) -> T:
        from ..models.res_planet import ResPlanet
        from ..models.res_ship_garbage import ResShipGarbage

        d = src_dict.copy()
        capacity_x = d.pop("capacityX", UNSET)

        capacity_y = d.pop("capacityY", UNSET)

        fuel_used = d.pop("fuelUsed", UNSET)

        _garbage = d.pop("garbage", UNSET)
        garbage: Union[Unset, ResShipGarbage]
        if isinstance(_garbage, Unset):
            garbage = UNSET
        else:
            garbage = ResShipGarbage.from_dict(_garbage)

        _planet = d.pop("planet", UNSET)
        planet: Union[Unset, ResPlanet]
        if isinstance(_planet, Unset):
            planet = UNSET
        else:
            planet = ResPlanet.from_dict(_planet)

        res_ship = cls(
            capacity_x=capacity_x,
            capacity_y=capacity_y,
            fuel_used=fuel_used,
            garbage=garbage,
            planet=planet,
        )

        res_ship.additional_properties = d
        return res_ship

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
