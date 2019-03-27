# Locksmith

Who better to manage keys?

Locksmith bolts on a more "enterprise" friendly management layer on
top of the amazingly capable WireGuard VPN.  The basic ID here is that
keys can be registered via a self-service system, approved
asyncronously, and activated for limited time spans.

## Life of a Peer

The peer initially posts their key off to a server where they are
pre-approved.  The idea here is that pre-approval saves everyone time
by cancelling registrations outright that won't be approved.  Think
about the cases where someone who won't be approved VPN access has
access to the form to request a key.

After pre-approval, the key is considered staged.  In this stage keys
cannot be activated for use.  To move from staged to approved one of
two conditions must be satisfied.  Either the network must have
auto-approval enabled where keys are truly self service, or an
appropriately approved user must approve the key.

Once approved the key is ready for use and a peer can request
activation of their key.  To request activation the peer must send an
activation request to the server after which time the key will be
activated for a configured amount of time.  For some networks this
will be an infinite period of time, but the intention is that keys are
validated for a bit longer than the average work day.

## Want to Help?

This is an open project, feel free to send your PRs and to get
involved!
