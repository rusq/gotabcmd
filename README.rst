=================================
Tableau tabcmd wrapper for golang
=================================

.. image:: https://travis-ci.org/rusq/gotabcmd.svg?branch=master
    :target: https://travis-ci.org/rusq/gotabcmd
.. image:: https://codecov.io/gh/rusq/gotabcmd/branch/master/graph/badge.svg
    :target: https://codecov.io/gh/rusq/gotabcmd

Purpose
-------

Automation of tabcmd with golang.

Can be used for containerizing tabcmd related tasks, in conjunction
with the `tabcmd docker image`_.

Restrictions
------------

As this is merely a wrapper around the ``tabcmd`` executable, there
can be only one instance of `Tableau`, due to the `nature`_ of
``tabcmd login/logout`` - and this is guaranteed by ``NewTableau()``.

You can still initialize a variable with ``tb:=&Tableau{}`` if you
must, but `bear in mind`_ that you might face the undesired behaviour.

Current state
-------------

Currently (apart from login/logout functions) only ``refreshextracts``
is supported.

Contributing
------------

Contributions are welcomed.


.. _`tabcmd docker image`: https://github.com/tableau/tableau-docker-samples/blob/2549b9f44be148437602275c598db131b4caaac1/samples/tabcmd/Dockerfile#L1
.. _`nature`: https://onlinehelp.tableau.com/current/server/en-us/tabcmd_cmd.htm#id5fba51c9-5608-4520-8ceb-2caf4846a2be
.. _`bear in mind`: https://i.kym-cdn.com/photos/images/original/001/035/451/6c9.png
