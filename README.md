# log
golang log

对golang自带的原生log进行了封装，增加了如下三个功能: 
1、每条日志都需要带一个id
2、记录日志调用处所在的文件与行号
3、如果日志文件被删除，会自动恢复
