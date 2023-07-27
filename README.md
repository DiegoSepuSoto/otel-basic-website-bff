# Basic Website BFF

#### This API is intended to use in my OpenTelemetry POC project.

The Website BFF API will return the following JSON:

```json
{
  "message": "La orden de John Doe está IN TRANSIT actualmente"
}
```

depending on the customerID sent by the query param.

The request can be made like this:

```bash
curl --request GET \
  --url 'http://localhost:8082/order?orderID=ajnsdoasn123123'
```