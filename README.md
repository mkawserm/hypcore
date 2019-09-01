# HypCore
Hyper Core is a small reusable high performance websocket and graphql focused golang package




# NOTE
================================================
1. Update connection track table on linux 
    
    echo 100000000 >  /proc/sys/net/netfilter/nf_conntrack_max

2. Run HypCore executable from the root user
    to handle 1 million or more open connections