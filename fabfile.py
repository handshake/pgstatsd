"""
Fabfile that controls the deployment of pgstatsd.
"""

import os

from fabric.api import cd, env, local, task, sudo, put, run

PROJECT_ROOT = os.path.relpath( os.path.dirname( __file__ ) )
REMOTE_BIN_DIR = "/opt/pgstatsd"
BUILD_DIR = os.path.join( PROJECT_ROOT, "build" )

env.user = ""
env.hosts = [""]

@task()
def push():
    """
    Pushes a build to the remote server.
    """

    global BUILD_DIR
    global REMOTE_BIN_DIR

    # Compile
    with cd( PROJECT_ROOT ):
        local( "make clean && make GOOS=linux GOARCH=amd64" )

        # Push the binary
        put( "%s/pgstatsd" % BUILD_DIR, "/tmp" )
        sudo( "mv /tmp/pgstatsd %s/pgstatsd" % REMOTE_BIN_DIR )
        sudo( "chmod +x %s/pgstatsd" % REMOTE_BIN_DIR )

    # Restart the service
    #restart()

@task()
def start():
    """
    Start the pgstatsd daemon.
    """

    sudo( "supervisorctl start pgstatsd" )


@task()
def stop():
    """
    Stop the pgstatsd daemon.
    """

    sudo( "supervisorctl stop pgstatsd" )


@task()
def restart():
    """
    Restart the pgstatsd daemon.
    """

    sudo( "supervisorctl restart pgstatsd" )


