# Event Acknowledgement and Delivery Retry

https://github.com/knative/specs/blob/main/specs/eventing/data-plane.md#event-acknowledgement-and-delivery-retry

From the Spec, **Event Acknowledgement**: 

>> Event recipients MUST use the HTTP response code to indicate acceptance of an event. The recipient SHOULD NOT return a response accepting the event until it has handled the event (processed the event or stored it in stable storage). The following response codes are explicitly defined; event recipients MAY also respond with other response codes. A response code not in this table SHOULD be treated as a retriable error.

From the Spec, **Delivery Retry**:

>> Where possible, event senders SHOULD re-attempt delivery of events where the HTTP request returned a retryable status code. It is RECOMMENDED that event senders implement some form of congestion control (such as exponential backoff) and delivery throttling when managing retry timing. Congestion control MAY cause event delivery to fail or MAY include not retrying failed delivery attempts. This specification does not document any specific congestion control algorithm or parameters. Brokers and Channels MUST implement congestion control and MUST implement retries.


# Testing Event Acknowledgement for a Broker

This is a data plane test, so we will be sending events to Knative components, in this case a Broker.

## [Pre] Creating a Broker and obtain the Broker URL

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

# Testing Delivery Retry for a Broker

In order to test Delivery Retry for a Broker, we need the following compoents: 
- A Broker with `.spec.delivery.retry` set to `1`
- A Trigger which refernece the Broker and points to a specific url (`events-fail-once`)
- An event consumer that can fail to cause the event redelivery

## [Pre] Run Event Consumer Kubernetes Service
 
 - Clone `https://github.com/salaboy/knative-conformance-experiment` and `cd` to `events-counter` and then run `ko apply -f config/` 

## [Pre] Creating a Broker, obtain the Broker URL and create a Trigger

```
kubectl apply -f data-plane/event-ack-and-retry/broker.yaml
```

Run the following command to obtain the address of the Broker:

```
kubectl get broker conformance-broker -ojson | jq .status.address.url
```

Create a Trigger that points to a service that fail to accept the event the first time we send it: 

```
kubectl apply -f data-plane/event-ack-and-retry/trigger.yaml
```


## [Test] Emit and Event and Observe Redelivery

```
curl -i -X POST -H "Content-Type: application/json" \
  -H "ce-specversion: 1.0" \
  -H "ce-source: curl-command" \
  -H "ce-type: ConformanceTested" \
  -H "ce-id: 123-abc" \
  -d '{"name":"Salaboy testing conformance"}' \
  http://broker-ingress.knative-eventing.127.0.0.1.nip.io/default/conformance-broker
```

Now I should observe the event failing in the `events-counter` pod and an inmedite redelivery should arrive. 


Congratulations you have tested the **Event Acknowledgement and Delivery Retry Conformance** :metal: !