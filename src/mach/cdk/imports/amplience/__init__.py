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


class AmplienceProvider(
    cdktf.TerraformProvider,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_amplience.AmplienceProvider",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        client_id: builtins.str,
        client_secret: builtins.str,
        hub_id: builtins.str,
        alias: typing.Optional[builtins.str] = None,
        auth_url: typing.Optional[builtins.str] = None,
        content_api_url: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param client_id: The OAuth Client ID for the Amplience management API https://amplience_provider.com/docs/api/dynamic-content/management/index.html#section/Authentication.
        :param client_secret: The OAuth Client Secret for Amplience management API. https://amplience_provider.com/docs/api/dynamic-content/management/index.html#section/Authentication
        :param hub_id: The Hub ID of the Amplience Hub to use this provider instance with.
        :param alias: Alias name.
        :param auth_url: The Amplience authentication URL.
        :param content_api_url: The base URL path for the Amplience Content API.
        '''
        config = AmplienceProviderConfig(
            client_id=client_id,
            client_secret=client_secret,
            hub_id=hub_id,
            alias=alias,
            auth_url=auth_url,
            content_api_url=content_api_url,
        )

        jsii.create(AmplienceProvider, self, [scope, id, config])

    @jsii.member(jsii_name="resetAlias")
    def reset_alias(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetAlias", []))

    @jsii.member(jsii_name="resetAuthUrl")
    def reset_auth_url(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetAuthUrl", []))

    @jsii.member(jsii_name="resetContentApiUrl")
    def reset_content_api_url(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetContentApiUrl", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="clientIdInput")
    def client_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "clientIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="clientSecretInput")
    def client_secret_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "clientSecretInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="hubIdInput")
    def hub_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "hubIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="aliasInput")
    def alias_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "aliasInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="authUrlInput")
    def auth_url_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "authUrlInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="contentApiUrlInput")
    def content_api_url_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "contentApiUrlInput"))

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
    @jsii.member(jsii_name="hubId")
    def hub_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "hubId"))

    @hub_id.setter
    def hub_id(self, value: builtins.str) -> None:
        jsii.set(self, "hubId", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="alias")
    def alias(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "alias"))

    @alias.setter
    def alias(self, value: typing.Optional[builtins.str]) -> None:
        jsii.set(self, "alias", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="authUrl")
    def auth_url(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "authUrl"))

    @auth_url.setter
    def auth_url(self, value: typing.Optional[builtins.str]) -> None:
        jsii.set(self, "authUrl", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="contentApiUrl")
    def content_api_url(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "contentApiUrl"))

    @content_api_url.setter
    def content_api_url(self, value: typing.Optional[builtins.str]) -> None:
        jsii.set(self, "contentApiUrl", value)


@jsii.data_type(
    jsii_type="labd_amplience.AmplienceProviderConfig",
    jsii_struct_bases=[],
    name_mapping={
        "client_id": "clientId",
        "client_secret": "clientSecret",
        "hub_id": "hubId",
        "alias": "alias",
        "auth_url": "authUrl",
        "content_api_url": "contentApiUrl",
    },
)
class AmplienceProviderConfig:
    def __init__(
        self,
        *,
        client_id: builtins.str,
        client_secret: builtins.str,
        hub_id: builtins.str,
        alias: typing.Optional[builtins.str] = None,
        auth_url: typing.Optional[builtins.str] = None,
        content_api_url: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param client_id: The OAuth Client ID for the Amplience management API https://amplience_provider.com/docs/api/dynamic-content/management/index.html#section/Authentication.
        :param client_secret: The OAuth Client Secret for Amplience management API. https://amplience_provider.com/docs/api/dynamic-content/management/index.html#section/Authentication
        :param hub_id: The Hub ID of the Amplience Hub to use this provider instance with.
        :param alias: Alias name.
        :param auth_url: The Amplience authentication URL.
        :param content_api_url: The base URL path for the Amplience Content API.
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "client_id": client_id,
            "client_secret": client_secret,
            "hub_id": hub_id,
        }
        if alias is not None:
            self._values["alias"] = alias
        if auth_url is not None:
            self._values["auth_url"] = auth_url
        if content_api_url is not None:
            self._values["content_api_url"] = content_api_url

    @builtins.property
    def client_id(self) -> builtins.str:
        '''The OAuth Client ID for the Amplience management API https://amplience_provider.com/docs/api/dynamic-content/management/index.html#section/Authentication.'''
        result = self._values.get("client_id")
        assert result is not None, "Required property 'client_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def client_secret(self) -> builtins.str:
        '''The OAuth Client Secret for Amplience management API.

        https://amplience_provider.com/docs/api/dynamic-content/management/index.html#section/Authentication
        '''
        result = self._values.get("client_secret")
        assert result is not None, "Required property 'client_secret' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def hub_id(self) -> builtins.str:
        '''The Hub ID of the Amplience Hub to use this provider instance with.'''
        result = self._values.get("hub_id")
        assert result is not None, "Required property 'hub_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def alias(self) -> typing.Optional[builtins.str]:
        '''Alias name.'''
        result = self._values.get("alias")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def auth_url(self) -> typing.Optional[builtins.str]:
        '''The Amplience authentication URL.'''
        result = self._values.get("auth_url")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def content_api_url(self) -> typing.Optional[builtins.str]:
        '''The base URL path for the Amplience Content API.'''
        result = self._values.get("content_api_url")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "AmplienceProviderConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class ContentRepository(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_amplience.ContentRepository",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        label: builtins.str,
        name: builtins.str,
        content_types: typing.Optional[typing.List["ContentRepositoryContentTypes"]] = None,
        features: typing.Optional[typing.List[builtins.str]] = None,
        item_locales: typing.Optional[typing.List[builtins.str]] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param label: -
        :param name: -
        :param content_types: content_types block.
        :param features: -
        :param item_locales: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = ContentRepositoryConfig(
            label=label,
            name=name,
            content_types=content_types,
            features=features,
            item_locales=item_locales,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(ContentRepository, self, [scope, id, config])

    @jsii.member(jsii_name="resetContentTypes")
    def reset_content_types(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetContentTypes", []))

    @jsii.member(jsii_name="resetFeatures")
    def reset_features(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetFeatures", []))

    @jsii.member(jsii_name="resetItemLocales")
    def reset_item_locales(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetItemLocales", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="labelInput")
    def label_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "labelInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="status")
    def status(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "status"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="type")
    def type(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "type"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="contentTypesInput")
    def content_types_input(
        self,
    ) -> typing.Optional[typing.List["ContentRepositoryContentTypes"]]:
        return typing.cast(typing.Optional[typing.List["ContentRepositoryContentTypes"]], jsii.get(self, "contentTypesInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="featuresInput")
    def features_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "featuresInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="itemLocalesInput")
    def item_locales_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "itemLocalesInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="contentTypes")
    def content_types(self) -> typing.List["ContentRepositoryContentTypes"]:
        return typing.cast(typing.List["ContentRepositoryContentTypes"], jsii.get(self, "contentTypes"))

    @content_types.setter
    def content_types(
        self,
        value: typing.List["ContentRepositoryContentTypes"],
    ) -> None:
        jsii.set(self, "contentTypes", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="features")
    def features(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "features"))

    @features.setter
    def features(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "features", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="itemLocales")
    def item_locales(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "itemLocales"))

    @item_locales.setter
    def item_locales(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "itemLocales", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="label")
    def label(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "label"))

    @label.setter
    def label(self, value: builtins.str) -> None:
        jsii.set(self, "label", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)


@jsii.data_type(
    jsii_type="labd_amplience.ContentRepositoryConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "label": "label",
        "name": "name",
        "content_types": "contentTypes",
        "features": "features",
        "item_locales": "itemLocales",
    },
)
class ContentRepositoryConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        label: builtins.str,
        name: builtins.str,
        content_types: typing.Optional[typing.List["ContentRepositoryContentTypes"]] = None,
        features: typing.Optional[typing.List[builtins.str]] = None,
        item_locales: typing.Optional[typing.List[builtins.str]] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param label: -
        :param name: -
        :param content_types: content_types block.
        :param features: -
        :param item_locales: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "label": label,
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
        if content_types is not None:
            self._values["content_types"] = content_types
        if features is not None:
            self._values["features"] = features
        if item_locales is not None:
            self._values["item_locales"] = item_locales

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
    def label(self) -> builtins.str:
        result = self._values.get("label")
        assert result is not None, "Required property 'label' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def content_types(
        self,
    ) -> typing.Optional[typing.List["ContentRepositoryContentTypes"]]:
        '''content_types block.'''
        result = self._values.get("content_types")
        return typing.cast(typing.Optional[typing.List["ContentRepositoryContentTypes"]], result)

    @builtins.property
    def features(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("features")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def item_locales(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("item_locales")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContentRepositoryConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_amplience.ContentRepositoryContentTypes",
    jsii_struct_bases=[],
    name_mapping={
        "content_type_uri": "contentTypeUri",
        "hub_content_type_id": "hubContentTypeId",
    },
)
class ContentRepositoryContentTypes:
    def __init__(
        self,
        *,
        content_type_uri: builtins.str,
        hub_content_type_id: builtins.str,
    ) -> None:
        '''
        :param content_type_uri: -
        :param hub_content_type_id: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "content_type_uri": content_type_uri,
            "hub_content_type_id": hub_content_type_id,
        }

    @builtins.property
    def content_type_uri(self) -> builtins.str:
        result = self._values.get("content_type_uri")
        assert result is not None, "Required property 'content_type_uri' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def hub_content_type_id(self) -> builtins.str:
        result = self._values.get("hub_content_type_id")
        assert result is not None, "Required property 'hub_content_type_id' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContentRepositoryContentTypes(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class Webhook(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_amplience.Webhook",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        label: builtins.str,
        method: builtins.str,
        active: typing.Optional[builtins.bool] = None,
        custom_payload: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        events: typing.Optional[typing.List[builtins.str]] = None,
        filter: typing.Optional[typing.List["WebhookFilter"]] = None,
        handlers: typing.Optional[typing.List[builtins.str]] = None,
        header: typing.Optional[typing.List["WebhookHeader"]] = None,
        notifications: typing.Optional[typing.List["WebhookNotifications"]] = None,
        secret: typing.Optional[builtins.str] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param label: -
        :param method: -
        :param active: -
        :param custom_payload: -
        :param events: -
        :param filter: filter block.
        :param handlers: -
        :param header: header block.
        :param notifications: notifications block.
        :param secret: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = WebhookConfig(
            label=label,
            method=method,
            active=active,
            custom_payload=custom_payload,
            events=events,
            filter=filter,
            handlers=handlers,
            header=header,
            notifications=notifications,
            secret=secret,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(Webhook, self, [scope, id, config])

    @jsii.member(jsii_name="resetActive")
    def reset_active(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetActive", []))

    @jsii.member(jsii_name="resetCustomPayload")
    def reset_custom_payload(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetCustomPayload", []))

    @jsii.member(jsii_name="resetEvents")
    def reset_events(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetEvents", []))

    @jsii.member(jsii_name="resetFilter")
    def reset_filter(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetFilter", []))

    @jsii.member(jsii_name="resetHandlers")
    def reset_handlers(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetHandlers", []))

    @jsii.member(jsii_name="resetHeader")
    def reset_header(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetHeader", []))

    @jsii.member(jsii_name="resetNotifications")
    def reset_notifications(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetNotifications", []))

    @jsii.member(jsii_name="resetSecret")
    def reset_secret(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetSecret", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="labelInput")
    def label_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "labelInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="methodInput")
    def method_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "methodInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="activeInput")
    def active_input(self) -> typing.Optional[builtins.bool]:
        return typing.cast(typing.Optional[builtins.bool], jsii.get(self, "activeInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="customPayloadInput")
    def custom_payload_input(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "customPayloadInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="eventsInput")
    def events_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "eventsInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="filterInput")
    def filter_input(self) -> typing.Optional[typing.List["WebhookFilter"]]:
        return typing.cast(typing.Optional[typing.List["WebhookFilter"]], jsii.get(self, "filterInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="handlersInput")
    def handlers_input(self) -> typing.Optional[typing.List[builtins.str]]:
        return typing.cast(typing.Optional[typing.List[builtins.str]], jsii.get(self, "handlersInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="headerInput")
    def header_input(self) -> typing.Optional[typing.List["WebhookHeader"]]:
        return typing.cast(typing.Optional[typing.List["WebhookHeader"]], jsii.get(self, "headerInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="notificationsInput")
    def notifications_input(
        self,
    ) -> typing.Optional[typing.List["WebhookNotifications"]]:
        return typing.cast(typing.Optional[typing.List["WebhookNotifications"]], jsii.get(self, "notificationsInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="secretInput")
    def secret_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "secretInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="active")
    def active(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "active"))

    @active.setter
    def active(self, value: builtins.bool) -> None:
        jsii.set(self, "active", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="customPayload")
    def custom_payload(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "customPayload"))

    @custom_payload.setter
    def custom_payload(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "customPayload", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="events")
    def events(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "events"))

    @events.setter
    def events(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "events", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="filter")
    def filter(self) -> typing.List["WebhookFilter"]:
        return typing.cast(typing.List["WebhookFilter"], jsii.get(self, "filter"))

    @filter.setter
    def filter(self, value: typing.List["WebhookFilter"]) -> None:
        jsii.set(self, "filter", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="handlers")
    def handlers(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "handlers"))

    @handlers.setter
    def handlers(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "handlers", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="header")
    def header(self) -> typing.List["WebhookHeader"]:
        return typing.cast(typing.List["WebhookHeader"], jsii.get(self, "header"))

    @header.setter
    def header(self, value: typing.List["WebhookHeader"]) -> None:
        jsii.set(self, "header", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="label")
    def label(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "label"))

    @label.setter
    def label(self, value: builtins.str) -> None:
        jsii.set(self, "label", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="method")
    def method(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "method"))

    @method.setter
    def method(self, value: builtins.str) -> None:
        jsii.set(self, "method", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="notifications")
    def notifications(self) -> typing.List["WebhookNotifications"]:
        return typing.cast(typing.List["WebhookNotifications"], jsii.get(self, "notifications"))

    @notifications.setter
    def notifications(self, value: typing.List["WebhookNotifications"]) -> None:
        jsii.set(self, "notifications", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="secret")
    def secret(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "secret"))

    @secret.setter
    def secret(self, value: builtins.str) -> None:
        jsii.set(self, "secret", value)


@jsii.data_type(
    jsii_type="labd_amplience.WebhookConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "label": "label",
        "method": "method",
        "active": "active",
        "custom_payload": "customPayload",
        "events": "events",
        "filter": "filter",
        "handlers": "handlers",
        "header": "header",
        "notifications": "notifications",
        "secret": "secret",
    },
)
class WebhookConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        label: builtins.str,
        method: builtins.str,
        active: typing.Optional[builtins.bool] = None,
        custom_payload: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        events: typing.Optional[typing.List[builtins.str]] = None,
        filter: typing.Optional[typing.List["WebhookFilter"]] = None,
        handlers: typing.Optional[typing.List[builtins.str]] = None,
        header: typing.Optional[typing.List["WebhookHeader"]] = None,
        notifications: typing.Optional[typing.List["WebhookNotifications"]] = None,
        secret: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param label: -
        :param method: -
        :param active: -
        :param custom_payload: -
        :param events: -
        :param filter: filter block.
        :param handlers: -
        :param header: header block.
        :param notifications: notifications block.
        :param secret: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "label": label,
            "method": method,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if active is not None:
            self._values["active"] = active
        if custom_payload is not None:
            self._values["custom_payload"] = custom_payload
        if events is not None:
            self._values["events"] = events
        if filter is not None:
            self._values["filter"] = filter
        if handlers is not None:
            self._values["handlers"] = handlers
        if header is not None:
            self._values["header"] = header
        if notifications is not None:
            self._values["notifications"] = notifications
        if secret is not None:
            self._values["secret"] = secret

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
    def label(self) -> builtins.str:
        result = self._values.get("label")
        assert result is not None, "Required property 'label' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def method(self) -> builtins.str:
        result = self._values.get("method")
        assert result is not None, "Required property 'method' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def active(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("active")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def custom_payload(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("custom_payload")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def events(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("events")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def filter(self) -> typing.Optional[typing.List["WebhookFilter"]]:
        '''filter block.'''
        result = self._values.get("filter")
        return typing.cast(typing.Optional[typing.List["WebhookFilter"]], result)

    @builtins.property
    def handlers(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("handlers")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def header(self) -> typing.Optional[typing.List["WebhookHeader"]]:
        '''header block.'''
        result = self._values.get("header")
        return typing.cast(typing.Optional[typing.List["WebhookHeader"]], result)

    @builtins.property
    def notifications(self) -> typing.Optional[typing.List["WebhookNotifications"]]:
        '''notifications block.'''
        result = self._values.get("notifications")
        return typing.cast(typing.Optional[typing.List["WebhookNotifications"]], result)

    @builtins.property
    def secret(self) -> typing.Optional[builtins.str]:
        result = self._values.get("secret")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "WebhookConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_amplience.WebhookFilter",
    jsii_struct_bases=[],
    name_mapping={"arguments": "arguments", "type": "type"},
)
class WebhookFilter:
    def __init__(
        self,
        *,
        arguments: typing.List["WebhookFilterArguments"],
        type: builtins.str,
    ) -> None:
        '''
        :param arguments: arguments block.
        :param type: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "arguments": arguments,
            "type": type,
        }

    @builtins.property
    def arguments(self) -> typing.List["WebhookFilterArguments"]:
        '''arguments block.'''
        result = self._values.get("arguments")
        assert result is not None, "Required property 'arguments' is missing"
        return typing.cast(typing.List["WebhookFilterArguments"], result)

    @builtins.property
    def type(self) -> builtins.str:
        result = self._values.get("type")
        assert result is not None, "Required property 'type' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "WebhookFilter(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_amplience.WebhookFilterArguments",
    jsii_struct_bases=[],
    name_mapping={"json_path": "jsonPath", "value": "value"},
)
class WebhookFilterArguments:
    def __init__(
        self,
        *,
        json_path: builtins.str,
        value: typing.List[builtins.str],
    ) -> None:
        '''
        :param json_path: -
        :param value: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "json_path": json_path,
            "value": value,
        }

    @builtins.property
    def json_path(self) -> builtins.str:
        result = self._values.get("json_path")
        assert result is not None, "Required property 'json_path' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def value(self) -> typing.List[builtins.str]:
        result = self._values.get("value")
        assert result is not None, "Required property 'value' is missing"
        return typing.cast(typing.List[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "WebhookFilterArguments(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_amplience.WebhookHeader",
    jsii_struct_bases=[],
    name_mapping={"key": "key", "secret_value": "secretValue", "value": "value"},
)
class WebhookHeader:
    def __init__(
        self,
        *,
        key: builtins.str,
        secret_value: typing.Optional[builtins.str] = None,
        value: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param key: -
        :param secret_value: -
        :param value: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "key": key,
        }
        if secret_value is not None:
            self._values["secret_value"] = secret_value
        if value is not None:
            self._values["value"] = value

    @builtins.property
    def key(self) -> builtins.str:
        result = self._values.get("key")
        assert result is not None, "Required property 'key' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def secret_value(self) -> typing.Optional[builtins.str]:
        result = self._values.get("secret_value")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def value(self) -> typing.Optional[builtins.str]:
        result = self._values.get("value")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "WebhookHeader(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_amplience.WebhookNotifications",
    jsii_struct_bases=[],
    name_mapping={"email": "email"},
)
class WebhookNotifications:
    def __init__(self, *, email: builtins.str) -> None:
        '''
        :param email: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "email": email,
        }

    @builtins.property
    def email(self) -> builtins.str:
        result = self._values.get("email")
        assert result is not None, "Required property 'email' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "WebhookNotifications(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


__all__ = [
    "AmplienceProvider",
    "AmplienceProviderConfig",
    "ContentRepository",
    "ContentRepositoryConfig",
    "ContentRepositoryContentTypes",
    "Webhook",
    "WebhookConfig",
    "WebhookFilter",
    "WebhookFilterArguments",
    "WebhookHeader",
    "WebhookNotifications",
]

publication.publish()
