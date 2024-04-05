from typing import TYPE_CHECKING, Any, Dict, List, Type, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.res_planet_garbage import ResPlanetGarbage


T = TypeVar("T", bound="ResPlanet")


@_attrs_define
class ResPlanet:
    """
    Attributes:
        garbage (Union[Unset, ResPlanetGarbage]):  Example: {'6fTzQid': [[0, 0], [0, 1], [1, 1]], 'RVnTkM59': [[0, 0],
            [0, 1], [1, 1], [2, 1], [1, 2]]}.
        name (Union[Unset, str]):
    """

    garbage: Union[Unset, "ResPlanetGarbage"] = UNSET
    name: Union[Unset, str] = UNSET
    additional_properties: Dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> Dict[str, Any]:
        garbage: Union[Unset, Dict[str, Any]] = UNSET
        if not isinstance(self.garbage, Unset):
            garbage = self.garbage.to_dict()

        name = self.name

        field_dict: Dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if garbage is not UNSET:
            field_dict["garbage"] = garbage
        if name is not UNSET:
            field_dict["name"] = name

        return field_dict

    @classmethod
    def from_dict(cls: Type[T], src_dict: Dict[str, Any]) -> T:
        from ..models.res_planet_garbage import ResPlanetGarbage

        d = src_dict.copy()
        _garbage = d.pop("garbage", UNSET)
        garbage: Union[Unset, ResPlanetGarbage]
        if isinstance(_garbage, Unset):
            garbage = UNSET
        else:
            garbage = ResPlanetGarbage.from_dict(_garbage)

        name = d.pop("name", UNSET)

        res_planet = cls(
            garbage=garbage,
            name=name,
        )

        res_planet.additional_properties = d
        return res_planet

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
