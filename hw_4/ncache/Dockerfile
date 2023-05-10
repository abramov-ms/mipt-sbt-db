FROM alachisoft/ncache

USER root
RUN apt-get -y update && apt-get -y upgrade
RUN apt-get install -y python3 python-is-python3 pip
RUN pip install ncache-client
