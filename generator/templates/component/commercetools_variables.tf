ct_project_key    = {{ site.Commercetools.ProjectKey|tf }}
ct_api_url        = {{ site.Commercetools.ApiURL|tf }}
ct_auth_url       = {{ site.Commercetools.TokenURL|tf }}

ct_stores = {
    {% for store in site.Commercetools.Stores %}
    {{ store.Key }} =  {
    key = {{ store.Key|tf }}
    variables = {
        {% for key, value in component.store_variables|get:store.key %}
        {{ key }} = {{ value|tfvalue }}
        {% endfor %}
    }
    secrets = {
        {% for key, value in component.store_secrets|get:store.key %}
        {{ key }} = {{ value|tfvalue }}
        {% endfor %}
    }
    {% endfor %}
}
