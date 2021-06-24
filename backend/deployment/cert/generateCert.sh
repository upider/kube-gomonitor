#!/bin/sh
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=www -hostname=gomonitor-manager.gomonitor,gomonitor-manager.gomonitor.svc,gomonitor-manager.gomonitor.svc.cluster,gomonitor-manager.gomonitor.svc.cluster.local  webhook-csr.json | cfssljson -bare webhook
caBundle=$(cat ca.pem | base64 | tr -d '\n')
sed -i "s/caBundle.*/caBundle: $caBundle/g" ../deployment.yaml
crt=$(cat webhook.pem | base64 | tr -d '\n')
sed -i "s/tls.crt.*/tls.crt: $crt/g" ../deployment.yaml
key=$(cat webhook-key.pem | base64 | tr -d '\n')
sed -i "s/tls.key.*/tls.key: $key/g" ../deployment.yaml