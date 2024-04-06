from typing import Any, Dict, List, Type, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="ResPlanetDiff")


@_attrs_define
class ResPlanetDiff:
    """
    Attributes:
        from_ (Union[Unset, str]):  Example: Earth.
        fuel (Union[Unset, int]):  Example: 100.
        to (Union[Unset, str]):  Example: Reinger.
    """

    from_: Union[Unset, str] = UNSET
    fuel: Union[Unset, int] = UNSET
    to: Union[Unset, str] = UNSET
    additional_properties: Dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> Dict[str, Any]:
        from_ = self.from_

        fuel = self.fuel

        to = self.to

        field_dict: Dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if from_ is not UNSET:
            field_dict["from"] = from_
        if fuel is not UNSET:
            field_dict["fuel"] = fuel
        if to is not UNSET:
            field_dict["to"] = to

        return field_dict

    @classmethod
    def from_dict(cls: Type[T], src_dict: Dict[str, Any]) -> T:
        d = src_dict.copy()
        from_ = d.pop("from", UNSET)

        fuel = d.pop("fuel", UNSET)

        to = d.pop("to", UNSET)

        res_planet_diff = cls(
            from_=from_,
            fuel=fuel,
            to=to,
        )

        res_planet_diff.additional_properties = d
        return res_planet_diff

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
