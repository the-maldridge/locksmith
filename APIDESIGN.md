# API Design

The design of locksmith is meant to be fairly straightfowards and
meant to be quick and easy to reason about.  This document tries to
explain what's going on and why it works.

## Getting Information on Networks - `/v1/networks/:id`

This endpoint will return a JSON structure containing information
about a particular network.  Since this allows you to tie key
identities to IPs, which would allow you to observe traffic flows, you
must posses a token with the `root` capability for the given network.

Such a request and response looks like this:

```
$ curl -X GET -H "Authorization: Bearer $TOKEN" http://localhost:1323/v1/networks/default | jq .
{
  "Name": "Default Network",
  "ID": "default",
  "Interface": "wg0",
  "Driver": "",
  "ApproveMode": "",
  "ApproveExpiry": 60000000000,
  "ActivateMode": "",
  "ActivateExpiry": 30000000000,
  "PreApproveHooks": null,
  "IPAM": [
    "dummy"
  ],
  "DNS": [
    "8.8.8.8"
  ],
  "AllowedIPs": [
    "10.0.0.0/23"
  ],
  "ApprovalExpirations": null,
  "ActivationExpirations": null,
  "AddressTable": null,
  "StagedPeers": {
    "TXO8ENOkmLCQsM3MlJsnO4q9fyBST7Diu1mtMo1/7zo=": {
      "Owner": "maldridge",
      "OwnerLabel": "",
      "PubKey": "TXO8ENOkmLCQsM3MlJsnO4q9fyBST7Diu1mtMo1/7zo=",
      "Addresses": null,
      "DNS": null,
      "AllowedIPs": null,
      "NetworkPubKey": ""
    }
  },
  "ApprovedPeers": null,
  "ActivePeers": null
}
```

## Adding a New Peer - `/v1/networks/:id/peers`

Adding a peer makes it available for the network to use.  Depending on
the setup of the network, anyone may be able to add a key, or it may
require the `add` permission for the given network.  Once you have
added the key it will by default be staged; again, the network may be
configured to automatically approve and/or activate the new key.

```
$ curl -X POST \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"PubKey":"F3sSyEZ/VSHVurBN3oAAL+Vt5+6/zlnbeiUJwy4kbwU=", "Owner":"maldridge"}' \
    http://localhost:1323/v1/networks/default/peers/disapprove
```

In the above example, the `Owner` is specified for the key.  This
allows the request to submit keys for other users, but doing so
requires the `other` permission for the network.

## Approving an Existing Peer - `/v1/networks/:id/peers/approve`

Prior to being used, a peer must be approved for use.  This will issue
configuration information to the peer, but not return this information
to the caller.  Approving a peer for non-automatic approval networks
requires the `approve` permision.

```
$ curl -X POST \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"PubKey":"F3sSyEZ/VSHVurBN3oAAL+Vt5+6/zlnbeiUJwy4kbwU="}' \
    http://localhost:1323/v1/networks/default/peers/approve
```

## Disapproving an Existing Peer - `/v1/networks/:id/peers/disapprove`

An existing peer's power to use the network can be revoked at any time
by disapproving it.  If the peer is active, disapproving it will force
the activation state to inactive.  Disapproved peers will still be
listed as staged on the network.

```
$ curl -X POST \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"PubKey":"F3sSyEZ/VSHVurBN3oAAL+Vt5+6/zlnbeiUJwy4kbwU="}' \
    http://localhost:1323/v1/networks/default/peers/disapprove
```

## Activating an Approved Peer - `/v1/networks/:id/peers/activate`

Activating a peer adds it to the wireguard interface and it is able to
be used.  Prior to activating the peer, it must be approved for use.
On some networks, peers are automatically approved and automatically
activated.  If you need to manually activate your peer, you must
either be on a network where self activation is enabled, or posess the
`activate` permission.

```
$ curl -X POST \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"PubKey":"F3sSyEZ/VSHVurBN3oAAL+Vt5+6/zlnbeiUJwy4kbwU="}' \
    http://localhost:1323/v1/networks/default/peers/activate
```


## Deactivating an Active Peer - `/v1/networks/:id/peers/deactivate`

Deactivating a peer removes it from the wireguard interface and it
will no longer be able to pass traffic.  The peer will stay approved.

```
$ curl -X POST \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"PubKey":"F3sSyEZ/VSHVurBN3oAAL+Vt5+6/zlnbeiUJwy4kbwU="}' \
    http://localhost:1323/v1/networks/default/peers/deactivate
```
