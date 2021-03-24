import abc
import builtins
import datetime
import enum
import typing

import jsii
import publication
import typing_extensions

from ._jsii import *

import cdktf
import constructs


class ApiClient(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.ApiClient",
):
    def __init__(
        self,
        scope_: constructs.Construct,
        id: builtins.str,
        *,
        name: builtins.str,
        scope: typing.List[builtins.str],
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope_: -
        :param id: -
        :param name: -
        :param scope: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = ApiClientConfig(
            name=name,
            scope=scope,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(ApiClient, self, [scope_, id, config])

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="scopeInput")
    def scope_input(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "scopeInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="secret")
    def secret(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "secret"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="scope")
    def scope(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "scope"))

    @scope.setter
    def scope(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "scope", value)


@jsii.data_type(
    jsii_type="labd_commercetools.ApiClientConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "name": "name",
        "scope": "scope",
    },
)
class ApiClientConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        name: builtins.str,
        scope: typing.List[builtins.str],
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param name: -
        :param scope: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "name": name,
            "scope": scope,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def scope(self) -> typing.List[builtins.str]:
        result = self._values.get("scope")
        assert result is not None, "Required property 'scope' is missing"
        return typing.cast(typing.List[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ApiClientConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class ApiExtension(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.ApiExtension",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        destination: typing.Mapping[builtins.str, builtins.str],
        trigger: typing.List["ApiExtensionTrigger"],
        key: typing.Optional[builtins.str] = None,
        timeout_in_ms: typing.Optional[jsii.Number] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param destination: -
        :param trigger: trigger block.
        :param key: -
        :param timeout_in_ms: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = ApiExtensionConfig(
            destination=destination,
            trigger=trigger,
            key=key,
            timeout_in_ms=timeout_in_ms,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(ApiExtension, self, [scope, id, config])

    @jsii.member(jsii_name="resetKey")
    def reset_key(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetKey", []))

    @jsii.member(jsii_name="resetTimeoutInMs")
    def reset_timeout_in_ms(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetTimeoutInMs", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="destinationInput")
    def destination_input(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "destinationInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="triggerInput")
    def trigger_input(self) -> typing.List["ApiExtensionTrigger"]:
        return typing.cast(typing.List["ApiExtensionTrigger"], jsii.get(self, "triggerInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="keyInput")
    def key_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "keyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="timeoutInMsInput")
    def timeout_in_ms_input(self) -> typing.Optional[jsii.Number]:
        return typing.cast(typing.Optional[jsii.Number], jsii.get(self, "timeoutInMsInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="destination")
    def destination(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "destination"))

    @destination.setter
    def destination(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "destination", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @key.setter
    def key(self, value: builtins.str) -> None:
        jsii.set(self, "key", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="timeoutInMs")
    def timeout_in_ms(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "timeoutInMs"))

    @timeout_in_ms.setter
    def timeout_in_ms(self, value: jsii.Number) -> None:
        jsii.set(self, "timeoutInMs", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="trigger")
    def trigger(self) -> typing.List["ApiExtensionTrigger"]:
        return typing.cast(typing.List["ApiExtensionTrigger"], jsii.get(self, "trigger"))

    @trigger.setter
    def trigger(self, value: typing.List["ApiExtensionTrigger"]) -> None:
        jsii.set(self, "trigger", value)


@jsii.data_type(
    jsii_type="labd_commercetools.ApiExtensionConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "destination": "destination",
        "trigger": "trigger",
        "key": "key",
        "timeout_in_ms": "timeoutInMs",
    },
)
class ApiExtensionConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        destination: typing.Mapping[builtins.str, builtins.str],
        trigger: typing.List["ApiExtensionTrigger"],
        key: typing.Optional[builtins.str] = None,
        timeout_in_ms: typing.Optional[jsii.Number] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param destination: -
        :param trigger: trigger block.
        :param key: -
        :param timeout_in_ms: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "destination": destination,
            "trigger": trigger,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if key is not None:
            self._values["key"] = key
        if timeout_in_ms is not None:
            self._values["timeout_in_ms"] = timeout_in_ms

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def destination(self) -> typing.Mapping[builtins.str, builtins.str]:
        result = self._values.get("destination")
        assert result is not None, "Required property 'destination' is missing"
        return typing.cast(typing.Mapping[builtins.str, builtins.str], result)

    @builtins.property
    def trigger(self) -> typing.List["ApiExtensionTrigger"]:
        '''trigger block.'''
        result = self._values.get("trigger")
        assert result is not None, "Required property 'trigger' is missing"
        return typing.cast(typing.List["ApiExtensionTrigger"], result)

    @builtins.property
    def key(self) -> typing.Optional[builtins.str]:
        result = self._values.get("key")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def timeout_in_ms(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("timeout_in_ms")
        return typing.cast(typing.Optional[jsii.Number], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ApiExtensionConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.ApiExtensionTrigger",
    jsii_struct_bases=[],
    name_mapping={"actions": "actions", "resource_type_id": "resourceTypeId"},
)
class ApiExtensionTrigger:
    def __init__(
        self,
        *,
        actions: typing.List[builtins.str],
        resource_type_id: builtins.str,
    ) -> None:
        '''
        :param actions: -
        :param resource_type_id: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "actions": actions,
            "resource_type_id": resource_type_id,
        }

    @builtins.property
    def actions(self) -> typing.List[builtins.str]:
        result = self._values.get("actions")
        assert result is not None, "Required property 'actions' is missing"
        return typing.cast(typing.List[builtins.str], result)

    @builtins.property
    def resource_type_id(self) -> builtins.str:
        result = self._values.get("resource_type_id")
        assert result is not None, "Required property 'resource_type_id' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ApiExtensionTrigger(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class CartDiscount(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.CartDiscount",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        name: typing.Mapping[builtins.str, builtins.str],
        predicate: builtins.str,
        sort_order: builtins.str,
        value: typing.List["CartDiscountValue"],
        description: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        is_active: typing.Optional[builtins.bool] = None,
        key: typing.Optional[builtins.str] = None,
        requires_discount_code: typing.Optional[builtins.bool] = None,
        stacking_mode: typing.Optional[builtins.str] = None,
        target: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        valid_from: typing.Optional[builtins.str] = None,
        valid_until: typing.Optional[builtins.str] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param name: -
        :param predicate: -
        :param sort_order: -
        :param value: value block.
        :param description: -
        :param is_active: -
        :param key: -
        :param requires_discount_code: -
        :param stacking_mode: -
        :param target: -
        :param valid_from: -
        :param valid_until: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = CartDiscountConfig(
            name=name,
            predicate=predicate,
            sort_order=sort_order,
            value=value,
            description=description,
            is_active=is_active,
            key=key,
            requires_discount_code=requires_discount_code,
            stacking_mode=stacking_mode,
            target=target,
            valid_from=valid_from,
            valid_until=valid_until,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(CartDiscount, self, [scope, id, config])

    @jsii.member(jsii_name="resetDescription")
    def reset_description(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDescription", []))

    @jsii.member(jsii_name="resetIsActive")
    def reset_is_active(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetIsActive", []))

    @jsii.member(jsii_name="resetKey")
    def reset_key(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetKey", []))

    @jsii.member(jsii_name="resetRequiresDiscountCode")
    def reset_requires_discount_code(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetRequiresDiscountCode", []))

    @jsii.member(jsii_name="resetStackingMode")
    def reset_stacking_mode(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetStackingMode", []))

    @jsii.member(jsii_name="resetTarget")
    def reset_target(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetTarget", []))

    @jsii.member(jsii_name="resetValidFrom")
    def reset_valid_from(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetValidFrom", []))

    @jsii.member(jsii_name="resetValidUntil")
    def reset_valid_until(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetValidUntil", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="predicateInput")
    def predicate_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "predicateInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="sortOrderInput")
    def sort_order_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "sortOrderInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="valueInput")
    def value_input(self) -> typing.List["CartDiscountValue"]:
        return typing.cast(typing.List["CartDiscountValue"], jsii.get(self, "valueInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="descriptionInput")
    def description_input(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "descriptionInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="isActiveInput")
    def is_active_input(self) -> typing.Optional[builtins.bool]:
        return typing.cast(typing.Optional[builtins.bool], jsii.get(self, "isActiveInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="keyInput")
    def key_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "keyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="requiresDiscountCodeInput")
    def requires_discount_code_input(self) -> typing.Optional[builtins.bool]:
        return typing.cast(typing.Optional[builtins.bool], jsii.get(self, "requiresDiscountCodeInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="stackingModeInput")
    def stacking_mode_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "stackingModeInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="targetInput")
    def target_input(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "targetInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="validFromInput")
    def valid_from_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "validFromInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="validUntilInput")
    def valid_until_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "validUntilInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="description")
    def description(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "description"))

    @description.setter
    def description(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "description", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="isActive")
    def is_active(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "isActive"))

    @is_active.setter
    def is_active(self, value: builtins.bool) -> None:
        jsii.set(self, "isActive", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @key.setter
    def key(self, value: builtins.str) -> None:
        jsii.set(self, "key", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "name"))

    @name.setter
    def name(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="predicate")
    def predicate(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "predicate"))

    @predicate.setter
    def predicate(self, value: builtins.str) -> None:
        jsii.set(self, "predicate", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="requiresDiscountCode")
    def requires_discount_code(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "requiresDiscountCode"))

    @requires_discount_code.setter
    def requires_discount_code(self, value: builtins.bool) -> None:
        jsii.set(self, "requiresDiscountCode", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="sortOrder")
    def sort_order(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "sortOrder"))

    @sort_order.setter
    def sort_order(self, value: builtins.str) -> None:
        jsii.set(self, "sortOrder", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="stackingMode")
    def stacking_mode(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "stackingMode"))

    @stacking_mode.setter
    def stacking_mode(self, value: builtins.str) -> None:
        jsii.set(self, "stackingMode", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="target")
    def target(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "target"))

    @target.setter
    def target(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "target", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="validFrom")
    def valid_from(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "validFrom"))

    @valid_from.setter
    def valid_from(self, value: builtins.str) -> None:
        jsii.set(self, "validFrom", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="validUntil")
    def valid_until(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "validUntil"))

    @valid_until.setter
    def valid_until(self, value: builtins.str) -> None:
        jsii.set(self, "validUntil", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="value")
    def value(self) -> typing.List["CartDiscountValue"]:
        return typing.cast(typing.List["CartDiscountValue"], jsii.get(self, "value"))

    @value.setter
    def value(self, value: typing.List["CartDiscountValue"]) -> None:
        jsii.set(self, "value", value)


@jsii.data_type(
    jsii_type="labd_commercetools.CartDiscountConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "name": "name",
        "predicate": "predicate",
        "sort_order": "sortOrder",
        "value": "value",
        "description": "description",
        "is_active": "isActive",
        "key": "key",
        "requires_discount_code": "requiresDiscountCode",
        "stacking_mode": "stackingMode",
        "target": "target",
        "valid_from": "validFrom",
        "valid_until": "validUntil",
    },
)
class CartDiscountConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        name: typing.Mapping[builtins.str, builtins.str],
        predicate: builtins.str,
        sort_order: builtins.str,
        value: typing.List["CartDiscountValue"],
        description: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        is_active: typing.Optional[builtins.bool] = None,
        key: typing.Optional[builtins.str] = None,
        requires_discount_code: typing.Optional[builtins.bool] = None,
        stacking_mode: typing.Optional[builtins.str] = None,
        target: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        valid_from: typing.Optional[builtins.str] = None,
        valid_until: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param name: -
        :param predicate: -
        :param sort_order: -
        :param value: value block.
        :param description: -
        :param is_active: -
        :param key: -
        :param requires_discount_code: -
        :param stacking_mode: -
        :param target: -
        :param valid_from: -
        :param valid_until: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "name": name,
            "predicate": predicate,
            "sort_order": sort_order,
            "value": value,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if description is not None:
            self._values["description"] = description
        if is_active is not None:
            self._values["is_active"] = is_active
        if key is not None:
            self._values["key"] = key
        if requires_discount_code is not None:
            self._values["requires_discount_code"] = requires_discount_code
        if stacking_mode is not None:
            self._values["stacking_mode"] = stacking_mode
        if target is not None:
            self._values["target"] = target
        if valid_from is not None:
            self._values["valid_from"] = valid_from
        if valid_until is not None:
            self._values["valid_until"] = valid_until

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def name(self) -> typing.Mapping[builtins.str, builtins.str]:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(typing.Mapping[builtins.str, builtins.str], result)

    @builtins.property
    def predicate(self) -> builtins.str:
        result = self._values.get("predicate")
        assert result is not None, "Required property 'predicate' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def sort_order(self) -> builtins.str:
        result = self._values.get("sort_order")
        assert result is not None, "Required property 'sort_order' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def value(self) -> typing.List["CartDiscountValue"]:
        '''value block.'''
        result = self._values.get("value")
        assert result is not None, "Required property 'value' is missing"
        return typing.cast(typing.List["CartDiscountValue"], result)

    @builtins.property
    def description(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("description")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def is_active(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("is_active")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def key(self) -> typing.Optional[builtins.str]:
        result = self._values.get("key")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def requires_discount_code(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("requires_discount_code")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def stacking_mode(self) -> typing.Optional[builtins.str]:
        result = self._values.get("stacking_mode")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def target(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("target")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def valid_from(self) -> typing.Optional[builtins.str]:
        result = self._values.get("valid_from")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def valid_until(self) -> typing.Optional[builtins.str]:
        result = self._values.get("valid_until")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "CartDiscountConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.CartDiscountValue",
    jsii_struct_bases=[],
    name_mapping={
        "type": "type",
        "distribution_channel_id": "distributionChannelId",
        "money": "money",
        "permyriad": "permyriad",
        "product_id": "productId",
        "supply_channel_id": "supplyChannelId",
        "variant": "variant",
    },
)
class CartDiscountValue:
    def __init__(
        self,
        *,
        type: builtins.str,
        distribution_channel_id: typing.Optional[builtins.str] = None,
        money: typing.Optional[typing.List["CartDiscountValueMoney"]] = None,
        permyriad: typing.Optional[jsii.Number] = None,
        product_id: typing.Optional[builtins.str] = None,
        supply_channel_id: typing.Optional[builtins.str] = None,
        variant: typing.Optional[jsii.Number] = None,
    ) -> None:
        '''
        :param type: -
        :param distribution_channel_id: -
        :param money: money block.
        :param permyriad: -
        :param product_id: -
        :param supply_channel_id: -
        :param variant: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "type": type,
        }
        if distribution_channel_id is not None:
            self._values["distribution_channel_id"] = distribution_channel_id
        if money is not None:
            self._values["money"] = money
        if permyriad is not None:
            self._values["permyriad"] = permyriad
        if product_id is not None:
            self._values["product_id"] = product_id
        if supply_channel_id is not None:
            self._values["supply_channel_id"] = supply_channel_id
        if variant is not None:
            self._values["variant"] = variant

    @builtins.property
    def type(self) -> builtins.str:
        result = self._values.get("type")
        assert result is not None, "Required property 'type' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def distribution_channel_id(self) -> typing.Optional[builtins.str]:
        result = self._values.get("distribution_channel_id")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def money(self) -> typing.Optional[typing.List["CartDiscountValueMoney"]]:
        '''money block.'''
        result = self._values.get("money")
        return typing.cast(typing.Optional[typing.List["CartDiscountValueMoney"]], result)

    @builtins.property
    def permyriad(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("permyriad")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def product_id(self) -> typing.Optional[builtins.str]:
        result = self._values.get("product_id")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def supply_channel_id(self) -> typing.Optional[builtins.str]:
        result = self._values.get("supply_channel_id")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def variant(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("variant")
        return typing.cast(typing.Optional[jsii.Number], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "CartDiscountValue(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.CartDiscountValueMoney",
    jsii_struct_bases=[],
    name_mapping={"cent_amount": "centAmount", "currency_code": "currencyCode"},
)
class CartDiscountValueMoney:
    def __init__(
        self,
        *,
        cent_amount: jsii.Number,
        currency_code: builtins.str,
    ) -> None:
        '''
        :param cent_amount: -
        :param currency_code: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "cent_amount": cent_amount,
            "currency_code": currency_code,
        }

    @builtins.property
    def cent_amount(self) -> jsii.Number:
        result = self._values.get("cent_amount")
        assert result is not None, "Required property 'cent_amount' is missing"
        return typing.cast(jsii.Number, result)

    @builtins.property
    def currency_code(self) -> builtins.str:
        result = self._values.get("currency_code")
        assert result is not None, "Required property 'currency_code' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "CartDiscountValueMoney(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class Channel(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.Channel",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        key: builtins.str,
        roles: typing.List[builtins.str],
        description: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        name: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param key: -
        :param roles: -
        :param description: -
        :param name: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = ChannelConfig(
            key=key,
            roles=roles,
            description=description,
            name=name,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(Channel, self, [scope, id, config])

    @jsii.member(jsii_name="resetDescription")
    def reset_description(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDescription", []))

    @jsii.member(jsii_name="resetName")
    def reset_name(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetName", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="keyInput")
    def key_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "keyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="rolesInput")
    def roles_input(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "rolesInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="descriptionInput")
    def description_input(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "descriptionInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="description")
    def description(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "description"))

    @description.setter
    def description(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "description", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @key.setter
    def key(self, value: builtins.str) -> None:
        jsii.set(self, "key", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "name"))

    @name.setter
    def name(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="roles")
    def roles(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "roles"))

    @roles.setter
    def roles(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "roles", value)


@jsii.data_type(
    jsii_type="labd_commercetools.ChannelConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "key": "key",
        "roles": "roles",
        "description": "description",
        "name": "name",
    },
)
class ChannelConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        key: builtins.str,
        roles: typing.List[builtins.str],
        description: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        name: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param key: -
        :param roles: -
        :param description: -
        :param name: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "key": key,
            "roles": roles,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if description is not None:
            self._values["description"] = description
        if name is not None:
            self._values["name"] = name

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def key(self) -> builtins.str:
        result = self._values.get("key")
        assert result is not None, "Required property 'key' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def roles(self) -> typing.List[builtins.str]:
        result = self._values.get("roles")
        assert result is not None, "Required property 'roles' is missing"
        return typing.cast(typing.List[builtins.str], result)

    @builtins.property
    def description(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("description")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def name(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("name")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ChannelConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class CommercetoolsProvider(
    cdktf.TerraformProvider,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.CommercetoolsProvider",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        api_url: builtins.str,
        client_id: builtins.str,
        client_secret: builtins.str,
        project_key: builtins.str,
        scopes: builtins.str,
        token_url: builtins.str,
        alias: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param api_url: The API URL of the commercetools platform. https://docs.commercetools.com/http-api
        :param client_id: The OAuth Client ID for a commercetools platform project. https://docs.commercetools.com/http-api-authorization
        :param client_secret: The OAuth Client Secret for a commercetools platform project. https://docs.commercetools.com/http-api-authorization
        :param project_key: The project key of commercetools platform project. https://docs.commercetools.com/getting-started
        :param scopes: A list as string of OAuth scopes assigned to a project key, to access resources in a commercetools platform project. https://docs.commercetools.com/http-api-authorization
        :param token_url: The authentication URL of the commercetools platform. https://docs.commercetools.com/http-api-authorization
        :param alias: Alias name.
        '''
        config = CommercetoolsProviderConfig(
            api_url=api_url,
            client_id=client_id,
            client_secret=client_secret,
            project_key=project_key,
            scopes=scopes,
            token_url=token_url,
            alias=alias,
        )

        jsii.create(CommercetoolsProvider, self, [scope, id, config])

    @jsii.member(jsii_name="resetAlias")
    def reset_alias(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetAlias", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="apiUrlInput")
    def api_url_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "apiUrlInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="clientIdInput")
    def client_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "clientIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="clientSecretInput")
    def client_secret_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "clientSecretInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="projectKeyInput")
    def project_key_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "projectKeyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="scopesInput")
    def scopes_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "scopesInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="tokenUrlInput")
    def token_url_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "tokenUrlInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="aliasInput")
    def alias_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "aliasInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="apiUrl")
    def api_url(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "apiUrl"))

    @api_url.setter
    def api_url(self, value: builtins.str) -> None:
        jsii.set(self, "apiUrl", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="clientId")
    def client_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "clientId"))

    @client_id.setter
    def client_id(self, value: builtins.str) -> None:
        jsii.set(self, "clientId", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="clientSecret")
    def client_secret(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "clientSecret"))

    @client_secret.setter
    def client_secret(self, value: builtins.str) -> None:
        jsii.set(self, "clientSecret", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="projectKey")
    def project_key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "projectKey"))

    @project_key.setter
    def project_key(self, value: builtins.str) -> None:
        jsii.set(self, "projectKey", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="scopes")
    def scopes(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "scopes"))

    @scopes.setter
    def scopes(self, value: builtins.str) -> None:
        jsii.set(self, "scopes", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="tokenUrl")
    def token_url(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "tokenUrl"))

    @token_url.setter
    def token_url(self, value: builtins.str) -> None:
        jsii.set(self, "tokenUrl", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="alias")
    def alias(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "alias"))

    @alias.setter
    def alias(self, value: typing.Optional[builtins.str]) -> None:
        jsii.set(self, "alias", value)


@jsii.data_type(
    jsii_type="labd_commercetools.CommercetoolsProviderConfig",
    jsii_struct_bases=[],
    name_mapping={
        "api_url": "apiUrl",
        "client_id": "clientId",
        "client_secret": "clientSecret",
        "project_key": "projectKey",
        "scopes": "scopes",
        "token_url": "tokenUrl",
        "alias": "alias",
    },
)
class CommercetoolsProviderConfig:
    def __init__(
        self,
        *,
        api_url: builtins.str,
        client_id: builtins.str,
        client_secret: builtins.str,
        project_key: builtins.str,
        scopes: builtins.str,
        token_url: builtins.str,
        alias: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param api_url: The API URL of the commercetools platform. https://docs.commercetools.com/http-api
        :param client_id: The OAuth Client ID for a commercetools platform project. https://docs.commercetools.com/http-api-authorization
        :param client_secret: The OAuth Client Secret for a commercetools platform project. https://docs.commercetools.com/http-api-authorization
        :param project_key: The project key of commercetools platform project. https://docs.commercetools.com/getting-started
        :param scopes: A list as string of OAuth scopes assigned to a project key, to access resources in a commercetools platform project. https://docs.commercetools.com/http-api-authorization
        :param token_url: The authentication URL of the commercetools platform. https://docs.commercetools.com/http-api-authorization
        :param alias: Alias name.
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "api_url": api_url,
            "client_id": client_id,
            "client_secret": client_secret,
            "project_key": project_key,
            "scopes": scopes,
            "token_url": token_url,
        }
        if alias is not None:
            self._values["alias"] = alias

    @builtins.property
    def api_url(self) -> builtins.str:
        '''The API URL of the commercetools platform.

        https://docs.commercetools.com/http-api
        '''
        result = self._values.get("api_url")
        assert result is not None, "Required property 'api_url' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def client_id(self) -> builtins.str:
        '''The OAuth Client ID for a commercetools platform project.

        https://docs.commercetools.com/http-api-authorization
        '''
        result = self._values.get("client_id")
        assert result is not None, "Required property 'client_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def client_secret(self) -> builtins.str:
        '''The OAuth Client Secret for a commercetools platform project.

        https://docs.commercetools.com/http-api-authorization
        '''
        result = self._values.get("client_secret")
        assert result is not None, "Required property 'client_secret' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def project_key(self) -> builtins.str:
        '''The project key of commercetools platform project.

        https://docs.commercetools.com/getting-started
        '''
        result = self._values.get("project_key")
        assert result is not None, "Required property 'project_key' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def scopes(self) -> builtins.str:
        '''A list as string of OAuth scopes assigned to a project key, to access resources in a commercetools platform project.

        https://docs.commercetools.com/http-api-authorization
        '''
        result = self._values.get("scopes")
        assert result is not None, "Required property 'scopes' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def token_url(self) -> builtins.str:
        '''The authentication URL of the commercetools platform.

        https://docs.commercetools.com/http-api-authorization
        '''
        result = self._values.get("token_url")
        assert result is not None, "Required property 'token_url' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def alias(self) -> typing.Optional[builtins.str]:
        '''Alias name.'''
        result = self._values.get("alias")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "CommercetoolsProviderConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class CustomObject(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.CustomObject",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        container: builtins.str,
        key: builtins.str,
        value: builtins.str,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param container: -
        :param key: -
        :param value: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = CustomObjectConfig(
            container=container,
            key=key,
            value=value,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(CustomObject, self, [scope, id, config])

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="containerInput")
    def container_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "containerInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="keyInput")
    def key_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "keyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="valueInput")
    def value_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "valueInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="container")
    def container(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "container"))

    @container.setter
    def container(self, value: builtins.str) -> None:
        jsii.set(self, "container", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @key.setter
    def key(self, value: builtins.str) -> None:
        jsii.set(self, "key", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="value")
    def value(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "value"))

    @value.setter
    def value(self, value: builtins.str) -> None:
        jsii.set(self, "value", value)


@jsii.data_type(
    jsii_type="labd_commercetools.CustomObjectConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "container": "container",
        "key": "key",
        "value": "value",
    },
)
class CustomObjectConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        container: builtins.str,
        key: builtins.str,
        value: builtins.str,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param container: -
        :param key: -
        :param value: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "container": container,
            "key": key,
            "value": value,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def container(self) -> builtins.str:
        result = self._values.get("container")
        assert result is not None, "Required property 'container' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def key(self) -> builtins.str:
        result = self._values.get("key")
        assert result is not None, "Required property 'key' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def value(self) -> builtins.str:
        result = self._values.get("value")
        assert result is not None, "Required property 'value' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "CustomObjectConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class DiscountCode(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.DiscountCode",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        cart_discounts: typing.List[builtins.str],
        code: builtins.str,
        description: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        groups: typing.Optional[typing.List[builtins.str]] = None,
        is_active: typing.Optional[builtins.bool] = None,
        max_applications: typing.Optional[jsii.Number] = None,
        max_applications_per_customer: typing.Optional[jsii.Number] = None,
        name: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        predicate: typing.Optional[builtins.str] = None,
        valid_from: typing.Optional[builtins.str] = None,
        valid_until: typing.Optional[builtins.str] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param cart_discounts: -
        :param code: -
        :param description: -
        :param groups: -
        :param is_active: -
        :param max_applications: -
        :param max_applications_per_customer: -
        :param name: -
        :param predicate: -
        :param valid_from: -
        :param valid_until: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = DiscountCodeConfig(
            cart_discounts=cart_discounts,
            code=code,
            description=description,
            groups=groups,
            is_active=is_active,
            max_applications=max_applications,
            max_applications_per_customer=max_applications_per_customer,
            name=name,
            predicate=predicate,
            valid_from=valid_from,
            valid_until=valid_until,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(DiscountCode, self, [scope, id, config])

    @jsii.member(jsii_name="resetDescription")
    def reset_description(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDescription", []))

    @jsii.member(jsii_name="resetGroups")
    def reset_groups(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetGroups", []))

    @jsii.member(jsii_name="resetIsActive")
    def reset_is_active(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetIsActive", []))

    @jsii.member(jsii_name="resetMaxApplications")
    def reset_max_applications(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetMaxApplications", []))

    @jsii.member(jsii_name="resetMaxApplicationsPerCustomer")
    def reset_max_applications_per_customer(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetMaxApplicationsPerCustomer", []))

    @jsii.member(jsii_name="resetName")
    def reset_name(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetName", []))

    @jsii.member(jsii_name="resetPredicate")
    def reset_predicate(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetPredicate", []))

    @jsii.member(jsii_name="resetValidFrom")
    def reset_valid_from(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetValidFrom", []))

    @jsii.member(jsii_name="resetValidUntil")
    def reset_valid_until(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetValidUntil", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="cartDiscountsInput")
    def cart_discounts_input(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "cartDiscountsInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="codeInput")
    def code_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "codeInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="descriptionInput")
    def description_input(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "descriptionInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="groupsInput")
    def groups_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "groupsInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="isActiveInput")
    def is_active_input(self) -> typing.Optional[builtins.bool]:
        return typing.cast(typing.Optional[builtins.bool], jsii.get(self, "isActiveInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="maxApplicationsInput")
    def max_applications_input(self) -> typing.Optional[jsii.Number]:
        return typing.cast(typing.Optional[jsii.Number], jsii.get(self, "maxApplicationsInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="maxApplicationsPerCustomerInput")
    def max_applications_per_customer_input(self) -> typing.Optional[jsii.Number]:
        return typing.cast(typing.Optional[jsii.Number], jsii.get(self, "maxApplicationsPerCustomerInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="predicateInput")
    def predicate_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "predicateInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="validFromInput")
    def valid_from_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "validFromInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="validUntilInput")
    def valid_until_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "validUntilInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="cartDiscounts")
    def cart_discounts(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "cartDiscounts"))

    @cart_discounts.setter
    def cart_discounts(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "cartDiscounts", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="code")
    def code(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "code"))

    @code.setter
    def code(self, value: builtins.str) -> None:
        jsii.set(self, "code", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="description")
    def description(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "description"))

    @description.setter
    def description(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "description", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="groups")
    def groups(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "groups"))

    @groups.setter
    def groups(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "groups", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="isActive")
    def is_active(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "isActive"))

    @is_active.setter
    def is_active(self, value: builtins.bool) -> None:
        jsii.set(self, "isActive", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="maxApplications")
    def max_applications(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "maxApplications"))

    @max_applications.setter
    def max_applications(self, value: jsii.Number) -> None:
        jsii.set(self, "maxApplications", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="maxApplicationsPerCustomer")
    def max_applications_per_customer(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "maxApplicationsPerCustomer"))

    @max_applications_per_customer.setter
    def max_applications_per_customer(self, value: jsii.Number) -> None:
        jsii.set(self, "maxApplicationsPerCustomer", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "name"))

    @name.setter
    def name(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="predicate")
    def predicate(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "predicate"))

    @predicate.setter
    def predicate(self, value: builtins.str) -> None:
        jsii.set(self, "predicate", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="validFrom")
    def valid_from(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "validFrom"))

    @valid_from.setter
    def valid_from(self, value: builtins.str) -> None:
        jsii.set(self, "validFrom", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="validUntil")
    def valid_until(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "validUntil"))

    @valid_until.setter
    def valid_until(self, value: builtins.str) -> None:
        jsii.set(self, "validUntil", value)


@jsii.data_type(
    jsii_type="labd_commercetools.DiscountCodeConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "cart_discounts": "cartDiscounts",
        "code": "code",
        "description": "description",
        "groups": "groups",
        "is_active": "isActive",
        "max_applications": "maxApplications",
        "max_applications_per_customer": "maxApplicationsPerCustomer",
        "name": "name",
        "predicate": "predicate",
        "valid_from": "validFrom",
        "valid_until": "validUntil",
    },
)
class DiscountCodeConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        cart_discounts: typing.List[builtins.str],
        code: builtins.str,
        description: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        groups: typing.Optional[typing.List[builtins.str]] = None,
        is_active: typing.Optional[builtins.bool] = None,
        max_applications: typing.Optional[jsii.Number] = None,
        max_applications_per_customer: typing.Optional[jsii.Number] = None,
        name: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        predicate: typing.Optional[builtins.str] = None,
        valid_from: typing.Optional[builtins.str] = None,
        valid_until: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param cart_discounts: -
        :param code: -
        :param description: -
        :param groups: -
        :param is_active: -
        :param max_applications: -
        :param max_applications_per_customer: -
        :param name: -
        :param predicate: -
        :param valid_from: -
        :param valid_until: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "cart_discounts": cart_discounts,
            "code": code,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if description is not None:
            self._values["description"] = description
        if groups is not None:
            self._values["groups"] = groups
        if is_active is not None:
            self._values["is_active"] = is_active
        if max_applications is not None:
            self._values["max_applications"] = max_applications
        if max_applications_per_customer is not None:
            self._values["max_applications_per_customer"] = max_applications_per_customer
        if name is not None:
            self._values["name"] = name
        if predicate is not None:
            self._values["predicate"] = predicate
        if valid_from is not None:
            self._values["valid_from"] = valid_from
        if valid_until is not None:
            self._values["valid_until"] = valid_until

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def cart_discounts(self) -> typing.List[builtins.str]:
        result = self._values.get("cart_discounts")
        assert result is not None, "Required property 'cart_discounts' is missing"
        return typing.cast(typing.List[builtins.str], result)

    @builtins.property
    def code(self) -> builtins.str:
        result = self._values.get("code")
        assert result is not None, "Required property 'code' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def description(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("description")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def groups(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("groups")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def is_active(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("is_active")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def max_applications(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("max_applications")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def max_applications_per_customer(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("max_applications_per_customer")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def name(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("name")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def predicate(self) -> typing.Optional[builtins.str]:
        result = self._values.get("predicate")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def valid_from(self) -> typing.Optional[builtins.str]:
        result = self._values.get("valid_from")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def valid_until(self) -> typing.Optional[builtins.str]:
        result = self._values.get("valid_until")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "DiscountCodeConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class ProductType(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.ProductType",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        name: builtins.str,
        attribute: typing.Optional[typing.List["ProductTypeAttribute"]] = None,
        description: typing.Optional[builtins.str] = None,
        key: typing.Optional[builtins.str] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param name: -
        :param attribute: attribute block.
        :param description: -
        :param key: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = ProductTypeConfig(
            name=name,
            attribute=attribute,
            description=description,
            key=key,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(ProductType, self, [scope, id, config])

    @jsii.member(jsii_name="resetAttribute")
    def reset_attribute(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetAttribute", []))

    @jsii.member(jsii_name="resetDescription")
    def reset_description(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDescription", []))

    @jsii.member(jsii_name="resetKey")
    def reset_key(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetKey", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="attributeInput")
    def attribute_input(self) -> typing.Optional[typing.List["ProductTypeAttribute"]]:
        return typing.cast(typing.Optional[typing.List["ProductTypeAttribute"]], jsii.get(self, "attributeInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="descriptionInput")
    def description_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "descriptionInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="keyInput")
    def key_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "keyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="attribute")
    def attribute(self) -> typing.List["ProductTypeAttribute"]:
        return typing.cast(typing.List["ProductTypeAttribute"], jsii.get(self, "attribute"))

    @attribute.setter
    def attribute(self, value: typing.List["ProductTypeAttribute"]) -> None:
        jsii.set(self, "attribute", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="description")
    def description(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "description"))

    @description.setter
    def description(self, value: builtins.str) -> None:
        jsii.set(self, "description", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @key.setter
    def key(self, value: builtins.str) -> None:
        jsii.set(self, "key", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)


@jsii.data_type(
    jsii_type="labd_commercetools.ProductTypeAttribute",
    jsii_struct_bases=[],
    name_mapping={
        "label": "label",
        "name": "name",
        "type": "type",
        "constraint": "constraint",
        "input_hint": "inputHint",
        "input_tip": "inputTip",
        "required": "required",
        "searchable": "searchable",
    },
)
class ProductTypeAttribute:
    def __init__(
        self,
        *,
        label: typing.Mapping[builtins.str, builtins.str],
        name: builtins.str,
        type: typing.List["ProductTypeAttributeType"],
        constraint: typing.Optional[builtins.str] = None,
        input_hint: typing.Optional[builtins.str] = None,
        input_tip: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        required: typing.Optional[builtins.bool] = None,
        searchable: typing.Optional[builtins.bool] = None,
    ) -> None:
        '''
        :param label: -
        :param name: -
        :param type: type block.
        :param constraint: -
        :param input_hint: -
        :param input_tip: -
        :param required: -
        :param searchable: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "label": label,
            "name": name,
            "type": type,
        }
        if constraint is not None:
            self._values["constraint"] = constraint
        if input_hint is not None:
            self._values["input_hint"] = input_hint
        if input_tip is not None:
            self._values["input_tip"] = input_tip
        if required is not None:
            self._values["required"] = required
        if searchable is not None:
            self._values["searchable"] = searchable

    @builtins.property
    def label(self) -> typing.Mapping[builtins.str, builtins.str]:
        result = self._values.get("label")
        assert result is not None, "Required property 'label' is missing"
        return typing.cast(typing.Mapping[builtins.str, builtins.str], result)

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def type(self) -> typing.List["ProductTypeAttributeType"]:
        '''type block.'''
        result = self._values.get("type")
        assert result is not None, "Required property 'type' is missing"
        return typing.cast(typing.List["ProductTypeAttributeType"], result)

    @builtins.property
    def constraint(self) -> typing.Optional[builtins.str]:
        result = self._values.get("constraint")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def input_hint(self) -> typing.Optional[builtins.str]:
        result = self._values.get("input_hint")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def input_tip(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("input_tip")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def required(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("required")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def searchable(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("searchable")
        return typing.cast(typing.Optional[builtins.bool], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ProductTypeAttribute(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.ProductTypeAttributeType",
    jsii_struct_bases=[],
    name_mapping={
        "name": "name",
        "element_type": "elementType",
        "localized_value": "localizedValue",
        "reference_type_id": "referenceTypeId",
        "type_reference": "typeReference",
        "values": "values",
    },
)
class ProductTypeAttributeType:
    def __init__(
        self,
        *,
        name: builtins.str,
        element_type: typing.Optional[typing.List["ProductTypeAttributeTypeElementType"]] = None,
        localized_value: typing.Optional[typing.List["ProductTypeAttributeTypeLocalizedValue"]] = None,
        reference_type_id: typing.Optional[builtins.str] = None,
        type_reference: typing.Optional[builtins.str] = None,
        values: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
    ) -> None:
        '''
        :param name: -
        :param element_type: element_type block.
        :param localized_value: localized_value block.
        :param reference_type_id: -
        :param type_reference: -
        :param values: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "name": name,
        }
        if element_type is not None:
            self._values["element_type"] = element_type
        if localized_value is not None:
            self._values["localized_value"] = localized_value
        if reference_type_id is not None:
            self._values["reference_type_id"] = reference_type_id
        if type_reference is not None:
            self._values["type_reference"] = type_reference
        if values is not None:
            self._values["values"] = values

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def element_type(
        self,
    ) -> typing.Optional[typing.List["ProductTypeAttributeTypeElementType"]]:
        '''element_type block.'''
        result = self._values.get("element_type")
        return typing.cast(typing.Optional[typing.List["ProductTypeAttributeTypeElementType"]], result)

    @builtins.property
    def localized_value(
        self,
    ) -> typing.Optional[typing.List["ProductTypeAttributeTypeLocalizedValue"]]:
        '''localized_value block.'''
        result = self._values.get("localized_value")
        return typing.cast(typing.Optional[typing.List["ProductTypeAttributeTypeLocalizedValue"]], result)

    @builtins.property
    def reference_type_id(self) -> typing.Optional[builtins.str]:
        result = self._values.get("reference_type_id")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def type_reference(self) -> typing.Optional[builtins.str]:
        result = self._values.get("type_reference")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def values(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("values")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ProductTypeAttributeType(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.ProductTypeAttributeTypeElementType",
    jsii_struct_bases=[],
    name_mapping={
        "name": "name",
        "localized_value": "localizedValue",
        "reference_type_id": "referenceTypeId",
        "type_reference": "typeReference",
        "values": "values",
    },
)
class ProductTypeAttributeTypeElementType:
    def __init__(
        self,
        *,
        name: builtins.str,
        localized_value: typing.Optional[typing.List["ProductTypeAttributeTypeElementTypeLocalizedValue"]] = None,
        reference_type_id: typing.Optional[builtins.str] = None,
        type_reference: typing.Optional[builtins.str] = None,
        values: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
    ) -> None:
        '''
        :param name: -
        :param localized_value: localized_value block.
        :param reference_type_id: -
        :param type_reference: -
        :param values: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "name": name,
        }
        if localized_value is not None:
            self._values["localized_value"] = localized_value
        if reference_type_id is not None:
            self._values["reference_type_id"] = reference_type_id
        if type_reference is not None:
            self._values["type_reference"] = type_reference
        if values is not None:
            self._values["values"] = values

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def localized_value(
        self,
    ) -> typing.Optional[typing.List["ProductTypeAttributeTypeElementTypeLocalizedValue"]]:
        '''localized_value block.'''
        result = self._values.get("localized_value")
        return typing.cast(typing.Optional[typing.List["ProductTypeAttributeTypeElementTypeLocalizedValue"]], result)

    @builtins.property
    def reference_type_id(self) -> typing.Optional[builtins.str]:
        result = self._values.get("reference_type_id")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def type_reference(self) -> typing.Optional[builtins.str]:
        result = self._values.get("type_reference")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def values(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("values")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ProductTypeAttributeTypeElementType(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.ProductTypeAttributeTypeElementTypeLocalizedValue",
    jsii_struct_bases=[],
    name_mapping={"key": "key", "label": "label"},
)
class ProductTypeAttributeTypeElementTypeLocalizedValue:
    def __init__(
        self,
        *,
        key: builtins.str,
        label: typing.Mapping[builtins.str, builtins.str],
    ) -> None:
        '''
        :param key: -
        :param label: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "key": key,
            "label": label,
        }

    @builtins.property
    def key(self) -> builtins.str:
        result = self._values.get("key")
        assert result is not None, "Required property 'key' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def label(self) -> typing.Mapping[builtins.str, builtins.str]:
        result = self._values.get("label")
        assert result is not None, "Required property 'label' is missing"
        return typing.cast(typing.Mapping[builtins.str, builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ProductTypeAttributeTypeElementTypeLocalizedValue(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.ProductTypeAttributeTypeLocalizedValue",
    jsii_struct_bases=[],
    name_mapping={"key": "key", "label": "label"},
)
class ProductTypeAttributeTypeLocalizedValue:
    def __init__(
        self,
        *,
        key: builtins.str,
        label: typing.Mapping[builtins.str, builtins.str],
    ) -> None:
        '''
        :param key: -
        :param label: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "key": key,
            "label": label,
        }

    @builtins.property
    def key(self) -> builtins.str:
        result = self._values.get("key")
        assert result is not None, "Required property 'key' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def label(self) -> typing.Mapping[builtins.str, builtins.str]:
        result = self._values.get("label")
        assert result is not None, "Required property 'label' is missing"
        return typing.cast(typing.Mapping[builtins.str, builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ProductTypeAttributeTypeLocalizedValue(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.ProductTypeConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "name": "name",
        "attribute": "attribute",
        "description": "description",
        "key": "key",
    },
)
class ProductTypeConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        name: builtins.str,
        attribute: typing.Optional[typing.List[ProductTypeAttribute]] = None,
        description: typing.Optional[builtins.str] = None,
        key: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param name: -
        :param attribute: attribute block.
        :param description: -
        :param key: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "name": name,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if attribute is not None:
            self._values["attribute"] = attribute
        if description is not None:
            self._values["description"] = description
        if key is not None:
            self._values["key"] = key

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def attribute(self) -> typing.Optional[typing.List[ProductTypeAttribute]]:
        '''attribute block.'''
        result = self._values.get("attribute")
        return typing.cast(typing.Optional[typing.List[ProductTypeAttribute]], result)

    @builtins.property
    def description(self) -> typing.Optional[builtins.str]:
        result = self._values.get("description")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def key(self) -> typing.Optional[builtins.str]:
        result = self._values.get("key")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ProductTypeConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class ProjectSettings(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.ProjectSettings",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        countries: typing.Optional[typing.List[builtins.str]] = None,
        currencies: typing.Optional[typing.List[builtins.str]] = None,
        external_oauth: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        languages: typing.Optional[typing.List[builtins.str]] = None,
        messages: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        name: typing.Optional[builtins.str] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param countries: -
        :param currencies: -
        :param external_oauth: -
        :param languages: -
        :param messages: -
        :param name: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = ProjectSettingsConfig(
            countries=countries,
            currencies=currencies,
            external_oauth=external_oauth,
            languages=languages,
            messages=messages,
            name=name,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(ProjectSettings, self, [scope, id, config])

    @jsii.member(jsii_name="resetCountries")
    def reset_countries(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetCountries", []))

    @jsii.member(jsii_name="resetCurrencies")
    def reset_currencies(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetCurrencies", []))

    @jsii.member(jsii_name="resetExternalOauth")
    def reset_external_oauth(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetExternalOauth", []))

    @jsii.member(jsii_name="resetLanguages")
    def reset_languages(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetLanguages", []))

    @jsii.member(jsii_name="resetMessages")
    def reset_messages(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetMessages", []))

    @jsii.member(jsii_name="resetName")
    def reset_name(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetName", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="countriesInput")
    def countries_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "countriesInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="currenciesInput")
    def currencies_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "currenciesInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="externalOauthInput")
    def external_oauth_input(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "externalOauthInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="languagesInput")
    def languages_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "languagesInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="messagesInput")
    def messages_input(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "messagesInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="countries")
    def countries(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "countries"))

    @countries.setter
    def countries(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "countries", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="currencies")
    def currencies(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "currencies"))

    @currencies.setter
    def currencies(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "currencies", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="externalOauth")
    def external_oauth(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "externalOauth"))

    @external_oauth.setter
    def external_oauth(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "externalOauth", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="languages")
    def languages(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "languages"))

    @languages.setter
    def languages(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "languages", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="messages")
    def messages(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "messages"))

    @messages.setter
    def messages(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "messages", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)


@jsii.data_type(
    jsii_type="labd_commercetools.ProjectSettingsConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "countries": "countries",
        "currencies": "currencies",
        "external_oauth": "externalOauth",
        "languages": "languages",
        "messages": "messages",
        "name": "name",
    },
)
class ProjectSettingsConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        countries: typing.Optional[typing.List[builtins.str]] = None,
        currencies: typing.Optional[typing.List[builtins.str]] = None,
        external_oauth: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        languages: typing.Optional[typing.List[builtins.str]] = None,
        messages: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        name: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param countries: -
        :param currencies: -
        :param external_oauth: -
        :param languages: -
        :param messages: -
        :param name: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {}
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if countries is not None:
            self._values["countries"] = countries
        if currencies is not None:
            self._values["currencies"] = currencies
        if external_oauth is not None:
            self._values["external_oauth"] = external_oauth
        if languages is not None:
            self._values["languages"] = languages
        if messages is not None:
            self._values["messages"] = messages
        if name is not None:
            self._values["name"] = name

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def countries(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("countries")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def currencies(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("currencies")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def external_oauth(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("external_oauth")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def languages(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("languages")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def messages(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("messages")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def name(self) -> typing.Optional[builtins.str]:
        result = self._values.get("name")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ProjectSettingsConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class ShippingMethod(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.ShippingMethod",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        name: builtins.str,
        description: typing.Optional[builtins.str] = None,
        is_default: typing.Optional[builtins.bool] = None,
        key: typing.Optional[builtins.str] = None,
        predicate: typing.Optional[builtins.str] = None,
        tax_category_id: typing.Optional[builtins.str] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param name: -
        :param description: -
        :param is_default: -
        :param key: -
        :param predicate: -
        :param tax_category_id: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = ShippingMethodConfig(
            name=name,
            description=description,
            is_default=is_default,
            key=key,
            predicate=predicate,
            tax_category_id=tax_category_id,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(ShippingMethod, self, [scope, id, config])

    @jsii.member(jsii_name="resetDescription")
    def reset_description(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDescription", []))

    @jsii.member(jsii_name="resetIsDefault")
    def reset_is_default(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetIsDefault", []))

    @jsii.member(jsii_name="resetKey")
    def reset_key(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetKey", []))

    @jsii.member(jsii_name="resetPredicate")
    def reset_predicate(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetPredicate", []))

    @jsii.member(jsii_name="resetTaxCategoryId")
    def reset_tax_category_id(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetTaxCategoryId", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="descriptionInput")
    def description_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "descriptionInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="isDefaultInput")
    def is_default_input(self) -> typing.Optional[builtins.bool]:
        return typing.cast(typing.Optional[builtins.bool], jsii.get(self, "isDefaultInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="keyInput")
    def key_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "keyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="predicateInput")
    def predicate_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "predicateInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="taxCategoryIdInput")
    def tax_category_id_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "taxCategoryIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="description")
    def description(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "description"))

    @description.setter
    def description(self, value: builtins.str) -> None:
        jsii.set(self, "description", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="isDefault")
    def is_default(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "isDefault"))

    @is_default.setter
    def is_default(self, value: builtins.bool) -> None:
        jsii.set(self, "isDefault", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @key.setter
    def key(self, value: builtins.str) -> None:
        jsii.set(self, "key", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="predicate")
    def predicate(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "predicate"))

    @predicate.setter
    def predicate(self, value: builtins.str) -> None:
        jsii.set(self, "predicate", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="taxCategoryId")
    def tax_category_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "taxCategoryId"))

    @tax_category_id.setter
    def tax_category_id(self, value: builtins.str) -> None:
        jsii.set(self, "taxCategoryId", value)


@jsii.data_type(
    jsii_type="labd_commercetools.ShippingMethodConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "name": "name",
        "description": "description",
        "is_default": "isDefault",
        "key": "key",
        "predicate": "predicate",
        "tax_category_id": "taxCategoryId",
    },
)
class ShippingMethodConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        name: builtins.str,
        description: typing.Optional[builtins.str] = None,
        is_default: typing.Optional[builtins.bool] = None,
        key: typing.Optional[builtins.str] = None,
        predicate: typing.Optional[builtins.str] = None,
        tax_category_id: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param name: -
        :param description: -
        :param is_default: -
        :param key: -
        :param predicate: -
        :param tax_category_id: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "name": name,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if description is not None:
            self._values["description"] = description
        if is_default is not None:
            self._values["is_default"] = is_default
        if key is not None:
            self._values["key"] = key
        if predicate is not None:
            self._values["predicate"] = predicate
        if tax_category_id is not None:
            self._values["tax_category_id"] = tax_category_id

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def description(self) -> typing.Optional[builtins.str]:
        result = self._values.get("description")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def is_default(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("is_default")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def key(self) -> typing.Optional[builtins.str]:
        result = self._values.get("key")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def predicate(self) -> typing.Optional[builtins.str]:
        result = self._values.get("predicate")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def tax_category_id(self) -> typing.Optional[builtins.str]:
        result = self._values.get("tax_category_id")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ShippingMethodConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class ShippingZone(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.ShippingZone",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        description: typing.Optional[builtins.str] = None,
        key: typing.Optional[builtins.str] = None,
        location: typing.Optional[typing.List["ShippingZoneLocation"]] = None,
        name: typing.Optional[builtins.str] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param description: -
        :param key: -
        :param location: location block.
        :param name: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = ShippingZoneConfig(
            description=description,
            key=key,
            location=location,
            name=name,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(ShippingZone, self, [scope, id, config])

    @jsii.member(jsii_name="resetDescription")
    def reset_description(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDescription", []))

    @jsii.member(jsii_name="resetKey")
    def reset_key(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetKey", []))

    @jsii.member(jsii_name="resetLocation")
    def reset_location(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetLocation", []))

    @jsii.member(jsii_name="resetName")
    def reset_name(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetName", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="descriptionInput")
    def description_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "descriptionInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="keyInput")
    def key_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "keyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="locationInput")
    def location_input(self) -> typing.Optional[typing.List["ShippingZoneLocation"]]:
        return typing.cast(typing.Optional[typing.List["ShippingZoneLocation"]], jsii.get(self, "locationInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="description")
    def description(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "description"))

    @description.setter
    def description(self, value: builtins.str) -> None:
        jsii.set(self, "description", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @key.setter
    def key(self, value: builtins.str) -> None:
        jsii.set(self, "key", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="location")
    def location(self) -> typing.List["ShippingZoneLocation"]:
        return typing.cast(typing.List["ShippingZoneLocation"], jsii.get(self, "location"))

    @location.setter
    def location(self, value: typing.List["ShippingZoneLocation"]) -> None:
        jsii.set(self, "location", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)


@jsii.data_type(
    jsii_type="labd_commercetools.ShippingZoneConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "description": "description",
        "key": "key",
        "location": "location",
        "name": "name",
    },
)
class ShippingZoneConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        description: typing.Optional[builtins.str] = None,
        key: typing.Optional[builtins.str] = None,
        location: typing.Optional[typing.List["ShippingZoneLocation"]] = None,
        name: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param description: -
        :param key: -
        :param location: location block.
        :param name: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {}
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if description is not None:
            self._values["description"] = description
        if key is not None:
            self._values["key"] = key
        if location is not None:
            self._values["location"] = location
        if name is not None:
            self._values["name"] = name

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def description(self) -> typing.Optional[builtins.str]:
        result = self._values.get("description")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def key(self) -> typing.Optional[builtins.str]:
        result = self._values.get("key")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def location(self) -> typing.Optional[typing.List["ShippingZoneLocation"]]:
        '''location block.'''
        result = self._values.get("location")
        return typing.cast(typing.Optional[typing.List["ShippingZoneLocation"]], result)

    @builtins.property
    def name(self) -> typing.Optional[builtins.str]:
        result = self._values.get("name")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ShippingZoneConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.ShippingZoneLocation",
    jsii_struct_bases=[],
    name_mapping={"country": "country", "state": "state"},
)
class ShippingZoneLocation:
    def __init__(
        self,
        *,
        country: builtins.str,
        state: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param country: -
        :param state: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "country": country,
        }
        if state is not None:
            self._values["state"] = state

    @builtins.property
    def country(self) -> builtins.str:
        result = self._values.get("country")
        assert result is not None, "Required property 'country' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def state(self) -> typing.Optional[builtins.str]:
        result = self._values.get("state")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ShippingZoneLocation(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class ShippingZoneRate(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.ShippingZoneRate",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        price: typing.List["ShippingZoneRatePrice"],
        shipping_method_id: builtins.str,
        shipping_zone_id: builtins.str,
        free_above: typing.Optional[typing.List["ShippingZoneRateFreeAbove"]] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param price: price block.
        :param shipping_method_id: -
        :param shipping_zone_id: -
        :param free_above: free_above block.
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = ShippingZoneRateConfig(
            price=price,
            shipping_method_id=shipping_method_id,
            shipping_zone_id=shipping_zone_id,
            free_above=free_above,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(ShippingZoneRate, self, [scope, id, config])

    @jsii.member(jsii_name="resetFreeAbove")
    def reset_free_above(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetFreeAbove", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="priceInput")
    def price_input(self) -> typing.List["ShippingZoneRatePrice"]:
        return typing.cast(typing.List["ShippingZoneRatePrice"], jsii.get(self, "priceInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="shippingMethodIdInput")
    def shipping_method_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "shippingMethodIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="shippingZoneIdInput")
    def shipping_zone_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "shippingZoneIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="freeAboveInput")
    def free_above_input(
        self,
    ) -> typing.Optional[typing.List["ShippingZoneRateFreeAbove"]]:
        return typing.cast(typing.Optional[typing.List["ShippingZoneRateFreeAbove"]], jsii.get(self, "freeAboveInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="freeAbove")
    def free_above(self) -> typing.List["ShippingZoneRateFreeAbove"]:
        return typing.cast(typing.List["ShippingZoneRateFreeAbove"], jsii.get(self, "freeAbove"))

    @free_above.setter
    def free_above(self, value: typing.List["ShippingZoneRateFreeAbove"]) -> None:
        jsii.set(self, "freeAbove", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="price")
    def price(self) -> typing.List["ShippingZoneRatePrice"]:
        return typing.cast(typing.List["ShippingZoneRatePrice"], jsii.get(self, "price"))

    @price.setter
    def price(self, value: typing.List["ShippingZoneRatePrice"]) -> None:
        jsii.set(self, "price", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="shippingMethodId")
    def shipping_method_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "shippingMethodId"))

    @shipping_method_id.setter
    def shipping_method_id(self, value: builtins.str) -> None:
        jsii.set(self, "shippingMethodId", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="shippingZoneId")
    def shipping_zone_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "shippingZoneId"))

    @shipping_zone_id.setter
    def shipping_zone_id(self, value: builtins.str) -> None:
        jsii.set(self, "shippingZoneId", value)


@jsii.data_type(
    jsii_type="labd_commercetools.ShippingZoneRateConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "price": "price",
        "shipping_method_id": "shippingMethodId",
        "shipping_zone_id": "shippingZoneId",
        "free_above": "freeAbove",
    },
)
class ShippingZoneRateConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        price: typing.List["ShippingZoneRatePrice"],
        shipping_method_id: builtins.str,
        shipping_zone_id: builtins.str,
        free_above: typing.Optional[typing.List["ShippingZoneRateFreeAbove"]] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param price: price block.
        :param shipping_method_id: -
        :param shipping_zone_id: -
        :param free_above: free_above block.
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "price": price,
            "shipping_method_id": shipping_method_id,
            "shipping_zone_id": shipping_zone_id,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if free_above is not None:
            self._values["free_above"] = free_above

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def price(self) -> typing.List["ShippingZoneRatePrice"]:
        '''price block.'''
        result = self._values.get("price")
        assert result is not None, "Required property 'price' is missing"
        return typing.cast(typing.List["ShippingZoneRatePrice"], result)

    @builtins.property
    def shipping_method_id(self) -> builtins.str:
        result = self._values.get("shipping_method_id")
        assert result is not None, "Required property 'shipping_method_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def shipping_zone_id(self) -> builtins.str:
        result = self._values.get("shipping_zone_id")
        assert result is not None, "Required property 'shipping_zone_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def free_above(self) -> typing.Optional[typing.List["ShippingZoneRateFreeAbove"]]:
        '''free_above block.'''
        result = self._values.get("free_above")
        return typing.cast(typing.Optional[typing.List["ShippingZoneRateFreeAbove"]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ShippingZoneRateConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.ShippingZoneRateFreeAbove",
    jsii_struct_bases=[],
    name_mapping={"cent_amount": "centAmount", "currency_code": "currencyCode"},
)
class ShippingZoneRateFreeAbove:
    def __init__(
        self,
        *,
        cent_amount: jsii.Number,
        currency_code: builtins.str,
    ) -> None:
        '''
        :param cent_amount: -
        :param currency_code: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "cent_amount": cent_amount,
            "currency_code": currency_code,
        }

    @builtins.property
    def cent_amount(self) -> jsii.Number:
        result = self._values.get("cent_amount")
        assert result is not None, "Required property 'cent_amount' is missing"
        return typing.cast(jsii.Number, result)

    @builtins.property
    def currency_code(self) -> builtins.str:
        result = self._values.get("currency_code")
        assert result is not None, "Required property 'currency_code' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ShippingZoneRateFreeAbove(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.ShippingZoneRatePrice",
    jsii_struct_bases=[],
    name_mapping={"cent_amount": "centAmount", "currency_code": "currencyCode"},
)
class ShippingZoneRatePrice:
    def __init__(
        self,
        *,
        cent_amount: jsii.Number,
        currency_code: builtins.str,
    ) -> None:
        '''
        :param cent_amount: -
        :param currency_code: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "cent_amount": cent_amount,
            "currency_code": currency_code,
        }

    @builtins.property
    def cent_amount(self) -> jsii.Number:
        result = self._values.get("cent_amount")
        assert result is not None, "Required property 'cent_amount' is missing"
        return typing.cast(jsii.Number, result)

    @builtins.property
    def currency_code(self) -> builtins.str:
        result = self._values.get("currency_code")
        assert result is not None, "Required property 'currency_code' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ShippingZoneRatePrice(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class State(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.State",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        key: builtins.str,
        type: builtins.str,
        description: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        initial: typing.Optional[builtins.bool] = None,
        name: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        roles: typing.Optional[typing.List[builtins.str]] = None,
        transitions: typing.Optional[typing.List[builtins.str]] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param key: -
        :param type: -
        :param description: -
        :param initial: -
        :param name: -
        :param roles: -
        :param transitions: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = StateConfig(
            key=key,
            type=type,
            description=description,
            initial=initial,
            name=name,
            roles=roles,
            transitions=transitions,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(State, self, [scope, id, config])

    @jsii.member(jsii_name="resetDescription")
    def reset_description(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDescription", []))

    @jsii.member(jsii_name="resetInitial")
    def reset_initial(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetInitial", []))

    @jsii.member(jsii_name="resetName")
    def reset_name(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetName", []))

    @jsii.member(jsii_name="resetRoles")
    def reset_roles(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetRoles", []))

    @jsii.member(jsii_name="resetTransitions")
    def reset_transitions(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetTransitions", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="keyInput")
    def key_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "keyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="typeInput")
    def type_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "typeInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="descriptionInput")
    def description_input(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "descriptionInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="initialInput")
    def initial_input(self) -> typing.Optional[builtins.bool]:
        return typing.cast(typing.Optional[builtins.bool], jsii.get(self, "initialInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="rolesInput")
    def roles_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "rolesInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="transitionsInput")
    def transitions_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "transitionsInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="description")
    def description(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "description"))

    @description.setter
    def description(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "description", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="initial")
    def initial(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "initial"))

    @initial.setter
    def initial(self, value: builtins.bool) -> None:
        jsii.set(self, "initial", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @key.setter
    def key(self, value: builtins.str) -> None:
        jsii.set(self, "key", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "name"))

    @name.setter
    def name(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="roles")
    def roles(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "roles"))

    @roles.setter
    def roles(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "roles", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="transitions")
    def transitions(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "transitions"))

    @transitions.setter
    def transitions(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "transitions", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="type")
    def type(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "type"))

    @type.setter
    def type(self, value: builtins.str) -> None:
        jsii.set(self, "type", value)


@jsii.data_type(
    jsii_type="labd_commercetools.StateConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "key": "key",
        "type": "type",
        "description": "description",
        "initial": "initial",
        "name": "name",
        "roles": "roles",
        "transitions": "transitions",
    },
)
class StateConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        key: builtins.str,
        type: builtins.str,
        description: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        initial: typing.Optional[builtins.bool] = None,
        name: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        roles: typing.Optional[typing.List[builtins.str]] = None,
        transitions: typing.Optional[typing.List[builtins.str]] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param key: -
        :param type: -
        :param description: -
        :param initial: -
        :param name: -
        :param roles: -
        :param transitions: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "key": key,
            "type": type,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if description is not None:
            self._values["description"] = description
        if initial is not None:
            self._values["initial"] = initial
        if name is not None:
            self._values["name"] = name
        if roles is not None:
            self._values["roles"] = roles
        if transitions is not None:
            self._values["transitions"] = transitions

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def key(self) -> builtins.str:
        result = self._values.get("key")
        assert result is not None, "Required property 'key' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def type(self) -> builtins.str:
        result = self._values.get("type")
        assert result is not None, "Required property 'type' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def description(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("description")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def initial(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("initial")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def name(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("name")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def roles(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("roles")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def transitions(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("transitions")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "StateConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class Store(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.Store",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        key: builtins.str,
        distribution_channels: typing.Optional[typing.List[builtins.str]] = None,
        languages: typing.Optional[typing.List[builtins.str]] = None,
        name: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        supply_channels: typing.Optional[typing.List[builtins.str]] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param key: -
        :param distribution_channels: -
        :param languages: -
        :param name: -
        :param supply_channels: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = StoreConfig(
            key=key,
            distribution_channels=distribution_channels,
            languages=languages,
            name=name,
            supply_channels=supply_channels,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(Store, self, [scope, id, config])

    @jsii.member(jsii_name="resetDistributionChannels")
    def reset_distribution_channels(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDistributionChannels", []))

    @jsii.member(jsii_name="resetLanguages")
    def reset_languages(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetLanguages", []))

    @jsii.member(jsii_name="resetName")
    def reset_name(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetName", []))

    @jsii.member(jsii_name="resetSupplyChannels")
    def reset_supply_channels(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetSupplyChannels", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="keyInput")
    def key_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "keyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="distributionChannelsInput")
    def distribution_channels_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "distributionChannelsInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="languagesInput")
    def languages_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "languagesInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="supplyChannelsInput")
    def supply_channels_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "supplyChannelsInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="distributionChannels")
    def distribution_channels(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "distributionChannels"))

    @distribution_channels.setter
    def distribution_channels(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "distributionChannels", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @key.setter
    def key(self, value: builtins.str) -> None:
        jsii.set(self, "key", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="languages")
    def languages(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "languages"))

    @languages.setter
    def languages(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "languages", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "name"))

    @name.setter
    def name(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="supplyChannels")
    def supply_channels(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "supplyChannels"))

    @supply_channels.setter
    def supply_channels(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "supplyChannels", value)


@jsii.data_type(
    jsii_type="labd_commercetools.StoreConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "key": "key",
        "distribution_channels": "distributionChannels",
        "languages": "languages",
        "name": "name",
        "supply_channels": "supplyChannels",
    },
)
class StoreConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        key: builtins.str,
        distribution_channels: typing.Optional[typing.List[builtins.str]] = None,
        languages: typing.Optional[typing.List[builtins.str]] = None,
        name: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        supply_channels: typing.Optional[typing.List[builtins.str]] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param key: -
        :param distribution_channels: -
        :param languages: -
        :param name: -
        :param supply_channels: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "key": key,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if distribution_channels is not None:
            self._values["distribution_channels"] = distribution_channels
        if languages is not None:
            self._values["languages"] = languages
        if name is not None:
            self._values["name"] = name
        if supply_channels is not None:
            self._values["supply_channels"] = supply_channels

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def key(self) -> builtins.str:
        result = self._values.get("key")
        assert result is not None, "Required property 'key' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def distribution_channels(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("distribution_channels")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def languages(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("languages")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def name(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("name")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def supply_channels(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("supply_channels")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "StoreConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class Subscription(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.Subscription",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        changes: typing.Optional[typing.List["SubscriptionChanges"]] = None,
        destination: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        format: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        key: typing.Optional[builtins.str] = None,
        message: typing.Optional[typing.List["SubscriptionMessage"]] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param changes: changes block.
        :param destination: -
        :param format: -
        :param key: -
        :param message: message block.
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = SubscriptionConfig(
            changes=changes,
            destination=destination,
            format=format,
            key=key,
            message=message,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(Subscription, self, [scope, id, config])

    @jsii.member(jsii_name="resetChanges")
    def reset_changes(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetChanges", []))

    @jsii.member(jsii_name="resetDestination")
    def reset_destination(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDestination", []))

    @jsii.member(jsii_name="resetFormat")
    def reset_format(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetFormat", []))

    @jsii.member(jsii_name="resetKey")
    def reset_key(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetKey", []))

    @jsii.member(jsii_name="resetMessage")
    def reset_message(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetMessage", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="changesInput")
    def changes_input(self) -> typing.Optional[typing.List["SubscriptionChanges"]]:
        return typing.cast(typing.Optional[typing.List["SubscriptionChanges"]], jsii.get(self, "changesInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="destinationInput")
    def destination_input(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "destinationInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="formatInput")
    def format_input(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "formatInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="keyInput")
    def key_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "keyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="messageInput")
    def message_input(self) -> typing.Optional[typing.List["SubscriptionMessage"]]:
        return typing.cast(typing.Optional[typing.List["SubscriptionMessage"]], jsii.get(self, "messageInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="changes")
    def changes(self) -> typing.List["SubscriptionChanges"]:
        return typing.cast(typing.List["SubscriptionChanges"], jsii.get(self, "changes"))

    @changes.setter
    def changes(self, value: typing.List["SubscriptionChanges"]) -> None:
        jsii.set(self, "changes", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="destination")
    def destination(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "destination"))

    @destination.setter
    def destination(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "destination", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="format")
    def format(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "format"))

    @format.setter
    def format(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "format", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @key.setter
    def key(self, value: builtins.str) -> None:
        jsii.set(self, "key", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="message")
    def message(self) -> typing.List["SubscriptionMessage"]:
        return typing.cast(typing.List["SubscriptionMessage"], jsii.get(self, "message"))

    @message.setter
    def message(self, value: typing.List["SubscriptionMessage"]) -> None:
        jsii.set(self, "message", value)


@jsii.data_type(
    jsii_type="labd_commercetools.SubscriptionChanges",
    jsii_struct_bases=[],
    name_mapping={"resource_type_ids": "resourceTypeIds"},
)
class SubscriptionChanges:
    def __init__(
        self,
        *,
        resource_type_ids: typing.Optional[typing.List[builtins.str]] = None,
    ) -> None:
        '''
        :param resource_type_ids: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if resource_type_ids is not None:
            self._values["resource_type_ids"] = resource_type_ids

    @builtins.property
    def resource_type_ids(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("resource_type_ids")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "SubscriptionChanges(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.SubscriptionConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "changes": "changes",
        "destination": "destination",
        "format": "format",
        "key": "key",
        "message": "message",
    },
)
class SubscriptionConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        changes: typing.Optional[typing.List[SubscriptionChanges]] = None,
        destination: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        format: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        key: typing.Optional[builtins.str] = None,
        message: typing.Optional[typing.List["SubscriptionMessage"]] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param changes: changes block.
        :param destination: -
        :param format: -
        :param key: -
        :param message: message block.
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {}
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if changes is not None:
            self._values["changes"] = changes
        if destination is not None:
            self._values["destination"] = destination
        if format is not None:
            self._values["format"] = format
        if key is not None:
            self._values["key"] = key
        if message is not None:
            self._values["message"] = message

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def changes(self) -> typing.Optional[typing.List[SubscriptionChanges]]:
        '''changes block.'''
        result = self._values.get("changes")
        return typing.cast(typing.Optional[typing.List[SubscriptionChanges]], result)

    @builtins.property
    def destination(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("destination")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def format(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("format")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def key(self) -> typing.Optional[builtins.str]:
        result = self._values.get("key")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def message(self) -> typing.Optional[typing.List["SubscriptionMessage"]]:
        '''message block.'''
        result = self._values.get("message")
        return typing.cast(typing.Optional[typing.List["SubscriptionMessage"]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "SubscriptionConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.SubscriptionMessage",
    jsii_struct_bases=[],
    name_mapping={"resource_type_id": "resourceTypeId", "types": "types"},
)
class SubscriptionMessage:
    def __init__(
        self,
        *,
        resource_type_id: typing.Optional[builtins.str] = None,
        types: typing.Optional[typing.List[builtins.str]] = None,
    ) -> None:
        '''
        :param resource_type_id: -
        :param types: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if resource_type_id is not None:
            self._values["resource_type_id"] = resource_type_id
        if types is not None:
            self._values["types"] = types

    @builtins.property
    def resource_type_id(self) -> typing.Optional[builtins.str]:
        result = self._values.get("resource_type_id")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def types(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("types")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "SubscriptionMessage(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class TaxCategory(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.TaxCategory",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        name: builtins.str,
        description: typing.Optional[builtins.str] = None,
        key: typing.Optional[builtins.str] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param name: -
        :param description: -
        :param key: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = TaxCategoryConfig(
            name=name,
            description=description,
            key=key,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(TaxCategory, self, [scope, id, config])

    @jsii.member(jsii_name="resetDescription")
    def reset_description(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDescription", []))

    @jsii.member(jsii_name="resetKey")
    def reset_key(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetKey", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="descriptionInput")
    def description_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "descriptionInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="keyInput")
    def key_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "keyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="description")
    def description(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "description"))

    @description.setter
    def description(self, value: builtins.str) -> None:
        jsii.set(self, "description", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @key.setter
    def key(self, value: builtins.str) -> None:
        jsii.set(self, "key", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)


@jsii.data_type(
    jsii_type="labd_commercetools.TaxCategoryConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "name": "name",
        "description": "description",
        "key": "key",
    },
)
class TaxCategoryConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        name: builtins.str,
        description: typing.Optional[builtins.str] = None,
        key: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param name: -
        :param description: -
        :param key: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "name": name,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if description is not None:
            self._values["description"] = description
        if key is not None:
            self._values["key"] = key

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def description(self) -> typing.Optional[builtins.str]:
        result = self._values.get("description")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def key(self) -> typing.Optional[builtins.str]:
        result = self._values.get("key")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "TaxCategoryConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class TaxCategoryRate(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.TaxCategoryRate",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        country: builtins.str,
        included_in_price: builtins.bool,
        name: builtins.str,
        tax_category_id: builtins.str,
        amount: typing.Optional[jsii.Number] = None,
        state: typing.Optional[builtins.str] = None,
        sub_rate: typing.Optional[typing.List["TaxCategoryRateSubRate"]] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param country: -
        :param included_in_price: -
        :param name: -
        :param tax_category_id: -
        :param amount: -
        :param state: -
        :param sub_rate: sub_rate block.
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = TaxCategoryRateConfig(
            country=country,
            included_in_price=included_in_price,
            name=name,
            tax_category_id=tax_category_id,
            amount=amount,
            state=state,
            sub_rate=sub_rate,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(TaxCategoryRate, self, [scope, id, config])

    @jsii.member(jsii_name="resetAmount")
    def reset_amount(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetAmount", []))

    @jsii.member(jsii_name="resetState")
    def reset_state(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetState", []))

    @jsii.member(jsii_name="resetSubRate")
    def reset_sub_rate(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetSubRate", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="countryInput")
    def country_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "countryInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="includedInPriceInput")
    def included_in_price_input(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "includedInPriceInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="taxCategoryIdInput")
    def tax_category_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "taxCategoryIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="amountInput")
    def amount_input(self) -> typing.Optional[jsii.Number]:
        return typing.cast(typing.Optional[jsii.Number], jsii.get(self, "amountInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="stateInput")
    def state_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "stateInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="subRateInput")
    def sub_rate_input(self) -> typing.Optional[typing.List["TaxCategoryRateSubRate"]]:
        return typing.cast(typing.Optional[typing.List["TaxCategoryRateSubRate"]], jsii.get(self, "subRateInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="amount")
    def amount(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "amount"))

    @amount.setter
    def amount(self, value: jsii.Number) -> None:
        jsii.set(self, "amount", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="country")
    def country(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "country"))

    @country.setter
    def country(self, value: builtins.str) -> None:
        jsii.set(self, "country", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="includedInPrice")
    def included_in_price(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "includedInPrice"))

    @included_in_price.setter
    def included_in_price(self, value: builtins.bool) -> None:
        jsii.set(self, "includedInPrice", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="state")
    def state(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "state"))

    @state.setter
    def state(self, value: builtins.str) -> None:
        jsii.set(self, "state", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="subRate")
    def sub_rate(self) -> typing.List["TaxCategoryRateSubRate"]:
        return typing.cast(typing.List["TaxCategoryRateSubRate"], jsii.get(self, "subRate"))

    @sub_rate.setter
    def sub_rate(self, value: typing.List["TaxCategoryRateSubRate"]) -> None:
        jsii.set(self, "subRate", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="taxCategoryId")
    def tax_category_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "taxCategoryId"))

    @tax_category_id.setter
    def tax_category_id(self, value: builtins.str) -> None:
        jsii.set(self, "taxCategoryId", value)


@jsii.data_type(
    jsii_type="labd_commercetools.TaxCategoryRateConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "country": "country",
        "included_in_price": "includedInPrice",
        "name": "name",
        "tax_category_id": "taxCategoryId",
        "amount": "amount",
        "state": "state",
        "sub_rate": "subRate",
    },
)
class TaxCategoryRateConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        country: builtins.str,
        included_in_price: builtins.bool,
        name: builtins.str,
        tax_category_id: builtins.str,
        amount: typing.Optional[jsii.Number] = None,
        state: typing.Optional[builtins.str] = None,
        sub_rate: typing.Optional[typing.List["TaxCategoryRateSubRate"]] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param country: -
        :param included_in_price: -
        :param name: -
        :param tax_category_id: -
        :param amount: -
        :param state: -
        :param sub_rate: sub_rate block.
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "country": country,
            "included_in_price": included_in_price,
            "name": name,
            "tax_category_id": tax_category_id,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if amount is not None:
            self._values["amount"] = amount
        if state is not None:
            self._values["state"] = state
        if sub_rate is not None:
            self._values["sub_rate"] = sub_rate

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def country(self) -> builtins.str:
        result = self._values.get("country")
        assert result is not None, "Required property 'country' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def included_in_price(self) -> builtins.bool:
        result = self._values.get("included_in_price")
        assert result is not None, "Required property 'included_in_price' is missing"
        return typing.cast(builtins.bool, result)

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def tax_category_id(self) -> builtins.str:
        result = self._values.get("tax_category_id")
        assert result is not None, "Required property 'tax_category_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def amount(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("amount")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def state(self) -> typing.Optional[builtins.str]:
        result = self._values.get("state")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def sub_rate(self) -> typing.Optional[typing.List["TaxCategoryRateSubRate"]]:
        '''sub_rate block.'''
        result = self._values.get("sub_rate")
        return typing.cast(typing.Optional[typing.List["TaxCategoryRateSubRate"]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "TaxCategoryRateConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.TaxCategoryRateSubRate",
    jsii_struct_bases=[],
    name_mapping={"amount": "amount", "name": "name"},
)
class TaxCategoryRateSubRate:
    def __init__(self, *, amount: jsii.Number, name: builtins.str) -> None:
        '''
        :param amount: -
        :param name: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "amount": amount,
            "name": name,
        }

    @builtins.property
    def amount(self) -> jsii.Number:
        result = self._values.get("amount")
        assert result is not None, "Required property 'amount' is missing"
        return typing.cast(jsii.Number, result)

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "TaxCategoryRateSubRate(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class Type(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_commercetools.Type",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        key: builtins.str,
        name: typing.Mapping[builtins.str, builtins.str],
        resource_type_ids: typing.List[builtins.str],
        description: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        field: typing.Optional[typing.List["TypeField"]] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param key: -
        :param name: -
        :param resource_type_ids: -
        :param description: -
        :param field: field block.
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = TypeConfig(
            key=key,
            name=name,
            resource_type_ids=resource_type_ids,
            description=description,
            field=field,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(Type, self, [scope, id, config])

    @jsii.member(jsii_name="resetDescription")
    def reset_description(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDescription", []))

    @jsii.member(jsii_name="resetField")
    def reset_field(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetField", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="keyInput")
    def key_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "keyInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="resourceTypeIdsInput")
    def resource_type_ids_input(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "resourceTypeIdsInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="descriptionInput")
    def description_input(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "descriptionInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="fieldInput")
    def field_input(self) -> typing.Optional[typing.List["TypeField"]]:
        return typing.cast(typing.Optional[typing.List["TypeField"]], jsii.get(self, "fieldInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="description")
    def description(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "description"))

    @description.setter
    def description(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "description", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="field")
    def field(self) -> typing.List["TypeField"]:
        return typing.cast(typing.List["TypeField"], jsii.get(self, "field"))

    @field.setter
    def field(self, value: typing.List["TypeField"]) -> None:
        jsii.set(self, "field", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="key")
    def key(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "key"))

    @key.setter
    def key(self, value: builtins.str) -> None:
        jsii.set(self, "key", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "name"))

    @name.setter
    def name(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="resourceTypeIds")
    def resource_type_ids(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "resourceTypeIds"))

    @resource_type_ids.setter
    def resource_type_ids(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "resourceTypeIds", value)


@jsii.data_type(
    jsii_type="labd_commercetools.TypeConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "key": "key",
        "name": "name",
        "resource_type_ids": "resourceTypeIds",
        "description": "description",
        "field": "field",
    },
)
class TypeConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        key: builtins.str,
        name: typing.Mapping[builtins.str, builtins.str],
        resource_type_ids: typing.List[builtins.str],
        description: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        field: typing.Optional[typing.List["TypeField"]] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param key: -
        :param name: -
        :param resource_type_ids: -
        :param description: -
        :param field: field block.
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "key": key,
            "name": name,
            "resource_type_ids": resource_type_ids,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if description is not None:
            self._values["description"] = description
        if field is not None:
            self._values["field"] = field

    @builtins.property
    def count(self) -> typing.Optional[jsii.Number]:
        '''
        :stability: experimental
        '''
        result = self._values.get("count")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def depends_on(self) -> typing.Optional[typing.List[cdktf.ITerraformDependable]]:
        '''
        :stability: experimental
        '''
        result = self._values.get("depends_on")
        return typing.cast(typing.Optional[typing.List[cdktf.ITerraformDependable]], result)

    @builtins.property
    def lifecycle(self) -> typing.Optional[cdktf.TerraformResourceLifecycle]:
        '''
        :stability: experimental
        '''
        result = self._values.get("lifecycle")
        return typing.cast(typing.Optional[cdktf.TerraformResourceLifecycle], result)

    @builtins.property
    def provider(self) -> typing.Optional[cdktf.TerraformProvider]:
        '''
        :stability: experimental
        '''
        result = self._values.get("provider")
        return typing.cast(typing.Optional[cdktf.TerraformProvider], result)

    @builtins.property
    def key(self) -> builtins.str:
        result = self._values.get("key")
        assert result is not None, "Required property 'key' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def name(self) -> typing.Mapping[builtins.str, builtins.str]:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(typing.Mapping[builtins.str, builtins.str], result)

    @builtins.property
    def resource_type_ids(self) -> typing.List[builtins.str]:
        result = self._values.get("resource_type_ids")
        assert result is not None, "Required property 'resource_type_ids' is missing"
        return typing.cast(typing.List[builtins.str], result)

    @builtins.property
    def description(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("description")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def field(self) -> typing.Optional[typing.List["TypeField"]]:
        '''field block.'''
        result = self._values.get("field")
        return typing.cast(typing.Optional[typing.List["TypeField"]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "TypeConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.TypeField",
    jsii_struct_bases=[],
    name_mapping={
        "label": "label",
        "name": "name",
        "type": "type",
        "input_hint": "inputHint",
        "required": "required",
    },
)
class TypeField:
    def __init__(
        self,
        *,
        label: typing.Mapping[builtins.str, builtins.str],
        name: builtins.str,
        type: typing.List["TypeFieldType"],
        input_hint: typing.Optional[builtins.str] = None,
        required: typing.Optional[builtins.bool] = None,
    ) -> None:
        '''
        :param label: -
        :param name: -
        :param type: type block.
        :param input_hint: -
        :param required: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "label": label,
            "name": name,
            "type": type,
        }
        if input_hint is not None:
            self._values["input_hint"] = input_hint
        if required is not None:
            self._values["required"] = required

    @builtins.property
    def label(self) -> typing.Mapping[builtins.str, builtins.str]:
        result = self._values.get("label")
        assert result is not None, "Required property 'label' is missing"
        return typing.cast(typing.Mapping[builtins.str, builtins.str], result)

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def type(self) -> typing.List["TypeFieldType"]:
        '''type block.'''
        result = self._values.get("type")
        assert result is not None, "Required property 'type' is missing"
        return typing.cast(typing.List["TypeFieldType"], result)

    @builtins.property
    def input_hint(self) -> typing.Optional[builtins.str]:
        result = self._values.get("input_hint")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def required(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("required")
        return typing.cast(typing.Optional[builtins.bool], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "TypeField(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.TypeFieldType",
    jsii_struct_bases=[],
    name_mapping={
        "name": "name",
        "element_type": "elementType",
        "localized_value": "localizedValue",
        "reference_type_id": "referenceTypeId",
        "values": "values",
    },
)
class TypeFieldType:
    def __init__(
        self,
        *,
        name: builtins.str,
        element_type: typing.Optional[typing.List["TypeFieldTypeElementType"]] = None,
        localized_value: typing.Optional[typing.List["TypeFieldTypeLocalizedValue"]] = None,
        reference_type_id: typing.Optional[builtins.str] = None,
        values: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
    ) -> None:
        '''
        :param name: -
        :param element_type: element_type block.
        :param localized_value: localized_value block.
        :param reference_type_id: -
        :param values: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "name": name,
        }
        if element_type is not None:
            self._values["element_type"] = element_type
        if localized_value is not None:
            self._values["localized_value"] = localized_value
        if reference_type_id is not None:
            self._values["reference_type_id"] = reference_type_id
        if values is not None:
            self._values["values"] = values

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def element_type(self) -> typing.Optional[typing.List["TypeFieldTypeElementType"]]:
        '''element_type block.'''
        result = self._values.get("element_type")
        return typing.cast(typing.Optional[typing.List["TypeFieldTypeElementType"]], result)

    @builtins.property
    def localized_value(
        self,
    ) -> typing.Optional[typing.List["TypeFieldTypeLocalizedValue"]]:
        '''localized_value block.'''
        result = self._values.get("localized_value")
        return typing.cast(typing.Optional[typing.List["TypeFieldTypeLocalizedValue"]], result)

    @builtins.property
    def reference_type_id(self) -> typing.Optional[builtins.str]:
        result = self._values.get("reference_type_id")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def values(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("values")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "TypeFieldType(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.TypeFieldTypeElementType",
    jsii_struct_bases=[],
    name_mapping={
        "name": "name",
        "localized_value": "localizedValue",
        "reference_type_id": "referenceTypeId",
        "values": "values",
    },
)
class TypeFieldTypeElementType:
    def __init__(
        self,
        *,
        name: builtins.str,
        localized_value: typing.Optional[typing.List["TypeFieldTypeElementTypeLocalizedValue"]] = None,
        reference_type_id: typing.Optional[builtins.str] = None,
        values: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
    ) -> None:
        '''
        :param name: -
        :param localized_value: localized_value block.
        :param reference_type_id: -
        :param values: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "name": name,
        }
        if localized_value is not None:
            self._values["localized_value"] = localized_value
        if reference_type_id is not None:
            self._values["reference_type_id"] = reference_type_id
        if values is not None:
            self._values["values"] = values

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def localized_value(
        self,
    ) -> typing.Optional[typing.List["TypeFieldTypeElementTypeLocalizedValue"]]:
        '''localized_value block.'''
        result = self._values.get("localized_value")
        return typing.cast(typing.Optional[typing.List["TypeFieldTypeElementTypeLocalizedValue"]], result)

    @builtins.property
    def reference_type_id(self) -> typing.Optional[builtins.str]:
        result = self._values.get("reference_type_id")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def values(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("values")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "TypeFieldTypeElementType(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.TypeFieldTypeElementTypeLocalizedValue",
    jsii_struct_bases=[],
    name_mapping={"key": "key", "label": "label"},
)
class TypeFieldTypeElementTypeLocalizedValue:
    def __init__(
        self,
        *,
        key: builtins.str,
        label: typing.Mapping[builtins.str, builtins.str],
    ) -> None:
        '''
        :param key: -
        :param label: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "key": key,
            "label": label,
        }

    @builtins.property
    def key(self) -> builtins.str:
        result = self._values.get("key")
        assert result is not None, "Required property 'key' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def label(self) -> typing.Mapping[builtins.str, builtins.str]:
        result = self._values.get("label")
        assert result is not None, "Required property 'label' is missing"
        return typing.cast(typing.Mapping[builtins.str, builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "TypeFieldTypeElementTypeLocalizedValue(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_commercetools.TypeFieldTypeLocalizedValue",
    jsii_struct_bases=[],
    name_mapping={"key": "key", "label": "label"},
)
class TypeFieldTypeLocalizedValue:
    def __init__(
        self,
        *,
        key: builtins.str,
        label: typing.Mapping[builtins.str, builtins.str],
    ) -> None:
        '''
        :param key: -
        :param label: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "key": key,
            "label": label,
        }

    @builtins.property
    def key(self) -> builtins.str:
        result = self._values.get("key")
        assert result is not None, "Required property 'key' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def label(self) -> typing.Mapping[builtins.str, builtins.str]:
        result = self._values.get("label")
        assert result is not None, "Required property 'label' is missing"
        return typing.cast(typing.Mapping[builtins.str, builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "TypeFieldTypeLocalizedValue(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


__all__ = [
    "ApiClient",
    "ApiClientConfig",
    "ApiExtension",
    "ApiExtensionConfig",
    "ApiExtensionTrigger",
    "CartDiscount",
    "CartDiscountConfig",
    "CartDiscountValue",
    "CartDiscountValueMoney",
    "Channel",
    "ChannelConfig",
    "CommercetoolsProvider",
    "CommercetoolsProviderConfig",
    "CustomObject",
    "CustomObjectConfig",
    "DiscountCode",
    "DiscountCodeConfig",
    "ProductType",
    "ProductTypeAttribute",
    "ProductTypeAttributeType",
    "ProductTypeAttributeTypeElementType",
    "ProductTypeAttributeTypeElementTypeLocalizedValue",
    "ProductTypeAttributeTypeLocalizedValue",
    "ProductTypeConfig",
    "ProjectSettings",
    "ProjectSettingsConfig",
    "ShippingMethod",
    "ShippingMethodConfig",
    "ShippingZone",
    "ShippingZoneConfig",
    "ShippingZoneLocation",
    "ShippingZoneRate",
    "ShippingZoneRateConfig",
    "ShippingZoneRateFreeAbove",
    "ShippingZoneRatePrice",
    "State",
    "StateConfig",
    "Store",
    "StoreConfig",
    "Subscription",
    "SubscriptionChanges",
    "SubscriptionConfig",
    "SubscriptionMessage",
    "TaxCategory",
    "TaxCategoryConfig",
    "TaxCategoryRate",
    "TaxCategoryRateConfig",
    "TaxCategoryRateSubRate",
    "Type",
    "TypeConfig",
    "TypeField",
    "TypeFieldType",
    "TypeFieldTypeElementType",
    "TypeFieldTypeElementTypeLocalizedValue",
    "TypeFieldTypeLocalizedValue",
]

publication.publish()
