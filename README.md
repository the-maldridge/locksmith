# Locksmith

Who better to manage keys?

Locksmith bolts on a more "enterprise" friendly management layer on
top of the amazingly capable WireGuard VPN.  The basic idea here is
that keys can be registered via a self-service system, approved
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
appropriately approved user must approve the key.  A network may
specify an approval lifetime which can be used to drive periodic
review of keys.

Once approved the key is ready for use and a peer can request
activation of their key.  A network may also specify that activation
is automatic after approval.  To request activation the peer must send
an activation request to the server after which time the key will be
activated for a configured amount of time.  For some networks this
will be an infinite period of time, but the intention is that keys are
validated for a bit longer than the average work day.

## How it Works

There are two components, there is the locksmith service which handles
lifecycle management for keys and generally provides the access
control, and there is the keyhole service which is a privileged
process that actually modifies the wireguard interface.

Keyhole runs as root on Linux systems in order to modify the
interface, but Locksmith can and should run as an unprivileged user.

## State of the project

All of the critical parts work, but this is not production ready code,
see below.

## Want to Help?

This is an open project, feel free to send your PRs and to get
involved!

The following areas are currently highly desired modules.  Please open
an issue and discuss if you're interested in working on any of these:

  * Tests.  This is arguably security critical code, and so test
    coverage is important.  I do not have the bandwidth right now to
    write good tests as this started out as a toy project to see what
    was possible.  Now that it is worth working on and finishing, some
    tests should be written.  I would accept changes that alter
    interfaces if they do not break overal functionality if they
    improve test coverage.
  * HA Storage backend.  It would be nice to have multiple locksmith
    instances, which would use the same storage backend.
  * Keyhole authentication.  Currently keyhole blindly accepts
    changes.  This is obviously not right.  Some means for locksmith
    to authenticate to keyhole would be ideal.  This is likely just a
    simple token, but I will entertain all options.
  * Frontend webapp.  For the vast majority of users, a frontend
    webapp is critical as it is how they will interact with the
    system.  This will need to be able to generate configuration
    fragments for users, perform administrative functions, and provide
    an interface to request key activation.  I'd like something simple
    here, but I'm not opposed to a VueJS or EmberJS style application.
  * CLI Tool.  Anything the webapp can do the CLI should be able to do
    as well.  I'd like this to be built with cobra so that it can read
    the configuration values in the same way.
