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


class Apikey(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_contentful.Apikey",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        name: builtins.str,
        space_id: builtins.str,
        description: typing.Optional[builtins.str] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param name: -
        :param space_id: -
        :param description: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = ApikeyConfig(
            name=name,
            space_id=space_id,
            description=description,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(Apikey, self, [scope, id, config])

    @jsii.member(jsii_name="resetDescription")
    def reset_description(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDescription", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="accessToken")
    def access_token(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "accessToken"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="spaceIdInput")
    def space_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="descriptionInput")
    def description_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "descriptionInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="description")
    def description(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "description"))

    @description.setter
    def description(self, value: builtins.str) -> None:
        jsii.set(self, "description", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="spaceId")
    def space_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceId"))

    @space_id.setter
    def space_id(self, value: builtins.str) -> None:
        jsii.set(self, "spaceId", value)


@jsii.data_type(
    jsii_type="labd_contentful.ApikeyConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "name": "name",
        "space_id": "spaceId",
        "description": "description",
    },
)
class ApikeyConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        name: builtins.str,
        space_id: builtins.str,
        description: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param name: -
        :param space_id: -
        :param description: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "name": name,
            "space_id": space_id,
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
    def space_id(self) -> builtins.str:
        result = self._values.get("space_id")
        assert result is not None, "Required property 'space_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def description(self) -> typing.Optional[builtins.str]:
        result = self._values.get("description")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ApikeyConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class Asset(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_contentful.Asset",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        archived: builtins.bool,
        asset_id: builtins.str,
        fields: typing.List["AssetFields"],
        locale: builtins.str,
        published: builtins.bool,
        space_id: builtins.str,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param archived: -
        :param asset_id: -
        :param fields: fields block.
        :param locale: -
        :param published: -
        :param space_id: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = AssetConfig(
            archived=archived,
            asset_id=asset_id,
            fields=fields,
            locale=locale,
            published=published,
            space_id=space_id,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(Asset, self, [scope, id, config])

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="archivedInput")
    def archived_input(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "archivedInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="assetIdInput")
    def asset_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "assetIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="fieldsInput")
    def fields_input(self) -> typing.List["AssetFields"]:
        return typing.cast(typing.List["AssetFields"], jsii.get(self, "fieldsInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="localeInput")
    def locale_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "localeInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="publishedInput")
    def published_input(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "publishedInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="spaceIdInput")
    def space_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="archived")
    def archived(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "archived"))

    @archived.setter
    def archived(self, value: builtins.bool) -> None:
        jsii.set(self, "archived", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="assetId")
    def asset_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "assetId"))

    @asset_id.setter
    def asset_id(self, value: builtins.str) -> None:
        jsii.set(self, "assetId", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="fields")
    def fields(self) -> typing.List["AssetFields"]:
        return typing.cast(typing.List["AssetFields"], jsii.get(self, "fields"))

    @fields.setter
    def fields(self, value: typing.List["AssetFields"]) -> None:
        jsii.set(self, "fields", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="locale")
    def locale(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "locale"))

    @locale.setter
    def locale(self, value: builtins.str) -> None:
        jsii.set(self, "locale", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="published")
    def published(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "published"))

    @published.setter
    def published(self, value: builtins.bool) -> None:
        jsii.set(self, "published", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="spaceId")
    def space_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceId"))

    @space_id.setter
    def space_id(self, value: builtins.str) -> None:
        jsii.set(self, "spaceId", value)


@jsii.data_type(
    jsii_type="labd_contentful.AssetConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "archived": "archived",
        "asset_id": "assetId",
        "fields": "fields",
        "locale": "locale",
        "published": "published",
        "space_id": "spaceId",
    },
)
class AssetConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        archived: builtins.bool,
        asset_id: builtins.str,
        fields: typing.List["AssetFields"],
        locale: builtins.str,
        published: builtins.bool,
        space_id: builtins.str,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param archived: -
        :param asset_id: -
        :param fields: fields block.
        :param locale: -
        :param published: -
        :param space_id: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "archived": archived,
            "asset_id": asset_id,
            "fields": fields,
            "locale": locale,
            "published": published,
            "space_id": space_id,
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
    def archived(self) -> builtins.bool:
        result = self._values.get("archived")
        assert result is not None, "Required property 'archived' is missing"
        return typing.cast(builtins.bool, result)

    @builtins.property
    def asset_id(self) -> builtins.str:
        result = self._values.get("asset_id")
        assert result is not None, "Required property 'asset_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def fields(self) -> typing.List["AssetFields"]:
        '''fields block.'''
        result = self._values.get("fields")
        assert result is not None, "Required property 'fields' is missing"
        return typing.cast(typing.List["AssetFields"], result)

    @builtins.property
    def locale(self) -> builtins.str:
        result = self._values.get("locale")
        assert result is not None, "Required property 'locale' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def published(self) -> builtins.bool:
        result = self._values.get("published")
        assert result is not None, "Required property 'published' is missing"
        return typing.cast(builtins.bool, result)

    @builtins.property
    def space_id(self) -> builtins.str:
        result = self._values.get("space_id")
        assert result is not None, "Required property 'space_id' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "AssetConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.AssetFields",
    jsii_struct_bases=[],
    name_mapping={"description": "description", "file": "file", "title": "title"},
)
class AssetFields:
    def __init__(
        self,
        *,
        description: typing.List["AssetFieldsDescription"],
        file: typing.Mapping[builtins.str, builtins.str],
        title: typing.List["AssetFieldsTitle"],
    ) -> None:
        '''
        :param description: description block.
        :param file: -
        :param title: title block.
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "description": description,
            "file": file,
            "title": title,
        }

    @builtins.property
    def description(self) -> typing.List["AssetFieldsDescription"]:
        '''description block.'''
        result = self._values.get("description")
        assert result is not None, "Required property 'description' is missing"
        return typing.cast(typing.List["AssetFieldsDescription"], result)

    @builtins.property
    def file(self) -> typing.Mapping[builtins.str, builtins.str]:
        result = self._values.get("file")
        assert result is not None, "Required property 'file' is missing"
        return typing.cast(typing.Mapping[builtins.str, builtins.str], result)

    @builtins.property
    def title(self) -> typing.List["AssetFieldsTitle"]:
        '''title block.'''
        result = self._values.get("title")
        assert result is not None, "Required property 'title' is missing"
        return typing.cast(typing.List["AssetFieldsTitle"], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "AssetFields(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.AssetFieldsDescription",
    jsii_struct_bases=[],
    name_mapping={"content": "content", "locale": "locale"},
)
class AssetFieldsDescription:
    def __init__(self, *, content: builtins.str, locale: builtins.str) -> None:
        '''
        :param content: -
        :param locale: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "content": content,
            "locale": locale,
        }

    @builtins.property
    def content(self) -> builtins.str:
        result = self._values.get("content")
        assert result is not None, "Required property 'content' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def locale(self) -> builtins.str:
        result = self._values.get("locale")
        assert result is not None, "Required property 'locale' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "AssetFieldsDescription(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.AssetFieldsTitle",
    jsii_struct_bases=[],
    name_mapping={"content": "content", "locale": "locale"},
)
class AssetFieldsTitle:
    def __init__(self, *, content: builtins.str, locale: builtins.str) -> None:
        '''
        :param content: -
        :param locale: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "content": content,
            "locale": locale,
        }

    @builtins.property
    def content(self) -> builtins.str:
        result = self._values.get("content")
        assert result is not None, "Required property 'content' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def locale(self) -> builtins.str:
        result = self._values.get("locale")
        assert result is not None, "Required property 'locale' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "AssetFieldsTitle(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class ContentfulProvider(
    cdktf.TerraformProvider,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_contentful.ContentfulProvider",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        cma_token: builtins.str,
        organization_id: builtins.str,
        alias: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param cma_token: The Contentful Management API token.
        :param organization_id: The organization ID.
        :param alias: Alias name.
        '''
        config = ContentfulProviderConfig(
            cma_token=cma_token, organization_id=organization_id, alias=alias
        )

        jsii.create(ContentfulProvider, self, [scope, id, config])

    @jsii.member(jsii_name="resetAlias")
    def reset_alias(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetAlias", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="cmaTokenInput")
    def cma_token_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "cmaTokenInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="organizationIdInput")
    def organization_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "organizationIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="aliasInput")
    def alias_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "aliasInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="cmaToken")
    def cma_token(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "cmaToken"))

    @cma_token.setter
    def cma_token(self, value: builtins.str) -> None:
        jsii.set(self, "cmaToken", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="organizationId")
    def organization_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "organizationId"))

    @organization_id.setter
    def organization_id(self, value: builtins.str) -> None:
        jsii.set(self, "organizationId", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="alias")
    def alias(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "alias"))

    @alias.setter
    def alias(self, value: typing.Optional[builtins.str]) -> None:
        jsii.set(self, "alias", value)


@jsii.data_type(
    jsii_type="labd_contentful.ContentfulProviderConfig",
    jsii_struct_bases=[],
    name_mapping={
        "cma_token": "cmaToken",
        "organization_id": "organizationId",
        "alias": "alias",
    },
)
class ContentfulProviderConfig:
    def __init__(
        self,
        *,
        cma_token: builtins.str,
        organization_id: builtins.str,
        alias: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param cma_token: The Contentful Management API token.
        :param organization_id: The organization ID.
        :param alias: Alias name.
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "cma_token": cma_token,
            "organization_id": organization_id,
        }
        if alias is not None:
            self._values["alias"] = alias

    @builtins.property
    def cma_token(self) -> builtins.str:
        '''The Contentful Management API token.'''
        result = self._values.get("cma_token")
        assert result is not None, "Required property 'cma_token' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def organization_id(self) -> builtins.str:
        '''The organization ID.'''
        result = self._values.get("organization_id")
        assert result is not None, "Required property 'organization_id' is missing"
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
        return "ContentfulProviderConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class Contenttype(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_contentful.Contenttype",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        display_field: builtins.str,
        field: typing.List["ContenttypeField"],
        name: builtins.str,
        space_id: builtins.str,
        description: typing.Optional[builtins.str] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param display_field: -
        :param field: field block.
        :param name: -
        :param space_id: -
        :param description: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = ContenttypeConfig(
            display_field=display_field,
            field=field,
            name=name,
            space_id=space_id,
            description=description,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(Contenttype, self, [scope, id, config])

    @jsii.member(jsii_name="resetDescription")
    def reset_description(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDescription", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="displayFieldInput")
    def display_field_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "displayFieldInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="fieldInput")
    def field_input(self) -> typing.List["ContenttypeField"]:
        return typing.cast(typing.List["ContenttypeField"], jsii.get(self, "fieldInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="spaceIdInput")
    def space_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="descriptionInput")
    def description_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "descriptionInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="description")
    def description(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "description"))

    @description.setter
    def description(self, value: builtins.str) -> None:
        jsii.set(self, "description", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="displayField")
    def display_field(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "displayField"))

    @display_field.setter
    def display_field(self, value: builtins.str) -> None:
        jsii.set(self, "displayField", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="field")
    def field(self) -> typing.List["ContenttypeField"]:
        return typing.cast(typing.List["ContenttypeField"], jsii.get(self, "field"))

    @field.setter
    def field(self, value: typing.List["ContenttypeField"]) -> None:
        jsii.set(self, "field", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="spaceId")
    def space_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceId"))

    @space_id.setter
    def space_id(self, value: builtins.str) -> None:
        jsii.set(self, "spaceId", value)


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "display_field": "displayField",
        "field": "field",
        "name": "name",
        "space_id": "spaceId",
        "description": "description",
    },
)
class ContenttypeConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        display_field: builtins.str,
        field: typing.List["ContenttypeField"],
        name: builtins.str,
        space_id: builtins.str,
        description: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param display_field: -
        :param field: field block.
        :param name: -
        :param space_id: -
        :param description: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "display_field": display_field,
            "field": field,
            "name": name,
            "space_id": space_id,
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
    def display_field(self) -> builtins.str:
        result = self._values.get("display_field")
        assert result is not None, "Required property 'display_field' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def field(self) -> typing.List["ContenttypeField"]:
        '''field block.'''
        result = self._values.get("field")
        assert result is not None, "Required property 'field' is missing"
        return typing.cast(typing.List["ContenttypeField"], result)

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def space_id(self) -> builtins.str:
        result = self._values.get("space_id")
        assert result is not None, "Required property 'space_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def description(self) -> typing.Optional[builtins.str]:
        result = self._values.get("description")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeField",
    jsii_struct_bases=[],
    name_mapping={
        "id": "id",
        "name": "name",
        "type": "type",
        "disabled": "disabled",
        "items": "items",
        "link_type": "linkType",
        "localized": "localized",
        "omitted": "omitted",
        "required": "required",
        "validation": "validation",
    },
)
class ContenttypeField:
    def __init__(
        self,
        *,
        id: builtins.str,
        name: builtins.str,
        type: builtins.str,
        disabled: typing.Optional[builtins.bool] = None,
        items: typing.Optional[typing.List["ContenttypeFieldItems"]] = None,
        link_type: typing.Optional[builtins.str] = None,
        localized: typing.Optional[builtins.bool] = None,
        omitted: typing.Optional[builtins.bool] = None,
        required: typing.Optional[builtins.bool] = None,
        validation: typing.Optional[typing.List["ContenttypeFieldValidation"]] = None,
    ) -> None:
        '''
        :param id: -
        :param name: -
        :param type: -
        :param disabled: -
        :param items: items block.
        :param link_type: -
        :param localized: -
        :param omitted: -
        :param required: -
        :param validation: validation block.
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "id": id,
            "name": name,
            "type": type,
        }
        if disabled is not None:
            self._values["disabled"] = disabled
        if items is not None:
            self._values["items"] = items
        if link_type is not None:
            self._values["link_type"] = link_type
        if localized is not None:
            self._values["localized"] = localized
        if omitted is not None:
            self._values["omitted"] = omitted
        if required is not None:
            self._values["required"] = required
        if validation is not None:
            self._values["validation"] = validation

    @builtins.property
    def id(self) -> builtins.str:
        result = self._values.get("id")
        assert result is not None, "Required property 'id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def type(self) -> builtins.str:
        result = self._values.get("type")
        assert result is not None, "Required property 'type' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def disabled(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("disabled")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def items(self) -> typing.Optional[typing.List["ContenttypeFieldItems"]]:
        '''items block.'''
        result = self._values.get("items")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldItems"]], result)

    @builtins.property
    def link_type(self) -> typing.Optional[builtins.str]:
        result = self._values.get("link_type")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def localized(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("localized")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def omitted(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("omitted")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def required(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("required")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def validation(self) -> typing.Optional[typing.List["ContenttypeFieldValidation"]]:
        '''validation block.'''
        result = self._values.get("validation")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldValidation"]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeField(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldItems",
    jsii_struct_bases=[],
    name_mapping={"link_type": "linkType", "type": "type", "validation": "validation"},
)
class ContenttypeFieldItems:
    def __init__(
        self,
        *,
        link_type: builtins.str,
        type: builtins.str,
        validation: typing.Optional[typing.List["ContenttypeFieldItemsValidation"]] = None,
    ) -> None:
        '''
        :param link_type: -
        :param type: -
        :param validation: validation block.
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "link_type": link_type,
            "type": type,
        }
        if validation is not None:
            self._values["validation"] = validation

    @builtins.property
    def link_type(self) -> builtins.str:
        result = self._values.get("link_type")
        assert result is not None, "Required property 'link_type' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def type(self) -> builtins.str:
        result = self._values.get("type")
        assert result is not None, "Required property 'type' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def validation(
        self,
    ) -> typing.Optional[typing.List["ContenttypeFieldItemsValidation"]]:
        '''validation block.'''
        result = self._values.get("validation")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldItemsValidation"]], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldItems(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldItemsValidation",
    jsii_struct_bases=[],
    name_mapping={
        "date": "date",
        "dimension": "dimension",
        "file_size": "fileSize",
        "link": "link",
        "mime_type": "mimeType",
        "range": "range",
        "regex": "regex",
        "size": "size",
        "unique": "unique",
    },
)
class ContenttypeFieldItemsValidation:
    def __init__(
        self,
        *,
        date: typing.Optional[typing.List["ContenttypeFieldItemsValidationDate"]] = None,
        dimension: typing.Optional[typing.List["ContenttypeFieldItemsValidationDimension"]] = None,
        file_size: typing.Optional[typing.List["ContenttypeFieldItemsValidationFileSize"]] = None,
        link: typing.Optional[typing.List[builtins.str]] = None,
        mime_type: typing.Optional[typing.List[builtins.str]] = None,
        range: typing.Optional[typing.List["ContenttypeFieldItemsValidationRange"]] = None,
        regex: typing.Optional[typing.List["ContenttypeFieldItemsValidationRegex"]] = None,
        size: typing.Optional[typing.List["ContenttypeFieldItemsValidationSize"]] = None,
        unique: typing.Optional[builtins.bool] = None,
    ) -> None:
        '''
        :param date: date block.
        :param dimension: dimension block.
        :param file_size: file_size block.
        :param link: -
        :param mime_type: -
        :param range: range block.
        :param regex: regex block.
        :param size: size block.
        :param unique: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if date is not None:
            self._values["date"] = date
        if dimension is not None:
            self._values["dimension"] = dimension
        if file_size is not None:
            self._values["file_size"] = file_size
        if link is not None:
            self._values["link"] = link
        if mime_type is not None:
            self._values["mime_type"] = mime_type
        if range is not None:
            self._values["range"] = range
        if regex is not None:
            self._values["regex"] = regex
        if size is not None:
            self._values["size"] = size
        if unique is not None:
            self._values["unique"] = unique

    @builtins.property
    def date(
        self,
    ) -> typing.Optional[typing.List["ContenttypeFieldItemsValidationDate"]]:
        '''date block.'''
        result = self._values.get("date")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldItemsValidationDate"]], result)

    @builtins.property
    def dimension(
        self,
    ) -> typing.Optional[typing.List["ContenttypeFieldItemsValidationDimension"]]:
        '''dimension block.'''
        result = self._values.get("dimension")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldItemsValidationDimension"]], result)

    @builtins.property
    def file_size(
        self,
    ) -> typing.Optional[typing.List["ContenttypeFieldItemsValidationFileSize"]]:
        '''file_size block.'''
        result = self._values.get("file_size")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldItemsValidationFileSize"]], result)

    @builtins.property
    def link(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("link")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def mime_type(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("mime_type")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def range(
        self,
    ) -> typing.Optional[typing.List["ContenttypeFieldItemsValidationRange"]]:
        '''range block.'''
        result = self._values.get("range")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldItemsValidationRange"]], result)

    @builtins.property
    def regex(
        self,
    ) -> typing.Optional[typing.List["ContenttypeFieldItemsValidationRegex"]]:
        '''regex block.'''
        result = self._values.get("regex")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldItemsValidationRegex"]], result)

    @builtins.property
    def size(
        self,
    ) -> typing.Optional[typing.List["ContenttypeFieldItemsValidationSize"]]:
        '''size block.'''
        result = self._values.get("size")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldItemsValidationSize"]], result)

    @builtins.property
    def unique(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("unique")
        return typing.cast(typing.Optional[builtins.bool], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldItemsValidation(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldItemsValidationDate",
    jsii_struct_bases=[],
    name_mapping={"err_message": "errMessage", "max": "max", "min": "min"},
)
class ContenttypeFieldItemsValidationDate:
    def __init__(
        self,
        *,
        err_message: typing.Optional[builtins.str] = None,
        max: typing.Optional[builtins.str] = None,
        min: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param err_message: -
        :param max: -
        :param min: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if err_message is not None:
            self._values["err_message"] = err_message
        if max is not None:
            self._values["max"] = max
        if min is not None:
            self._values["min"] = min

    @builtins.property
    def err_message(self) -> typing.Optional[builtins.str]:
        result = self._values.get("err_message")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def max(self) -> typing.Optional[builtins.str]:
        result = self._values.get("max")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def min(self) -> typing.Optional[builtins.str]:
        result = self._values.get("min")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldItemsValidationDate(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldItemsValidationDimension",
    jsii_struct_bases=[],
    name_mapping={
        "err_message": "errMessage",
        "max_height": "maxHeight",
        "max_width": "maxWidth",
        "min_height": "minHeight",
        "min_width": "minWidth",
    },
)
class ContenttypeFieldItemsValidationDimension:
    def __init__(
        self,
        *,
        err_message: builtins.str,
        max_height: jsii.Number,
        max_width: jsii.Number,
        min_height: jsii.Number,
        min_width: jsii.Number,
    ) -> None:
        '''
        :param err_message: -
        :param max_height: -
        :param max_width: -
        :param min_height: -
        :param min_width: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "err_message": err_message,
            "max_height": max_height,
            "max_width": max_width,
            "min_height": min_height,
            "min_width": min_width,
        }

    @builtins.property
    def err_message(self) -> builtins.str:
        result = self._values.get("err_message")
        assert result is not None, "Required property 'err_message' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def max_height(self) -> jsii.Number:
        result = self._values.get("max_height")
        assert result is not None, "Required property 'max_height' is missing"
        return typing.cast(jsii.Number, result)

    @builtins.property
    def max_width(self) -> jsii.Number:
        result = self._values.get("max_width")
        assert result is not None, "Required property 'max_width' is missing"
        return typing.cast(jsii.Number, result)

    @builtins.property
    def min_height(self) -> jsii.Number:
        result = self._values.get("min_height")
        assert result is not None, "Required property 'min_height' is missing"
        return typing.cast(jsii.Number, result)

    @builtins.property
    def min_width(self) -> jsii.Number:
        result = self._values.get("min_width")
        assert result is not None, "Required property 'min_width' is missing"
        return typing.cast(jsii.Number, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldItemsValidationDimension(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldItemsValidationFileSize",
    jsii_struct_bases=[],
    name_mapping={"err_message": "errMessage", "max": "max", "min": "min"},
)
class ContenttypeFieldItemsValidationFileSize:
    def __init__(
        self,
        *,
        err_message: typing.Optional[builtins.str] = None,
        max: typing.Optional[jsii.Number] = None,
        min: typing.Optional[jsii.Number] = None,
    ) -> None:
        '''
        :param err_message: -
        :param max: -
        :param min: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if err_message is not None:
            self._values["err_message"] = err_message
        if max is not None:
            self._values["max"] = max
        if min is not None:
            self._values["min"] = min

    @builtins.property
    def err_message(self) -> typing.Optional[builtins.str]:
        result = self._values.get("err_message")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def max(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("max")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def min(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("min")
        return typing.cast(typing.Optional[jsii.Number], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldItemsValidationFileSize(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldItemsValidationRange",
    jsii_struct_bases=[],
    name_mapping={"err_message": "errMessage", "max": "max", "min": "min"},
)
class ContenttypeFieldItemsValidationRange:
    def __init__(
        self,
        *,
        err_message: typing.Optional[builtins.str] = None,
        max: typing.Optional[jsii.Number] = None,
        min: typing.Optional[jsii.Number] = None,
    ) -> None:
        '''
        :param err_message: -
        :param max: -
        :param min: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if err_message is not None:
            self._values["err_message"] = err_message
        if max is not None:
            self._values["max"] = max
        if min is not None:
            self._values["min"] = min

    @builtins.property
    def err_message(self) -> typing.Optional[builtins.str]:
        result = self._values.get("err_message")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def max(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("max")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def min(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("min")
        return typing.cast(typing.Optional[jsii.Number], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldItemsValidationRange(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldItemsValidationRegex",
    jsii_struct_bases=[],
    name_mapping={"err_message": "errMessage", "flags": "flags", "pattern": "pattern"},
)
class ContenttypeFieldItemsValidationRegex:
    def __init__(
        self,
        *,
        err_message: typing.Optional[builtins.str] = None,
        flags: typing.Optional[builtins.str] = None,
        pattern: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param err_message: -
        :param flags: -
        :param pattern: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if err_message is not None:
            self._values["err_message"] = err_message
        if flags is not None:
            self._values["flags"] = flags
        if pattern is not None:
            self._values["pattern"] = pattern

    @builtins.property
    def err_message(self) -> typing.Optional[builtins.str]:
        result = self._values.get("err_message")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def flags(self) -> typing.Optional[builtins.str]:
        result = self._values.get("flags")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def pattern(self) -> typing.Optional[builtins.str]:
        result = self._values.get("pattern")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldItemsValidationRegex(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldItemsValidationSize",
    jsii_struct_bases=[],
    name_mapping={"err_message": "errMessage", "max": "max", "min": "min"},
)
class ContenttypeFieldItemsValidationSize:
    def __init__(
        self,
        *,
        err_message: typing.Optional[builtins.str] = None,
        max: typing.Optional[jsii.Number] = None,
        min: typing.Optional[jsii.Number] = None,
    ) -> None:
        '''
        :param err_message: -
        :param max: -
        :param min: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if err_message is not None:
            self._values["err_message"] = err_message
        if max is not None:
            self._values["max"] = max
        if min is not None:
            self._values["min"] = min

    @builtins.property
    def err_message(self) -> typing.Optional[builtins.str]:
        result = self._values.get("err_message")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def max(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("max")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def min(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("min")
        return typing.cast(typing.Optional[jsii.Number], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldItemsValidationSize(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldValidation",
    jsii_struct_bases=[],
    name_mapping={
        "date": "date",
        "dimension": "dimension",
        "file_size": "fileSize",
        "link": "link",
        "mime_type": "mimeType",
        "range": "range",
        "regex": "regex",
        "size": "size",
        "unique": "unique",
    },
)
class ContenttypeFieldValidation:
    def __init__(
        self,
        *,
        date: typing.Optional[typing.List["ContenttypeFieldValidationDate"]] = None,
        dimension: typing.Optional[typing.List["ContenttypeFieldValidationDimension"]] = None,
        file_size: typing.Optional[typing.List["ContenttypeFieldValidationFileSize"]] = None,
        link: typing.Optional[typing.List[builtins.str]] = None,
        mime_type: typing.Optional[typing.List[builtins.str]] = None,
        range: typing.Optional[typing.List["ContenttypeFieldValidationRange"]] = None,
        regex: typing.Optional[typing.List["ContenttypeFieldValidationRegex"]] = None,
        size: typing.Optional[typing.List["ContenttypeFieldValidationSize"]] = None,
        unique: typing.Optional[builtins.bool] = None,
    ) -> None:
        '''
        :param date: date block.
        :param dimension: dimension block.
        :param file_size: file_size block.
        :param link: -
        :param mime_type: -
        :param range: range block.
        :param regex: regex block.
        :param size: size block.
        :param unique: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if date is not None:
            self._values["date"] = date
        if dimension is not None:
            self._values["dimension"] = dimension
        if file_size is not None:
            self._values["file_size"] = file_size
        if link is not None:
            self._values["link"] = link
        if mime_type is not None:
            self._values["mime_type"] = mime_type
        if range is not None:
            self._values["range"] = range
        if regex is not None:
            self._values["regex"] = regex
        if size is not None:
            self._values["size"] = size
        if unique is not None:
            self._values["unique"] = unique

    @builtins.property
    def date(self) -> typing.Optional[typing.List["ContenttypeFieldValidationDate"]]:
        '''date block.'''
        result = self._values.get("date")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldValidationDate"]], result)

    @builtins.property
    def dimension(
        self,
    ) -> typing.Optional[typing.List["ContenttypeFieldValidationDimension"]]:
        '''dimension block.'''
        result = self._values.get("dimension")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldValidationDimension"]], result)

    @builtins.property
    def file_size(
        self,
    ) -> typing.Optional[typing.List["ContenttypeFieldValidationFileSize"]]:
        '''file_size block.'''
        result = self._values.get("file_size")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldValidationFileSize"]], result)

    @builtins.property
    def link(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("link")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def mime_type(self) -> typing.Optional[typing.List[builtins.str]]:
        result = self._values.get("mime_type")
        return typing.cast(typing.Optional[typing.List[builtins.str]], result)

    @builtins.property
    def range(self) -> typing.Optional[typing.List["ContenttypeFieldValidationRange"]]:
        '''range block.'''
        result = self._values.get("range")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldValidationRange"]], result)

    @builtins.property
    def regex(self) -> typing.Optional[typing.List["ContenttypeFieldValidationRegex"]]:
        '''regex block.'''
        result = self._values.get("regex")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldValidationRegex"]], result)

    @builtins.property
    def size(self) -> typing.Optional[typing.List["ContenttypeFieldValidationSize"]]:
        '''size block.'''
        result = self._values.get("size")
        return typing.cast(typing.Optional[typing.List["ContenttypeFieldValidationSize"]], result)

    @builtins.property
    def unique(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("unique")
        return typing.cast(typing.Optional[builtins.bool], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldValidation(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldValidationDate",
    jsii_struct_bases=[],
    name_mapping={"err_message": "errMessage", "max": "max", "min": "min"},
)
class ContenttypeFieldValidationDate:
    def __init__(
        self,
        *,
        err_message: typing.Optional[builtins.str] = None,
        max: typing.Optional[builtins.str] = None,
        min: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param err_message: -
        :param max: -
        :param min: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if err_message is not None:
            self._values["err_message"] = err_message
        if max is not None:
            self._values["max"] = max
        if min is not None:
            self._values["min"] = min

    @builtins.property
    def err_message(self) -> typing.Optional[builtins.str]:
        result = self._values.get("err_message")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def max(self) -> typing.Optional[builtins.str]:
        result = self._values.get("max")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def min(self) -> typing.Optional[builtins.str]:
        result = self._values.get("min")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldValidationDate(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldValidationDimension",
    jsii_struct_bases=[],
    name_mapping={
        "err_message": "errMessage",
        "max_height": "maxHeight",
        "max_width": "maxWidth",
        "min_height": "minHeight",
        "min_width": "minWidth",
    },
)
class ContenttypeFieldValidationDimension:
    def __init__(
        self,
        *,
        err_message: builtins.str,
        max_height: jsii.Number,
        max_width: jsii.Number,
        min_height: jsii.Number,
        min_width: jsii.Number,
    ) -> None:
        '''
        :param err_message: -
        :param max_height: -
        :param max_width: -
        :param min_height: -
        :param min_width: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "err_message": err_message,
            "max_height": max_height,
            "max_width": max_width,
            "min_height": min_height,
            "min_width": min_width,
        }

    @builtins.property
    def err_message(self) -> builtins.str:
        result = self._values.get("err_message")
        assert result is not None, "Required property 'err_message' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def max_height(self) -> jsii.Number:
        result = self._values.get("max_height")
        assert result is not None, "Required property 'max_height' is missing"
        return typing.cast(jsii.Number, result)

    @builtins.property
    def max_width(self) -> jsii.Number:
        result = self._values.get("max_width")
        assert result is not None, "Required property 'max_width' is missing"
        return typing.cast(jsii.Number, result)

    @builtins.property
    def min_height(self) -> jsii.Number:
        result = self._values.get("min_height")
        assert result is not None, "Required property 'min_height' is missing"
        return typing.cast(jsii.Number, result)

    @builtins.property
    def min_width(self) -> jsii.Number:
        result = self._values.get("min_width")
        assert result is not None, "Required property 'min_width' is missing"
        return typing.cast(jsii.Number, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldValidationDimension(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldValidationFileSize",
    jsii_struct_bases=[],
    name_mapping={"err_message": "errMessage", "max": "max", "min": "min"},
)
class ContenttypeFieldValidationFileSize:
    def __init__(
        self,
        *,
        err_message: typing.Optional[builtins.str] = None,
        max: typing.Optional[jsii.Number] = None,
        min: typing.Optional[jsii.Number] = None,
    ) -> None:
        '''
        :param err_message: -
        :param max: -
        :param min: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if err_message is not None:
            self._values["err_message"] = err_message
        if max is not None:
            self._values["max"] = max
        if min is not None:
            self._values["min"] = min

    @builtins.property
    def err_message(self) -> typing.Optional[builtins.str]:
        result = self._values.get("err_message")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def max(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("max")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def min(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("min")
        return typing.cast(typing.Optional[jsii.Number], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldValidationFileSize(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldValidationRange",
    jsii_struct_bases=[],
    name_mapping={"err_message": "errMessage", "max": "max", "min": "min"},
)
class ContenttypeFieldValidationRange:
    def __init__(
        self,
        *,
        err_message: typing.Optional[builtins.str] = None,
        max: typing.Optional[jsii.Number] = None,
        min: typing.Optional[jsii.Number] = None,
    ) -> None:
        '''
        :param err_message: -
        :param max: -
        :param min: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if err_message is not None:
            self._values["err_message"] = err_message
        if max is not None:
            self._values["max"] = max
        if min is not None:
            self._values["min"] = min

    @builtins.property
    def err_message(self) -> typing.Optional[builtins.str]:
        result = self._values.get("err_message")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def max(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("max")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def min(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("min")
        return typing.cast(typing.Optional[jsii.Number], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldValidationRange(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldValidationRegex",
    jsii_struct_bases=[],
    name_mapping={"err_message": "errMessage", "flags": "flags", "pattern": "pattern"},
)
class ContenttypeFieldValidationRegex:
    def __init__(
        self,
        *,
        err_message: typing.Optional[builtins.str] = None,
        flags: typing.Optional[builtins.str] = None,
        pattern: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param err_message: -
        :param flags: -
        :param pattern: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if err_message is not None:
            self._values["err_message"] = err_message
        if flags is not None:
            self._values["flags"] = flags
        if pattern is not None:
            self._values["pattern"] = pattern

    @builtins.property
    def err_message(self) -> typing.Optional[builtins.str]:
        result = self._values.get("err_message")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def flags(self) -> typing.Optional[builtins.str]:
        result = self._values.get("flags")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def pattern(self) -> typing.Optional[builtins.str]:
        result = self._values.get("pattern")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldValidationRegex(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.ContenttypeFieldValidationSize",
    jsii_struct_bases=[],
    name_mapping={"err_message": "errMessage", "max": "max", "min": "min"},
)
class ContenttypeFieldValidationSize:
    def __init__(
        self,
        *,
        err_message: typing.Optional[builtins.str] = None,
        max: typing.Optional[jsii.Number] = None,
        min: typing.Optional[jsii.Number] = None,
    ) -> None:
        '''
        :param err_message: -
        :param max: -
        :param min: -
        '''
        self._values: typing.Dict[str, typing.Any] = {}
        if err_message is not None:
            self._values["err_message"] = err_message
        if max is not None:
            self._values["max"] = max
        if min is not None:
            self._values["min"] = min

    @builtins.property
    def err_message(self) -> typing.Optional[builtins.str]:
        result = self._values.get("err_message")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def max(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("max")
        return typing.cast(typing.Optional[jsii.Number], result)

    @builtins.property
    def min(self) -> typing.Optional[jsii.Number]:
        result = self._values.get("min")
        return typing.cast(typing.Optional[jsii.Number], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "ContenttypeFieldValidationSize(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class Entry(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_contentful.Entry",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        archived: builtins.bool,
        contenttype_id: builtins.str,
        entry_id: builtins.str,
        field: typing.List["EntryField"],
        locale: builtins.str,
        published: builtins.bool,
        space_id: builtins.str,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param archived: -
        :param contenttype_id: -
        :param entry_id: -
        :param field: field block.
        :param locale: -
        :param published: -
        :param space_id: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = EntryConfig(
            archived=archived,
            contenttype_id=contenttype_id,
            entry_id=entry_id,
            field=field,
            locale=locale,
            published=published,
            space_id=space_id,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(Entry, self, [scope, id, config])

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="archivedInput")
    def archived_input(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "archivedInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="contenttypeIdInput")
    def contenttype_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "contenttypeIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="entryIdInput")
    def entry_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "entryIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="fieldInput")
    def field_input(self) -> typing.List["EntryField"]:
        return typing.cast(typing.List["EntryField"], jsii.get(self, "fieldInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="localeInput")
    def locale_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "localeInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="publishedInput")
    def published_input(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "publishedInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="spaceIdInput")
    def space_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="archived")
    def archived(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "archived"))

    @archived.setter
    def archived(self, value: builtins.bool) -> None:
        jsii.set(self, "archived", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="contenttypeId")
    def contenttype_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "contenttypeId"))

    @contenttype_id.setter
    def contenttype_id(self, value: builtins.str) -> None:
        jsii.set(self, "contenttypeId", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="entryId")
    def entry_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "entryId"))

    @entry_id.setter
    def entry_id(self, value: builtins.str) -> None:
        jsii.set(self, "entryId", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="field")
    def field(self) -> typing.List["EntryField"]:
        return typing.cast(typing.List["EntryField"], jsii.get(self, "field"))

    @field.setter
    def field(self, value: typing.List["EntryField"]) -> None:
        jsii.set(self, "field", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="locale")
    def locale(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "locale"))

    @locale.setter
    def locale(self, value: builtins.str) -> None:
        jsii.set(self, "locale", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="published")
    def published(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "published"))

    @published.setter
    def published(self, value: builtins.bool) -> None:
        jsii.set(self, "published", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="spaceId")
    def space_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceId"))

    @space_id.setter
    def space_id(self, value: builtins.str) -> None:
        jsii.set(self, "spaceId", value)


@jsii.data_type(
    jsii_type="labd_contentful.EntryConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "archived": "archived",
        "contenttype_id": "contenttypeId",
        "entry_id": "entryId",
        "field": "field",
        "locale": "locale",
        "published": "published",
        "space_id": "spaceId",
    },
)
class EntryConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        archived: builtins.bool,
        contenttype_id: builtins.str,
        entry_id: builtins.str,
        field: typing.List["EntryField"],
        locale: builtins.str,
        published: builtins.bool,
        space_id: builtins.str,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param archived: -
        :param contenttype_id: -
        :param entry_id: -
        :param field: field block.
        :param locale: -
        :param published: -
        :param space_id: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "archived": archived,
            "contenttype_id": contenttype_id,
            "entry_id": entry_id,
            "field": field,
            "locale": locale,
            "published": published,
            "space_id": space_id,
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
    def archived(self) -> builtins.bool:
        result = self._values.get("archived")
        assert result is not None, "Required property 'archived' is missing"
        return typing.cast(builtins.bool, result)

    @builtins.property
    def contenttype_id(self) -> builtins.str:
        result = self._values.get("contenttype_id")
        assert result is not None, "Required property 'contenttype_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def entry_id(self) -> builtins.str:
        result = self._values.get("entry_id")
        assert result is not None, "Required property 'entry_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def field(self) -> typing.List["EntryField"]:
        '''field block.'''
        result = self._values.get("field")
        assert result is not None, "Required property 'field' is missing"
        return typing.cast(typing.List["EntryField"], result)

    @builtins.property
    def locale(self) -> builtins.str:
        result = self._values.get("locale")
        assert result is not None, "Required property 'locale' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def published(self) -> builtins.bool:
        result = self._values.get("published")
        assert result is not None, "Required property 'published' is missing"
        return typing.cast(builtins.bool, result)

    @builtins.property
    def space_id(self) -> builtins.str:
        result = self._values.get("space_id")
        assert result is not None, "Required property 'space_id' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "EntryConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


@jsii.data_type(
    jsii_type="labd_contentful.EntryField",
    jsii_struct_bases=[],
    name_mapping={"content": "content", "id": "id", "locale": "locale"},
)
class EntryField:
    def __init__(
        self,
        *,
        content: builtins.str,
        id: builtins.str,
        locale: builtins.str,
    ) -> None:
        '''
        :param content: -
        :param id: -
        :param locale: -
        '''
        self._values: typing.Dict[str, typing.Any] = {
            "content": content,
            "id": id,
            "locale": locale,
        }

    @builtins.property
    def content(self) -> builtins.str:
        result = self._values.get("content")
        assert result is not None, "Required property 'content' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def id(self) -> builtins.str:
        result = self._values.get("id")
        assert result is not None, "Required property 'id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def locale(self) -> builtins.str:
        result = self._values.get("locale")
        assert result is not None, "Required property 'locale' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "EntryField(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class Environment(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_contentful.Environment",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        name: builtins.str,
        space_id: builtins.str,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param name: -
        :param space_id: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = EnvironmentConfig(
            name=name,
            space_id=space_id,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(Environment, self, [scope, id, config])

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
    @jsii.member(jsii_name="spaceIdInput")
    def space_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="spaceId")
    def space_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceId"))

    @space_id.setter
    def space_id(self, value: builtins.str) -> None:
        jsii.set(self, "spaceId", value)


@jsii.data_type(
    jsii_type="labd_contentful.EnvironmentConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "name": "name",
        "space_id": "spaceId",
    },
)
class EnvironmentConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        name: builtins.str,
        space_id: builtins.str,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param name: -
        :param space_id: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "name": name,
            "space_id": space_id,
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
    def space_id(self) -> builtins.str:
        result = self._values.get("space_id")
        assert result is not None, "Required property 'space_id' is missing"
        return typing.cast(builtins.str, result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "EnvironmentConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class Locale(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_contentful.Locale",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        code: builtins.str,
        name: builtins.str,
        space_id: builtins.str,
        cda: typing.Optional[builtins.bool] = None,
        cma: typing.Optional[builtins.bool] = None,
        fallback_code: typing.Optional[builtins.str] = None,
        optional: typing.Optional[builtins.bool] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param code: -
        :param name: -
        :param space_id: -
        :param cda: -
        :param cma: -
        :param fallback_code: -
        :param optional: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = LocaleConfig(
            code=code,
            name=name,
            space_id=space_id,
            cda=cda,
            cma=cma,
            fallback_code=fallback_code,
            optional=optional,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(Locale, self, [scope, id, config])

    @jsii.member(jsii_name="resetCda")
    def reset_cda(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetCda", []))

    @jsii.member(jsii_name="resetCma")
    def reset_cma(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetCma", []))

    @jsii.member(jsii_name="resetFallbackCode")
    def reset_fallback_code(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetFallbackCode", []))

    @jsii.member(jsii_name="resetOptional")
    def reset_optional(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetOptional", []))

    @jsii.member(jsii_name="synthesizeAttributes")
    def _synthesize_attributes(self) -> typing.Mapping[builtins.str, typing.Any]:
        return typing.cast(typing.Mapping[builtins.str, typing.Any], jsii.invoke(self, "synthesizeAttributes", []))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="codeInput")
    def code_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "codeInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="id")
    def id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "id"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="nameInput")
    def name_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "nameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="spaceIdInput")
    def space_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="cdaInput")
    def cda_input(self) -> typing.Optional[builtins.bool]:
        return typing.cast(typing.Optional[builtins.bool], jsii.get(self, "cdaInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="cmaInput")
    def cma_input(self) -> typing.Optional[builtins.bool]:
        return typing.cast(typing.Optional[builtins.bool], jsii.get(self, "cmaInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="fallbackCodeInput")
    def fallback_code_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "fallbackCodeInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="optionalInput")
    def optional_input(self) -> typing.Optional[builtins.bool]:
        return typing.cast(typing.Optional[builtins.bool], jsii.get(self, "optionalInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="cda")
    def cda(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "cda"))

    @cda.setter
    def cda(self, value: builtins.bool) -> None:
        jsii.set(self, "cda", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="cma")
    def cma(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "cma"))

    @cma.setter
    def cma(self, value: builtins.bool) -> None:
        jsii.set(self, "cma", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="code")
    def code(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "code"))

    @code.setter
    def code(self, value: builtins.str) -> None:
        jsii.set(self, "code", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="fallbackCode")
    def fallback_code(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "fallbackCode"))

    @fallback_code.setter
    def fallback_code(self, value: builtins.str) -> None:
        jsii.set(self, "fallbackCode", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="optional")
    def optional(self) -> builtins.bool:
        return typing.cast(builtins.bool, jsii.get(self, "optional"))

    @optional.setter
    def optional(self, value: builtins.bool) -> None:
        jsii.set(self, "optional", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="spaceId")
    def space_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceId"))

    @space_id.setter
    def space_id(self, value: builtins.str) -> None:
        jsii.set(self, "spaceId", value)


@jsii.data_type(
    jsii_type="labd_contentful.LocaleConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "code": "code",
        "name": "name",
        "space_id": "spaceId",
        "cda": "cda",
        "cma": "cma",
        "fallback_code": "fallbackCode",
        "optional": "optional",
    },
)
class LocaleConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        code: builtins.str,
        name: builtins.str,
        space_id: builtins.str,
        cda: typing.Optional[builtins.bool] = None,
        cma: typing.Optional[builtins.bool] = None,
        fallback_code: typing.Optional[builtins.str] = None,
        optional: typing.Optional[builtins.bool] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param code: -
        :param name: -
        :param space_id: -
        :param cda: -
        :param cma: -
        :param fallback_code: -
        :param optional: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "code": code,
            "name": name,
            "space_id": space_id,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if cda is not None:
            self._values["cda"] = cda
        if cma is not None:
            self._values["cma"] = cma
        if fallback_code is not None:
            self._values["fallback_code"] = fallback_code
        if optional is not None:
            self._values["optional"] = optional

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
    def code(self) -> builtins.str:
        result = self._values.get("code")
        assert result is not None, "Required property 'code' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def name(self) -> builtins.str:
        result = self._values.get("name")
        assert result is not None, "Required property 'name' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def space_id(self) -> builtins.str:
        result = self._values.get("space_id")
        assert result is not None, "Required property 'space_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def cda(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("cda")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def cma(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("cma")
        return typing.cast(typing.Optional[builtins.bool], result)

    @builtins.property
    def fallback_code(self) -> typing.Optional[builtins.str]:
        result = self._values.get("fallback_code")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def optional(self) -> typing.Optional[builtins.bool]:
        result = self._values.get("optional")
        return typing.cast(typing.Optional[builtins.bool], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "LocaleConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class Space(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_contentful.Space",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        name: builtins.str,
        default_locale: typing.Optional[builtins.str] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param name: -
        :param default_locale: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = SpaceConfig(
            name=name,
            default_locale=default_locale,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(Space, self, [scope, id, config])

    @jsii.member(jsii_name="resetDefaultLocale")
    def reset_default_locale(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetDefaultLocale", []))

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
    @jsii.member(jsii_name="defaultLocaleInput")
    def default_locale_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "defaultLocaleInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="defaultLocale")
    def default_locale(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "defaultLocale"))

    @default_locale.setter
    def default_locale(self, value: builtins.str) -> None:
        jsii.set(self, "defaultLocale", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)


@jsii.data_type(
    jsii_type="labd_contentful.SpaceConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "name": "name",
        "default_locale": "defaultLocale",
    },
)
class SpaceConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        name: builtins.str,
        default_locale: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param name: -
        :param default_locale: -
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
        if default_locale is not None:
            self._values["default_locale"] = default_locale

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
    def default_locale(self) -> typing.Optional[builtins.str]:
        result = self._values.get("default_locale")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "SpaceConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


class WebhookA(
    cdktf.TerraformResource,
    metaclass=jsii.JSIIMeta,
    jsii_type="labd_contentful.WebhookA",
):
    def __init__(
        self,
        scope: constructs.Construct,
        id: builtins.str,
        *,
        name: builtins.str,
        space_id: builtins.str,
        topics: typing.List[builtins.str],
        url: builtins.str,
        headers: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        http_basic_auth_password: typing.Optional[builtins.str] = None,
        http_basic_auth_username: typing.Optional[builtins.str] = None,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
    ) -> None:
        '''
        :param scope: -
        :param id: -
        :param name: -
        :param space_id: -
        :param topics: -
        :param url: -
        :param headers: -
        :param http_basic_auth_password: -
        :param http_basic_auth_username: -
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        '''
        config = WebhookAConfig(
            name=name,
            space_id=space_id,
            topics=topics,
            url=url,
            headers=headers,
            http_basic_auth_password=http_basic_auth_password,
            http_basic_auth_username=http_basic_auth_username,
            count=count,
            depends_on=depends_on,
            lifecycle=lifecycle,
            provider=provider,
        )

        jsii.create(WebhookA, self, [scope, id, config])

    @jsii.member(jsii_name="resetHeaders")
    def reset_headers(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetHeaders", []))

    @jsii.member(jsii_name="resetHttpBasicAuthPassword")
    def reset_http_basic_auth_password(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetHttpBasicAuthPassword", []))

    @jsii.member(jsii_name="resetHttpBasicAuthUsername")
    def reset_http_basic_auth_username(self) -> None:
        return typing.cast(None, jsii.invoke(self, "resetHttpBasicAuthUsername", []))

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
    @jsii.member(jsii_name="spaceIdInput")
    def space_id_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceIdInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="topicsInput")
    def topics_input(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "topicsInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="urlInput")
    def url_input(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "urlInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="version")
    def version(self) -> jsii.Number:
        return typing.cast(jsii.Number, jsii.get(self, "version"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="headersInput")
    def headers_input(
        self,
    ) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], jsii.get(self, "headersInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="httpBasicAuthPasswordInput")
    def http_basic_auth_password_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "httpBasicAuthPasswordInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="httpBasicAuthUsernameInput")
    def http_basic_auth_username_input(self) -> typing.Optional[builtins.str]:
        return typing.cast(typing.Optional[builtins.str], jsii.get(self, "httpBasicAuthUsernameInput"))

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="headers")
    def headers(self) -> typing.Mapping[builtins.str, builtins.str]:
        return typing.cast(typing.Mapping[builtins.str, builtins.str], jsii.get(self, "headers"))

    @headers.setter
    def headers(self, value: typing.Mapping[builtins.str, builtins.str]) -> None:
        jsii.set(self, "headers", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="httpBasicAuthPassword")
    def http_basic_auth_password(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "httpBasicAuthPassword"))

    @http_basic_auth_password.setter
    def http_basic_auth_password(self, value: builtins.str) -> None:
        jsii.set(self, "httpBasicAuthPassword", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="httpBasicAuthUsername")
    def http_basic_auth_username(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "httpBasicAuthUsername"))

    @http_basic_auth_username.setter
    def http_basic_auth_username(self, value: builtins.str) -> None:
        jsii.set(self, "httpBasicAuthUsername", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="name")
    def name(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "name"))

    @name.setter
    def name(self, value: builtins.str) -> None:
        jsii.set(self, "name", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="spaceId")
    def space_id(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "spaceId"))

    @space_id.setter
    def space_id(self, value: builtins.str) -> None:
        jsii.set(self, "spaceId", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="topics")
    def topics(self) -> typing.List[builtins.str]:
        return typing.cast(typing.List[builtins.str], jsii.get(self, "topics"))

    @topics.setter
    def topics(self, value: typing.List[builtins.str]) -> None:
        jsii.set(self, "topics", value)

    @builtins.property # type: ignore[misc]
    @jsii.member(jsii_name="url")
    def url(self) -> builtins.str:
        return typing.cast(builtins.str, jsii.get(self, "url"))

    @url.setter
    def url(self, value: builtins.str) -> None:
        jsii.set(self, "url", value)


@jsii.data_type(
    jsii_type="labd_contentful.WebhookAConfig",
    jsii_struct_bases=[cdktf.TerraformMetaArguments],
    name_mapping={
        "count": "count",
        "depends_on": "dependsOn",
        "lifecycle": "lifecycle",
        "provider": "provider",
        "name": "name",
        "space_id": "spaceId",
        "topics": "topics",
        "url": "url",
        "headers": "headers",
        "http_basic_auth_password": "httpBasicAuthPassword",
        "http_basic_auth_username": "httpBasicAuthUsername",
    },
)
class WebhookAConfig(cdktf.TerraformMetaArguments):
    def __init__(
        self,
        *,
        count: typing.Optional[jsii.Number] = None,
        depends_on: typing.Optional[typing.List[cdktf.ITerraformDependable]] = None,
        lifecycle: typing.Optional[cdktf.TerraformResourceLifecycle] = None,
        provider: typing.Optional[cdktf.TerraformProvider] = None,
        name: builtins.str,
        space_id: builtins.str,
        topics: typing.List[builtins.str],
        url: builtins.str,
        headers: typing.Optional[typing.Mapping[builtins.str, builtins.str]] = None,
        http_basic_auth_password: typing.Optional[builtins.str] = None,
        http_basic_auth_username: typing.Optional[builtins.str] = None,
    ) -> None:
        '''
        :param count: 
        :param depends_on: 
        :param lifecycle: 
        :param provider: 
        :param name: -
        :param space_id: -
        :param topics: -
        :param url: -
        :param headers: -
        :param http_basic_auth_password: -
        :param http_basic_auth_username: -
        '''
        if isinstance(lifecycle, dict):
            lifecycle = cdktf.TerraformResourceLifecycle(**lifecycle)
        self._values: typing.Dict[str, typing.Any] = {
            "name": name,
            "space_id": space_id,
            "topics": topics,
            "url": url,
        }
        if count is not None:
            self._values["count"] = count
        if depends_on is not None:
            self._values["depends_on"] = depends_on
        if lifecycle is not None:
            self._values["lifecycle"] = lifecycle
        if provider is not None:
            self._values["provider"] = provider
        if headers is not None:
            self._values["headers"] = headers
        if http_basic_auth_password is not None:
            self._values["http_basic_auth_password"] = http_basic_auth_password
        if http_basic_auth_username is not None:
            self._values["http_basic_auth_username"] = http_basic_auth_username

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
    def space_id(self) -> builtins.str:
        result = self._values.get("space_id")
        assert result is not None, "Required property 'space_id' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def topics(self) -> typing.List[builtins.str]:
        result = self._values.get("topics")
        assert result is not None, "Required property 'topics' is missing"
        return typing.cast(typing.List[builtins.str], result)

    @builtins.property
    def url(self) -> builtins.str:
        result = self._values.get("url")
        assert result is not None, "Required property 'url' is missing"
        return typing.cast(builtins.str, result)

    @builtins.property
    def headers(self) -> typing.Optional[typing.Mapping[builtins.str, builtins.str]]:
        result = self._values.get("headers")
        return typing.cast(typing.Optional[typing.Mapping[builtins.str, builtins.str]], result)

    @builtins.property
    def http_basic_auth_password(self) -> typing.Optional[builtins.str]:
        result = self._values.get("http_basic_auth_password")
        return typing.cast(typing.Optional[builtins.str], result)

    @builtins.property
    def http_basic_auth_username(self) -> typing.Optional[builtins.str]:
        result = self._values.get("http_basic_auth_username")
        return typing.cast(typing.Optional[builtins.str], result)

    def __eq__(self, rhs: typing.Any) -> builtins.bool:
        return isinstance(rhs, self.__class__) and rhs._values == self._values

    def __ne__(self, rhs: typing.Any) -> builtins.bool:
        return not (rhs == self)

    def __repr__(self) -> str:
        return "WebhookAConfig(%s)" % ", ".join(
            k + "=" + repr(v) for k, v in self._values.items()
        )


__all__ = [
    "Apikey",
    "ApikeyConfig",
    "Asset",
    "AssetConfig",
    "AssetFields",
    "AssetFieldsDescription",
    "AssetFieldsTitle",
    "ContentfulProvider",
    "ContentfulProviderConfig",
    "Contenttype",
    "ContenttypeConfig",
    "ContenttypeField",
    "ContenttypeFieldItems",
    "ContenttypeFieldItemsValidation",
    "ContenttypeFieldItemsValidationDate",
    "ContenttypeFieldItemsValidationDimension",
    "ContenttypeFieldItemsValidationFileSize",
    "ContenttypeFieldItemsValidationRange",
    "ContenttypeFieldItemsValidationRegex",
    "ContenttypeFieldItemsValidationSize",
    "ContenttypeFieldValidation",
    "ContenttypeFieldValidationDate",
    "ContenttypeFieldValidationDimension",
    "ContenttypeFieldValidationFileSize",
    "ContenttypeFieldValidationRange",
    "ContenttypeFieldValidationRegex",
    "ContenttypeFieldValidationSize",
    "Entry",
    "EntryConfig",
    "EntryField",
    "Environment",
    "EnvironmentConfig",
    "Locale",
    "LocaleConfig",
    "Space",
    "SpaceConfig",
    "WebhookA",
    "WebhookAConfig",
]

publication.publish()
