#!/bin/bash
upload_server="jmeter-collector.servicecomb"
rm -rf ./dashboard

if [ "$REPORT_UPLOAD_SERVER" != "" ]; then
	echo "Set upload_server to " $REPORT_UPLOAD_SERVER
	upload_server=$REPORT_UPLOAD_SERVER
fi

echo $upload_server

now_str=$(date +"%Y-%m-%d-%H-%M-%S")
echo $now_str
jmeter -g ./result.jtl -o dashboard
tar zcf ${now_str}_upload.tgz ./result.jtl ./dashboard/
curl -X POST --data-binary @./"${now_str}"_upload.tgz -H 'testname: spring-demo' http://$upload_server/upload
