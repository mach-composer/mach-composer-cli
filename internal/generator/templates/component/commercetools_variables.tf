ct_project_key    = {{ site.Commercetools.ProjectKey|tf }}
ct_api_url        = {{ site.Commercetools.APIURL|tf }}
ct_auth_url       = {{ site.Commercetools.TokenURL|tf }}

ct_stores = {
    {% for store in site.Commercetools.Stores %}
    {{ store.Key }} =  {
        key = {{ store.Key|tf }}
        variables = {
            {% for key, value in component.StoreVariables|get:store.Key %}
            {{ key }} = {{ value|tfvalue }}
            {% endfor %}
        }
        secrets = {
            {% for key, value in component.StoreSecrets|get:store.Key %}
            {{ key }} = {{ value|tfvalue }}
            {% endfor %}
        }
    }
    {% endfor %}
}
