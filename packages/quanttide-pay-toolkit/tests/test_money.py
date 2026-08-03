from decimal import Decimal

import pytest

from quanttide_pay_toolkit.money import Money


def test_construct_from_decimal() -> None:
    money = Money(Decimal("12.34"))
    assert money.amount == Decimal("12.34")
    assert money.currency == "CNY"


def test_construct_from_str_and_int() -> None:
    assert Money("9.99").amount == Decimal("9.99")
    assert Money(100).amount == Decimal("100.00")


def test_quantize_to_cents() -> None:
    assert Money("0.005").amount == Decimal("0.01")
    assert Money("0.004").amount == Decimal("0.00")


def test_from_cents() -> None:
    money = Money.from_cents(1234)
    assert money.amount == Decimal("12.34")
    assert money.cents == 1234


def test_from_cents_rejects_non_int() -> None:
    with pytest.raises(TypeError, match="cents 必须为 int"):
        Money.from_cents("1234")  # type: ignore[arg-type]


def test_cents() -> None:
    assert Money("1.00").cents == 100
    assert Money("0.01").cents == 1


def test_add() -> None:
    result = Money("1.10") + Money("2.20")
    assert result == Money("3.30")
    assert result.currency == "CNY"


def test_sub() -> None:
    result = Money("3.30") - Money("1.10")
    assert result == Money("2.20")


def test_add_different_currency_raises() -> None:
    with pytest.raises(ValueError, match="币种不一致"):
        _ = Money("1.00") + Money("1.00", currency="USD")


def test_sub_negative_raises() -> None:
    with pytest.raises(ValueError, match="不能为负数"):
        _ = Money("0.10") - Money("0.20")


def test_negative_amount_raises() -> None:
    with pytest.raises(ValueError, match="不能为负数"):
        Money("-1.00")


def test_invalid_currency_raises() -> None:
    with pytest.raises(ValueError, match="币种代码不合法"):
        Money("1.00", currency="cny")
