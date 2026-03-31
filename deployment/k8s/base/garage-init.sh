garagepod="garage-0"
garagesvc="garage-0.garage-service"
garagecmd="kubectl exec -i $garagepod -- /garage"
k8secretname="garage-app-secret"
BUCKETNAME="movies"
secretKey=""
accessKey=""


isGarageExist=$(kubectl get pods $garagepod --no-headers 2>&1 | grep $garagepod | grep Ready -i | wc -l)

while [[ "$isGarageExist" < 0 ]] 
do
  sleep 1
done

servcerId=$($garagecmd status | tail -n 1 | awk '{print $1}')
echo $servcerId
echo $garagecmd bucket list
verifyLayout=$($garagecmd bucket list 2>&1 | grep 500 | grep "Layout not ready" )
echo $verifyLayout
if [[ -n $verifyLayout ]]; 
then
  $garagecmd layout assign -z dc1 -c 1G $servcerId 
  $garagecmd layout show
  $garagecmd layout apply --version 1
fi

verifyLayoutAgain=$($garagecmd bucket list 2>&1 | grep 500 | grep "Layout not ready")
if [[ -n $verifyLayoutAgain ]];
then
  echo "error during creating the layout"
  exit 1
fi

verifyBucket=$($garagecmd bucket list 2>&1 | grep $BUCKETNAME -i)
echo $verifyBucket
if [[ -z $verifyBucket ]];
then
  $garagecmd bucket create $BUCKETNAME
  secretKey=$($garagecmd key create movies-key | grep Secret | awk '{print $3}')
  accessKey=$($garagecmd bucket allow --read --write --owner $BUCKETNAME --key movies-key | tail -n 1 | awk '{print $2}')
fi

verifyKey=$($garagecmd key list | tail -n 1 | grep Expiration)
if [[ -z $verifyKey ]]
then
  secretKey=$($garagecmd key info movies-key --show-secret | grep Secret | awk '{print $3}')
  accessKey=$($garagecmd key info movies-key | grep "Key ID" | awk '{print $3}')
fi

isSecretExists=$(kubectl get secrets $k8secretname --no-headers 2>&1 | grep 'not found' )
if [[ -n $isSecretExists ]] 
then
  touch /tmp/.secrets 
  echo "AWS_ACCESS_KEY: $accessKey" > /tmp/.secrets 
  echo "AWS_SECRET_KEY: $secretKey" >> /tmp/.secrets 
  echo "AWS_URL: http://$garagesvc:3900" >> /tmp/.secrets 
  echo "BUCKETNAME: $BUCKETNAME" >> /tmp/.secrets 

  kubectl create secret generic $k8secretname --from-file=/tmp/.secrets 
else
  touch /tmp/.secrets  
  echo "AWS_ACCESS_KEY= $accessKey" > /tmp/.secrets 
  echo "AWS_SECRET_KEY= $secretKey" >> /tmp/.secrets 
  echo "AWS_URL= http://$garagesvc:3900" >> /tmp/.secrets 
  echo "BUCKETNAME= $BUCKETNAME" >> /tmp/.secrets 
  kubectl create secret generic $k8secretname --from-file=/tmp/.secrets --dry-run=client -o yaml | kubectl apply -f -
fi 


echo $secretKey
echo $accessKey