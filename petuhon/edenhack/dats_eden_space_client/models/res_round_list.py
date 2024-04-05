from typing import TYPE_CHECKING, Any, Dict, List, Type, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

if TYPE_CHECKING:
    from ..models.res_round import ResRound


T = TypeVar("T", bound="ResRoundList")


@_attrs_define
class ResRoundList:
    """
    Attributes:
        rounds (Union[Unset, List['ResRound']]):
    """

    rounds: Union[Unset, List["ResRound"]] = UNSET
    additional_properties: Dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> Dict[str, Any]:
        rounds: Union[Unset, List[Dict[str, Any]]] = UNSET
        if not isinstance(self.rounds, Unset):
            rounds = []
            for rounds_item_data in self.rounds:
                rounds_item = rounds_item_data.to_dict()
                rounds.append(rounds_item)

        field_dict: Dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if rounds is not UNSET:
            field_dict["rounds"] = rounds

        return field_dict

    @classmethod
    def from_dict(cls: Type[T], src_dict: Dict[str, Any]) -> T:
        from ..models.res_round import ResRound

        d = src_dict.copy()
        rounds = []
        _rounds = d.pop("rounds", UNSET)
        for rounds_item_data in _rounds or []:
            rounds_item = ResRound.from_dict(rounds_item_data)

            rounds.append(rounds_item)

        res_round_list = cls(
            rounds=rounds,
        )

        res_round_list.additional_properties = d
        return res_round_list

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
