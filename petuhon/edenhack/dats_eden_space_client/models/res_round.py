from typing import Any, Dict, List, Type, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="ResRound")


@_attrs_define
class ResRound:
    """
    Attributes:
        start_at (Union[Unset, str]):  Example: 2024-04-04 14:00:00.
        end_at (Union[Unset, str]):  Example: 2024-04-04 14:30:00.
        is_current (Union[Unset, bool]):
        name (Union[Unset, str]):  Example: round 1.
        planet_count (Union[Unset, int]):  Example: 100.
    """

    start_at: Union[Unset, str] = UNSET
    end_at: Union[Unset, str] = UNSET
    is_current: Union[Unset, bool] = UNSET
    name: Union[Unset, str] = UNSET
    planet_count: Union[Unset, int] = UNSET
    additional_properties: Dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> Dict[str, Any]:
        start_at = self.start_at

        end_at = self.end_at

        is_current = self.is_current

        name = self.name

        planet_count = self.planet_count

        field_dict: Dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if start_at is not UNSET:
            field_dict["startAt"] = start_at
        if end_at is not UNSET:
            field_dict["endAt"] = end_at
        if is_current is not UNSET:
            field_dict["isCurrent"] = is_current
        if name is not UNSET:
            field_dict["name"] = name
        if planet_count is not UNSET:
            field_dict["planetCount"] = planet_count

        return field_dict

    @classmethod
    def from_dict(cls: Type[T], src_dict: Dict[str, Any]) -> T:
        d = src_dict.copy()
        start_at = d.pop("startAt", UNSET)

        end_at = d.pop("endAt", UNSET)

        is_current = d.pop("isCurrent", UNSET)

        name = d.pop("name", UNSET)

        planet_count = d.pop("planetCount", UNSET)

        res_round = cls(
            start_at=start_at,
            end_at=end_at,
            is_current=is_current,
            name=name,
            planet_count=planet_count,
        )

        res_round.additional_properties = d
        return res_round

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
