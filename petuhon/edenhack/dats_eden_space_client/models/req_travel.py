from typing import Any, Dict, List, Type, TypeVar, Union, cast

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="ReqTravel")


@_attrs_define
class ReqTravel:
    """
    Attributes:
        planets (Union[Unset, List[str]]):  Example: ['Reinger-77', 'Earth'].
    """

    planets: Union[Unset, List[str]] = UNSET
    additional_properties: Dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> Dict[str, Any]:
        planets: Union[Unset, List[str]] = UNSET
        if not isinstance(self.planets, Unset):
            planets = self.planets

        field_dict: Dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if planets is not UNSET:
            field_dict["planets"] = planets

        return field_dict

    @classmethod
    def from_dict(cls: Type[T], src_dict: Dict[str, Any]) -> T:
        d = src_dict.copy()
        planets = cast(List[str], d.pop("planets", UNSET))

        req_travel = cls(
            planets=planets,
        )

        req_travel.additional_properties = d
        return req_travel

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
