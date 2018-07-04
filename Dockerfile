FROM centos:7
ARG config 
ADD ./logs /home/leisure/logs
ADD ./conf /home/leisure/conf
ADD ./leisure /home/leisure/leisure
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
WORKDIR /home/leisure

RUN mv ${config} ./conf/active.conf

VOLUME ["/home/leisure/logs"]

CMD ["/home/leisure/leisure", "--config=./conf/active.conf"]
