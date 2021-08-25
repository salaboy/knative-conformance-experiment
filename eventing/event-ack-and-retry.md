# Event Acknowledgement and Delivery Retry

https://github.com/knative/specs/blob/main/specs/eventing/data-plane.md#event-acknowledgement-and-delivery-retry

"Event recipients MUST use the HTTP response code to indicate acceptance of an event. The recipient SHOULD NOT return a response accepting the event until it has handled the event (processed the event or stored it in stable storage). The following response codes are explicitly defined; event recipients MAY also respond with other response codes. A response code not in this table SHOULD be treated as a retriable error."




## Testing Event Acknowledgement for a Broker

This is a data plane test, so we will be sending events to Knative components, in this case a Broker

## [Pre] Creating a Broker: 

```
kubectl apply -f data-plane/event-ack-and-retry/broker.yaml
```

Run the following command to obtain the address of the Broker:

```
kubectl get broker conformance-broker -ojson | jq .status.address.url
```

## [Test] Response Code to a Valid CloudEvent

Emit a CloudEvent to the broker, replace the URL of the Broker with the obtained in the previous step:

```
curl -i -X POST -H "Content-Type: application/json" \
  -H "ce-specversion: 1.0" \
  -H "ce-source: curl-command" \
  -H "ce-type: ConformanceTested" \
  -H "ce-id: 123-abc" \
  -d '{"name":"Salaboy testing conformance"}' \
  http://broker-ingress.knative-eventing.127.0.0.1.nip.io/default/conformance-broker 
```

You should obtain a `202 Accepted` HTTP code: 

```
HTTP/1.1 202 Accepted
date: Wed, 25 Aug 2021 11:37:33 GMT
content-length: 0
x-envoy-upstream-service-time: 1
server: envoy
```

## [Test] Response Code to wrong CloudEvent Spec Version

Emit a CloudEvent to the broker with a wrong `ce-specversion` value:

```
curl -i -X POST -H "Content-Type: application/json" \
  -H "ce-specversion: 3.0" \
  -H "ce-source: curl-command" \
  -H "ce-type: ConformanceTested" \
  -H "ce-id: 123-abc" \
  -d '{"name":"Salaboy testing conformance"}' \
  http://broker-ingress.knative-eventing.127.0.0.1.nip.io/default/conformance-broker
 ``` 

You should obtain a `400 bad Request` HTTP Code: 

```
HTTP/1.1 400 Bad Request
date: Wed, 25 Aug 2021 11:39:38 GMT
content-length: 0
x-envoy-upstream-service-time: 0
server: envoy
```



Congratulations you have tested the **Event Acknowledgement and Delivery Retry Conformance** :metal: !