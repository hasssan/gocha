{% macro message_list(scope, messages) %}{% for msg in messages %}
{% if scope != "none" %}    {% endif %}- {{msg.Subject}} ([{{msg.ID|slice:":10"}}]({{url}}/commit/{{msg.ID}})){% endfor %}{% endmacro %}
# {{appName}} {{version}}
---
{% for type, group in message_group %}
## {{type|title}}
{% for scope, messages in group %}
{% if scope != "none" %}- **{{scope}}:**{% endif %}{{message_list(scope, messages)}}{% endfor %}
{% endfor %}