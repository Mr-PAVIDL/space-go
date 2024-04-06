from typing import TYPE_CHECKING, Any, Dict, List, Type, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.req_collect_garbage import ReqCollectGarbage


T = TypeVar("T", bound="ReqCollect")


@_attrs_define
class ReqCollect:
    """
    Attributes:
        garbage (Union[Unset, ReqCollectGarbage]):  Example: {'71B2XMi': [[2, 10], [2, 9], [2, 8], [3, 8]]}.
    """

    garbage: Union[Unset, "ReqCollectGarbage"] = UNSET
    additional_properties: Dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> Dict[str, Any]:
        garbage: Union[Unset, Dict[str, Any]] = UNSET
        if not isinstance(self.garbage, Unset):
            garbage = self.garbage.to_dict()

        field_dict: Dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if garbage is not UNSET:
            field_dict["garbage"] = garbage

        return field_dict

    @classmethod
    def from_dict(cls: Type[T], src_dict: Dict[str, Any]) -> T:
        from ..models.req_collect_garbage import ReqCollectGarbage

        d = src_dict.copy()
        _garbage = d.pop("garbage", UNSET)
        garbage: Union[Unset, ReqCollectGarbage]
        if isinstance(_garbage, Unset):
            garbage = UNSET
        else:
            garbage = ReqCollectGarbage.from_dict(_garbage)

        req_collect = cls(
            garbage=garbage,
        )

        req_collect.additional_properties = d
        return req_collect

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
